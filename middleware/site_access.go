package middleware

import (
	"github.com/W1ndys/easy-qfnu-api-go/internal/config"
	"github.com/W1ndys/easy-qfnu-api-go/internal/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SiteTokenCookie  = "site_token"
	AdminTokenCookie = "admin_token"
)

// SiteAccessRequired 前端访问验证中间件
func SiteAccessRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否开启访问验证
		if !config.IsSiteAccessEnabled() {
			c.Next()
			return
		}

		// 从 Cookie 获取 Token
		token, err := c.Cookie(SiteTokenCookie)
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/access")
			c.Abort()
			return
		}

		// 验证 Token
		if !crypto.ValidateToken(token, "site") {
			c.SetCookie(SiteTokenCookie, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/access")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuthRequired 管理员认证中间件
func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(AdminTokenCookie)
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		if !crypto.ValidateToken(token, "admin") {
			c.SetCookie(AdminTokenCookie, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
