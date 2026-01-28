# QFNU-API-GO 项目介绍

## 项目概述

这是一个使用 Go 语言封装的曲阜师范大学教务系统 API 中间件项目。作为我学习 Go 语言的练手项目，它不仅实现了教务系统数据的结构化获取，更承载了我对 Web API 设计的思考和实践。

## 设计哲学

我将这个项目的设计哲学命名为 **"极简分层架构"(Minimalist Layered Architecture)**，核心理念是：

> **够用就好，简洁至上，一目了然**

### 1. 文件命名极简化

传统 Go 项目常见 `grade_handler.go`、`grade_service.go` 这样的命名。我选择直接用 `grade.go`，通过**目录**而非**后缀**来区分职责：

```
api/v1/grade.go    ← 这就是 Handler
service/grade.go   ← 这就是 Service
model/grade.go     ← 这就是 Model
```

文件名只说明"这是什么业务"，目录名说明"这是什么角色"。

### 2. 分层清晰但不繁琐

采用经典的四层架构，但刻意保持轻量：

```
┌─────────────────────────────────────┐
│           main.go (入口)            │
│      路由注册 + 静态资源 + 启动      │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│       middleware/ (中间件层)         │
│      CORS / Auth / Logger           │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│         api/v1/ (接口层)            │
│     参数绑定 → 调用Service → 响应    │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│        service/ (业务逻辑层)         │
│    HTTP客户端 + 爬虫 + HTML解析      │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│         model/ (数据模型层)          │
│      结构体 + 请求参数 + 响应结构     │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│        common/ (公共工具层)          │
│    统一响应 / 日志 / 请求工具        │
└─────────────────────────────────────┘
```

没有 DAO 层、没有 Repository 层——因为这个项目不需要数据库，所有数据来自教务系统的实时爬取。**不为抽象而抽象**。

### 3. 统一响应封装 + 泛型

Go 1.18 引入泛型后，我立刻将其应用于响应封装：

```go
type Response[T any] struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data T      `json:"data"`
}

// 一行代码返回成功响应，类型自动推导
response.Success(c, gradeList)
```

所有 API 响应结构一致，前端只需一套解析逻辑。

### 4. 单文件部署

使用 Go 1.16 的 `embed` 特性，将前端资源编译进二进制文件：

```go
//go:embed web
var webFS embed.FS
```

最终产物只有一个可执行文件，`scp` 到服务器直接运行，无需 Docker、无需 Nginx。

### 5. 中间件解耦

鉴权、日志、跨域处理全部抽离为独立中间件：

```go
r.Use(middleware.RequestLogger())    // 全局请求日志
r.Use(middleware.Cors())             // 全局跨域
apiGroup.Use(middleware.AuthRequired())  // API 路由鉴权
```

新增功能只需实现业务逻辑，无需关心这些横切关注点。

### 6. 哨兵错误 + 精确捕获

定义明确的错误类型，Handler 层精确处理：

```go
// service 层定义
var ErrCookieExpired = errors.New("cookie_expired_or_invalid")

// handler 层捕获
if errors.Is(err, service.ErrCookieExpired) {
    response.CookieExpired(c)
    return
}
```

### 7. 容错优先

课程类型支持中文名或 ID 双向输入，未知类型直接透传：

```go
func GetCourseTypeID(input string) string {
    if id, ok := CourseTypeNameToID[input]; ok {
        return id
    }
    return input  // 未知类型不报错，假设是 ID 直接返回
}
```

教务系统更新新增课程类型时，本系统不会因此失效。

## 技术选型

| 层级        | 技术        | 选择理由             |
| ----------- | ----------- | -------------------- |
| Web 框架    | Gin         | 性能优秀，生态成熟   |
| HTTP 客户端 | Resty       | 链式调用，响应拦截器 |
| HTML 解析   | goquery     | jQuery 风格，上手快  |
| 日志        | slog + tint | 标准库 + 彩色输出    |
| 日志切割    | lumberjack  | 自动切割，自动压缩   |
| 静态资源    | embed       | 编译嵌入，单文件部署 |

## 目录结构

```
easy-qfnu-api-go/
├── api/v1/          # 接口层 - Handler
├── common/          # 公共层 - 工具函数
│   ├── logger/      # 日志配置
│   ├── request/     # 请求辅助
│   └── response/    # 统一响应
├── middleware/      # 中间件层
├── model/           # 模型层 - 结构体定义
├── service/         # 逻辑层 - 业务实现
├── web/             # 前端资源 (编译时嵌入)
├── logs/            # 运行时日志 (自动生成)
├── docs/            # 文档
├── main.go          # 程序入口
└── go.mod           # 依赖管理
```

## 设计原则总结

1. **简洁优先** - 文件命名简化，目录即职责
2. **分层清晰** - Handler → Service → Model 职责分明
3. **中间件解耦** - 横切关注点独立处理
4. **统一响应** - 泛型封装 + 错误码体系
5. **单文件部署** - embed 嵌入前端资源
6. **容错设计** - 允许未知值通过，不轻易报错
7. **够用就好** - 不为抽象而抽象，不为扩展而过度设计

## 写在最后

这个项目是我学习 Go 的产物，也是我对 "什么是好的 API 设计" 的一次回答。

它不追求大而全的企业级架构，而是探索在保持代码质量的前提下，能有多简洁。

**简单的代码更容易理解，容易理解的代码更容易维护。**

---

_作者：W1ndys_
_项目地址：https://github.com/W1ndys/easy-qfnu-api-go_
