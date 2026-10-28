package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type PollController struct {
	pollService service.PollService
}

func NewPollController(pollService service.PollService) *PollController {
	return &PollController{
		pollService: pollService,
	}
}

// FindPolls			godoc
// @ID					findCommunityPolls
// @Summary				Gets a page of a community's polls.
// @Param				communityId path int true "Community ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns a page of polls, newest first. In a private community, only members may see them (403 otherwise).
// @Produce				application/json
// @Tags				Poll
// @Success				200 {object} transport.PollsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/polls [get]
func (p *PollController) FindPolls(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)

	polls, err := p.pollService.FindPolls(user.ID, communityId, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, polls)
}

// CreatePoll			godoc
// @ID					createPoll
// @Summary				Creates a poll in a community.
// @Param				communityId path int true "Community ID"
// @Param				Poll body transport.CreatePollRequest true "The poll"
// @Description			Only members may (403 otherwise). Takes 2 to 10 different options; closesAt, when given, must be in the future.
// @Accept				application/json
// @Produce				application/json
// @Tags				Poll
// @Success				201 {object} transport.Response{data=transport.CreatePollResponse}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/polls [post]
func (p *PollController) CreatePoll(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	request := transport.CreatePollRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	pollId, err := p.pollService.CreatePoll(user.ID, communityId, service.CreatePollInput{
		Question:    request.Question,
		Description: request.Description,
		Options:     request.Options,
		ClosesAt:    request.ClosesAt,
	})
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{
		Data: transport.CreatePollResponse{ID: pollId},
	})
}

// FindPoll				godoc
// @ID					findPoll
// @Summary				Gets a poll with its options.
// @Param				pollId path int true "Poll ID"
// @Description			Each option's votes are null until the logged user voted or the poll closed. In a private community, only members may see it (403 otherwise).
// @Produce				application/json
// @Tags				Poll
// @Success				200 {object} transport.PollResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/polls/{pollId} [get]
func (p *PollController) FindPoll(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	pollId, ok := idParam(ctx, "pollId")
	if !ok {
		return
	}

	poll, err := p.pollService.FindPoll(user.ID, pollId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, poll)
}

// Vote					godoc
// @ID					votePoll
// @Summary				Votes in a poll.
// @Param				pollId path int true "Poll ID"
// @Param				Vote body transport.VotePollRequest true "The chosen option"
// @Description			Only members vote (403 otherwise), once (409 after that), while the poll is open (409 once closed).
// @Accept				application/json
// @Produce				application/json
// @Tags				Poll
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/polls/{pollId}/vote [post]
func (p *PollController) Vote(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	pollId, ok := idParam(ctx, "pollId")
	if !ok {
		return
	}

	request := transport.VotePollRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := p.pollService.Vote(user.ID, pollId, request.OptionID); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// Close				godoc
// @ID					closePoll
// @Summary				Closes a poll.
// @Param				pollId path int true "Poll ID"
// @Description			Ends voting now. Only the poll's author or the community's owner may (403 otherwise); 409 when already closed.
// @Produce				application/json
// @Tags				Poll
// @Success				204
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/polls/{pollId}/close [post]
func (p *PollController) Close(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	pollId, ok := idParam(ctx, "pollId")
	if !ok {
		return
	}

	if err := p.pollService.ClosePoll(user.ID, pollId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
