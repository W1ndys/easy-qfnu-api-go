package model

// {
//   "currentWeekRaw": "第18周/20周",
//   "courses": [
//     {
//       "index": 1,
//       "name": "网络管理",
//       "credit": "3",
//       "category": "任选",
//       "location": "嵌入式实验室204",
//       "classes": "23网安班,22网安班",
//       "rawTimeString": "第18周 星期一 [02-03-04]节",
//       "timeParsed": {
//         "week": 18,
//         "dayOfWeek": 1,
//         "periodArray": [2, 3, 4]
//       }
//     }
//   ]
// }

// ClassScheduleResponse 课程表响应结构
type ClassScheduleResponse struct {
	CurrentWeekRaw string           `json:"currentWeekRaw"` // 当前周次原始字符串
	Courses        []ClassSchedules `json:"courses"`        // 课程列表
}

// ClassSchedules 课程表信息
type ClassSchedules struct {
	Index         int            `json:"index"`         // 课程索引
	Name          string         `json:"name"`          // 课程名称
	Credit        string         `json:"credit"`        // 学分
	Category      string         `json:"category"`      // 课程类别
	Location      string         `json:"location"`      // 上课地点
	Classes       string         `json:"classes"`       // 上课班级
	RawTimeString string         `json:"rawTimeString"` // 原始时间字符串
	TimeParsed    ClassTimeParse `json:"timeParsed"`    // 解析后的时间信息
}

// ClassTimeParse 课程时间解析信息
type ClassTimeParse struct {
	Week        int   `json:"week"`        // 周次
	DayOfWeek   int   `json:"dayOfWeek"`   // 星期几 (1-7)
	PeriodArray []int `json:"periodArray"` // 节次数组
}

// ClassSchedulesRequest 课程表请求结构
type ClassSchedulesRequest struct {
	Date string `form:"date"` // 日期 (e.g., 2026-01-01)
}
