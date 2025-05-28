package controller

import (
	"net/http"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
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

// CreateTopic		godoc
// @Summary			Creates a new Topic
// @Description		Creates a new Topic inside a community's forum
// @Tags			Topics
// @Success			201 {object} dto.Response{}
// @Router			/topic [post]
func (tc *TopicController) CreateTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	req := dto.CreateTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.topicService.CreateTopic(user.ID, req.CommunityID, req.Title, req.Text); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{})
}

// ReplyTopic		godoc
// @Summary			Creates a new Topic Reply
// @Description		Replies to an existing Topic inside a community's forum
// @Tags			Topics
// @Success			201 {object} dto.Response{}
// @Router			/topic/{topicId} [post]
func (tc *TopicController) ReplyTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	parentTopicId := helpers.MustParseUint(ctx.Param("topicId"))

	req := dto.ReplyTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.topicService.ReplyTopic(user.ID, req.CommunityID, parentTopicId, req.Title, req.Text); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{})
}
