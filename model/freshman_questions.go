package model

// FreshmanQuestion 对应数据库中的 questions 表
type FreshmanQuestion struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`
	QuestionText string `json:"question"`
	OptionA      string `json:"option_a"`
	OptionB      string `json:"option_b"`
	OptionC      string `json:"option_c"`
	OptionD      string `json:"option_d"`
	OptionAnswer string `json:"option_answer"`
}
