package v1

import (
	"net/http"

	"github.com/W1ndys/qfnu-api-go/service"
	"github.com/gin-gonic/gin"
)

// GetGradeList 是给 Gin 用的处理函数
func GetGradeList(c *gin.Context) {
	// 1. 获取参数
	token := c.GetHeader("X-Token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "缺少 Token"})
		return
	}
	term := c.DefaultQuery("term", "")
	courseType := c.DefaultQuery("course_type", "")
	courseName := c.DefaultQuery("course_name", "")
	displayType := c.DefaultQuery("display_type", "all")

	// 2. 调用业务逻辑 (Service 层)
	// 这里的 FetchGrades 首字母是大写，所以能被跨包调用
	data, err := service.FetchGrades(token, term, courseType, courseName, displayType)
	// 3. 处理业务结果
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回 JSON
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}
