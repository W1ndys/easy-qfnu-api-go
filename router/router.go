package router

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/W1ndys/easy-qfnu-api-go/api/v1/questions"
	zhjw "github.com/W1ndys/easy-qfnu-api-go/api/v1/zhjw"
	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/middleware"
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

		// 新生试题库搜索
		apiV1.GET("/questions/search", questions.Search)
	}

	// 【受保护接口组】 (Protected)
	// zhjw 教务系统相关接口
	zhjwGroup := apiV1.Group("/zhjw")
	zhjwGroup.Use(middleware.AuthRequired())
	{
		// 成绩相关接口
		zhjwGroup.GET("/grade", zhjw.GetGradeList)
		// 教学计划/培养方案
		zhjwGroup.GET("/course-plan", zhjw.GetCoursePlan)
		// 考试安排相关接口
		zhjwGroup.GET("/exam", zhjw.GetExamSchedules)
		// 选课结果相关接口
		zhjwGroup.GET("/selection", zhjw.GetSelectionResults)
		// 课程表相关接口
		zhjwGroup.GET("/schedule", zhjw.GetClassSchedules)
	}
}

func installStaticRoutes(r *gin.Engine, webFS embed.FS) {
	// 1. 加载 HTML 模板
	r.HTMLRender = loadTemplates(webFS)

	// 2. 注册静态资源路由
	staticFiles, _ := fs.Sub(webFS, "web/static")
	r.StaticFS("/static", http.FS(staticFiles))

	// 3. 注册页面路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/grade", func(c *gin.Context) {
		c.HTML(http.StatusOK, "grade.html", nil)
	})
	r.GET("/schedule", func(c *gin.Context) {
		c.HTML(http.StatusOK, "schedule.html", nil)
	})
	r.GET("/course-plan", func(c *gin.Context) {
		c.HTML(http.StatusOK, "course-plan.html", nil)
	})
	r.GET("/exam", func(c *gin.Context) {
		c.HTML(http.StatusOK, "exam.html", nil)
	})
	r.GET("/selection", func(c *gin.Context) {
		c.HTML(http.StatusOK, "selection.html", nil)
	})
	r.GET("/questions", func(c *gin.Context) {
		c.HTML(http.StatusOK, "questions.html", nil)
	})
}
