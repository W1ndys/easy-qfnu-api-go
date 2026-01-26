package main

import (
	"embed"
	"log"

	"github.com/W1ndys/qfnu-api-go/common/logger"
	"github.com/W1ndys/qfnu-api-go/router"
	"github.com/fatih/color"
)

// ---------------------------------------------------------
// 1. 嵌入 web 目录下的所有文件
// ---------------------------------------------------------
//
//go:embed web
var webFS embed.FS

func main() {
	// 初始化日志
	logger.InitLogger("./logs", "easy-qfnu-api", "info")

	// 初始化路由 (注入 webFS)
	r := router.InitRouter(webFS)

	// ---------------------------------------------------------
	// 启动提示
	// ---------------------------------------------------------
	printBanner()

	r.Run("0.0.0.0:8080")
}

func printBanner() {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	log.Println(green("√ 服务器启动成功！"))
	log.Println(cyan("➜ 接口地址: http://localhost:8080/api/v1/"))
	log.Println(cyan("➜ 网页首页: http://localhost:8080/"))
	log.Println(red("! 注意: 请勿关闭此窗口"))
}
