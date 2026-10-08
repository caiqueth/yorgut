package controller

import (
	"errors"
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

// topicErrorStatus maps the service's sentinel errors to HTTP statuses.
func topicErrorStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrTopicNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrNotCommunityMember):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// FindTopicAndReplies	godoc
// @Summary				Lists all the replies of a given topic
// @Description			Returns a list of topics
// @Produce				application/json
// @Tags				Topics
// @Success				200 {object} transport.TopicRepliesResponse{}
// @Router				/topics/{topicId}/replies [get]
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
// @Description		Creates a new Topic inside a community's forum. Only members of the community may post (403 otherwise).
// @Param			Topic body transport.CreateTopicRequest true "Create Topic"
// @Tags			Topics
// @Success			201 {object} transport.Response{}
// @Router			/topics [post]
func (tc *TopicController) CreateTopic(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	req := transport.CreateTopicRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.topicService.CreateTopic(user.ID, req.CommunityID, req.Title, req.Text); err != nil {
		ctx.JSON(topicErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}

// ReplyTopic		godoc
// @Summary			Creates a new Topic Reply
// @Description		Replies to an existing Topic inside a community's forum. Only members of the topic's community may reply (403 otherwise).
// @Param			topicId path string true "Topic ID"
// @Param			Reply body transport.ReplyTopicRequest true "Reply"
// @Tags			Topics
// @Success			201 {object} transport.Response{}
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
		ctx.JSON(topicErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}
