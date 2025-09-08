#!/bin/bash

# 确保在正确的目录
cd "$(dirname "$0")"

# 检查并安装 Go
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    sudo apt-get update
    sudo apt-get install -y golang-go git
fi

# 设置 GOPROXY
export GOPROXY=https://goproxy.cn,direct

# 确保目录结构正确
mkdir -p templates/{live,file}

# 检查是否已经在运行
if pgrep -f "template-server" > /dev/null; then
    echo "Stopping existing template server..."
    pkill -f "template-server"
    sleep 2
fi

# 编译服务
echo "Building server..."
go build -o template-server .

# 创建日志目录
mkdir -p logs

# 启动服务
echo "Starting server..."
nohup ./template-server > logs/server.log 2>&1 &

# 等待服务启动
sleep 2

# 检查服务状态
if pgrep -f "template-server" > /dev/null; then
    echo "Server started successfully!"
    echo "Process ID: $(pgrep -f "template-server")"
    echo "Log file: $(pwd)/logs/server.log"
    echo "You can check logs with: tail -f logs/server.log"
else
    echo "Failed to start server. Check logs/server.log for details."
    exit 1
fi

# 显示最新日志
echo -e "\nLatest logs:"
tail -n 10 logs/server.log
