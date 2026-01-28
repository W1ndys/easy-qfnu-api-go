package zhjw

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-go/common/request"
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/model"
	zhjwService "github.com/W1ndys/easy-qfnu-api-go/services/zhjw"
	"github.com/gin-gonic/gin"
)

// GetClassSchedules 是给 Gin 用的处理函数
func GetClassSchedules(c *gin.Context) {

	// 获取参数，能放行到这里，说明已经通过鉴权中间件检查
	Authorization := request.GetCurrentUserAuthorization(c)

	// 绑定查询参数到结构体
	var req model.ClassSchedulesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "查询参数错误，请检查后重试")
		return
	}

	// 调用业务逻辑 (Service 层)
	// 这里的 FetchClassSchedules 首字母是大写，所以能被跨包调用
	data, err := zhjwService.FetchClassSchedules(Authorization, req.Date)
	// 处理业务结果
	// 如果有错误，返回错误信息
	if errors.Is(err, zhjwService.ErrCookieExpired) {
		response.CookieExpired(c)
		return
	} else if errors.Is(err, zhjwService.ErrResourceNotFound) {
		response.ResourceNotFound(c)
		return
	} else if err != nil {
		response.FailWithCode(c, 1, "获取课程表失败: "+err.Error())
		return
	}
	response.Success(c, data)

}
