package admin

import (
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/internal/config"
	"github.com/W1ndys/easy-qfnu-api-go/internal/crypto"
	"github.com/W1ndys/easy-qfnu-api-go/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// Login 管理员登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if !config.VerifyAdminPassword(req.Password) {
		response.Fail(c, "密码错误")
		return
	}

	expireHours := 24
	token := crypto.GenerateToken("admin", expireHours)
	expireAt := time.Now().Add(time.Hour * time.Duration(expireHours)).Unix()

	c.SetCookie(middleware.AdminTokenCookie, token, expireHours*3600, "/", "", false, true)

	response.Success(c, gin.H{
		"token":     token,
		"expire_at": expireAt,
	})
}

// Logout 管理员登出
func Logout(c *gin.Context) {
	c.SetCookie(middleware.AdminTokenCookie, "", -1, "/", "", false, true)
	response.Success(c, gin.H{})
}

// CheckInit 检查是否需要初始化
func CheckInit(c *gin.Context) {
	adminPwd := config.Get(config.KeyAdminPassword)
	sitePwd := config.Get(config.KeySiteAccessPassword)
	response.Success(c, gin.H{
		"need_init":       adminPwd == "",
		"site_pwd_set":    sitePwd != "",
	})
}

type InitRequest struct {
	AdminPassword string `json:"admin_password" binding:"required"`
	SitePassword  string `json:"site_password" binding:"required"`
}

// Init 初始化密码
func Init(c *gin.Context) {
	// 检查是否已初始化
	if config.Get(config.KeyAdminPassword) != "" {
		response.Fail(c, "已完成初始化")
		return
	}

	var req InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := config.SetAdminPassword(req.AdminPassword); err != nil {
		response.Fail(c, "设置管理员密码失败")
		return
	}

	if err := config.SetSitePassword(req.SitePassword); err != nil {
		response.Fail(c, "设置访问密码失败")
		return
	}

	response.Success(c, gin.H{})
}
