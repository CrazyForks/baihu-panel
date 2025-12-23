#!/bin/sh

# ============================
# 激活 Python 虚拟环境
# ============================
if [ -d /opt/venv ]; then
    export PATH="/opt/venv/bin:$PATH"
fi

# ============================
# 创建必要目录
# ============================
mkdir -p /app/data /app/data/scripts /app/configs

# ============================
# 设置不可变属性
# ============================
if [ -f /app/data/ql.db ]; then
    chattr +i /app/data/ql.db 2>/dev/null || true
fi

if [ -d /app/configs ]; then
    chattr +i /app/configs 2>/dev/null || true
fi

if [ -d /app/data/scripts ]; then
    chattr +i /app/data/scripts 2>/dev/null || true
fi

# ============================
# 启动应用
# ============================
exec ./baihu
