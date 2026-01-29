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
			slog.String("method", method),     // 请求方法
			slog.String("path", path),         // 请求路径
			slog.Int("status", statusCode),    // 响应状态码
			slog.String("ip", clientIP),       // 客户端 IP
			slog.Duration("latency", latency), // 请求耗时
		}
		if raw != "" {
			args = append(args, slog.String("query", raw))
		}
		if errorMessage != "" {
			args = append(args, slog.String("error", errorMessage))
		}

		// 根据状态码决定日志级别
		if statusCode >= 500 {
			slog.Error("[请求服务器错误]", args...)
		} else if statusCode >= 400 {
			slog.Warn("[请求客户端错误]", args...)
		} else {
			slog.Info("[响应成功]", args...)
		}
	}
}
