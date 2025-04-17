# 🧪 ROS 包维护者测试说明文档

## 📌 背景说明

为确保 ROS 多项目构建环境在 CI/CD 中稳定运行，请各 ROS 包维护者在合入主构建系统前，**自行测试项目是否可在统一构建镜像中成功编译**。

测试方式基于本仓库提供的容器化环境与自动化脚本，测试结果请通过 YAML 文件进行记录并提交。

---

## ✅ 测试流程说明

请参考主仓库中的 [`README.md`](./README.md) 进行操作，包括：

- 初始化 `.env` 文件
- 添加你的仓库地址至 `projects.txt`
- 执行构建脚本 `./build.sh`

构建完成后请根据测试情况填写 `build_info.yaml`。

---

## 📄 提交编译信息（build_info.yaml）

测试完成后，请在你的 ROS 项目仓库根目录下新建并填写 `build_info.yaml`，内容包括依赖、编译状态等信息，格式如下：

```yaml
# build_info.yaml

project:
  name: your_ros_package
  repo_url: http://192.168.0.108/yourteam/your_ros_package.git
  branch: main

dependencies:
  apt:
    - libopencv-dev
    - libpcl-dev
  ros:
    - rclcpp
    - sensor_msgs
    - tf2_ros

build:
  status: success  # 可选值：success | fail
  errors: []        # 如果失败，写出关键报错信息（可多条）
  warnings: []      # 可选：构建过程的主要警告

notes: >
  项目依赖本地相机 SDK 编译时请连接设备；
  若无设备可在 CMake 中添加 OFF 开关；
  已在 ROS humble + Ubuntu 22.04 下验证通过。
```

---

## 📬 提交说明

- 请将 `build_info.yaml` 提交到你的项目仓库根目录
- CI 将自动读取该文件，统计依赖项与构建状态

---

## 🧩 进阶计划

- 支持扩展为 `projects.yaml`，配置分支、路径等参数
- 支持集成 GitLab CI / GitHub Actions 实现自动构建
- 使用多阶段构建优化镜像大小与构建速度
- 可接入本地或远程私有 Docker Registry 实现构建缓存

---

## 🔐 安全说明

- `.env` 文件中包含 Git 用户名和密码，请勿提交到 Git 仓库
- `.env` 默认已在 `.gitignore` 中忽略
- 推荐使用 Git Token 代替明文密码进行访问

---

## 📮 联系与支持

由 [infinitykernel.com](http://infinitykernel.com) 提供支持。  
如需定制 ROS 构建环境或 CI/CD 集成请联系：

📧 0x@infinitykernel.com
