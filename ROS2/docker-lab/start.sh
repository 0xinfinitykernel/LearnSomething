#!/bin/bash

# Step 1: 检查 .env 是否已存在
if [ -f .env ]; then
  echo "[✓] .env 文件已存在，跳过生成。"
else
  echo "[*] 生成 .env 文件..."

  # 可交互式设置用户名和密码（默认值也支持自动化）
  read -p "请输入 Git 用户名 [默认: gituser]: " input_user
  read -s -p "请输入 Git 密码: " input_pass
  echo ""

  GIT_USER=${input_user:-gituser}
  GIT_PASS=$input_pass

  {
    echo "USER=$(whoami)"
    echo "USER_ID=$(id -u)"
    echo "GROUP_ID=$(id -g)"
    echo "HOME_DIR=$HOME"
    echo "DISPLAY=$DISPLAY"
    echo "GIT_USER=${GIT_USER}"
    echo "GIT_PASS=${GIT_PASS}"
  } > .env

  echo "[✓] .env 文件生成完毕。"
fi

# Step 1.5: 复制 ros_env2.sh 到用户主目录（如果不存在）
ROS_ENV_SRC="./ros_env2.sh"
ROS_ENV_DEST="${HOME}/ros_env2.sh"

if [ ! -f "$ROS_ENV_DEST" ]; then
  if [ -f "$ROS_ENV_SRC" ]; then
    echo "[*] 复制 ros_env2.sh 到 ${ROS_ENV_DEST} ..."
    cp "$ROS_ENV_SRC" "$ROS_ENV_DEST"
    chmod +x "$ROS_ENV_DEST"
    echo "[✓] 已复制 ros_env2.sh"
  else
    echo "[!] 未找到 ros_env2.sh，跳过复制。"
  fi
else
  echo "[✓] ${ROS_ENV_DEST} 已存在，跳过复制。"
fi

# Step 2: 启动 docker-compose
echo "[*] 启动 Docker Compose..."
docker-compose up -d

