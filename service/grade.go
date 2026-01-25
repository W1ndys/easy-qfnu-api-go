package service

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/W1ndys/qfnu-api-go/model"
	"github.com/go-resty/resty/v2"
)

// FetchGrades 抓取并解析成绩
// 返回值：([]model.Grade, error) -> Go 函数支持多返回值
func FetchGrades(cookie string, term string, courseType string, courseName string, displayType string) ([]model.Grade, error) {

	// 课程类型：支持中文名称或ID，统一转换为ID
	courseType = model.GetCourseTypeID(courseType)

	// 准备请求
	client := resty.New()
	targetURL := "http://zhjw.qfnu.edu.cn/jsxsd/kscj/cjcx_list"
	formData := map[string]string{
		"kksj": strings.TrimSpace(term),        // 开课时间
		"kcxz": strings.TrimSpace(courseType),  // 课程性质
		"kcmc": strings.TrimSpace(courseName),  // 课程名称
		"xsfs": strings.TrimSpace(displayType), // 显示方式
	}

	// 发起 POST 请求
	resp, err := client.R().
		SetHeader("Cookie", cookie).
		SetHeader("User-Agent", "Mozilla/5.0...").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(formData).
		Post(targetURL)

	// cookie失效检测，响应代码不是200或者包含“用户登录”视为失效
	if err == nil && (resp.StatusCode() != 200 || strings.Contains(string(resp.Body()), "用户登录")) {
		err = fmt.Errorf("Cookie 失效，请重新获取 Cookie")
		return nil, err
	}

	// 语法点 3: 错误处理习惯
	if err != nil {
		return nil, err // 遇到错误立刻返回
	}

	// 3. 解析 HTML (调用内部私有函数)
	return parseHtml(resp.Body())
}

// parseHtml 是私有函数(小写p)，只在这个文件内部使用，外部不需要知道解析细节
func parseHtml(htmlBody []byte) ([]model.Grade, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var grades []model.Grade // 使用 model.Grade

	doc.Find("#dataList tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		} // 跳过表头
		tds := s.Find("td")
		if tds.Length() < 10 {
			return
		}

		// 组装数据
		g := model.Grade{
			Semester:   strings.TrimSpace(tds.Eq(1).Text()),
			CourseCode: strings.TrimSpace(tds.Eq(2).Text()),
			CourseName: strings.TrimSpace(tds.Eq(3).Text()),
			Score:      strings.TrimSpace(tds.Eq(5).Text()),
			Credit:     strings.TrimSpace(tds.Eq(7).Text()),
			GPA:        strings.TrimSpace(tds.Eq(9).Text()),
			ExamType:   strings.TrimSpace(tds.Eq(11).Text()),
			CourseProp: strings.TrimSpace(tds.Eq(14).Text()),
		}
		grades = append(grades, g)
	})

	// 如果包含"未查询到数据"字样，说明结果为空
	if strings.Contains(string(htmlBody), "未查询到数据") {
		return nil, fmt.Errorf("未查询到数据，请检查查询条件")
	}

	if len(grades) == 0 {
		// 语法点 4: 自定义错误
		return nil, fmt.Errorf("解析结果为空，可能是Cookie失效或页面结构变更")
	}

	return grades, nil
}
