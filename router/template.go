package router

import (
	"embed"
	"html/template"
	"strings"

	"github.com/gin-contrib/multitemplate"
)

// loadTemplates 加载 HTML 模板
// 使用 gin-contrib/multitemplate 实现模板继承
// 自动扫描 web/templates 目录下的所有 .html 文件（排除 layouts 子目录）
func loadTemplates(webFS embed.FS) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 解析布局模板
	baseBytes, _ := webFS.ReadFile("web/templates/layouts/base.html")
	baseContent := string(baseBytes)

	// 动态扫描 templates 目录下的页面文件
	entries, err := webFS.ReadDir("web/templates")
	if err != nil {
		panic("failed to read templates directory: " + err.Error())
	}

	for _, entry := range entries {
		// 跳过目录（如 layouts）
		if entry.IsDir() {
			continue
		}

		// 只处理 .html 文件
		name := entry.Name()
		if !strings.HasSuffix(name, ".html") {
			continue
		}

		// 读取页面模板 (embed.FS 始终使用正斜杠)
		pagePath := "web/templates/" + name
		pageBytes, err := webFS.ReadFile(pagePath)
		if err != nil {
			continue
		}

		// 组合 base + page
		tmpl := template.Must(template.New("base").Parse(baseContent))
		template.Must(tmpl.Parse(string(pageBytes)))
		r.Add(name, tmpl)
	}

	return r
}
