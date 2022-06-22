# LVM引导不起来UUID故障解决

## 问题
重新引导后，系统将无法启动或激活其中一个 LVM 卷组。尝试手动激活组会产生以下消息：
Couldn't find device with uuid '[UUID]'

# 解决方法
* 元数据位置处的字符串 LABELONE 将设备标记为 LVM 设备的一部分。如果没有此字符串，lvm3 将不会尝试将设备用作物理卷。物理设备元数据可能会由于系统错误或故意的手动操作而被覆盖。
* 对于 Linux 版本 4 或 5，默认元数据区域为 192Kb。对于 Linux 版本 6，默认元数据为 1Mb。在尝试任何恢复或修复活动之前，强烈建议您制作此区域的备份副本
```
# /bin/dd if=/dev/xvdd of=/root/xvdd.metadata bs=1K count=[192_or_1024]
```
要检查 LVM 签名，请执行以下操作：
```
# /bin/strings /root/xvdd.metadata | /bin/fgrep LABELONE
```
如果未生成任何输出，则元数据已损坏

查看以下信息：

/etc/lvm/backup

/etc/lvm/archives

如果这些内容已更改或不一致，则还原活动可能会损坏整个卷数据。验证物理设备或 LUN 是否仍然可用，或者仍提供给此服务器。

请努力为所有 LVM 卷进行定期备份。虽然它们可能是可恢复的，但配置错误也可能完全损坏整个数据集。

1) 请注意，仅当您的系统上还有任何多路径设备时，此项目才重要。如果不在同一台服务器上使用 LVM 和多路径，则可以安全地跳过此项。

在系统启动期间，每次块设备（如磁盘驱动器或 LUN）时，LVM 子系统都会收到通知，以便该设备可用于构建 LVM 卷。这是一个异步过程;无法保证每次启动系统时都会以相同的顺序发现设备。这意味着多路径设备的物理路径很可能在复合逻辑设备完成之前被发现，从而导致在为多路径子系统提供设备之前，LVM 子系统声明了物理路径。这种情况有两个不良结果：

 a. LVM 仅使用和声明到多路径设备的单个路径，使系统容易受到单点故障的影响，从而导致与存储的连接灾难性丢失。

 b. 由于 LVM 获得物理路径的独占所有权，因此多路径层将设备报告为忙碌，无法构建多路径设备。这会使存储受到单点故障的影响，从而阻止对存储的访问。解决方案是强制 LVM 仅考虑那些实际保存 LVM 卷一部分的块设备。执行此操作的方法是查看 /etc/lvm/lvm.conf 文件：
```
# /bin/fgrep -n -e 's/#.*//' -e '/filter/p' /etc/lvm/lvm.conf
filter = [ "a/.*/" ]
```
如果您的输出与上述内容类似，则表示您遇到了 LVM 配置问题，如果您在系统上使用多路径设备，则可能会破坏这些设备。更改此参数超出了本说明的范围。假设您已更正此值，我们将继续进行。

2) 有时，保存 LVM 数据的存储识别速度很慢，如果在系统稳定后访问，则可能会成功挂载。首先，我们将清点可用的块设备并确定UUID：
```
# /sbin/vgscan
Reading all physical volumes.  This may take a while...
Couldn't fine device with uuid  70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
Found volume group "data_vg" using metadata type lvm2
```
现在我们有导致问题的UUID，我们必须找到关联的设备：
```
# /sbin/pvs -o +uuid
Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  PV             VG      Fmt  Attr PSize   PFree   PV UUID
  /dev/xvdc      data_vg lvm2 a--  996.00M      0  VrVT1L-CTcT-9Nn9-oIAx-BnEA-X7sv-vJO6RE
  /dev/xvde      data_vg lvm2 a--  996.00M 428.00M tGIqvd-lsYv-7JmV-1bfD-t7BL-HaGi-rmIYW0
  unknown device data_vg lvm2 a-m  996.00M      0  70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt
```
3) 我们可以查看逻辑卷在物理设备之间的分布情况，如下所示：
```
# /sbin/lvs -o +devices
Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  LV         VG      Attr   LSize Origin Snap%  Move Log Copy%  Convert Devices
  data_vg_lv data_vg -wi-a- 2.50G                                       /dev/xvdc(0)
  data_vg_lv data_vg -wi-a- 2.50G                                       unknown device(0)
  data_vg_lv data_vg -wi-a- 2.50G                                       /dev/xvde(0) =
```
4) 尝试激活卷组：
```
# /sbin/vgchange -a y data_vg
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Refusing activation of partial LV data_vg_lv. Use --partial to override.
  1 logical volume(s) in volume group "data_vg" now active
```
5) 尝试缩小卷组并删除丢失的设备：
```
# /sbin/vgreduce --removemissing data_vg
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  WARNING: Partial LV data_vg_lv needs to be repaired or removed.
  WARNING: There are still partial LVs in VG data_vg.
  To remove them unconditionally use: vgreduce --removemissing --force.
  Proceeding to remove empty missing PVs.
# /sbin/vgreduce --removemissing data_vg --force
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Removing partial LV data_vg_lv.
  Logical volume "data_vg_lv" successfully removed
  Wrote out consistent volume group data_vg
```
6) 从组中消除丢失的设备后，LVM 设备应激活：
```
# /sbin/pvs
  PV         VG      Fmt  Attr PSize   PFree
  /dev/xvdc  data_vg lvm2 a--  996.00M 996.00M
  /dev/xvde  data_vg lvm2 a--  996.00M 996.00M
# /sbin/lvs -o +devices
#
# /sbin/vgscan
  Reading all physical volumes.  This may take a while...
  Found volume group "data_vg" using metadata type lvm2
```

```
# /sbin/vgdisplay
  --- Volume group ---
  VG Name               data_vg
  System ID
  Format                lvm2
  Metadata Areas        2
  Metadata Sequence No  5
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                0
  Open LV               0
  Max PV                0
  Cur PV                2
  Act PV                2
  VG Size               1.95 GB
  PE Size               4.00 MB
  Total PE              498
  Alloc PE / Size       0 / 0
  Free  PE / Size       498 / 1.95 GB
  VG UUID               yTOvvd-ZjUe-gXP0-41BT-qUIk-8uPR-lpr9Pw
```
7) 我们可以尝试使用存储在 /etc/lvm/archive/ 目录中的信息来还原卷组：
```
# /sbin/vgcfgrestore -f data_vg_00003-1023778751.vg data_vg
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Cannot restore Volume Group data_vg with 1 PVs marked as missing.
  Restore failed.
```
8) 尝试根据从卷组信息派生的 UUID 设置覆盖或重新查询设备信息：
```
# /sbin/pvcreate --restorefile /etc/lvm/archive/data_vg_00003-1023778751.vg --uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt /dev/xvdd
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Writing physical volume data to disk "/dev/xvdd"
  Physical volume "/dev/xvdd" successfully created
```
如果上面操作失败，执行下面命令清除元数据后再执行8操作：
```
wipefs -a data_vg
```
9) 使用文本编辑器打开卷组信息，例如 /etc/lvm/archive/data_vg_00003-1023778751.vg，然后从标志条目中删除"MISSING"字符串，使其如下所示：
```
flags = [ ]
```
10) 使用以下修改后的条目恢复 LVM：
```
# /sbin/ vgcfgrestore -f /etc/lvm/archive/data_vg_00003-1023778751.vg data_vg
  Restored volume group data_vg
```
11) 验证环境：
```
# /sbin/vgscan
  Reading all physical volumes.  This may take a while...
  Found volume group "data_vg" using metadata type lvm2
```

```
# /sbin/ pvs -o +uuid
  PV         VG      Fmt  Attr PSize   PFree   PV UUID
  /dev/xvdc  data_vg lvm2 a--  996.00M      0  VrVT1L-CTcT-9Nn9-oIAx-BnEA-X7sv-vJO6RE
  /dev/xvdd  data_vg lvm2 a--  996.00M      0  70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt
  /dev/xvde  data_vg lvm2 a--  996.00M 428.00M tGIqvd-lsYv-7JmV-1bfD-t7BL-HaGi-rmIYW0
```  
12) 验证 LVM 可用性：
```
# /sbin/lvs -o +devices
  LV         VG      Attr   LSize Origin Snap%  Move Log Copy%  Convert Devices
  data_vg_lv data_vg -wi--- 2.50G                                       /dev/xvdc(0)
  data_vg_lv data_vg -wi--- 2.50G                                       /dev/xvdd(0)
  data_vg_lv data_vg -wi--- 2.50G                                       /dev/xvde(0)
```