package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	// 假设这是你原本的包路径，保持不变

	v1 "github.com/W1ndys/qfnu-api-go/api/v1"
	"github.com/W1ndys/qfnu-api-go/common/logger"
	"github.com/W1ndys/qfnu-api-go/middleware"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------
// 1. 嵌入 web 目录下的所有文件
// ---------------------------------------------------------
//
//go:embed web
var webFS embed.FS

func main() {
	// 初始化日志
	logger.InitLogger("./logs", "qfnu-api.log", "info")

	r := gin.Default()

	// 注册中间件
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Cors())

	// 注册 API 路由 (优先匹配)
	// 创建一个带鉴权中间件的路由组
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthRequired())

	{
		// 这里的 Handler 不需要再写检查 Token 的代码了
		apiGroup.GET("/grades", v1.GetGradeList)
		// 以后加的新接口也不用写：
		// apiGroup.GET("/schedule", v1.GetSchedule)
	}

	//  核心：实现根目录挂载静态资源 (作为兜底逻辑)
	// 第一步：剥离 "web" 这一层目录
	// 这样访问时不需要带 /web 前缀，直接对应 web 目录内部结构
	staticFiles, _ := fs.Sub(webFS, "web")

	// 第二步：创建标准的文件服务器
	// http.FileServer 具备以下自动功能：
	// 1. 访问 / -> 自动寻找 index.html
	// 2. 访问 /about.html -> 寻找 about.html
	// 3. 访问 /css/style.css -> 寻找 css/style.css
	fileServer := http.FileServer(http.FS(staticFiles))

	// 第三步：使用 NoRoute 作为静态资源入口
	// 逻辑：如果请求没有命中上面的 /api 路由，就会进入这里
	r.NoRoute(func(c *gin.Context) {
		// 出于安全考虑，可以拦截一下非 GET/HEAD 请求
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusMethodNotAllowed)
			return
		}

		// 将请求转交给 Go 原生的文件服务器处理
		// 它会自动处理 Content-Type、Content-Length 和 404
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	// ---------------------------------------------------------
	// 启动提示
	// ---------------------------------------------------------
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	log.Println(green("√ 服务器启动成功！"))
	log.Println(cyan("➜ 接口地址: http://localhost:8080/api/grades"))
	log.Println(cyan("➜ 网页首页: http://localhost:8080/")) // 直接访问根路径
	log.Println(red("! 注意: 请勿关闭此窗口"))

	r.Run("0.0.0.0:8080")
}
