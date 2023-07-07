# 实时补丁更新方法
安装依赖
```depending on
sudo apt-get update
sudo apt-get install build-essential bc curl ca-certificates fakeroot gnupg2 libssl-dev lsb-release libelf-dev bison flex rsync zstd dwarves libncurses-dev
```

# 获取内核原代码
```linux & patch
wget https://mirrors.edge.kernel.org/pub/linux/kernel/v5.x/linux-5.15.119.tar.xz
wget https://mirrors.edge.kernel.org/pub/linux/kernel/projects/rt/5.15/patch-5.15.119-rt65.patch.gz
```

# 解压下载的文件
```
tar -xvf linux-5.15.119.tar.xz
gzip -d patch-5.15.119-rt65.patch.gz
cd linux-5.15.119
```

# 安装补丁
```
patch -p1 < ../patch-5.15.119-rt65.patch
```

# 复制系统当前内核的.config文件
```
cp /boot/config-5.11.0-41-generic .config
```

# 调用图形化界面，设置.config文件
```
make menuconfig
```

# ARM开启 Fully Preemptible Kernel(RT)
```
vi .config
# CONFIG_KVM is not set

-----------------------WHY DO IT-----------------------
Changes since v5.6.19-rt12:

  - Rebase to v5.9-rc2

  - The seqcount related patches have been replaced on top of the
    seqcount series by Ahmed S. Darwis which landed mainline. 

  - The posix-timer patches have been dropped because upstream changes
    cover all of was needed on RT's side. As a result RT relies on
    HAVE_POSIX_CPU_TIMERS_TASK_WORK. This is provided only by x86.
    The RT patch provides this option for ARM/ARM64/POWERPC as long as
    KVM is disabled. The reason is that the task work must be handled
    before KVM returns to guest.
https://lore.kernel.org/linux-rt-users/20200824154605.v66t2rsxobt3r5jg@linutronix.de/
```

# 界面改动的地方
```
General setup

Preemption Model (Voluntary Kernel Preemption(Desktop))
—[x] Fully Preemptible Kernel(RT)

```

# gedit .config
```
修改.config文件，搜索关键词，将

CONFIG_MODULE_SIG_ALL

CONFIG_MODULE_SIG_KEY

CONFIG_SYSTEM_TRUSTED_KEYS

CONFIG_SYSTEM_REVOCATION_LIST

CONFIG_SYSTEM_REVOCATION_KEYS

五项注释掉，最后把CONFIG_DEBUG_INFO=y去掉，不然新内核带debug信息超大。

```

# 编译和安装新内核
```
make -j$(nproc)
sudo make modules_install -j$(nproc)
sudo make install -j$(nproc)
```

# 更新引导加载器并重启
```
sudo update-grub
sudo reboot
```