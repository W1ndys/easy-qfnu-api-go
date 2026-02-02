package course_recommendation

import (
	"errors"
	"time"

	"github.com/W1ndys/easy-qfnu-api-go/internal/database"
	"github.com/W1ndys/easy-qfnu-api-go/model"
)

var (
	ErrNotFound = errors.New("推荐记录不存在")
)

// Query 根据关键词查询可见的课程推荐（匹配课程名称或教师姓名）
func Query(keyword string) ([]model.CourseRecommendationPublic, error) {
	db := database.GetCourseRecDB()
	if db == nil {
		return nil, errors.New("数据库连接失败")
	}

	pattern := "%" + keyword + "%"
	rows, err := db.Query(`
		SELECT course_name, teacher_name, recommendation_reason, recommender_nickname, recommendation_time
		FROM course_recommendations
		WHERE is_visible = 1 AND (course_name LIKE ? OR teacher_name LIKE ?)
		ORDER BY recommendation_time DESC
	`, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CourseRecommendationPublic
	for rows.Next() {
		var r model.CourseRecommendationPublic
		if err := rows.Scan(&r.CourseName, &r.TeacherName, &r.RecommendationReason, &r.RecommenderNickname, &r.RecommendationTime); err != nil {
			continue
		}
		list = append(list, r)
	}

	if list == nil {
		return []model.CourseRecommendationPublic{}, nil
	}
	return list, nil
}

// Recommend 提交课程推荐
func Recommend(req model.CourseRecommendationRecommendRequest) (int64, error) {
	db := database.GetCourseRecDB()
	if db == nil {
		return 0, errors.New("数据库连接失败")
	}

	now := time.Now().Unix()
	_, err := db.Exec(`
		INSERT INTO course_recommendations (course_name, teacher_name, recommendation_reason, recommender_nickname, recommendation_time, is_visible)
		VALUES (?, ?, ?, ?, ?, 0)
	`, req.CourseName, req.TeacherName, req.RecommendationReason, req.RecommenderNickname, now)
	if err != nil {
		return 0, err
	}

	return now, nil
}

// Review 审核课程推荐（设置是否可见）
func Review(recommendationID int64, isVisible bool) error {
	db := database.GetCourseRecDB()
	if db == nil {
		return errors.New("数据库连接失败")
	}

	visibleInt := 0
	if isVisible {
		visibleInt = 1
	}

	result, err := db.Exec(`
		UPDATE course_recommendations SET is_visible = ? WHERE id = ?
	`, visibleInt, recommendationID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetAll 获取所有课程推荐（管理员用，包含不可见的）
func GetAll() ([]model.CourseRecommendation, error) {
	db := database.GetCourseRecDB()
	if db == nil {
		return nil, errors.New("数据库连接失败")
	}

	rows, err := db.Query(`
		SELECT id, course_name, teacher_name, recommendation_reason, recommender_nickname, recommendation_time, is_visible
		FROM course_recommendations
		ORDER BY recommendation_time DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CourseRecommendation
	for rows.Next() {
		var r model.CourseRecommendation
		if err := rows.Scan(&r.ID, &r.CourseName, &r.TeacherName, &r.RecommendationReason, &r.RecommenderNickname, &r.RecommendationTime, &r.IsVisible); err != nil {
			continue
		}
		list = append(list, r)
	}

	if list == nil {
		return []model.CourseRecommendation{}, nil
	}
	return list, nil
}

// Delete 删除课程推荐（管理员用）
func Delete(recommendationID int64) error {
	db := database.GetCourseRecDB()
	if db == nil {
		return errors.New("数据库连接失败")
	}

	result, err := db.Exec(`DELETE FROM course_recommendations WHERE id = ?`, recommendationID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
