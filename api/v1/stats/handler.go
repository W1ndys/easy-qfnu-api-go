package stats

import (
	"strconv"

	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	statsService "github.com/W1ndys/easy-qfnu-api-go/services/stats"
	"github.com/gin-gonic/gin"
)

// GetDashboard 获取大屏统计数据
func GetDashboard(c *gin.Context) {
	data, err := statsService.GetDashboardData()
	if err != nil {
		response.Fail(c, "获取统计数据失败: "+err.Error())
		return
	}
	response.Success(c, data)
}

// GetTrend 获取调用趋势
func GetTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, _ := strconv.Atoi(daysStr)

	data, err := statsService.GetTrendData(days)
	if err != nil {
		response.Fail(c, "获取趋势数据失败: "+err.Error())
		return
	}
	response.Success(c, data)
}
