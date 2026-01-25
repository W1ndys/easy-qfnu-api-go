package request

import "github.com/gin-gonic/gin"

// GetCurrentUserToken 从上下文中安全获取 Token
func GetCurrentUserToken(c *gin.Context) string {
	// MustGet 取不到会 panic，GetString 取不到返回空字符串
	// 因为中间件已经保证了 token 存在，这里可以用 GetString
	return c.GetString("token")
}
