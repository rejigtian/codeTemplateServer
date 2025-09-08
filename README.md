# Template Server 部署文档

## 前置要求
- Ubuntu 服务器（推荐 Ubuntu 22.04 或以上）
- Docker（如未安装会自动安装）
- Git（如未安装会自动安装）
- 配置好的 SSH key（用于访问 GitHub）

## 配置步骤

### 1. 配置 SSH Key
1. 生成 SSH key（如果没有）：
```bash
ssh-keygen -t ed25519 -C "your_email@example.com"
```

2. 将公钥添加到 GitHub：
```bash
cat ~/.ssh/id_ed25519.pub
# 复制输出内容到 GitHub -> Settings -> SSH and GPG keys -> New SSH key
```

3. 测试 GitHub 连接：
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

### 2. 配置 API Keys
修改 `auth.go` 文件中的 API keys：
```go
var (
    readerKey = "your_reader_key_here"    // 只读权限
    writerKey = "your_writer_key_here"    // 读写权限
    adminKey  = "your_admin_key_here"     // 管理员权限
)
```

建议使用随机生成的 UUID 作为 key：
```bash
# 生成随机 key
uuidgen | tr -d '-'
```

### 3. 配置 Git 仓库
修改 `config.json` 文件：
```json
{
    "gitRemote": "git@github.com:your-org/your-template-repo.git"
}
```

确保：
1. 仓库已经创建
2. SSH key 有仓库的访问权限
3. 仓库为私有仓库（如果包含私密模板）

### 4. 部署服务

1. 克隆代码：
```bash
git clone https://github.com/your-org/codeTemplateServer.git
cd codeTemplateServer
```

2. 运行启动脚本：
```bash
chmod +x start.sh
./start.sh
```

启动脚本会：
- 安装必要的依赖（Go、Docker）
- 创建必要的目录
- 构建并启动 Docker 容器

### 5. 验证部署

1. 检查服务状态：
```bash
docker ps | grep template-server
docker logs template-server
```

2. 测试 API：
```bash
# 测试读取权限
curl -H "X-API-Key: your_reader_key" http://localhost:8080/api/templates/list

# 测试上传权限
curl -H "X-API-Key: your_writer_key" http://localhost:8080/api/templates/list
```

## 目录结构
```
/codeTemplateServer/
├── templates/          # 模板存储目录
│   ├── live/          # Live Templates
│   └── file/          # File Templates
├── config.json        # 配置文件
├── auth.go           # API Key 配置
└── start.sh          # 启动脚本
```

## 常见问题

### 1. GitHub 连接问题
如果无法连接 GitHub，检查：
- SSH key 是否正确配置
- 是否使用了 443 端口配置
- 服务器防火墙是否允许 443 端口

### 2. Docker 镜像下载慢
已配置国内镜像源（腾讯云），如果还是很慢，可以修改 `/etc/docker/daemon.json`：
```json
{
    "registry-mirrors": [
        "https://mirror.ccs.tencentyun.com",
        "https://docker.mirrors.ustc.edu.cn"
    ]
}
```

### 3. 权限问题
确保：
- templates 目录有正确的读写权限
- SSH key 有正确的文件权限（600）
- Docker 用户有权限访问挂载的目录

## 安全建议
1. 定期更换 API keys
2. 使用强密码策略生成 keys
3. 确保 Git 仓库为私有
4. 定期检查访问日志
5. 配置服务器防火墙，只开放必要端口

## 维护
- 日志位置：`docker logs template-server`
- 配置文件：`config.json`
- 模板备份：自动通过 Git 同步
- 服务自动重启：已配置 `--restart always`

## 升级
1. 拉取最新代码：
```bash
git pull
```

2. 重新启动服务：
```bash
./start.sh
```

## 联系方式
如有问题，请联系：
- 邮件：tianruijie@wepie.com
- 内部群：WP Coder 技术支持群
