package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type FriendshipController struct {
	friendshipService service.FriendshipService
}

func NewFriendshipController(service service.FriendshipService) *FriendshipController {
	return &FriendshipController{
		friendshipService: service,
	}
}

// AddFriend			godoc
// @ID					addFriend
// @Summary				Creates a Friendship request.
// @Description			Logged user makes a new friendship request.
// @Param				Friendship body transport.CreateFriendshipRequest true "Create Friendship"
// @Accept				application/json
// @Produce				application/json
// @Tags				Friendship
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/friendships [post]
func (f *FriendshipController) AddFriend(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.CreateFriendshipRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := f.friendshipService.CreateFriendship(user.ID, request.UserRequesteeID); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// ReplyFriendship		godoc
// @ID					replyFriendship
// @Summary				Reply a Friendship request.
// @Description			Repliess a Friendship request.
// @Param				FriendshipRequest body transport.ReplyFriendshipRequest true "Accept Friendship"
// @Accept				application/json
// @Produce				application/json
// @Tags				Friendship
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/friendships/reply [post]
func (f *FriendshipController) ReplyFriendship(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.ReplyFriendshipRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := f.friendshipService.ReplyFriendship(user.ID, request.FriendshipID, request.Accept); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// RemoveFriendByUserID		godoc
// @ID					removeFriend
// @Summary					Deletes a Friendship by the target userId Id.
// @Param					userId path int true "Target user ID"
// @Description				Logged user deletes a friendship.
// @Produce					application/json
// @Tags					Friendship
// @Success					204
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/friendships/friend/{userId} [delete]
func (f *FriendshipController) RemoveFriendByUserID(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	if err := f.friendshipService.RemoveFriend(user.ID, userId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
