package site

import (
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/internal/config"
	"github.com/W1ndys/easy-qfnu-api-go/internal/crypto"
	"github.com/W1ndys/easy-qfnu-api-go/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

type VerifyRequest struct {
	Password string `json:"password" binding:"required"`
}

// Verify 验证访问密码
func Verify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	// 验证密码
	if !config.VerifySitePassword(req.Password) {
		response.Fail(c, "密码错误")
		return
	}

	// 生成 Token
	expireHours := config.GetTokenExpireHours()
	token := crypto.GenerateToken("site", expireHours)
	expireAt := time.Now().Add(time.Hour * time.Duration(expireHours)).Unix()

	// 设置 Cookie
	c.SetCookie(middleware.SiteTokenCookie, token, expireHours*3600, "/", "", false, true)

	response.Success(c, gin.H{
		"token":     token,
		"expire_at": expireAt,
	})
}

// CheckToken 检查 Token 是否有效
func CheckToken(c *gin.Context) {
	token, err := c.Cookie(middleware.SiteTokenCookie)
	if err != nil || token == "" {
		response.Fail(c, "未登录")
		return
	}

	if !crypto.ValidateToken(token, "site") {
		response.Fail(c, "Token 已过期")
		return
	}

	response.Success(c, gin.H{"valid": true})
}

// GetAnnouncements 获取启用的公告列表
func GetAnnouncements(c *gin.Context) {
	announcements := config.GetActiveAnnouncements()
	response.Success(c, announcements)
}
