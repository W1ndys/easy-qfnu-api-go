package middleware

import (
	"strings"
	"time"

	"github.com/W1ndys/easy-qfnu-api-go/common/stats"
	"github.com/gin-gonic/gin"
)

// StatsCollector 统计收集中间件
func StatsCollector() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 只统计 API 路由，排除静态资源、页面、统计接口和站点公共接口
		shouldCollect := strings.HasPrefix(path, "/api/") &&
			!strings.HasPrefix(path, "/api/v1/stats/") &&
			!strings.HasPrefix(path, "/api/v1/site/")

		var start time.Time
		if shouldCollect {
			start = time.Now()
		}

		// 处理请求
		c.Next()

		// 收集日志
		if shouldCollect {
			latency := time.Since(start).Milliseconds()
			stats.Collect(stats.RequestLog{
				Path:       path,
				Method:     c.Request.Method,
				StatusCode: c.Writer.Status(),
				LatencyMs:  latency,
				ClientIP:   c.ClientIP(),
				UserAgent:  c.Request.UserAgent(),
				CreatedAt:  time.Now().Unix(),
			})
		}
	}
}
