# The Problem
After a reboot the system will not start or activate one of the LVM volume groups. Trying to activate the group manually produces the message:

Couldn't find device with uuid '[UUID]'
# The Solution
* The string LABELONE at the metadata location marks the device as being part of an LVM device. Without this string, lvm3 will not attempt to utilize the device as a physical volume. The physical device metadata can be overwritten because of a system error or a deliberate manual action.

* For Linux versions 4 or 5, the default metadata area is 192Kb. For Linux version 6, the default metadata is 1Mb. Prior to attempting any recovery or repair activities, you are strongly encouraged to make a backup copy of this area:

```
# /bin/dd if=/dev/xvdd of=/root/xvdd.metadata bs=1K count=[192_or_1024]
```
To check for the LVM signature, do this:

```
# /bin/strings /root/xvdd.metadata | /bin/fgrep LABELONE
```
If no output is produced, the metadata is corrupt.

Review the information in:

/etc/lvm/backup
/etc/lvm/archives
for changes. If these have been changed or are inconsistent, it is possible for the restoration activity to corrupt the entire volume data. Verify that the physical device or LUN is still available, or still presented to this server.

Be diligent to make regular backup for all LVM volumes. While they may be recoverable, it is also possible for a misconfiguration to completely corrupt the entire dataset.

1. Note that this item is important only if you also have any multipathed devices on your system. If you do not use LVM and multipathing on the same server, you may safely skip this item.

During system startup, the LVM subsystem is notified each time a block device, such as a disk drive or LUN, so that the device can be used to construct an LVM volume. This is an asynchronous process; there is no guarantee that the devices will be discovered in the same order each time the sytem is booted. This means that the physical paths of a multipathed device are likely to be discovered before the composite logical device is complete, leading to the physical path being claimed by the LVM subsystem before the multipath subsystem is offered the device. There are two undesirable results of this condition:

a. Only a single path to the multipath device gets utilized and claimed by LVM, leaving the system vulnerable to a single-point failure causing catastrophic loss of connectivity to the storage.

b. Since LVM gains exclusive ownership of the physical path, the multipathing layer reports the device as busy and cannot build the multipath device. This leaves the storage subject to a single-point failure preventing access to storage. The solution is to force the LVM to consider only those block devices which are to actually hold part of an LVM volume. The way to do this is to look into the /etc/lvm/lvm.conf file:

```
# /bin/fgrep -n -e 's/#.*//' -e '/filter/p' /etc/lvm/lvm.conf
filter = [ "a/.*/" ]
```
If your output looks like the above, you have an LVM configuration problem that will likely break any multipath devices, if you use them on your system. Changing this parameter is outside the scope of this note. We will proceed assuming you have corrected this value.

2. Sometimes the storage which holds the LVM data is slow to be recognized and may be successfully mounted if accessed after the system has stabilized. To begin, we will inventory the available block devices and determine the UUID’s:

```
# /sbin/vgscan
Reading all physical volumes.  This may take a while...
Couldn't fine device with uuid  70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
Found volume group "data_vg" using metadata type lvm2
```
Now that we have the UUID causing the issue, we must find the associated device:

```
# /sbin/pvs -o +uuid
Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  PV             VG      Fmt  Attr PSize   PFree   PV UUID
  /dev/xvdc      data_vg lvm2 a--  996.00M      0  VrVT1L-CTcT-9Nn9-oIAx-BnEA-X7sv-vJO6RE
  /dev/xvde      data_vg lvm2 a--  996.00M 428.00M tGIqvd-lsYv-7JmV-1bfD-t7BL-HaGi-rmIYW0
  unknown device data_vg lvm2 a-m  996.00M      0  70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt
```
3. We can view the distribution of the logical volumes across the physical devices like this:

```
# /sbin/lvs -o +devices
Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  LV         VG      Attr   LSize Origin Snap%  Move Log Copy%  Convert Devices
  data_vg_lv data_vg -wi-a- 2.50G                                       /dev/xvdc(0)
  data_vg_lv data_vg -wi-a- 2.50G                                       unknown device(0)
  data_vg_lv data_vg -wi-a- 2.50G                                       /dev/xvde(0) =
 ```
4. Try to activate the volume group:

```
# /sbin/vgchange -a y data_vg
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Refusing activation of partial LV data_vg_lv. Use --partial to override.
  1 logical volume(s) in volume group "data_vg" now active
```
5. Try to reduce the volume group and remove the missing device:

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
6. With the missing device eliminated from the group, the LVM device should activate:

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
7. We can try to restore the volume group using the information stored in the /etc/lvm/archive/ directory:

```
# /sbin/vgcfgrestore -f data_vg_00003-1023778751.vg data_vg
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Cannot restore Volume Group data_vg with 1 PVs marked as missing.
  Restore failed.
```
8. Trying to overwrite or resqore the device information based on the UUID settings derived from the volume group information:

```
# /sbin/pvcreate --restorefile /etc/lvm/archive/data_vg_00003-1023778751.vg --uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt /dev/xvdd
  Couldn't find device with uuid 70FBaa-3QKh-HTAF-gUzZ-u3mu-2RRs-hI3BIt.
  Writing physical volume data to disk "/dev/xvdd"
  Physical volume "/dev/xvdd" successfully created

```  
<u> If the above operation fails, execute the following command to clear the metadata and then perform the operation in step 8：</u>

```
wipefs -a data_vg
```

9. Open the volume group information, for example /etc/lvm/archive/data_vg_00003-1023778751.vg using a text editor and remove the “MISSING” string from the flags entry so it looks like this:

```
flags = [ ]
```
10. Restore the LVM using this modified entry:

```
# /sbin/ vgcfgrestore -f /etc/lvm/archive/data_vg_00003-1023778751.vg data_vg
  Restored volume group data_vg
```
11. Verify the environment:

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
12. Verify the LVM availability:

```
# /sbin/lvs -o +devices
  LV         VG      Attr   LSize Origin Snap%  Move Log Copy%  Convert Devices
  data_vg_lv data_vg -wi--- 2.50G                                       /dev/xvdc(0)
  data_vg_lv data_vg -wi--- 2.50G                                       /dev/xvdd(0)
  data_vg_lv data_vg -wi--- 2.50G                                       /dev/xvde(0)
```
