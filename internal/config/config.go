package config

import (
	"time"

	"github.com/W1ndys/easy-qfnu-api-go/internal/crypto"
	"github.com/W1ndys/easy-qfnu-api-go/internal/database"
)

const (
	KeySiteAccessEnabled  = "site_access_enabled"
	KeySiteAccessPassword = "site_access_password"
	KeyAdminPassword      = "admin_password"
	KeyTokenExpireHours   = "token_expire_hours"
)

// Get 获取配置值
func Get(key string) string {
	db := database.GetAppDB()
	if db == nil {
		return ""
	}

	var value string
	err := db.QueryRow("SELECT value FROM system_config WHERE key = ?", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

// Set 设置配置值
func Set(key, value string) error {
	db := database.GetAppDB()
	if db == nil {
		return nil
	}

	now := time.Now().Unix()
	_, err := db.Exec(`
		INSERT INTO system_config (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = ?
	`, key, value, now, value, now)
	return err
}

// IsSiteAccessEnabled 是否开启访问验证
func IsSiteAccessEnabled() bool {
	// 1. 如果没设置密码，视为无需验证
	if Get(KeySiteAccessPassword) == "" {
		return false
	}

	// 2. 检查开关状态
	value := Get(KeySiteAccessEnabled)
	// 如果没初始化（为空）或明确关闭（false），视为无需验证
	if value == "" || value == "false" {
		return false
	}

	return true
}

// GetTokenExpireHours 获取 Token 过期时间（小时）
func GetTokenExpireHours() int {
	value := Get(KeyTokenExpireHours)
	if value == "" {
		return 24 // 默认 24 小时
	}
	var hours int
	if _, err := time.ParseDuration(value + "h"); err == nil {
		hours = 24
	}
	return hours
}

// VerifySitePassword 验证访问密码
func VerifySitePassword(password string) bool {
	hash := Get(KeySiteAccessPassword)
	if hash == "" {
		return false
	}
	return crypto.CheckPassword(password, hash)
}

// VerifyAdminPassword 验证管理员密码
func VerifyAdminPassword(password string) bool {
	hash := Get(KeyAdminPassword)
	if hash == "" {
		return false
	}
	return crypto.CheckPassword(password, hash)
}

// SetSitePassword 设置访问密码
func SetSitePassword(password string) error {
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}
	return Set(KeySiteAccessPassword, hash)
}

// SetAdminPassword 设置管理员密码
func SetAdminPassword(password string) error {
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}
	return Set(KeyAdminPassword, hash)
}

// Announcement 公告结构
type Announcement struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	SortOrder int    `json:"sort_order"`
}

// GetActiveAnnouncements 获取启用的公告列表
func GetActiveAnnouncements() []Announcement {
	db := database.GetAppDB()
	if db == nil {
		return []Announcement{}
	}

	rows, err := db.Query(`
		SELECT id, title, content, type, sort_order
		FROM announcements
		WHERE is_active = 1
		ORDER BY sort_order ASC, id DESC
	`)
	if err != nil {
		return []Announcement{}
	}
	defer rows.Close()

	var list []Announcement
	for rows.Next() {
		var a Announcement
		rows.Scan(&a.ID, &a.Title, &a.Content, &a.Type, &a.SortOrder)
		list = append(list, a)
	}

	if list == nil {
		return []Announcement{}
	}
	return list
}
