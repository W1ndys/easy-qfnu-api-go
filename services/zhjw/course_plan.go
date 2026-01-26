package zhjw

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/W1ndys/qfnu-api-go/model"
)

// FetchCoursePlan 获取培养方案
func FetchCoursePlan(token string) (*model.CoursePlanResponse, error) {
	// 1. 请求页面
	url := "http://zhjw.qfnu.edu.cn/jsxsd/pyfa/topyfamx"
	client := NewClient(token)
	resp, err := client.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// 2. 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, fmt.Errorf("parse html failed: %w", err)
	}

	response := &model.CoursePlanResponse{}

	// 3. 提取培养目标 (第一个 span#pymb)
	response.Objectives = strings.TrimSpace(doc.Find("span#pymb").First().Text())

	// 4. 提取详细说明 (第二个 span#pymb，或者根据上下文查找)
	// 因为页面上有两个 id="pymb" 的元素，goquery 的 Find("#id") 可能只会返回第一个或所有。
	// 这里尝试获取所有并取第二个，如果只有一个则为空
	pymbSelection := doc.Find("span#pymb")
	if pymbSelection.Length() >= 2 {
		response.Description = strings.TrimSpace(pymbSelection.Eq(1).Text())
	}

	// 5. 解析课程列表
	response.Groups = parseCourseGroups(doc)

	return response, nil
}

func parseCourseGroups(doc *goquery.Document) []model.CourseGroup {
	var groups []model.CourseGroup
	var currentGroup *model.CourseGroup

	// 查找表格 #mxh 的 tbody 下的所有 tr
	doc.Find("table#mxh tbody tr").Each(func(i int, s *goquery.Selection) {
		// 跳过表头行 (通常有 class="Nsb_r_list_thb" 或者没有数据单元格)
		if s.Find("th").Length() > 0 {
			return
		}

		// 跳过小计行 (含有 "小计" 文本)
		if strings.Contains(s.Text(), "小计") {
			return
		}

		cells := s.Find("td")
		// 根据单元格数量判断是否是新的一组
		// 新的一组的第一行通常包含 rowspan 的课程体系/选课组名称
		// 观察 HTML 结构:
		// 第一行: [课程体系(rowspan), 选课组(rowspan?), 课程编号, 课程名称, 完成情况, 课程性质, 课程属性, 学分, ..., 开设学期]
		// 后续行: [选课组(如果上面没合并?), 课程编号, 课程名称, ...]
		// 注意：HTML 示例中，"课程体系" 和 "选课组" 实际上是在第一列合并显示的，格式如 "通识课-身心健康课组2-国家安全教育 (应修 1 / 已修 1)"

		// 逻辑：检查第一列是否包含 "应修" 和 "已修" 关键字，如果是，则视为新组的开始
		firstCellText := strings.TrimSpace(cells.First().Text())
		isNewGroup := strings.Contains(firstCellText, "应修") && strings.Contains(firstCellText, "已修")

		var course model.CoursePlanInfo
		var cellOffset int //用于调整后续单元格的索引

		if isNewGroup {
			// 保存上一组
			if currentGroup != nil {
				groups = append(groups, *currentGroup)
			}

			// 解析新组信息
			groupName, req, earned := parseGroupHeader(firstCellText)
			currentGroup = &model.CourseGroup{
				GroupName:       groupName,
				RequiredCredits: req,
				EarnedCredits:   earned,
				Courses:         []model.CoursePlanInfo{},
			}

			// 新组的第一行，课程信息从第3个单元格开始 (索引2) ?
			// 让我们看 HTML:
			// Cell 0: Group Name (rowspan)
			// Cell 1: empty (&nbsp;) or Checkbox?
			// Cell 2: Course Code (580001)
			// Cell 3: Course Name
			// ...
			cellOffset = 2
		} else {
			// 不是新组，说明是该组的后续课程
			// 这种行的第一列通常是空的或者 Checkbox占位?
			// 看 HTML 例子:
			// <TR> <TD>&nbsp;</TD> <TD>306061</TD> ... </TR>
			// 所以索引 0 是空列，索引 1 是 Course Code
			cellOffset = 1
		}

		// 提取课程信息
		// 注意：如果 isNewGroup，code 在 cells[2]；否则在 cells[1]
		if cells.Length() > cellOffset+5 {
			course.CourseCode = strings.TrimSpace(cells.Eq(cellOffset).Text())
			course.CourseName = strings.TrimSpace(cells.Eq(cellOffset + 1).Text())
			course.Status = strings.TrimSpace(cells.Eq(cellOffset + 2).Text())
			course.CourseProp = strings.TrimSpace(cells.Eq(cellOffset + 3).Text())
			course.CourseAttr = strings.TrimSpace(cells.Eq(cellOffset + 4).Text())

			creditsStr := strings.TrimSpace(cells.Eq(cellOffset + 5).Text())
			course.Credits, _ = strconv.ParseFloat(creditsStr, 64)

			// 总学时是倒数第二列 (HTML 中是 cells.Length()-2 ? 或者是固定位置?)
			// HTML 表头有: 讲课, 实践, 讲座, 实验, 设计, 上机, 讨论, 课外, 网络, 总学时(rowspan or separate?)
			// 总学时在 HTML 中是倒数第2个 td (倒数第1个是开设学期)
			// 为了稳健，我们取 cells.Last().Prev().Text() ?
			// 让我们数一下列数：
			// Base columns after CourseAttr(idx+4): Credits(idx+5), 9 types of hours, Total Hours, Term
			// So Total Hours is at index + 5 + 9 + 1 = index + 15
			// Term is at index + 16 (Last)

			// 让我们用相对位置：倒数第一个是学期，倒数第二个是总学时
			course.Term = strings.TrimSpace(cells.Last().Text())

			// 处理总学时里面可能包含 input 标签的情况
			hoursText := strings.TrimSpace(cells.Eq(cells.Length() - 2).Text())
			course.Hours = hoursText
		}

		if currentGroup != nil && course.CourseName != "" {
			currentGroup.Courses = append(currentGroup.Courses, course)
		}
	})

	// 添加最后一组
	if currentGroup != nil {
		groups = append(groups, *currentGroup)
	}

	return groups
}

// parseGroupHeader 解析 "通识课... (应修 1 / 已修 1)"
func parseGroupHeader(text string) (name string, required float64, earned float64) {
	// 正则匹配: (.*)\(应修\s*([\d\.]+)\s*/\s*已修\s*([\d\.]+)\)
	re := regexp.MustCompile(`(.*?)\s*[\(（]应修\s*([\d\.]+)\s*/\s*已修\s*([\d\.]+)[）\)]`)
	matches := re.FindStringSubmatch(text)

	if len(matches) == 4 {
		name = strings.TrimSpace(matches[1])
		required, _ = strconv.ParseFloat(matches[2], 64)
		earned, _ = strconv.ParseFloat(matches[3], 64)
	} else {
		name = text // Fallback
	}
	return
}
