# Background

```
Welcome to Ubuntu 20.04.6 LTS (GNU/Linux 5.10.104-tegra aarch64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
jtop 4.2.2 - (c) 2023, Raffaello Bonghi [raffaello@rnext.it]
Website: https://rnext.it/jetson_stats

 Platform                          Serial Number: [s|XX CLICK TO READ XXX]
  Machine: aarch64                 Hardware
  System: Linux                     Model: Jetson AGX Orin
  Distribution: Ubuntu 20.04 focal  699-level Part Number: 699-13701-0005-500 M.0
  Release: 5.10.104-tegra           P-Number: p3701-0005
  Python: 3.8.10                    Module: NVIDIA Jetson AGX Orin (64GB ram)
                                    SoC: tegra23x
 Libraries                          CUDA Arch BIN: 8.7
  CUDA: 11.4.315                    Codename: Concord
  cuDNN: 8.6.0.166                  L4T: 35.3.1
  TensorRT: 8.5.2.2                 Jetpack: 5.1.1
  VPI: 2.2.7
  Vulkan: 1.3.204                  Hostname: ubuntu
  OpenCV: 4.5.4 with CUDA: YES     Interfaces
                                    eth0: 10.10.0.5
                                    docker0: 172.17.0.1
```
## Compiled CUDA is No
```
$ sudo jtop

* OpenCV:	4.1.1	compiled CUDA:	NO
```

## Uninstall default opencv packages
```
$ sudo apt purge libopencv*
$ sudo apt autoremove

```

## Update && Upgrade

```
$ sudo apt update
$ sudo apt upgrade

```

## Install dependencies
* Generic tools
```
$ sudo apt install build-essential cmake pkg-config unzip yasm git checkinstall
```
* Image I/O libs
```
$ sudo apt install libjpeg-dev libpng-dev libtiff-dev
```
* Video/Audio Libs - FFMPEG, GSTREAMER, x264 and so on
```
$ sudo apt install libavcodec-dev libavformat-dev libswscale-dev libavresample-dev
$ sudo apt install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev
$ sudo apt install libxvidcore-dev x264 libx264-dev libfaac-dev libmp3lame-dev libtheora-dev 
$ sudo apt install libfaac-dev libmp3lame-dev libvorbis-dev

```
* OpenCore - Adaptive Multi Rate Narrow Band(AMRNB) and Wide Band(AMRWB) speech codec
```
$ sudo apt install libopencore-amrnb-dev libopencore-amrwb-dev

```
* Cameras programming interface libs
```
$ sudo apt-get install libdc1394-22 libdc1394-22-dev libxine2-dev libv4l-dev v4l-utils
$ cd /usr/include/linux
$ sudo ln -s -f ../libv4l1-videodev.h videodev.h
$ cd ~

```
* GTK lib for the graphical user functionalites coming from OpenCV highghui module
```
$ sudo apt-get install libgtk-3-dev

```
* Python libraries for python3
```
$ sudo apt-get install python3-dev python3-pip
$ sudo -H pip3 install -U pip numpy
$ sudo apt install python3-testresources

```
* Parallelism library C++ for CPU
```
$ sudo apt-get install libtbb-dev
```
* Optimization libraries for OpenCV
```
$ sudo apt-get install libatlas-base-dev gfortran
```
* Optional libraries
```
$ sudo apt-get install libprotobuf-dev protobuf-compiler
$ sudo apt-get install libgoogle-glog-dev libgflags-dev
$ sudo apt-get install libgphoto2-dev libeigen3-dev libhdf5-dev doxygen

```
* Download OpenCV && unpack 
```
$ cd ~/Downloads
$ wget -O opencv.zip https://github.com/opencv/opencv/archive/refs/tags/4.5.4.zip
$ wget -O opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/refs/tags/4.5.4.zip
$ unzip opencv.zip
$ unzip opencv_contrib.zip

```
* build opencv
```
$ cd opencv-4.5.4
$ mkdir build
$ cd build

```
* Cmake
```
## CUDA_TOOLKIT_ROOT_DIR应改为自己开发板上CUDA的根目录
## CUDA_ARCH_BIN应改为GPU计算能力，当前所用开发板支持的CUDA版本为11.4，计算能力8.7
## OPENCV_EXTRA_MODULES_PATH为opencv_contrib的路径
$ cmake -D CMAKE_BUILD_TYPE=RELEASE -D CMAKE_INSTALL_PREFIX=/usr/local \
-D BUILD_opencv_python2=1 -D BUILD_opencv_python3=1 -D WITH_FFMPEG=1 \
-D CUDA_TOOLKIT_ROOT_DIR=/usr/local/cuda-11.4 \
-D WITH_TBB=ON -D ENABLE_FAST_MATH=1 -D CUDA_FAST_MATH=1 -D WITH_CUBLAS=1 \
-D WITH_CUDA=ON -D BUILD_opencv_cudacodec=OFF -D WITH_CUDNN=ON \
-D OPENCV_DNN_CUDA=ON \
-D CUDA_ARCH_BIN=8.7 \
-D WITH_V4L=ON -D WITH_QT=OFF -D WITH_OPENGL=ON -D WITH_GSTREAMER=ON \
-D OPENCV_GENERATE_PKGCONFIG=ON -D OPENCV_PC_FILE_NAME=opencv.pc \
-D OPENCV_ENABLE_NONFREE=ON \
-D OPENCV_EXTRA_MODULES_PATH=/home/nvidia/Downloads/opencv_contrib-4.5.4/modules ..

```
* Make
```
$ make -j$(nproc)
$ sudo make install

```
