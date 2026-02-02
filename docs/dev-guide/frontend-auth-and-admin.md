# 前端访问验证与后台管理面板开发文档

> 创建时间：2026-01-30

## 一、功能概述

### 1.1 需求说明

| 功能模块         | 描述                                                                       |
| ---------------- | -------------------------------------------------------------------------- |
| **前端访问验证** | 用户访问子功能页面时需要输入访问密码，验证通过后才能使用。主页不需要验证。 |
| **后台管理面板** | 单一管理员通过密码登录，可管理访问密码、公告内容等                         |

### 1.2 功能范围

**前端访问验证**：

- 保护范围：子功能页面（成绩查询、课表、考试安排等）
- 不保护：主页 `/`、数据大屏 `/dashboard`
- 验证方式：输入访问密码，验证通过后存储 Token

**后台管理面板**：

- 管理员认证：单一管理员，仅需密码登录
- 访问密码管理：设置/修改前端访问密码
- 公告管理：增删改查公告，公告显示在每个页面顶部
- 无需配置站点名称

---

## 二、架构重构

### 2.1 数据库重构

当前所有数据存储在 `stats.db`，职责不清晰。重构为：

| 数据库文件      | 用途     | 存储内容                             |
| --------------- | -------- | ------------------------------------ |
| `data/stats.db` | 统计数据 | API 请求日志、搜索热词、系统运行信息 |
| `data/app.db`   | 应用数据 | 系统配置、公告、管理员信息           |

### 2.2 目录结构重构

**现有结构问题**：

- `common/stats/` 混合了数据库连接和业务逻辑
- 缺少统一的数据库管理层

**新增目录**：

```
├── internal/                     # 内部包（新增）
│   ├── database/                # 数据库管理
│   │   ├── database.go         # 数据库连接管理
│   │   ├── stats.go            # 统计数据库操作
│   │   └── app.go              # 应用数据库操作
│   ├── config/                  # 系统配置
│   │   └── config.go           # 配置读写
│   └── crypto/                  # 加密工具
│       └── crypto.go           # 密码加密/Token生成
```

**新增 API 和中间件**：

```
├── api/v1/
│   ├── admin/                   # 管理后台 API（新增）
│   │   ├── handler.go          # 登录/登出
│   │   ├── config.go           # 配置管理
│   │   └── announcement.go     # 公告管理
│   └── site/                    # 站点公共 API（新增）
│       └── handler.go          # 访问验证、获取公告
├── middleware/
│   ├── site_access.go          # 前端访问验证（新增）
│   └── admin_auth.go           # 管理员认证（新增）
```

**新增前端文件**：

```
├── web/templates/
│   ├── access.html             # 访问验证页面（新增）
│   ├── components/
│   │   └── announcement.html   # 公告组件（新增）
│   └── admin/                   # 管理后台页面（新增）
│       ├── login.html
│       ├── index.html
│       ├── config.html
│       └── announcements.html
└── web/static/js/
    ├── access.js               # 访问验证逻辑（新增）
    └── admin/                   # 管理后台 JS（新增）
```

---

## 三、数据库设计

### 3.1 应用数据库 `app.db`

**系统配置表**：

```sql
CREATE TABLE IF NOT EXISTS system_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at INTEGER NOT NULL
);
```

**配置项说明**：

| key                    | 说明                     | 默认值 |
| ---------------------- | ------------------------ | ------ |
| `site_access_enabled`  | 是否开启访问验证         | `true` |
| `site_access_password` | 访问密码（bcrypt加密）   | -      |
| `admin_password`       | 管理员密码（bcrypt加密） | -      |
| `token_expire_hours`   | Token 过期时间(小时)     | `24`   |

**公告表**：

```sql
CREATE TABLE IF NOT EXISTS announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT DEFAULT 'info',      -- info/warning/error
    is_active INTEGER DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
```

---

## 四、API 设计

### 4.1 站点公共 API `/api/v1/site`

| 方法 | 路径             | 说明                    | 认证 |
| ---- | ---------------- | ----------------------- | ---- |
| POST | `/verify`        | 验证访问密码            | 无   |
| GET  | `/announcements` | 获取启用的公告列表      | 无   |
| GET  | `/check-token`   | 检查访问 Token 是否有效 | 无   |

### 4.2 管理后台 API `/api/v1/admin`

| 方法   | 路径                 | 说明         | 认证 |
| ------ | -------------------- | ------------ | ---- |
| POST   | `/login`             | 管理员登录   | 无   |
| POST   | `/logout`            | 管理员登出   | 需要 |
| GET    | `/config`            | 获取所有配置 | 需要 |
| PUT    | `/config`            | 更新配置     | 需要 |
| GET    | `/announcements`     | 获取所有公告 | 需要 |
| POST   | `/announcements`     | 创建公告     | 需要 |
| PUT    | `/announcements/:id` | 更新公告     | 需要 |
| DELETE | `/announcements/:id` | 删除公告     | 需要 |

### 4.3 API 请求/响应格式

**访问验证请求**：

```json
POST /api/v1/site/verify
{
    "password": "访问密码"
}
```

**验证成功响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "xxx",
    "expire_at": 1706745600
  }
}
```

**管理员登录请求**：

```json
POST /api/v1/admin/login
{
    "password": "管理员密码"
}
```

---

## 五、路由配置

### 5.1 页面路由保护

**需要保护的页面**：

- `/grade` - 成绩查询
- `/schedule` - 课表查询
- `/course-plan` - 培养方案
- `/exam` - 考试安排
- `/selection` - 选课结果
- `/questions` - 题库

**不需要保护的页面**：

- `/` - 主页
- `/dashboard` - 数据大屏

### 5.2 路由配置示例

```go
// router/router.go

// 公开页面
r.GET("/", renderIndex)
r.GET("/dashboard", renderDashboard)

// 受保护页面（需要访问验证）
protected := r.Group("")
protected.Use(middleware.SiteAccessRequired())
{
    protected.GET("/grade", renderGrade)
    protected.GET("/schedule", renderSchedule)
    protected.GET("/course-plan", renderCoursePlan)
    protected.GET("/exam", renderExam)
    protected.GET("/selection", renderSelection)
    protected.GET("/questions", renderQuestions)
}

// 管理后台页面
r.GET("/admin/login", renderAdminLogin)
adminPages := r.Group("/admin")
adminPages.Use(middleware.AdminAuthRequired())
{
    adminPages.GET("/", renderAdminIndex)
    adminPages.GET("/config", renderAdminConfig)
    adminPages.GET("/announcements", renderAdminAnnouncements)
}
```

---

## 六、开发阶段划分

### 阶段一：基础设施搭建

- [ ] 创建 `internal/database/` 目录和数据库管理模块
- [ ] 创建 `data/app.db` 应用数据库
- [ ] 创建 `internal/crypto/` 加密工具模块
- [ ] 创建 `internal/config/` 配置管理模块

### 阶段二：前端访问验证

- [ ] 创建 `middleware/site_access.go` 访问验证中间件
- [ ] 创建 `api/v1/site/handler.go` 站点公共 API
- [ ] 创建 `web/templates/access.html` 访问验证页面
- [ ] 集成到路由配置

### 阶段三：管理后台认证

- [ ] 创建 `middleware/admin_auth.go` 管理员认证中间件
- [ ] 创建 `api/v1/admin/handler.go` 登录/登出 API
- [ ] 创建 `web/templates/admin/login.html` 登录页面

### 阶段四：配置管理

- [ ] 创建 `api/v1/admin/config.go` 配置管理 API
- [ ] 创建 `web/templates/admin/config.html` 配置管理页面
- [ ] 实现访问密码设置功能

### 阶段五：公告管理

- [ ] 创建 `api/v1/admin/announcement.go` 公告 CRUD API
- [ ] 创建 `web/templates/admin/announcements.html` 公告管理页面
- [ ] 创建 `web/templates/components/announcement.html` 公告展示组件
- [ ] 在各页面集成公告展示

### 阶段六：测试与优化

- [ ] 功能测试
- [ ] 安全性检查
- [ ] UI/UX 优化

---

## 七、安全考虑

1. **密码存储**：使用 bcrypt 加密存储
2. **Token 安全**：使用 HMAC-SHA256 签名，包含过期时间
3. **防暴力破解**：登录失败次数限制
4. **HTTPS**：生产环境强制 HTTPS

---

## 八、Token 实现方案（无状态）

### 8.1 Token 结构

```
Token = Base64(Payload) + "." + Signature
```

**Payload 内容**：

```json
{
  "type": "site|admin",
  "exp": 1706745600
}
```

| 字段   | 说明                                               |
| ------ | -------------------------------------------------- |
| `type` | Token 类型：`site`（访问验证）或 `admin`（管理员） |
| `exp`  | 过期时间（Unix 时间戳）                            |

### 8.2 Token 生成流程

```go
func GenerateToken(tokenType string, expireHours int) string {
    // 1. 构建 Payload
    payload := map[string]interface{}{
        "type": tokenType,
        "exp":  time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
    }

    // 2. JSON 序列化并 Base64 编码
    payloadJSON, _ := json.Marshal(payload)
    payloadBase64 := base64.URLEncoding.EncodeToString(payloadJSON)

    // 3. 使用密钥生成签名
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(payloadBase64))
    signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

    // 4. 拼接 Token
    return payloadBase64 + "." + signature
}
```

### 8.3 Token 验证流程

```go
func ValidateToken(token string, expectedType string) bool {
    // 1. 分割 Token
    parts := strings.Split(token, ".")
    if len(parts) != 2 {
        return false
    }

    // 2. 验证签名
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(parts[0]))
    expectedSig := base64.URLEncoding.EncodeToString(h.Sum(nil))
    if parts[1] != expectedSig {
        return false
    }

    // 3. 解析 Payload
    payloadJSON, _ := base64.URLEncoding.DecodeString(parts[0])
    var payload map[string]interface{}
    json.Unmarshal(payloadJSON, &payload)

    // 4. 验证类型
    if payload["type"] != expectedType {
        return false
    }

    // 5. 验证过期时间
    exp := int64(payload["exp"].(float64))
    if time.Now().Unix() > exp {
        return false
    }

    return true
}
```

### 8.4 密钥管理

**密钥来源**：从环境变量 `TOKEN_SECRET` 读取，若未设置则自动生成随机密钥。

```go
var secretKey string

func init() {
    secretKey = os.Getenv("TOKEN_SECRET")
    if secretKey == "" {
        // 生成 32 字节随机密钥
        key := make([]byte, 32)
        rand.Read(key)
        secretKey = base64.StdEncoding.EncodeToString(key)
    }
}
```

**注意事项**：

- 生产环境建议在 `.env` 中配置固定密钥
- 若未配置，每次重启服务会生成新密钥，导致所有 Token 失效

---

## 九、流程图

### 9.1 前端访问验证流程

```
用户访问受保护页面
        ↓
检查 Cookie 中是否有 Token
        ↓
   ┌────┴────┐
   ↓         ↓
  有        无
   ↓         ↓
验证Token  显示密码输入页
   ↓         ↓
 ┌─┴─┐    输入密码
 ↓   ↓       ↓
有效 无效  调用验证API
 ↓   ↓       ↓
正常 清除  ┌──┴──┐
访问 Token ↓    ↓
     ↓   成功  失败
   重定向  ↓    ↓
   验证页 存储  提示
         Token 错误
           ↓
         正常访问
```

### 9.2 管理员登录流程

```
访问 /admin/*
      ↓
检查管理员 Token
      ↓
  ┌───┴───┐
  ↓       ↓
 有效    无效
  ↓       ↓
正常访问 重定向到
        /admin/login
           ↓
       输入密码
           ↓
      调用登录API
           ↓
       ┌──┴──┐
       ↓    ↓
      成功  失败
       ↓    ↓
     存储  提示
     Token 错误
       ↓
    重定向到
    /admin/
```
