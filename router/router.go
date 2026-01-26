package router

import (
	"embed"
	"io/fs"
	"net/http"

	zhjw "github.com/W1ndys/qfnu-api-go/api/v1/zhjw"
	"github.com/W1ndys/qfnu-api-go/common/response"
	"github.com/W1ndys/qfnu-api-go/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由引擎
func InitRouter(webFS embed.FS) *gin.Engine {
	r := gin.Default()

	// 1. 注册中间件
	installMiddlewares(r)

	// 2. 注册 API 路由
	installAPIRoutes(r)

	// 3. 注册静态资源 (Web)
	installStaticRoutes(r, webFS)

	return r
}

func installMiddlewares(r *gin.Engine) {
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Cors())
}

func installAPIRoutes(r *gin.Engine) {
	// 创建/api根路由组
	apiRoot := r.Group("/api")
	{
		// 健康检查接口
		apiRoot.GET("/health", func(c *gin.Context) {
			response.Success(c, "API is healthy")
		})
	}

	// 创建 v1 根组 (仅用于统一前缀 /api/v1)
	apiV1 := apiRoot.Group("/v1")

	// 【公开接口组】 (Public)
	// 特点：不挂载 AuthRequired 中间件
	{
		// apiV1.GET("/news", v1.GetNewsList)
		// apiV1.GET("/calendar", v1.GetCalendar)
	}

	// 【受保护接口组】 (Protected)
	// zhjw 教务系统相关接口
	zhjwGroup := apiV1.Group("/zhjw")
	zhjwGroup.Use(middleware.AuthRequired())
	{
		// 成绩相关接口
		zhjwGroup.GET("/grades", zhjw.GetGradeList)
		// 教学计划/培养方案
		zhjwGroup.GET("/course-plan", zhjw.GetCoursePlan)
		// 考试安排相关接口
		zhjwGroup.GET("/exam-schedules", zhjw.GetExamSchedules)
	}
}

func installStaticRoutes(r *gin.Engine, webFS embed.FS) {
	// 第一步：剥离 "web" 这一层目录
	staticFiles, _ := fs.Sub(webFS, "web")

	// 第二步：创建标准的文件服务器
	fileServer := http.FileServer(http.FS(staticFiles))

	// 第三步：使用 NoRoute 作为静态资源入口
	r.NoRoute(func(c *gin.Context) {
		// 出于安全考虑，拦截非 GET/HEAD 请求
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusMethodNotAllowed)
			return
		}

		// 将请求转交给 Go 原生的文件服务器处理
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
