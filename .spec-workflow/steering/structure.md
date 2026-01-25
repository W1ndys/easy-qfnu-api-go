# 项目结构

## 目录组织

```
qfnu-api-go/
├── api/                    # 接口层
│   └── v1/                 # API v1 版本
│       ├── grade.go        # 成绩相关接口
│       ├── login.go        # 登录相关接口
│       └── ...             # 其他业务接口
├── common/                 # 公共工具层
│   ├── logger/             # 日志配置
│   │   └── logger.go       # slog + tint 配置
│   ├── request/            # 请求辅助工具
│   │   └── request.go      # HTTP 请求封装
│   └── response/           # 统一响应封装
│       └── response.go     # 泛型响应结构
├── middleware/             # 中间件层
│   ├── auth.go             # 鉴权中间件
│   ├── cors.go             # 跨域中间件
│   └── logger.go           # 请求日志中间件
├── model/                  # 数据模型层
│   ├── grade.go            # 成绩相关结构体
│   ├── common.go           # 公共结构体
│   └── ...                 # 其他业务模型
├── service/                # 业务逻辑层
│   ├── grade.go            # 成绩业务逻辑（爬虫+解析）
│   ├── login.go            # 登录业务逻辑
│   └── ...                 # 其他业务服务
├── web/                    # 前端资源（编译时嵌入）
│   ├── index.html          # 首页
│   ├── grade.html          # 成绩查询页面
│   └── ...                 # 其他页面和静态资源
├── docs/                   # 文档目录
│   ├── API.md              # API 接口文档
│   └── PROJECT.md          # 项目介绍
├── logs/                   # 运行时日志（自动生成，已 gitignore）
├── tmp/                    # Air 热重载临时目录（已 gitignore）
├── .air.toml               # Air 热重载配置
├── .gitignore              # Git 忽略文件
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验
├── main.go                 # 程序入口
├── LICENSE                 # 开源协议
└── README.md               # 项目说明
```

## 命名规范

### 文件命名

- **业务文件**: 直接使用业务名，如 `grade.go`、`login.go`
- **不加后缀**: 通过目录区分职责，不使用 `_handler`、`_service` 后缀
- **小写蛇形**: 多单词使用下划线，如 `course_type.go`

### 代码命名

- **结构体**: PascalCase，如 `GradeInfo`、`LoginRequest`
- **函数/方法**: PascalCase（导出）或 camelCase（私有）
- **常量**: PascalCase 或 UPPER_SNAKE_CASE
- **变量**: camelCase
- **包名**: 小写单词，如 `response`、`middleware`

## 导入规范

### 导入顺序

```go
import (
    // 1. 标准库
    "fmt"
    "net/http"

    // 2. 第三方库
    "github.com/gin-gonic/gin"
    "github.com/go-resty/resty/v2"

    // 3. 项目内部包
    "qfnu-api-go/common/response"
    "qfnu-api-go/service"
)
```

### 包引用规则

- 使用绝对路径导入项目包
- 避免循环依赖
- 依赖方向：api → service → model，common 被所有层引用

## 代码结构模式

### Handler 层（api/v1/）

```go
// 1. 参数绑定
var req model.GradeRequest
if err := c.ShouldBindJSON(&req); err != nil {
    response.BadRequest(c, "参数错误")
    return
}

// 2. 调用 Service
result, err := service.GetGrade(req)

// 3. 错误处理
if err != nil {
    // 精确捕获哨兵错误
    if errors.Is(err, service.ErrCookieExpired) {
        response.CookieExpired(c)
        return
    }
    response.Error(c, err.Error())
    return
}

// 4. 返回响应
response.Success(c, result)
```

### Service 层（service/）

```go
// 1. 构造请求
client := resty.New()
resp, err := client.R().
    SetHeader("Cookie", authorization).
    Get(targetURL)

// 2. 错误检查
if err != nil {
    return nil, err
}

// 3. HTML 解析
doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

// 4. 数据提取和封装
result := &model.GradeResponse{}
doc.Find("selector").Each(func(i int, s *goquery.Selection) {
    // 解析逻辑
})

// 5. 返回结果
return result, nil
```

### Model 层（model/）

```go
// 请求参数结构体
type GradeRequest struct {
    Semester string `json:"semester" binding:"required"`
}

// 响应数据结构体
type GradeInfo struct {
    CourseName string  `json:"course_name"`
    Credit     float64 `json:"credit"`
    Score      string  `json:"score"`
}
```

## 代码组织原则

1. **目录即职责**: 通过目录名区分代码角色，文件名只表示业务领域
2. **单一职责**: 每个文件只处理一个业务领域
3. **扁平优先**: 避免过深的目录嵌套
4. **统一入口**: main.go 作为唯一入口，负责组装和启动

## 模块边界

### 层级依赖规则

```
main.go
   ↓
middleware/ ←→ common/
   ↓
api/v1/
   ↓
service/
   ↓
model/
```

- Handler 只调用 Service，不直接处理业务逻辑
- Service 不依赖 Handler，只依赖 Model 和 common
- Model 是纯数据结构，不包含业务逻辑
- common 是公共工具，被所有层引用但不引用业务层

### 错误传递规则

- Service 层定义哨兵错误（如 `ErrCookieExpired`）
- Handler 层使用 `errors.Is()` 精确捕获并返回对应响应
- 未知错误返回通用错误响应

## 代码规模指南

- **文件大小**: 单文件不超过 500 行
- **函数大小**: 单函数不超过 100 行
- **嵌套深度**: 不超过 4 层
- **参数数量**: 函数参数不超过 5 个

## 文档规范

- 导出的函数和类型需要添加注释
- 复杂逻辑添加行内注释说明意图
- API 接口在 docs/API.md 中统一文档化
- 项目级文档放在 docs/ 目录
