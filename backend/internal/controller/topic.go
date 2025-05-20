package controller

import (
	"net/http"
	"yorgut/internal/helpers"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type TopicController struct {
	topicService service.TopicService
}

func NewTopicController(service service.TopicService) *TopicController {
	return &TopicController{
		topicService: service,
	}
}

// FindHotTopics		godoc
// @Summary				Gets a list of most recently replied topics.
// @Description			Returns a list of topics
// @Produce				application/json
// @Tags				Topics
// @Success				200 {object} dto.TopicsResponse{}
// @Router				/topics/hot/{communityId} [get]
func (tc *TopicController) FindHotTopics(ctx *gin.Context) {
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	topics, err := tc.topicService.FindHotTopics(communityId, 5)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, topics)

}
