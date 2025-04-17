# 🤖 ROS 2 多项目自动构建运行环境

本项目提供一套完整的 Docker 化 ROS 2 构建与运行方案，支持从多个私有 Git 仓库自动拉取代码、安装依赖、构建并运行 ROS 2 节点。适用于开发调试、CI/CD 集成与可移植部署。

---

## 📁 项目结构说明

| 文件/目录             | 说明 |
|----------------------|------|
| `Dockerfile`         | 构建 ROS 2 容器镜像，支持多项目 clone + install + build |
| `projects.txt`       | 多项目 Git 仓库地址列表 |
| `build.sh`           | 构建镜像脚本，读取 `.env` 自动注入构建参数 |
| `start.sh`           | 启动容器脚本（首次执行可生成 `.env`） |
| `docker-compose.yml` | 启动容器、挂载目录、绑定 X11 等 |
| `.env.template`      | 示例配置文件，复制为 `.env` 使用 |

---

## 🚀 使用指南

### 1️⃣ 第一次使用：生成 `.env`

首次运行请执行：

```bash
./start.sh
```

该脚本会自动生成 `.env` 文件，并引导你输入 Git 用户名与密码（用于克隆私有仓库）。

也可手动复制：

```bash
cp .env.template .env
```

---

### 2️⃣ 编辑项目列表

将你要构建的 Git 仓库地址写入 `projects.txt`，一行一个，例如：

```
http://192.168.0.108/yourteam/lx_camera_ros2.git
http://192.168.0.108/yourteam/another_module.git
```

无需在 URL 中写入用户名和密码，系统会自动从 `.env` 注入。

---

### 3️⃣ 构建镜像

执行以下命令构建完整环境镜像：

```bash
./build.sh
```

构建过程将执行以下内容：

- 克隆所有仓库
- 自动执行每个项目中的 `install.sh`（如果存在）
- 使用 `colcon build` 构建全部 ROS 2 包
- 最终生成镜像 `ros2_multi_project_runtime`

---

### 4️⃣ 启动容器运行 ROS 系统

容器构建完成后，可直接启动：

```bash
./start.sh
```

该命令将基于 `docker-compose.yml` 启动容器，并进入你指定的 ROS 工作目录（如 `/home/youruser/ws`），自动执行设定的启动指令。

---

## 🔐 安全说明

- `.env` 包含敏感信息（如 Git 密码/Token），**请勿提交到仓库**
- `.env` 已在 `.gitignore` 中忽略
- 推荐在企业或 CI/CD 场景中使用 Git 访问 Token 替代密码

---

## 🧩 常见问题

### ❓ Q: 每次启动都会重新构建吗？

不会。只有在你修改项目代码或 `projects.txt` 后，手动执行 `./build.sh` 才会重新构建。平时只需执行 `./start.sh` 启动即可。

---

### ❓ Q: 如何指定启动哪个 ROS 节点？

你可以修改 `docker-compose.yml` 中的 `command:` 字段，例如：

```yaml
command: >
  bash -c "
  source /opt/ros/humble/setup.bash &&
  source install/setup.bash &&
  ros2 launch your_package your_launch.py"
```

---

## 🧠 扩展建议

- 支持从 `projects.yaml` 加载更复杂的项目配置（如指定分支/路径）
- 集成 GitLab CI / GitHub Actions，实现自动构建推送镜像
- 使用多阶段构建优化镜像体积与构建速度
- 接入本地或远程私有 Docker Registry

---

## 🧩 附加说明：ros_env2.sh 环境加载脚本

项目中包含一个环境初始化脚本 `ros_env2.sh`，用于在容器或交互终端中自动加载 ROS 2 环境。

### 🔧 文件结构示例：

```bash
#!/bin/bash
set -e

# 加载 ROS 2 主环境
source /opt/ros/humble/setup.bash

# 如需加载其他工作空间，可按需取消注释
# source "/home/parallels/catkin_ws/devel/setup.bash"
# source "/opt/cartographer_ros/setup.bash"

# 执行传入命令
exec "$@"
```
## 📮 联系与支持

由 [infinitykernel.com](https://infinitykernel.com) 提供支持  
如需定制 ROS 构建系统或 CI/CD 集成，请联系：

📧 **0x@infinitykernel.com**

