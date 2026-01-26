package model // 语法点 1: 包声明

// ExamSchedule 定义考试安排结构
// 每条数据返回的格式分别是：序号	校区	考试场次	课程编号	课程名称	授课教师	考试时间	考场	座位号	准考证号	备注	操作
type ExamSchedule struct {
	Index       string `json:"index"`        // 序号
	Campus      string `json:"campus"`       // 校区
	Session     string `json:"session"`      // 考试场次
	CourseId    string `json:"course_id"`    // 课程编号
	CourseName  string `json:"course_name"`  // 课程名称
	Instructor  string `json:"instructor"`   // 授课教师
	ExamTime    string `json:"exam_time"`    // 考试时间
	ExamRoom    string `json:"exam_room"`    // 考场
	SeatNumber  string `json:"seat_number"`  // 座位号
	AdmissionNo string `json:"admission_no"` // 准考证号
	Remarks     string `json:"remarks"`      // 备注
	Operation   string `json:"operation"`    // 操作
}

// ExamSchedulesRequest 定义前端查询参数
// Gin 使用 "form" tag 来解析 Query String (?term=...)
type ExamSchedulesRequest struct {
	Term string `form:"term"` // 学期，对应 upstream: kksj
}
