# WP Template Server

WP Template Server 是 WP Template Coder 插件的服务器端组件，提供模板存储、共享和版本控制功能。

## 功能特点

- 模板存储和管理
- 基于角色的访问控制
- Git 版本控制集成
- RESTful API 接口
- Docker 容器化部署
- 自动备份和同步

## 系统要求

### 服务器环境
- Ubuntu 22.04 或更高版本（推荐）
- 2GB RAM 以上
- 10GB 可用磁盘空间

### 必要软件
- Docker（如未安装会自动安装）
- Git（如未安装会自动安装）
- Go 1.21 或更高版本（用于本地开发）

## 快速开始

### 使用 Docker（推荐）

1. 克隆仓库：
```bash
git clone https://github.com/your-org/wp-template-server.git
cd wp-template-server
```

2. 配置环境：
```bash
cp config.example.json config.json
# 编辑 config.json 配置文件
```

3. 启动服务：
```bash
./start.sh
```

### 手动部署

1. 安装依赖：
```bash
# 安装 Go
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装其他依赖
sudo apt update
sudo apt install -y git make
```

2. 构建和运行：
```bash
go build
./start_native.sh
```

## 配置说明

### 基本配置
编辑 `config.json`：
```json
{
    "port": 8080,
    "dataDir": "./data",
    "gitRemote": "git@github.com:your-org/your-template-repo.git",
    "backupInterval": 3600,
    "maxUploadSize": 10485760
}
```

### API 密钥配置
编辑 `auth.go`：
```go
var (
    readerKey = "your_reader_key_here"    // 只读权限
    writerKey = "your_writer_key_here"    // 读写权限
    adminKey  = "your_admin_key_here"     // 管理员权限
)
```

建议使用 UUID 生成密钥：
```bash
uuidgen | tr -d '-'
```

### Git 配置

1. 生成 SSH 密钥：
```bash
ssh-keygen -t ed25519 -C "template-server@your-domain.com"
```

2. 配置 GitHub 访问：
```bash
# 配置 SSH 使用 443 端口（国内服务器需要）
mkdir -p ~/.ssh
echo "Host github.com
    Hostname ssh.github.com
    Port 443
    User git" > ~/.ssh/config

# 测试连接
ssh -T git@github.com
```

## API 接口

### 认证
所有 API 请求需要包含 API 密钥：
```bash
curl -H "X-API-Key: your_api_key" http://your-server:8080/api/...
```

### 主要端点

#### 模板管理
- `GET /api/templates/list` - 获取模板列表
- `POST /api/templates/upload` - 上传模板
- `GET /api/templates/download/{id}` - 下载模板
- `DELETE /api/templates/{id}` - 删除模板（需要管理员权限）

#### 系统管理
- `GET /api/system/status` - 获取系统状态
- `POST /api/system/backup` - 触发备份（需要管理员权限）
- `GET /api/system/logs` - 获取系统日志（需要管理员权限）

## 目录结构
```
/wp-template-server/
├── api/              # API 处理器
├── auth/             # 认证相关
├── config/           # 配置管理
├── data/             # 数据存储
│   ├── templates/    # 模板文件
│   └── backup/       # 备份文件
├── docker/           # Docker 相关文件
├── scripts/          # 辅助脚本
└── templates/        # 模板存储
    ├── live/         # Live Templates
    └── file/         # File Templates
```

## 监控和维护

### 日志查看
```bash
# Docker 部署
docker logs template-server

# 本地部署
tail -f /var/log/template-server.log
```

### 备份管理
```bash
# 手动触发备份
curl -H "X-API-Key: admin_key" -X POST http://localhost:8080/api/system/backup

# 查看备份状态
curl -H "X-API-Key: admin_key" http://localhost:8080/api/system/backup/status
```

### 健康检查
```bash
curl http://localhost:8080/health
```

## 常见问题

### 1. 连接问题
- 检查防火墙设置
- 验证 API 密钥
- 确认服务器状态

### 2. 上传失败
- 检查文件大小限制
- 验证用户权限
- 确认存储空间

### 3. Git 同步问题
- 检查 SSH 密钥配置
- 验证仓库访问权限
- 确认网络连接

## 安全建议

1. 定期更换 API 密钥
2. 使用强密码策略
3. 启用 HTTPS
4. 限制 IP 访问
5. 定期审计日志

## 性能优化

1. 配置合适的上传限制
2. 启用压缩
3. 使用 CDN（可选）
4. 优化数据库查询
5. 配置缓存策略

## 技术支持
- 邮箱：tianruijie@wepie.com
- 内部群：WP Coder 技术支持群

## 贡献
欢迎提交 Pull Request 或提出 Issue。

## 许可证
本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。