package model // 语法点 1: 包声明

// SelectionResult 定义选课结果结构
// 序号	课程名称	课程编号	上课老师	总学时	学分	课程属性	课程性质	选课操作人	选课时间
type SelectionResult struct {
	Index      string `json:"index"`       // 序号
	CourseName string `json:"course_name"` // 课程名称
	CourseId   string `json:"course_id"`   // 课程编号
	Teacher    string `json:"teacher"`     // 上课老师
	Hours      string `json:"hours"`       // 总学时
	Credit     string `json:"credit"`      // 学分
	CourseAttr string `json:"course_attr"` // 课程属性
	CourseProp string `json:"course_prop"` // 课程性质
	Operator   string `json:"operator"`    // 选课操作人
	SelectTime string `json:"select_time"` // 选课时间
}

//  SelectionResultsRequest 定义前端查询参数
// Gin 使用 "form" tag 来解析 Query String (?term=...)
type SelectionResultsRequest struct {
	Term string `form:"term"` // 学期，对应 upstream: kksj
}

// SelectionResultsResponse 选课结果查询响应结构
type SelectionResultsResponse struct {
	Results []SelectionResult `json:"results"` // 选课结果列表
}
