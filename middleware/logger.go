package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 记录每次 HTTP 请求的详细信息
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 请求处理完后，记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 组装 Log 属性
		args := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		}
		if raw != "" {
			args = append(args, slog.String("query", raw))
		}
		if errorMessage != "" {
			args = append(args, slog.String("error", errorMessage))
		}

		// 根据状态码决定日志级别
		if statusCode >= 500 {
			slog.Error("Request Failed", args...)
		} else if statusCode >= 400 {
			slog.Warn("Request Client Error", args...)
		} else {
			slog.Info("Request Success", args...)
		}
	}
}
