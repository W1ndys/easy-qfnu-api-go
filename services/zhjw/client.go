package zhjw

import (
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
)

// ErrCookieExpired 定义一个哨兵错误 (Sentinel Error)
// 这样在 API Handler 层可以通过 errors.Is() 精确捕获它，然后返回 401
var ErrCookieExpired = errors.New("cookie_expired_or_invalid")

// NewJwcClient 创建一个配置好“自动检查机制”的 Resty 客户端
func NewClient(Authorization string) *resty.Client {
	client := resty.New()

	// 1. 统一设置 Header (避免在每个请求里重复写)
	client.SetHeader("Cookie", Authorization)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	client.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	// 2. 核心：注册响应拦截器 (Middleware)
	// 每次请求回来，在你的业务代码运行之前，这个函数会先运行
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// 获取响应内容
		body := resp.String()

		// 统一检查规则：
		// 强智教务系统 Cookie 失效时，通常会重定向到登录页，或者 body 里包含特定文字
		// 这里的判断条件根据你的实际情况添加
		if resp.StatusCode() != 200 ||
			strings.Contains(body, "用户登录") {

			// 如果命中，返回特定错误
			// 注意：这里返回 error 会导致后续的 API 请求直接报错返回，
			// 不会再执行 fetchGrade 里的 parseHtml 逻辑
			return ErrCookieExpired
		}

		return nil // 正常，继续执行
	})

	return client
}
