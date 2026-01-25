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

// GradeRequest 定义前端查询参数
// Gin 使用 "form" tag 来解析 Query String (?term=...)
type GradeRequest struct {
	Term        string `form:"term"`         // 学期，对应 upstream: kksj
	CourseType  string `form:"course_type"`  // 课程性质，对应 upstream: kcxz
	CourseName  string `form:"course_name"`  // 课程名称，对应 upstream: kcmc
	DisplayType string `form:"display_type"` // 显示方式，对应 upstream: xsfs
}

// 课程性质 类型名字与ID的对应
// <option value="01">公共课</option>
// <option value="02">公共基础课</option>
// <option value="03">专业基础课</option>
// <option value="04">专业课</option>
// <option value="05">专业选修课</option>
// <option value="06">公共选修课</option>
// <option value="07">专业任选课</option>
// <option value="08">实践教学环节</option>
// <option value="09">公共任选课</option>
// <option value="10">教师教育基础课程（必修）</option>
// <option value="11">专业必修课</option>
// <option value="12">学科基础必修课</option>
// <option value="13">专业方向限选课</option>
// <option value="14">考试报名虚拟课程</option>
// <option value="15">教师教育选修课程</option>
// <option value="16">公共必修课</option>
type CourseType string

const (
	CourseTypePublic          CourseType = "01" // 公共课
	CourseTypePublicBasic     CourseType = "02" // 公共基础课
	CourseTypeMajorBasic      CourseType = "03" // 专业基础课
	CourseTypeMajor           CourseType = "04" // 专业课
	CourseTypeMajorElective   CourseType = "05" // 专业选修课
	CourseTypePublicElective  CourseType = "06" // 公共选修课
	CourseTypeMajorOptional   CourseType = "07" // 专业任选课
	CourseTypePractical       CourseType = "08" // 实践教学环节
	CourseTypePublicOptional  CourseType = "09" // 公共任选课
	CourseTypeTeacherEduReq   CourseType = "10" // 教师教育基础课程（必修）
	CourseTypeMajorRequired   CourseType = "11" // 专业必修课
	CourseTypeDisciplineBasic CourseType = "12" // 学科基础必修课
	CourseTypeMajorDirection  CourseType = "13" // 专业方向限选课
	CourseTypeExamRegVirtual  CourseType = "14" // 考试报名虚拟课程
	CourseTypeTeacherEduElec  CourseType = "15" // 教师教育选修课程
	CourseTypePublicRequired  CourseType = "16" // 公共必修课
)

// 显示方式 与 value 的对应
// <option value="all">显示全部成绩</option>
// <option value="max">显示最好成绩</option>

// CourseTypeNameToID 中文名称到ID的映射表
var CourseTypeNameToID = map[string]string{
	"公共课":          "01",
	"公共基础课":        "02",
	"专业基础课":        "03",
	"专业课":          "04",
	"专业选修课":        "05",
	"公共选修课":        "06",
	"专业任选课":        "07",
	"实践教学环节":       "08",
	"公共任选课":        "09",
	"教师教育基础课程（必修）": "10",
	"专业必修课":        "11",
	"学科基础必修课":      "12",
	"专业方向限选课":      "13",
	"考试报名虚拟课程":     "14",
	"教师教育选修课程":     "15",
	"公共必修课":        "16",
}

// GetCourseTypeID 根据输入返回课程类型ID
// 如果输入是中文名称则转换为ID，如果已经是ID则直接返回
func GetCourseTypeID(input string) string {
	if input == "" {
		return ""
	}
	// 尝试从映射表查找
	if id, ok := CourseTypeNameToID[input]; ok {
		return id
	}
	// 如果不在映射表中，假设已经是ID，直接返回，不做额外限制，防止后续教务系统更新新增类型而本系统未更新
	return input
}
