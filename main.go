package main

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// Grade 对应表格中的一行成绩数据
type Grade struct {
	Semester   string `json:"semester"`    // 开课学期
	CourseCode string `json:"course_code"` // 课程编号
	CourseName string `json:"course_name"` // 课程名称
	Score      string `json:"score"`       // 成绩
	Credit     string `json:"credit"`      // 学分
	GPA        string `json:"gpa"`         // 绩点
	ExamType   string `json:"exam_type"`   // 考核方式
	CourseProp string `json:"course_prop"` // 课程性质
}

// ParseGrades 解析 HTML
func ParseGrades(htmlBody []byte) ([]Grade, error) {
	// 1. 加载 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var grades []Grade

	// 2. 定位表格行
	doc.Find("#dataList tr").Each(func(i int, s *goquery.Selection) {
		// 跳过表头
		if i == 0 {
			return
		}

		tds := s.Find("td")
		// 容错处理：如果不满10列，可能是无效行
		if tds.Length() < 10 {
			return
		}

		// 3. 提取数据
		grade := Grade{
			Semester:   strings.TrimSpace(tds.Eq(1).Text()),
			CourseCode: strings.TrimSpace(tds.Eq(2).Text()),
			CourseName: strings.TrimSpace(tds.Eq(3).Text()),
			Score:      strings.TrimSpace(tds.Eq(5).Text()), // 重点关注：这里通常有很多空格
			Credit:     strings.TrimSpace(tds.Eq(7).Text()),
			GPA:        strings.TrimSpace(tds.Eq(9).Text()),
			ExamType:   strings.TrimSpace(tds.Eq(11).Text()),
			CourseProp: strings.TrimSpace(tds.Eq(14).Text()),
		}

		grades = append(grades, grade)
	})

	return grades, nil
}

// Cors 中间件：允许跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求来源

		if origin != "" {
			// 允许的前端来源，"*" 表示允许所有，生产环境建议指定具体域名
			c.Header("Access-Control-Allow-Origin", "*")
			// 允许的 Header，重点是 X-Token
			c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, X-Token")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 浏览器会在发送真实请求前发一个 OPTIONS 预检请求，直接返回 200 即可
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// 1. 应用跨域中间件 (必须放在路由定义之前)
	r.Use(Cors())

	client := resty.New() // 复用 HTTP Client

	r.GET("/api/grades", func(c *gin.Context) {
		// 1. 获取 Cookie
		cookieValue := c.GetHeader("X-Token")
		if cookieValue == "" {
			cookieValue = c.Query("cookie")
		}

		if cookieValue == "" {
			c.JSON(401, gin.H{"code": 401, "msg": "Cookie 缺失，请在 Header(X-Token) 或 URL 参数中提供"})
			return
		}

		// 2. 准备请求
		targetURL := "http://zhjw.qfnu.edu.cn/jsxsd/kscj/cjcx_list"
		formData := map[string]string{
			"kksj": "",
			"kcxz": "",
			"kcmc": "",
			"xsfs": "all",
		}

		// 3. 代理请求
		resp, err := client.R().
			SetHeader("Cookie", cookieValue).
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(formData).
			Post(targetURL)

		if err != nil {
			c.JSON(502, gin.H{"code": 502, "msg": "教务系统连接失败", "error": err.Error()})
			return
		}

		// 4. 检查 Cookie 是否有效
		// 如果教务系统跳转到了登录页，通常会包含特定的 Title 或 input 框
		htmlContent := resp.String() // Resty 自动转 String (默认 UTF-8)
		if strings.Contains(htmlContent, "用户登录") || strings.Contains(htmlContent, "login_btn") {
			c.JSON(401, gin.H{"code": 401, "msg": "Cookie 已失效，请重新登录"})
			return
		}

		// 5. 解析数据
		// 直接传入 Body 字节流
		grades, err := ParseGrades(resp.Body())
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "HTML 解析异常", "error": err.Error()})
			return
		}

		// 6. 返回 JSON
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"list":  grades,
				"count": len(grades),
			},
		})
	})

	r.Run(":8080")
}
