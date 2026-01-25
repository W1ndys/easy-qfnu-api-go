package v1

import (
	"errors"
	"net/http"

	"github.com/W1ndys/qfnu-api-go/common/response"
	"github.com/W1ndys/qfnu-api-go/model"
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

	// 绑定查询参数到结构体
	var req model.GradeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		// 如果绑定失败，可以忽略或报错，这里选择忽略使用默认值
	}

	// 调用业务逻辑 (Service 层)
	// 这里的 FetchGrades 首字母是大写，所以能被跨包调用
	data, err := service.FetchGrades(token, req.Term, req.CourseType, req.CourseName, req.DisplayType)
	// 处理业务结果
	// 如果有错误，返回错误信息
	if errors.Is(err, service.ErrCookieExpired) {
		response.CookieExpired(c)
		return
	}
	if err != nil {
		response.FailWithCode(c, 1, "获取成绩失败: "+err.Error())
		return
	}
	response.Success(c, data)

}
