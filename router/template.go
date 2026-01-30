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

	// 加载根目录模板
	loadTemplatesFromDir(webFS, r, baseContent, "web/templates", "")

	// 加载 admin 子目录模板
	loadTemplatesFromDir(webFS, r, baseContent, "web/templates/admin", "admin/")

	return r
}

// loadTemplatesFromDir 从指定目录加载模板
func loadTemplatesFromDir(webFS embed.FS, r multitemplate.Renderer, baseContent, dir, prefix string) {
	entries, err := webFS.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		// 跳过目录
		if entry.IsDir() {
			continue
		}

		// 只处理 .html 文件
		name := entry.Name()
		if !strings.HasSuffix(name, ".html") {
			continue
		}

		// 读取页面模板
		pagePath := dir + "/" + name
		pageBytes, err := webFS.ReadFile(pagePath)
		if err != nil {
			continue
		}

		// 组合 base + page
		tmpl := template.Must(template.New("base").Parse(baseContent))
		template.Must(tmpl.Parse(string(pageBytes)))
		r.Add(prefix+name, tmpl)
	}
}
