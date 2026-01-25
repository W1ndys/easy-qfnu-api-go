package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	v1 "github.com/W1ndys/qfnu-api-go/api/v1"
	"github.com/W1ndys/qfnu-api-go/common/logger"
	"github.com/W1ndys/qfnu-api-go/middleware"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

//go:embed web
var webFS embed.FS

func main() {
	// 初始化日志 (放在最前面)
	// 日志会保存在 ./logs/app.log
	logger.InitLogger("./logs", "qfnu-api.log", "info")

	r := gin.Default() // Default 默认带了 Logger 和 Recovery，你可以改用 gin.New() 手动添加

	// 注册自定义日志中间件 (它会打印漂亮的 JSON 请求日志)
	r.Use(middleware.RequestLogger())

	// 注册 CORS
	r.Use(middleware.Cors())
	// 注册 API 路由 (优先匹配)
	r.GET("/api/grades", v1.GetGradeList)

	// 处理静态资源 (解决路由冲突 + 支持 Vue History 模式)
	// 从 embed 中剥离 "web" 前缀
	staticFiles, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}

	// 创建文件服务器
	fileServer := http.FileServer(http.FS(staticFiles))

	// 使用 NoRoute 接管所有未定义的路由
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// A. 如果是 API 请求但没匹配到路由，返回 JSON 404 (而不是 HTML)
		if strings.HasPrefix(path, "/api") {
			c.JSON(404, gin.H{"code": 404, "msg": "API Not Found"})
			return
		}

		// B. 处理静态文件
		// 尝试打开请求的文件，检查是否存在
		file, err := staticFiles.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
			// 文件不存在 (说明是 Vue 的前端路由，比如 /grades)
			// 直接返回 index.html，让前端 Vue Router 去处理页面显示
			c.FileFromFS("index.html", http.FS(staticFiles))
			return
		}
		// 记得关闭文件句柄，避免资源泄露
		file.Close()

		// 文件存在 (比如 /static/js/app.js)，直接服务该文件
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	log.Println(green("√ 服务器启动成功！"))
	log.Println(cyan("➜ 前端地址: http://localhost:8080/"))
	log.Println(red("! 注意: 请勿关闭此窗口"))

	r.Run("0.0.0.0:8080")
}
