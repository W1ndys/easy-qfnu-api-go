package model // 语法点 1: 包声明

// Grade 定义成绩结构
// 语法点 2: 首字母大写 = Public (公开)
type Grade struct {
	Semester   string `json:"semester"`
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	Score      string `json:"score"`
	Credit     string `json:"credit"`
	GPA        string `json:"gpa"`
	ExamType   string `json:"exam_type"`
	CourseProp string `json:"course_prop"`
}
