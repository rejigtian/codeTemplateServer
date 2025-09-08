#!/bin/bash

# 确保在正确的目录
cd "$(dirname "$0")"

# 确保目录结构正确
mkdir -p templates/{live,file}

# 下载依赖
go mod download

# 停止并删除旧容器（如果存在）
docker rm -f template-server 2>/dev/null || true

# 重新构建镜像
echo "Building Docker image..."
docker build -t template-server . || {
    echo "Failed to build Docker image"
    exit 1
}

# 启动新容器
echo "Starting container..."
docker run -d \
  --name template-server \
  --restart always \
  -p 8080:8080 \
  -v $(pwd)/templates:/app/templates \
  -v ~/.ssh:/root/.ssh:ro \
  template-server || {
    echo "Failed to start container"
    exit 1
}

# 检查容器状态
echo "Checking container status..."
sleep 2
docker ps | grep template-server
echo "Container logs:"
docker logs template-server