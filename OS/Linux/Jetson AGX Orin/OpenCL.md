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
## 二进制安装的方式,多了SPIR，SPIR-V，如果不想自行再次编译的,推荐这种安装方式.
## make install 或二进制安装 只是安装库文件，不会复制icd到标准目录。如果导致clGetPlatform时获取不到平台时。手动拷贝：
sudo mkdir -p /etc/OpenCL/vendors
sudo cp /usr/local/etc/OpenCL/vendors/pocl.icd /etc/OpenCL/vendors/
```

More information [GitHub](https://github.com/pocl/pocl)

# Source Code Install
use llvm 10 compile pocl 3.0, if higher pocl version you should be higher llvm version
pocl 4.0 use llvm 15.0.7 build
```
sudo apt-get install clang

## 执行后这里自动安装的版本是clang10+llvm10，但是不会安装libclang-cpp-dev,需要自己安装对应的包（dpkg -l | grep 看一下,注意和自己安装的clang+vllm版本匹配）

sudo apt-get install libclang-cpp10-dev 

```

```
sudo apt-get install libncurses5
#安装后还是提示找不到，那就手动创建个软链接
sudo ln -s /lib/aarch64-linux-gnu/libtinfo.so.5 /lib/aarch64-linux-gnu/libtinfo.so

```

## ruby (编译ocl-icd用)
```
sudo apt install ruby

```
### 下载
ocl-icd：https://github.com/OCL-dev/ocl-icd/tags

### 解压
```
tar -zxvf ocl-icd-2.3.1.tar.gz

```
### 编译
```
cd ocl-icd-2.3.1
./bootstrap
./configure
make
make check
sudo make install
#把OpenCL的头文件复制到标准目录中
sudo cp ~/work/ocl-icd-2.3.1/khronos-headers/* /usr/include

```
## pocl

### 下载
pocl_3.0: http://portablecl.org/download.html

### 解压
```
tar -zxvf pocl-3.0.tar.gz
```

### 编译
* 安装依赖
```
sudo apt-get install valgrind
#以下依赖基本都有，除了hwloc和cmake还有clinfo。hwloc和cmake必不可少。
sudo apt-get install gcc patch hwloc cmake git pkg-config make clinfo

```
* 在解压的目录pocl-3.0中
```
mkdir build
cd build
cmake -DENABLE_CUDA=ON -DCLANG_MARCH_FLAG= -DCMAKE_BUILD_TYPE=Release -DLLC_HOST_CPU=cortex-a78 -DHOST_CPU_CACHELINE_SIZE=64 -DWITH_LLVM_CONFIG=/usr/bin/llvm-config-10 -DSINGLE_LLVM_LIB=1 -DENABLE_ICD=1 -DENABLECUDNN=1 ..
make
sudo make install
#make install 只是安装库文件，不会复制icd到标准目录。这将导致clGetPlatform时获取不到平台。手动拷贝：
sudo mkdir -p /etc/OpenCL/vendors
sudo cp /usr/local/etc/OpenCL/vendors/pocl.icd /etc/OpenCL/vendors/

```