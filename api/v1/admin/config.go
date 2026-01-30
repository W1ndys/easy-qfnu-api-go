package admin

import (
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/internal/config"
	"github.com/gin-gonic/gin"
)

// GetConfig 获取所有配置
func GetConfig(c *gin.Context) {
	response.Success(c, gin.H{
		"site_access_enabled": config.IsSiteAccessEnabled(),
		"token_expire_hours":  config.GetTokenExpireHours(),
	})
}

type UpdateConfigRequest struct {
	SiteAccessEnabled  *bool   `json:"site_access_enabled"`
	SiteAccessPassword *string `json:"site_access_password"`
	AdminPassword      *string `json:"admin_password"`
	TokenExpireHours   *string `json:"token_expire_hours"`
}

// UpdateConfig 更新配置
func UpdateConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if req.SiteAccessEnabled != nil {
		if *req.SiteAccessEnabled {
			config.Set(config.KeySiteAccessEnabled, "true")
		} else {
			config.Set(config.KeySiteAccessEnabled, "false")
		}
	}

	if req.SiteAccessPassword != nil && *req.SiteAccessPassword != "" {
		config.SetSitePassword(*req.SiteAccessPassword)
	}

	if req.AdminPassword != nil && *req.AdminPassword != "" {
		config.SetAdminPassword(*req.AdminPassword)
	}

	if req.TokenExpireHours != nil {
		config.Set(config.KeyTokenExpireHours, *req.TokenExpireHours)
	}

	response.Success(c, gin.H{})
}
