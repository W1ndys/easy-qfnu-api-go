package zhjw

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-go/common/request"
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	zhjwService "github.com/W1ndys/easy-qfnu-api-go/services/zhjw"
	"github.com/gin-gonic/gin"
)

// GetCoursePlan 是给 Gin 用的处理函数
func GetCoursePlan(c *gin.Context) {

	// 获取参数，能放行到这里，说明已经通过鉴权中间件检查
	Authorization := request.GetCurrentUserAuthorization(c)

	// 调用业务逻辑 (Service 层)
	// 这里的 FetchCoursePlan 首字母是大写，所以能被跨包调用
	data, err := zhjwService.FetchCoursePlan(Authorization)
	// 处理业务结果
	// 如果有错误，返回错误信息
	if errors.Is(err, zhjwService.ErrCookieExpired) {
		response.CookieExpired(c)
		return
	} else if errors.Is(err, zhjwService.ErrResourceNotFound) {
		response.ResourceNotFound(c)
		return
	} else if err != nil {
		response.FailWithCode(c, 1, "获取培养方案失败: "+err.Error())
		return
	}
	response.Success(c, data)

}
