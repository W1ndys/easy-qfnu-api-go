package main

import (
	// 引入 api 包

	v1 "github.com/W1ndys/qfnu-api-go/api/v1"
	"github.com/W1ndys/qfnu-api-go/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 核心修复点：注册 CORS 中间件
	// 必须放在所有路由注册之前！
	r.Use(middleware.Cors())
	// 注册路由
	// 我们把处理函数委托给了 v1 包里的 GetGradeList
	r.GET("/api/grades", v1.GetGradeList)

	r.Run(":8080")
}
