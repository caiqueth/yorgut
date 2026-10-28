package controller

import (
	"net/http"
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
// @Summary				Gets a topic and a page of its replies
// @Param				topicId path int true "Topic ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns the topic and a page of its replies, oldest first. In a private community, only members may read (403 otherwise).
// @Produce				application/json
// @Tags				Topics
// @Success				200 {object} transport.TopicRepliesResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/topics/{topicId}/replies [get]
func (t *TopicController) FindTopicAndReplies(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	topicId, ok := idParam(ctx, "topicId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)
	topics, err := t.topicService.FindTopicReplies(user.ID, topicId, page, pageSize)
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
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/topics [post]
func (t *TopicController) CreateTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	req := transport.CreateTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := t.topicService.CreateTopic(user.ID, req.CommunityID, req.Title, req.Text); err != nil {
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
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/topics/{topicId} [post]
func (t *TopicController) ReplyTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	parentTopicId, ok := idParam(ctx, "topicId")
	if !ok {
		return
	}

	req := transport.ReplyTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := t.topicService.ReplyTopic(user.ID, parentTopicId, req.Title, req.Text); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}
