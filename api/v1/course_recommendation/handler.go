package course_recommendation

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-go/common/response"
	"github.com/W1ndys/easy-qfnu-api-go/model"
	services "github.com/W1ndys/easy-qfnu-api-go/services/course_recommendation"
	"github.com/gin-gonic/gin"
)

// Query 查询课程推荐接口
func Query(c *gin.Context) {
	var req model.CourseRecommendationQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "请输入查询关键词")
		return
	}

	list, err := services.Query(req.Keyword)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}

	response.Success(c, list)
}

// Recommend 提交课程推荐接口
func Recommend(c *gin.Context) {
	var req model.CourseRecommendationRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "请填写完整的推荐信息")
		return
	}

	recommendationTime, err := services.Recommend(req)
	if err != nil {
		response.Fail(c, "提交失败: "+err.Error())
		return
	}

	response.Success(c, model.CourseRecommendationRecommendResponse{
		Message:            "推荐成功",
		RecommendationTime: recommendationTime,
	})
}

// Review 审核课程推荐接口（管理员）
func Review(c *gin.Context) {
	var req model.CourseRecommendationReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err := services.Review(req.RecommendationID, req.IsVisible)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			response.FailWithCode(c, response.CodeResourceNotFound, "推荐记录不存在")
			return
		}
		response.Fail(c, "审核失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "审核成功"})
}

// GetAll 获取所有课程推荐（管理员）
func GetAll(c *gin.Context) {
	list, err := services.GetAll()
	if err != nil {
		response.Fail(c, "获取失败: "+err.Error())
		return
	}

	response.Success(c, list)
}

// Delete 删除课程推荐（管理员）
func Delete(c *gin.Context) {
	var req struct {
		RecommendationID int64 `json:"recommendation_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err := services.Delete(req.RecommendationID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			response.FailWithCode(c, response.CodeResourceNotFound, "推荐记录不存在")
			return
		}
		response.Fail(c, "删除失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
