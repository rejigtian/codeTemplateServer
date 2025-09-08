FROM golang:1.21-alpine

# 安装 git
RUN apk add --no-cache git openssh

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY *.go ./

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# 创建模板目录
RUN mkdir -p /app/templates/live /app/templates/file

# 设置启动命令
CMD ["/app/server"]
