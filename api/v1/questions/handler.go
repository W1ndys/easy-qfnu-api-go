package questions

import (
	"github.com/W1ndys/qfnu-api-go/common/response"
	"github.com/W1ndys/qfnu-api-go/model"
	services "github.com/W1ndys/qfnu-api-go/services/questions"
	"github.com/gin-gonic/gin"
)

// Search 题目搜索接口
func Search(c *gin.Context) {
	keyword := c.Query("keyword")

	// 如果关键词为空，返回空列表
	if keyword == "" {
		response.Success(c, []model.FreshmanQuestion{})
		return
	}

	questions, err := services.SearchQuestions(keyword)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}

	// 确保返回空数组而不是 null
	if questions == nil {
		questions = []model.FreshmanQuestion{}
	}

	response.Success(c, questions)
}
