package model

// CourseRecommendation 课程推荐数据结构
type CourseRecommendation struct {
	ID                   int64  `json:"id"`
	CourseName           string `json:"course_name"`
	TeacherName          string `json:"teacher_name"`
	RecommendationReason string `json:"recommendation_reason"`
	RecommenderNickname  string `json:"recommender_nickname"`
	RecommendationTime   int64  `json:"recommendation_time"`
	IsVisible            bool   `json:"is_visible"`
	Campus               string `json:"campus"`              // 校区：曲阜/日照
	RecommendationYear   string `json:"recommendation_year"` // 推荐依据年份
}

// CourseRecommendationPublic 对外展示的课程推荐（不包含 is_visible 字段）
type CourseRecommendationPublic struct {
	CourseName           string `json:"course_name"`
	TeacherName          string `json:"teacher_name"`
	RecommendationReason string `json:"recommendation_reason"`
	RecommenderNickname  string `json:"recommender_nickname"`
	RecommendationTime   int64  `json:"recommendation_time"`
	Campus               string `json:"campus"`
	RecommendationYear   string `json:"recommendation_year"`
}

// CourseRecommendationQueryRequest 查询请求参数
type CourseRecommendationQueryRequest struct {
	Keyword string `form:"keyword" binding:"required"`
}

// CourseRecommendationRecommendRequest 推荐请求参数
type CourseRecommendationRecommendRequest struct {
	CourseName           string `json:"course_name" binding:"required"`
	TeacherName          string `json:"teacher_name" binding:"required"`
	RecommendationReason string `json:"recommendation_reason" binding:"required"`
	RecommenderNickname  string `json:"recommender_nickname"` // 允许为空
	Campus               string `json:"campus" binding:"required"`
	RecommendationYear   string `json:"recommendation_year" binding:"required"`
}

// CourseRecommendationReviewRequest 审核请求参数
type CourseRecommendationReviewRequest struct {
	RecommendationID int64 `json:"recommendation_id" binding:"required"`
	IsVisible        bool  `json:"is_visible"`
}

// CourseRecommendationUpdateRequest 更新请求参数（管理员）
type CourseRecommendationUpdateRequest struct {
	RecommendationID     int64  `json:"recommendation_id" binding:"required"`
	CourseName           string `json:"course_name" binding:"required"`
	TeacherName          string `json:"teacher_name" binding:"required"`
	RecommendationReason string `json:"recommendation_reason" binding:"required"`
	RecommenderNickname  string `json:"recommender_nickname"`
	IsVisible            bool   `json:"is_visible"`
	Campus               string `json:"campus" binding:"required"`
	RecommendationYear   string `json:"recommendation_year" binding:"required"`
}

// CourseRecommendationRecommendResponse 推荐成功响应
type CourseRecommendationRecommendResponse struct {
	Message            string `json:"message"`
	RecommendationTime int64  `json:"recommendation_time"`
}
