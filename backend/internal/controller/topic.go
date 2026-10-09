package controller

import (
	"net/http"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

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
// @ID					findTopicReplies
// @Summary				Lists all the replies of a given topic
// @Param				topicId path int true "Topic ID"
// @Description			Returns the topic and its replies
// @Produce				application/json
// @Tags				Topics
// @Success				200 {object} transport.TopicRepliesResponse{}
// @Failure			401
// @Failure			404
// @Security			BearerAuth
// @Router				/topics/{topicId}/replies [get]
func (tc *TopicController) FindTopicAndReplies(ctx *gin.Context) {
	topicId := helpers.MustParseUint(ctx.Param("topicId"))
	topics, err := tc.topicService.FindTopicReplies(topicId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, topics)

}

// CreateTopic		godoc
// @ID					createTopic
// @Summary			Creates a new Topic
// @Description		Creates a new Topic inside a community's forum. Only members of the community may post (403 otherwise).
// @Param			Topic body transport.CreateTopicRequest true "Create Topic"
// @Accept				application/json
// @Produce				application/json
// @Tags			Topics
// @Success			201 {object} transport.Response{}
// @Failure			400
// @Failure			401
// @Failure			403
// @Security			BearerAuth
// @Router			/topics [post]
func (tc *TopicController) CreateTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	req := transport.CreateTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.topicService.CreateTopic(user.ID, req.CommunityID, req.Title, req.Text); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}

// ReplyTopic		godoc
// @ID					replyTopic
// @Summary			Creates a new Topic Reply
// @Description		Replies to an existing Topic inside a community's forum. Only members of the topic's community may reply (403 otherwise).
// @Param			topicId path int true "Topic ID"
// @Param			Reply body transport.ReplyTopicRequest true "Reply"
// @Accept				application/json
// @Produce				application/json
// @Tags			Topics
// @Success			201 {object} transport.Response{}
// @Failure			400
// @Failure			401
// @Failure			403
// @Failure			404
// @Security			BearerAuth
// @Router			/topics/{topicId} [post]
func (tc *TopicController) ReplyTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	parentTopicId := helpers.MustParseUint(ctx.Param("topicId"))

	req := transport.ReplyTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.topicService.ReplyTopic(user.ID, parentTopicId, req.Title, req.Text); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}
