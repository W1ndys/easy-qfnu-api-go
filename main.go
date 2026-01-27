package main

import (
	"embed"
	"log"
	"os"

	"github.com/W1ndys/qfnu-api-go/common/logger"
	"github.com/W1ndys/qfnu-api-go/router"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// ---------------------------------------------------------
// 1. 嵌入 web 目录下的静态资源
// ---------------------------------------------------------
// 只嵌入运行时需要的文件，排除构建工具和源文件
//
//go:embed web/*.html web/static/css/tailwind.css web/static/js web/static/favico.ico
var webFS embed.FS

func main() {
	// 尝试加载 .env 文件，忽略错误（因为环境变量可能已经存在）
	_ = godotenv.Load()

	// 初始化日志
	logger.InitLogger("./logs", "easy-qfnu-api", "info")

	// 初始化路由 (注入 webFS)
	r := router.InitRouter(webFS)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ---------------------------------------------------------
	// 启动提示
	// ---------------------------------------------------------
	printBanner(port)

	r.Run("0.0.0.0:" + port)
}

func printBanner(port string) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	log.Println(green("√ 服务器启动成功！"))
	log.Println(cyan("➜ 接口地址: http://localhost:" + port + "/api/v1/"))
	log.Println(cyan("➜ 网页首页: http://localhost:" + port + "/"))
	log.Println(red("! 注意: 请勿关闭此窗口"))
}
