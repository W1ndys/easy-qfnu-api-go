package middleware

import (
	"time"

	"github.com/W1ndys/easy-qfnu-api-go/common/stats"
	"github.com/gin-gonic/gin"
)

// StatsCollector 统计收集中间件
func StatsCollector() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start).Milliseconds()

		// 收集日志
		stats.Collect(stats.RequestLog{
			Path:       c.Request.URL.Path,
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			LatencyMs:  latency,
			ClientIP:   c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			CreatedAt:  time.Now().Unix(),
		})
	}
}
