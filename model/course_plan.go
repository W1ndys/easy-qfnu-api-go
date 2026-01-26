package model

// CoursePlanInfo 课程详细信息
type CoursePlanInfo struct {
	CourseName string  `json:"course_name"` // 课程名称
	CourseCode string  `json:"course_code"` // 课程编号
	Status     string  `json:"status"`      // 完成情况 (e.g., 已修(优))
	CourseProp string  `json:"course_prop"` // 课程性质 (e.g., 公共必修课)
	CourseAttr string  `json:"course_attr"` // 课程属性 (e.g., 必修)
	Credits    float64 `json:"credits"`     // 学分
	Hours      string  `json:"hours"`       // 总学时
	Term       string  `json:"term"`        // 开设学期
}

// CourseGroup 选课组信息
type CourseGroup struct {
	GroupName       string           `json:"group_name"`       // 选课组名称
	RequiredCredits float64          `json:"required_credits"` // 应修学分
	EarnedCredits   float64          `json:"earned_credits"`   // 已修学分
	Courses         []CoursePlanInfo `json:"courses"`          // 组内课程列表
}

// CoursePlanResponse 培养方案响应结构
type CoursePlanResponse struct {
	Objectives  string        `json:"objectives"`  // 培养目标
	Description string        `json:"description"` // 详细说明
	Groups      []CourseGroup `json:"groups"`      // 课程组列表
}
