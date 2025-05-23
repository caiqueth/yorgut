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

// FindTopicAndReplies	godoc
// @Summary				Lists all the replies of a given topic
// @Description			Returns a list of topics
// @Produce				application/json
// @Tags				Topics
// @Success				200 {object} dto.TopicRepliesResponse{}
// @Router				/topic/{topicId}/replies [get]
func (tc *TopicController) FindTopicAndReplies(ctx *gin.Context) {
	topicId := helpers.MustParseUint(ctx.Param("topicId"))
	topics, err := tc.topicService.FindTopicReplies(topicId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, topics)

}
