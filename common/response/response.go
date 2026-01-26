package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构体
type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Result 返回通用响应
func Result[T any](c *gin.Context, httpStatus int, code int, msg string, data T) {
	c.JSON(httpStatus, Response[T]{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功响应 (自动推导 T)
func Success[T any](c *gin.Context, data T) {
	Result(c, http.StatusOK, 0, GetMsg(CodeSuccess), data)
}

// Fail 失败响应
// 修复点：显式指定泛型类型为 [any]
func Fail(c *gin.Context, msg string) {
	// 这里必须写 Result[any]，告诉编译器 T 是 interface{}
	Result[any](c, http.StatusOK, 1, msg, nil)
}

// Cookie 过期响应
func CookieExpired(c *gin.Context) {
	Result[any](c, http.StatusUnauthorized, CodeAuthExpired, GetMsg(CodeAuthExpired), nil)
}

// FailWithCode 自定义错误码响应
// 修复点：显式指定泛型类型为 [any]
func FailWithCode(c *gin.Context, code int, msg string) {
	// 这里必须写 Result[any]
	Result[any](c, http.StatusOK, code, msg, nil)
}
