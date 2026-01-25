package middleware

import (
	"github.com/W1ndys/qfnu-api-go/common/response"
	"github.com/gin-gonic/gin"
)

// AuthRequired 鉴权中间件
// 作用：强制要求请求必须带 Authorization，否则直接拦截
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization
		Authorization := c.GetHeader("Authorization")

		// 2. 检查是否存在
		if Authorization == "" {
			// 如果没有 Authorization，直接报错返回
			response.FailWithCode(c, response.CodeAuthExpired, "缺少鉴权 Authorization (Cookie)")

			// 🛑 核心步骤：Abort
			// 这一步非常重要！它告诉 Gin 停止执行后面的 Handler，直接返回响应。
			c.Abort()
			return
		}

		// 3. 将 Authorization 放入上下文 (Context)
		// 这样后续的 Handler 就可以直接取用，不用再读 Header 了
		c.Set("Authorization", Authorization)

		// 4. 放行，执行下一个 Handler
		c.Next()
	}
}
