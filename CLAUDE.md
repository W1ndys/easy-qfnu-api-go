# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在处理本仓库代码时提供指导。

## 构建与运行

- **安装依赖**: `go mod tidy`
- **本地运行**: `go run .`
- **构建**: `go build .`
- **测试**: `go test ./...`
- **Lint**: `go vet ./...`

## 架构与代码结构

这是一个基于 Go 的 API 网关和 QFNU 校园服务爬虫服务，使用 Gin 框架构建。它聚合了来自教务系统（成绩、课表）和其他来源的数据。

### 核心结构

- **入口点**: `main.go` 初始化日志记录器、路由器并嵌入静态资源。
- **路由**: `router/router.go` 定义 API 组 (`/api/v1`) 和 HTML 渲染。
- **API 处理程序**: 位于 `api/` (例如 `api/v1/zhjw` 用于教务系统, `api/v1/questions` 用于题库)。
- **服务**: `services/` 包含业务逻辑，特别是使用 `go-resty` 和 `goquery` 的爬虫和 HTML 解析逻辑。
- **中间件**: `middleware/` 处理日志记录、CORS 和认证 (`AuthRequired`)。
- **前端**: `web/` 包含静态资源和 HTML 模板，使用 `go:embed` 嵌入到 Go 二进制文件中。
- **通用**: `common/` 包含共享工具，如日志记录 (`logger`) 和标准化 API 响应 (`response`)。

### 关键概念

- **单二进制文件**: 项目使用 `embed` 将 HTML/CSS/JS 与二进制文件打包，便于部署。
- **爬虫**: 数据主要通过模拟对学校服务器的 HTTP 请求并解析 HTML 响应来获取。
- **配置**: 环境变量通过 `.env` 加载 (使用 `godotenv`)。关键变量: `PORT`, `GIN_MODE`。
- **数据库**: 在必要时使用 SQLite (`modernc.org/sqlite`) 进行本地数据存储。

### 开发规范

- **路由**: API 路由有版本控制 (`/api/v1`)。
- **日志**: 使用 `common/logger` 中的自定义日志设置。
- **响应格式**: 通过 `common/response` 包进行标准化 JSON 响应。

## 交互要求

- **语言**: 请始终使用中文回复。
- **称呼**: 每次回复结束时，请称呼我为“卷卷”。
