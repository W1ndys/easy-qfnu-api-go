package router

import (
	"embed"
	"html/template"

	"github.com/gin-contrib/multitemplate"
)

// loadTemplates 加载 HTML 模板
// 使用 gin-contrib/multitemplate 实现模板继承
func loadTemplates(webFS embed.FS) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 解析布局模板
	baseBytes, _ := webFS.ReadFile("web/templates/layouts/base.html")
	baseContent := string(baseBytes)

	// 为每个页面创建独立的模板集
	pages := []string{"index.html", "grade.html"}
	for _, page := range pages {
		pageBytes, _ := webFS.ReadFile("web/templates/" + page)
		// 组合 base + page
		tmpl := template.Must(template.New("base").Parse(baseContent))
		template.Must(tmpl.Parse(string(pageBytes)))
		r.Add(page, tmpl)
	}

	return r
}
