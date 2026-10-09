package controller

import (
	"net/http"
	"yorgut/internal/helpers"
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
// @Failure			400
// @Failure			401
// @Failure			404
// @Failure			409
// @Security			BearerAuth
// @Router				/friendships [post]
func (f *FriendshipController) AddFriend(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.CreateFriendshipRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Failure			400
// @Failure			401
// @Failure			403
// @Failure			404
// @Security			BearerAuth
// @Router				/friendships/reply [post]
func (f *FriendshipController) ReplyFriendship(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.ReplyFriendshipRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shouldThisUserBeReplying := f.friendshipService.ShouldThisUserBeReplying(request.FriendshipID, user.ID)
	if !shouldThisUserBeReplying {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "this friendship does not belong to this user"})
		return
	}

	if err := f.friendshipService.ReplyFriendship(request.FriendshipID, request.Accept); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// RemoveFriendByFriendshipID	godoc
// @ID					deleteFriendship
// @Summary						Deletes a Friendship by id.
// @Param						friendshipId path int true "Friendship ID"
// @Description					Deletes a Friendship.
// @Produce						application/json
// @Tags						Friendship
// @Success						204
// @Failure			401
// @Failure			403
// @Failure			404
// @Security			BearerAuth
// @Router						/friendships/{friendshipId} [delete]
func (f *FriendshipController) RemoveFriendByFriendshipID(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendshipId := helpers.MustParseUint(ctx.Param("friendshipId"))

	shouldThisUserBeReplying := f.friendshipService.ShouldThisUserBeDeleting(friendshipId, user.ID)
	if !shouldThisUserBeReplying {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "this friendship does not belong to this user"})
		return
	}

	if err := f.friendshipService.DeleteFriendship(friendshipId); err != nil {
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
// @Failure			401
// @Failure			404
// @Security			BearerAuth
// @Router					/friendships/friend/{userId} [delete]
func (f *FriendshipController) RemoveFriendByUserID(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	if err := f.friendshipService.RemoveFriend(user.ID, userId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)

}
