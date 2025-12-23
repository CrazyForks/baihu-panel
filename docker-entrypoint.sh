#!/bin/sh

PYTHON_VENV_DIR="/app/envs/python"
NODE_ENV_DIR="/app/envs/node"

# ============================
# 创建必要目录
# ============================
mkdir -p /app/data /app/data/scripts /app/configs /app/envs

# ============================
# 创建 Python 虚拟环境（如果不存在）
# ============================
if [ ! -d "$PYTHON_VENV_DIR" ]; then
    echo "Creating Python virtual environment..."
    python3 -m venv "$PYTHON_VENV_DIR"
    "$PYTHON_VENV_DIR/bin/pip" config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
    "$PYTHON_VENV_DIR/bin/pip" install --upgrade pip
    echo "Python virtual environment created at $PYTHON_VENV_DIR"
else
    echo "Python virtual environment already exists at $PYTHON_VENV_DIR"
fi

# ============================
# 创建 Node 环境目录（如果不存在）
# ============================
if [ ! -d "$NODE_ENV_DIR" ]; then
    echo "Creating Node environment directory..."
    mkdir -p "$NODE_ENV_DIR"
    # 设置 npm 全局安装目录
    npm config set prefix "$NODE_ENV_DIR"
    # 配置淘宝镜像
    npm config set registry https://registry.npmmirror.com
    echo "Node environment created at $NODE_ENV_DIR"
else
    echo "Node environment already exists at $NODE_ENV_DIR"
fi

# ============================
# 设置环境变量
# ============================
export PATH="$NODE_ENV_DIR/bin:$PYTHON_VENV_DIR/bin:$PATH"
export NODE_PATH="$NODE_ENV_DIR/lib/node_modules"

# ============================
# 启动应用
# ============================
exec ./baihu
