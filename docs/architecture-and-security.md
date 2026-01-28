# 曲奇API 架构与安全分析

本文档详细说明曲奇API服务如何与教务系统进行交互，以及相关的安全机制分析。

---

## 一、系统架构概述

曲奇API是一个**代理/转发型**服务，它不直接存储用户的账号密码，而是通过用户提供的会话凭证（Cookie）代为访问教务系统，并将返回的HTML数据解析为结构化的JSON格式。

```mermaid
graph TB
    subgraph 用户端
        A[用户浏览器/客户端]
    end

    subgraph 曲奇API服务
        B[Gin Web 框架]
        C[认证中间件]
        D[API Handler]
        E[Service 层]
        F[HTTP Client]
        G[HTML 解析器]
    end

    subgraph 教务系统
        H[教务系统<br/>zhjw.qfnu.edu.cn]
    end

    A -->|1. 携带 Authorization 请求| B
    B --> C
    C -->|2. 验证 Authorization 存在| D
    D --> E
    E --> F
    F -->|3. 转发 Cookie 请求| H
    H -->|4. 返回 HTML 页面| F
    F --> G
    G -->|5. 解析为结构化数据| E
    E -->|6. 返回 JSON 响应| A
```

---

## 二、认证流程详解

### 2.1 Cookie 获取流程（用户侧）

用户需要**自行登录教务系统**获取有效的会话Cookie，曲奇API不参与登录过程。

```mermaid
sequenceDiagram
    participant User as 用户
    participant Browser as 浏览器
    participant ZHJW as 教务系统
    participant API as 曲奇API

    Note over User,ZHJW: 第一阶段：用户自行获取Cookie
    User->>Browser: 1. 打开教务系统登录页
    Browser->>ZHJW: 2. 访问登录页面
    ZHJW-->>Browser: 3. 返回登录表单
    User->>Browser: 4. 输入学号/密码
    Browser->>ZHJW: 5. POST 登录请求
    ZHJW-->>Browser: 6. 登录成功，Set-Cookie: JSESSIONID=xxx
    User->>Browser: 7. F12 复制 Cookie

    Note over User,API: 第二阶段：使用Cookie调用API
    User->>API: 8. GET /api/v1/zhjw/grades<br/>Authorization: JSESSIONID=xxx
    API->>ZHJW: 9. 携带Cookie请求数据
    ZHJW-->>API: 10. 返回HTML数据
    API-->>User: 11. 返回JSON数据
```

### 2.2 请求处理流程

```mermaid
flowchart TD
    A[客户端请求] --> B{Authorization<br/>头存在?}
    B -->|否| C[返回 401<br/>Cookie已过期]
    B -->|是| D[存储到 Gin Context]
    D --> E[调用 Service 层]
    E --> F[创建 HTTP Client<br/>设置 Cookie 头]
    F --> G[请求教务系统]
    G --> H{响应检查}
    H -->|包含'用户登录'| I[返回 401<br/>Cookie已过期]
    H -->|包含'未查询到数据'| J[返回 404<br/>资源不存在]
    H -->|正常响应| K[解析 HTML]
    K --> L[返回 JSON 数据]
```

---

## 三、数据流向分析

### 3.1 完整请求链路

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant Middleware as 认证中间件
    participant Handler as API Handler
    participant Service as Service层
    participant HTTPClient as HTTP Client
    participant Parser as HTML解析器
    participant ZHJW as 教务系统

    Client->>Middleware: GET /api/v1/zhjw/grades<br/>Authorization: JSESSIONID=abc123

    Note over Middleware: 验证 Authorization 存在
    Middleware->>Handler: 传递请求 + Context

    Handler->>Service: GetGradeList(cookie, params)

    Service->>HTTPClient: NewClient(cookie)
    Note over HTTPClient: 设置请求头<br/>Cookie: JSESSIONID=abc123<br/>User-Agent: Mozilla/5.0...

    HTTPClient->>ZHJW: POST /jsxsd/kscj/cjcx_list<br/>Cookie: JSESSIONID=abc123

    ZHJW-->>HTTPClient: 200 OK<br/>HTML Table 数据

    Note over HTTPClient: 响应拦截器检查<br/>是否包含登录页面

    HTTPClient->>Parser: 传递 HTML
    Parser->>Parser: goquery 解析表格

    Parser-->>Service: 结构化数据
    Service-->>Handler: GradeResult
    Handler-->>Client: JSON Response<br/>{"code":200,"data":{...}}
```

### 3.2 核心组件职责

```mermaid
graph LR
    subgraph 路由层
        R[router.go<br/>路由注册]
    end

    subgraph 中间件层
        M[auth.go<br/>认证校验]
    end

    subgraph API层
        A1[grade.go]
        A2[exam_schedules.go]
        A3[class_schedules.go]
        A4[course_plan.go]
    end

    subgraph Service层
        S1[grade.go<br/>成绩查询]
        S2[exam_schedules.go<br/>考试安排]
        S3[class_schedules.go<br/>课程表]
        S4[course_plan.go<br/>培养方案]
    end

    subgraph 基础设施
        C[client.go<br/>HTTP客户端]
        P[HTML解析<br/>goquery]
    end

    R --> M
    M --> A1 & A2 & A3 & A4
    A1 --> S1
    A2 --> S2
    A3 --> S3
    A4 --> S4
    S1 & S2 & S3 & S4 --> C
    C --> P
```

---

## 四、安全机制分析

### 4.1 安全设计要点

```mermaid
mindmap
  root((安全设计))
    凭证处理
      不存储账号密码
      不存储Cookie
      Cookie仅在请求期间使用
      内存中短暂存在
    权限控制
      用户只能访问自己的数据
      权限由教务系统控制
      API无额外权限提升
    传输安全
      支持HTTPS部署
      Cookie不记录到日志
    会话管理
      自动检测Cookie过期
      过期立即返回401
```

### 4.2 安全边界划分

```mermaid
graph TB
    subgraph 用户责任区
        U1[保管好自己的Cookie]
        U2[不要泄露给他人]
        U3[定期更换/重新登录]
    end

    subgraph API服务责任区
        A1[不存储敏感凭证]
        A2[仅做数据转发和解析]
        A3[不记录Cookie到日志]
        A4[检测异常立即中断]
    end

    subgraph 教务系统责任区
        Z1[验证Cookie有效性]
        Z2[控制数据访问权限]
        Z3[管理会话生命周期]
        Z4[记录访问日志]
    end
```

### 4.3 威胁模型分析

```mermaid
flowchart TD
    subgraph 潜在威胁
        T1[Cookie泄露]
        T2[中间人攻击]
        T3[API服务被入侵]
        T4[请求频率过高]
    end

    subgraph 防护措施
        D1[用户自行保管<br/>Cookie有时效性]
        D2[建议HTTPS部署<br/>加密传输]
        D3[不存储凭证<br/>无持久化数据可窃取]
        D4[可配置限流中间件]
    end

    subgraph 影响范围
        I1[仅影响泄露者本人]
        I2[Cookie过期后失效]
        I3[不影响其他用户]
        I4[教务系统可封禁IP]
    end

    T1 --> D1 --> I1
    T2 --> D2 --> I2
    T3 --> D3 --> I3
    T4 --> D4 --> I4
```

---

## 五、对教务系统的影响评估

### 5.1 数据安全影响

| 风险点 | 风险等级 | 说明 |
|--------|----------|------|
| 数据泄露 | **低** | API仅能获取用户自己的数据，不能越权访问他人数据 |
| 数据篡改 | **无** | API为只读操作，不提供任何写入/修改接口 |
| 数据完整性 | **无影响** | 不对教务系统数据进行任何修改 |
| 凭证安全 | **低** | 不存储密码，Cookie有时效性且由用户自行管理 |

### 5.2 系统负载影响

```mermaid
graph LR
    subgraph 请求特征
        A[单用户低频请求]
        B[仅查询操作]
        C[无批量爬取设计]
    end

    subgraph 对教务系统影响
        D[等同于用户手动访问]
        E[不增加额外负担]
        F[可被教务系统限流]
    end

    A --> D
    B --> E
    C --> F
```

### 5.3 安全结论

```mermaid
graph TD
    Q[曲奇API会对教务系统<br/>造成数据安全隐患吗?]

    Q --> A1[不会造成数据泄露]
    Q --> A2[不会造成数据篡改]
    Q --> A3[不会造成系统负担]

    A1 --> R1[只能访问Cookie持有者自己的数据<br/>权限由教务系统控制]
    A2 --> R2[所有接口均为只读GET请求<br/>无任何写入操作]
    A3 --> R3[请求频率等同手动操作<br/>无批量爬取机制]

    R1 & R2 & R3 --> C[结论：安全风险极低<br/>本质是数据格式转换服务]
```

---

## 六、技术实现细节

### 6.1 Cookie 转发机制

```go
// services/zhjw/client.go
func NewClient(Authorization string) *resty.Client {
    client := resty.New()
    client.SetBaseURL("http://zhjw.qfnu.edu.cn")
    client.SetHeader("Cookie", Authorization)  // 关键：Cookie转发
    // ...
    return client
}
```

### 6.2 会话失效检测

```go
// services/zhjw/client.go - 响应拦截器
client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
    body := resp.String()
    // 检测是否被重定向到登录页
    if strings.Contains(body, "用户登录") {
        return ErrCookieExpired
    }
    return nil
})
```

### 6.3 错误码映射

| API 错误码 | HTTP 状态 | 含义 | 教务系统状态 |
|-----------|----------|------|-------------|
| 200 | 200 | 成功 | 正常响应 |
| 401 | 401 | Cookie已过期 | 重定向到登录页 |
| 404 | 200 | 无数据 | 返回"未查询到数据" |
| 502 | 502 | 教务系统不可用 | 连接超时/错误 |

---

## 七、部署安全建议

1. **启用 HTTPS**：确保客户端与API服务之间的通信加密
2. **配置限流**：防止单用户高频请求对教务系统造成压力
3. **日志脱敏**：确保日志中不记录 Cookie 等敏感信息
4. **定期审计**：检查是否有异常的请求模式

---

## 八、总结

曲奇API采用**透传代理**模式，具有以下安全特性：

- **最小权限原则**：仅转发用户自己的Cookie，无法越权
- **无状态设计**：不存储任何用户凭证，无持久化攻击面
- **只读操作**：所有接口均为查询，不修改任何数据
- **会话隔离**：每个请求独立，Cookie仅在单次请求中使用

**对教务系统的影响**：等同于用户通过浏览器手动访问，不会造成数据安全隐患。
