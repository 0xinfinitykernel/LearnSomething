#!/bin/bash
set -e

echo "[*] 加载 .env..."
set -a
source .env
set +a

echo "[*] 构建镜像 ros2_multi_project_runtime ..."
docker build -f ./Dockerfile -t ros2_multi_project_runtime \
  --build-arg USER=${USER} \
  --build-arg USER_ID=${USER_ID} \
  --build-arg GROUP_ID=${GROUP_ID} \
  --build-arg GIT_USER=${GIT_USER} \
  --build-arg GIT_PASS=${GIT_PASS} \
  ..

echo "[✓] 构建完成。"

