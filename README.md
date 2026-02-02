# easy-qfnu-api-go

QFNU 综合校园服务网关。基于 Go (Gin) 构建，不仅支持教务系统（成绩/课表）清洗，更聚合了通知公告检索、题库搜索等多元化服务，实现单文件部署的一站式校园数据中心。A comprehensive campus service gateway for QFNU built with Go. Aggregates educational data (grades/schedules), notification search, and question banks into a single-binary, high-performance API solution.

API文档地址：https://easy-qfnu-api.apifox.cn/

> 该项目仅供学习交流使用，请勿用于商业用途，本工具仅供学习HTTP协议和HTML解析技术，使用前应获得系统管理员授权。开发者不对用户违反服务条款的行为负责。本工具仅供个人学习研究，请遵守目标网站服务条款
>
> 该项目是作者学习Golang语言过程中练手的一个小项目，代码质量和设计模式可能不够完善，欢迎各位大佬指正和交流。

---

## 部署指南

### 环境要求

- Go 1.21+ (推荐 1.25+)
- Git (用于克隆源码)

```bash
# 1. 克隆仓库
git clone https://github.com/W1ndys/easy-qfnu-api-go.git
cd easy-qfnu-api-go

# 2. 安装依赖
go mod tidy

# 3. 编译
go build -o easy-qfnu-api-go

# 4. 运行
./easy-qfnu-api-go
```

### 配置说明

配置通过环境变量或 `.env` 文件设置：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `PORT` | `8141` | 服务监听端口 |

**示例 `.env` 文件：**

```env
PORT=8141
```

### 验证部署

启动成功后，访问以下地址验证：

| 地址 | 说明 |
|------|------|
| `http://127.0.0.1:8141/` | Web 首页 |
| `http://127.0.0.1:8141/grade` | 成绩查询页面 |
| `http://127.0.0.1:8141/api/health` | API 健康检查 |

### 生产环境部署建议

#### 使用 systemd（Linux）

创建服务文件 `/etc/systemd/system/easy-qfnu-api.service`：

```ini
[Unit]
Description=Easy QFNU API Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/easy-qfnu-api-go
ExecStart=/opt/easy-qfnu-api-go/easy-qfnu-api-go
Restart=always
RestartSec=5
Environment=GIN_MODE=release
Environment=PORT=8141

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable easy-qfnu-api
sudo systemctl start easy-qfnu-api
```

#### 使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://127.0.0.1:8141;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### 使用 Docker（可选）

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o easy-qfnu-api-go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/easy-qfnu-api-go .
EXPOSE 8141
ENV GIN_MODE=release
CMD ["./easy-qfnu-api-go"]
```

构建并运行：

```bash
docker build -t easy-qfnu-api-go .
docker run -d -p 8141:8141 --name qfnu-api easy-qfnu-api-go
```

---

## 架构说明

详细的架构设计和安全分析请参阅 [架构与安全分析文档](docs/architecture-and-security.md)。

---

## 开发文档

- [开发规范 v1](docs/dev-guide/v1.md)
- [开发规范 v2](docs/dev-guide/v2.md)
- [UI 设计规范](docs/UI.md)
