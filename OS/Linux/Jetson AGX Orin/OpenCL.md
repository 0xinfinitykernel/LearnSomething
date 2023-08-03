# Binary packages

## PoCL with CUDA driver support for Linux x86_64, aarch64 and ppc64le can be found on conda-forge distribution and can be installed with
```
wget "https://github.com/conda-forge/miniforge/releases/latest/download/Mambaforge-$(uname)-$(uname -m).sh"
bash Mambaforge-$(uname)-$(uname -m).sh   # install mambaforge
```
## To install pocl with cuda driver
```
mamba install pocl-cuda
```
##  To install all drivers
```
mamba install pocl
```
### Fix some problems
```
#make install 或二进制安装 只是安装库文件，不会复制icd到标准目录。如果导致clGetPlatform时获取不到平台时。手动拷贝：
sudo mkdir -p /etc/OpenCL/vendors
sudo cp /usr/local/etc/OpenCL/vendors/pocl.icd /etc/OpenCL/vendors/
```

More information [GitHub](https://github.com/pocl/pocl)