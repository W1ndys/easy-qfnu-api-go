package zhjw

import (
	"bytes"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/W1ndys/qfnu-api-go/model"
)

// FetchExamSchedules 抓取并解析成绩，返回包含统计信息的响应
func FetchExamSchedules(cookie string, term string) ([]model.ExamSchedule, error) {

	// 使用工厂函数创建 Client (自带检查功能)
	client := NewClient(cookie)

	targetURL := "http://zhjw.qfnu.edu.cn/jsxsd/xsks/xsksap_list"
	formData := map[string]string{
		"xnxqid": strings.TrimSpace(term), // 学期id

	}

	// 记录重要的业务行为
	slog.Info("开始抓取考试安排",
		"term", term,
		"cookie_len", len(cookie), // 不要记录完整 cookie，记录长度即可，保护隐私
	)
	// 发起 POST 请求
	resp, err := client.R().
		SetFormData(formData).
		Post(targetURL)

	// 错误处理
	if err != nil {
		return nil, err // 遇到错误立刻返回
	}

	// 解析 HTML (调用内部私有函数)
	examSchedules, err := parseExamSchedulesHtml(resp.Body())
	if err != nil {
		return nil, err
	}

	return examSchedules, nil
}

// parseExamSchedulesHtml 解析考试安排 HTML
func parseExamSchedulesHtml(htmlBody []byte) ([]model.ExamSchedule, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var schedules []model.ExamSchedule = make([]model.ExamSchedule, 0)

	doc.Find("#dataList tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return // 跳过表头
		}
		tds := s.Find("td")

		// 检查是否是“未查询到数据”的提示行
		if tds.Length() == 1 && strings.Contains(tds.Text(), "未查询到数据") {
			return
		}

		// 正常数据行应该有 12 列
		if tds.Length() < 12 {
			return
		}

		// 组装数据
		// 序号 校区 考试场次 课程编号 课程名称 授课教师 考试时间 考场 座位号 准考证号 备注 操作
		es := model.ExamSchedule{
			Index:       strings.TrimSpace(tds.Eq(0).Text()),
			Campus:      strings.TrimSpace(tds.Eq(1).Text()),
			Session:     strings.TrimSpace(tds.Eq(2).Text()),
			CourseId:    strings.TrimSpace(tds.Eq(3).Text()),
			CourseName:  strings.TrimSpace(tds.Eq(4).Text()),
			Instructor:  strings.TrimSpace(tds.Eq(5).Text()),
			ExamTime:    strings.TrimSpace(tds.Eq(6).Text()),
			ExamRoom:    strings.TrimSpace(tds.Eq(7).Text()),
			SeatNumber:  strings.TrimSpace(tds.Eq(8).Text()),
			AdmissionNo: strings.TrimSpace(tds.Eq(9).Text()),
			Remarks:     strings.TrimSpace(tds.Eq(10).Text()),
			Operation:   strings.TrimSpace(tds.Eq(11).Text()),
		}
		schedules = append(schedules, es)
	})

	return schedules, nil
}
