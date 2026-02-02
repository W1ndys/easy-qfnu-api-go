package main

import (
	"embed"
	"flag"
	"log"
	"os"

	"github.com/W1ndys/easy-qfnu-api-go/common/logger"
	"github.com/W1ndys/easy-qfnu-api-go/common/stats"
	"github.com/W1ndys/easy-qfnu-api-go/internal/config"
	"github.com/W1ndys/easy-qfnu-api-go/router"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// ---------------------------------------------------------
// 1. 嵌入 web 目录下的静态资源和模板
// ---------------------------------------------------------
// 只嵌入运行时需要的文件，排除构建工具和源文件
//
//go:embed web/templates web/static/css/tailwind.css web/static/js web/static/favico.ico
var webFS embed.FS

func main() {
	// 定义命令行参数
	resetPwd := flag.String("reset-password", "", "重置管理员密码")
	flag.Parse()

	// 处理重置密码逻辑
	if *resetPwd != "" {
		if err := config.SetAdminPassword(*resetPwd); err != nil {
			log.Fatalf("重置密码失败: %v", err)
		}
		log.Printf("管理员密码已成功重置为: %s", *resetPwd)
		return
	}

	// 尝试加载 .env 文件，忽略错误（因为环境变量可能已经存在）
	_ = godotenv.Load()

	// 设置 Gin 为 Release 模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化日志
	logger.InitLogger("./logs", "easy-qfnu-api", "info")

	// 初始化统计模块
	stats.InitCollector()
	stats.RecordStartTime()

	// 初始化路由 (注入 webFS)
	r := router.InitRouter(webFS)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8141"
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
	log.Println(cyan("➜ 网页首页: http://127.0.0.1:" + port + "/"))
	log.Println(cyan("➜ 后台地址: http://127.0.0.1:" + port + "/admin"))
	log.Println(red("! 注意: 请勿关闭此窗口"))
}
