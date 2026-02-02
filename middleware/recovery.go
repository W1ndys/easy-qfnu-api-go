package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/W1ndys/easy-qfnu-api-go/common/notify"
	"github.com/gin-gonic/gin"
)

// Recovery 自定义恢复中间件，捕获 panic 并发送飞书通知
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				errMsg := fmt.Sprintf("%v", err)
				path := c.Request.URL.Path
				method := c.Request.Method

				// 记录日志
				slog.Error("服务器 Panic",
					"error", errMsg,
					"path", path,
					"method", method,
					"stack", stack,
				)

				// 发送飞书通知
				notify.NotifyError(
					"Panic Recovery",
					fmt.Sprintf("[%s] %s - %s", method, path, errMsg),
					stack,
				)

				// 返回 500 错误
				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}
