#!/bin/bash

# Step 1: 检查 .env 是否已存在
if [ -f .env ]; then
  echo "[✓] .env 文件已存在，跳过生成。"
else
  echo "[*] 生成 .env 文件..."
  echo "USER_ID=$(id -u)" > .env
  echo "GROUP_ID=$(id -g)" >> .env
  echo "HOME_DIR=$HOME" >> .env
  echo "DISPLAY=$DISPLAY" >> .env
  echo "[✓] .env 文件生成完毕。"
fi

# Step 2: 启动 docker-compose
echo "[*] 启动 Docker Compose..."
docker-compose up -d
