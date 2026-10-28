package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type CommunityMembershipController struct {
	communityMembershipService service.CommunityMembershipService
}

func NewCommunityMembershipController(service service.CommunityMembershipService) *CommunityMembershipController {
	return &CommunityMembershipController{
		communityMembershipService: service,
	}
}

// FindJoinRequests		godoc
// @ID					findCommunityJoinRequests
// @Summary				Lists who asked to join a private community.
// @Param				communityId path int true "Community ID"
// @Description			Returns the pending join requests, oldest first. Only the owner may (403 otherwise).
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindJoinRequestsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/join-requests [get]
func (cm *CommunityMembershipController) FindJoinRequests(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	requests, err := cm.communityMembershipService.FindJoinRequests(user.ID, communityId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

// ReplyJoinRequest		godoc
// @ID					replyCommunityJoinRequest
// @Summary				Accepts or declines a request to join a private community.
// @Param				communityId path int true "Community ID"
// @Param				userId path int true "ID of the user who asked to join"
// @Param				Reply body transport.ReplyJoinRequestRequest true "Accept or decline"
// @Description			Accepting makes the user a member; declining drops the request. Only the owner may (403 otherwise).
// @Accept				application/json
// @Produce				application/json
// @Tags				Community
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/join-requests/{userId}/reply [post]
func (cm *CommunityMembershipController) ReplyJoinRequest(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	request := transport.ReplyJoinRequestRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	err := cm.communityMembershipService.ReplyJoinRequest(user.ID, communityId, userId, request.Accept)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// FindInvitableFriends	godoc
// @ID					findInvitableFriends
// @Summary				Gets a page of the logged user's friends to invite to a community.
// @Param				communityId path int true "Community ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns a page of friends, most recently seen first, each saying whether they already are a member or already invited. Only members may (403 otherwise).
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindInvitableFriendsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/invitable-friends [get]
func (cm *CommunityMembershipController) FindInvitableFriends(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)

	friends, err := cm.communityMembershipService.FindInvitableFriends(user.ID, communityId, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, friends)
}

// InviteToCommunity	godoc
// @ID					inviteToCommunity
// @Summary				Invites a friend to join a community.
// @Param				communityId path int true "Community ID"
// @Param				Invitation body transport.InviteToCommunityRequest true "The friend to invite"
// @Description			Only members may invite (403 otherwise), and only their friends (403). 409 when the friend is already a member or already invited.
// @Accept				application/json
// @Produce				application/json
// @Tags				Community
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/invitations [post]
func (cm *CommunityMembershipController) InviteToCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	request := transport.InviteToCommunityRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	err := cm.communityMembershipService.InviteToCommunity(user.ID, communityId, request.UserID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
