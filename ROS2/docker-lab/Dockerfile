FROM osrf/ros:humble-desktop-full

LABEL maintainer="0x@infinitykernel.com"

# ========= 构建时参数 =========
ARG USER
ARG USER_ID
ARG GROUP_ID
ARG GIT_USER
ARG GIT_PASS

# 安装 sudo、开发工具
RUN apt update && \
    apt install -y sudo git python3-colcon-common-extensions curl python3 && \
    rm -rf /var/lib/apt/lists/*

# 创建和宿主机一致的用户，避免挂载目录 root 权限问题
RUN groupadd -g ${GROUP_ID} ${USER} && \
    useradd -m -u ${USER_ID} -g ${GROUP_ID} ${USER} && \
    echo "${USER} ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# 切换用户
USER ${USER}
WORKDIR /home/${USER}

# 默认工作区结构
ENV ROS_WORKSPACE=/home/${USER}/ws
RUN mkdir -p $ROS_WORKSPACE/src
WORKDIR $ROS_WORKSPACE/src

# ======= 多项目支持 =======
# 添加项目列表文件（你提前准备好）
COPY projects.txt /tmp/projects.txt

# 克隆所有项目（项目地址写在 projects.txt，自动注入用户名密码）
RUN set -e; \
    ENCODED_PASS=$(python3 -c "import urllib.parse; print(urllib.parse.quote('${GIT_PASS}'))"); \
    while IFS= read -r raw_url; do \
      proto=$(echo $raw_url | cut -d/ -f1); \
      host_and_path=$(echo $raw_url | cut -d/ -f3-); \
      full_url="${proto}//${GIT_USER}:${ENCODED_PASS}@${host_and_path}"; \
      echo "⏬ Cloning $full_url ..."; \
      git clone "$full_url"; \
    done < /tmp/projects.txt

# 安装每个项目的 install.sh（如果有）
RUN set -e; \
    for d in */ ; do \
      if [ -f "$d/install.sh" ]; then \
        echo "🔧 Installing $d/install.sh ..."; \
        cd "$d" && chmod +x install.sh && sudo ./install.sh && cd ..; \
      fi \
    done

# 自动修复缺失宏定义
#RUN find . \( -name "*.cpp" -o -name "*.h" \) -type f | xargs sed -i 's/LX_INT_2D_UNDISTORT_SCALE/LX_INT_3D_UNDISTORT_SCALE/g'
# 自动修复缺失宏定义（仅当 SDK 中没有定义时才替换）
RUN if grep -qr "LX_INT_2D_UNDISTORT_SCALE" /opt/Lanxin-MRDVS/include; then \
      echo "✅ SDK defines LX_INT_2D_UNDISTORT_SCALE, no need to patch"; \
    else \
      echo "⚠️  LX_INT_2D_UNDISTORT_SCALE not found, patching to LX_INT_3D_UNDISTORT_SCALE"; \
      find . \( -name "*.cpp" -o -name "*.h" \) -type f -print0 \
        | xargs -0 --no-run-if-empty sed -i 's/LX_INT_2D_UNDISTORT_SCALE/LX_INT_3D_UNDISTORT_SCALE/g'; \
    fi

# 编译所有项目
WORKDIR $ROS_WORKSPACE
RUN bash -c "source /opt/ros/humble/setup.bash && colcon build"

# 默认入口由 docker-compose 控制
