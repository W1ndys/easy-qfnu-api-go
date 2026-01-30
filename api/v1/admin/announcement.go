package admin

import (
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/internal/database"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type AnnouncementRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Type      string `json:"type"`
	IsActive  *bool  `json:"is_active"`
	SortOrder *int   `json:"sort_order"`
}

// GetAnnouncements 获取所有公告
func GetAnnouncements(c *gin.Context) {
	db := database.GetAppDB()
	if db == nil {
		response.Fail(c, "数据库错误")
		return
	}

	rows, err := db.Query(`
		SELECT id, title, content, type, is_active, sort_order, created_at, updated_at
		FROM announcements
		ORDER BY sort_order ASC, id DESC
	`)
	if err != nil {
		response.Fail(c, "查询失败")
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, createdAt, updatedAt int64
		var title, content, aType string
		var isActive, sortOrder int
		rows.Scan(&id, &title, &content, &aType, &isActive, &sortOrder, &createdAt, &updatedAt)
		list = append(list, gin.H{
			"id":         id,
			"title":      title,
			"content":    content,
			"type":       aType,
			"is_active":  isActive == 1,
			"sort_order": sortOrder,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	if list == nil {
		list = []gin.H{}
	}
	response.Success(c, list)
}

// CreateAnnouncement 创建公告
func CreateAnnouncement(c *gin.Context) {
	var req AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	db := database.GetAppDB()
	if db == nil {
		response.Fail(c, "数据库错误")
		return
	}

	aType := "info"
	if req.Type != "" {
		aType = req.Type
	}
	isActive := 1
	if req.IsActive != nil && !*req.IsActive {
		isActive = 0
	}
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	now := time.Now().Unix()
	result, err := db.Exec(`
		INSERT INTO announcements (title, content, type, is_active, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.Title, req.Content, aType, isActive, sortOrder, now, now)

	if err != nil {
		response.Fail(c, "创建失败")
		return
	}

	id, _ := result.LastInsertId()
	response.Success(c, gin.H{"id": id})
}

// UpdateAnnouncement 更新公告
func UpdateAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, "无效的ID")
		return
	}

	var req AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	db := database.GetAppDB()
	if db == nil {
		response.Fail(c, "数据库错误")
		return
	}

	aType := "info"
	if req.Type != "" {
		aType = req.Type
	}
	isActive := 1
	if req.IsActive != nil && !*req.IsActive {
		isActive = 0
	}
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	now := time.Now().Unix()
	_, err = db.Exec(`
		UPDATE announcements
		SET title = ?, content = ?, type = ?, is_active = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, req.Title, req.Content, aType, isActive, sortOrder, now, id)

	if err != nil {
		response.Fail(c, "更新失败")
		return
	}

	response.Success(c, gin.H{})
}

// DeleteAnnouncement 删除公告
func DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, "无效的ID")
		return
	}

	db := database.GetAppDB()
	if db == nil {
		response.Fail(c, "数据库错误")
		return
	}

	_, err = db.Exec("DELETE FROM announcements WHERE id = ?", id)
	if err != nil {
		response.Fail(c, "删除失败")
		return
	}

	response.Success(c, gin.H{})
}
