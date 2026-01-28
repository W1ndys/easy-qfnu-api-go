package zhjw

import (
	"bytes"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/W1ndys/easy-qfnu-api-go/model"
)

// FetchClassSchedules 抓取并解析课程表
func FetchClassSchedules(cookie string, date string) (*model.ClassScheduleResponse, error) {

	// 使用工厂函数创建 Client (自带检查功能)
	client := NewClient(cookie)

	targetURL := "http://zhjw.qfnu.edu.cn/jsxsd/framework/main_index_loadkb.jsp"
	formData := map[string]string{
		"rq": strings.TrimSpace(date), // 日期
	}

	// 记录重要的业务行为
	slog.Info("开始抓取课程表",
		"date", date,
		"cookie_len", len(cookie), // 不要记录完整 cookie，记录长度即可，保护隐私
	)
	// 发起 POST 请求
	resp, err := client.R().
		SetFormData(formData).
		Post(targetURL)

	// log.Printf("响应内容：%s", resp.Body())

	// 错误处理
	if err != nil {
		return nil, err // 遇到错误立刻返回
	}

	// 解析 HTML
	return parseClassSchedulesHtml(resp.Body())
}

// parseClassSchedulesHtml 解析课程表 HTML
func parseClassSchedulesHtml(htmlBody []byte) (*model.ClassScheduleResponse, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	response := &model.ClassScheduleResponse{}
	htmlStr := string(htmlBody)

	// 1. 获取当前周次信息
	// 从 script 中提取: $("#li_showWeek").html("<span class=\"main_text main_color\">第18周</span>/20周");
	// 只匹配包含 span 标签的内容，排除 "当前登录已失效，请重新登录！" 这种纯文本
	// 匹配的内容：
	// - "<span class=\"main_text main_color\">第18周</span>/20周"
	// - "<span class=\"main_text main_color\">当前日期不在教学周历内</span>"
	reWeek := regexp.MustCompile(`\$\("#li_showWeek"\)\.html\("(<span[^>]*>.*?</span>.*?)"\);`)
	matches := reWeek.FindStringSubmatch(htmlStr)
	if len(matches) > 1 {
		// 清除 html 标签
		reTag := regexp.MustCompile(`<[^>]+>`)
		cleanWeek := reTag.ReplaceAllString(matches[1], "")
		// 处理转义字符
		cleanWeek = strings.ReplaceAll(cleanWeek, `\"`, `"`)
		response.CurrentWeekRaw = cleanWeek
	}

	var courses []model.ClassSchedules
	index := 1

	// 2. 遍历表格中的 p 标签 (带有 title 属性)
	// 格式: 课程学分：3<br/>课程属性：任选<br/>课程名称：网络管理<br/>上课时间：第18周 星期一 [02-03-04]节<br/>上课地点：嵌入式实验室204<br/>课堂名称：23网安班,22网安班
	doc.Find(".kb_table p[title]").Each(func(i int, s *goquery.Selection) {
		titleVal, exists := s.Attr("title")
		if !exists {
			return
		}

		parts := strings.Split(titleVal, "<br/>")
		course := model.ClassSchedules{
			Index: index,
		}

		for _, part := range parts {
			kv := strings.SplitN(part, "：", 2)
			if len(kv) != 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])

			switch key {
			case "课程学分":
				course.Credit = val
			case "课程属性":
				course.Category = val
			case "课程名称":
				course.Name = val
			case "上课时间":
				course.RawTimeString = val
				course.TimeParsed = parseTime(val)
			case "上课地点":
				course.Location = val
			case "课堂名称":
				course.Classes = val
			}
		}
		courses = append(courses, course)
		index++
	})

	response.Courses = courses
	return response, nil
}

func parseTime(raw string) model.ClassTimeParse {
	// 示例: 第18周 星期一 [02-03-04]节
	result := model.ClassTimeParse{}

	// 解析周
	reWeek := regexp.MustCompile(`第(\d+)周`)
	weekMatch := reWeek.FindStringSubmatch(raw)
	if len(weekMatch) > 1 {
		w, _ := strconv.Atoi(weekMatch[1])
		result.Week = w
	}

	// 解析星期
	if strings.Contains(raw, "星期一") {
		result.DayOfWeek = 1
	} else if strings.Contains(raw, "星期二") {
		result.DayOfWeek = 2
	} else if strings.Contains(raw, "星期三") {
		result.DayOfWeek = 3
	} else if strings.Contains(raw, "星期四") {
		result.DayOfWeek = 4
	} else if strings.Contains(raw, "星期五") {
		result.DayOfWeek = 5
	} else if strings.Contains(raw, "星期六") {
		result.DayOfWeek = 6
	} else if strings.Contains(raw, "星期日") {
		result.DayOfWeek = 7
	}

	// 解析节次 [02-03-04]节
	rePeriod := regexp.MustCompile(`\[([\d-]+)\]`)
	periodMatch := rePeriod.FindStringSubmatch(raw)
	if len(periodMatch) > 1 {
		pStr := periodMatch[1] // 02-03-04
		pParts := strings.Split(pStr, "-")
		var periods []int
		for _, p := range pParts {
			val, err := strconv.Atoi(p)
			if err == nil {
				periods = append(periods, val)
			}
		}
		result.PeriodArray = periods
	}

	return result
}
