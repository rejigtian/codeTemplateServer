#!/bin/bash

# 检查参数
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <server_host> <ssh_key_path>"
    echo "Example: $0 user@your-server.com ~/.ssh/id_rsa"
    exit 1
fi

SERVER_HOST=$1
SSH_KEY=$2

# 创建远程目录
ssh -i $SSH_KEY $SERVER_HOST "mkdir -p ~/template-server"

# 复制文件到服务器
scp -i $SSH_KEY Dockerfile docker-compose.yml *.go $SERVER_HOST:~/template-server/

# 在服务器上构建和启动服务
ssh -i $SSH_KEY $SERVER_HOST "cd ~/template-server && docker-compose up -d --build"

echo "Deployment completed!"
