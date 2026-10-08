package controller

import (
	"errors"
	"fmt"
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
// @Summary				Creates a Friendship request.
// @Description			Logged user makes a new friendship request.
// @Param				Friendship body transport.CreateFriendshipRequest true "Create Friendship"
// @Produce				application/json
// @Tags				Friendship
// @Success				204
// @Router				/friendships [post]
func (f *FriendshipController) AddFriend(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.CreateFriendshipRequest{}
	ctx.ShouldBindJSON(&request)

	if err := f.friendshipService.CreateFriendship(user.ID, request.UserRequesteeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create friendship: %v", err.Error())})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// ReplyFriendship		godoc
// @Summary				Reply a Friendship request.
// @Description			Repliess a Friendship request.
// @Param				FriendshipRequest body transport.ReplyFriendshipRequest true "Accept Friendship"
// @Produce				application/json
// @Tags				Friendship
// @Success				204
// @Router				/friendships/reply [post]
func (f *FriendshipController) ReplyFriendship(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.ReplyFriendshipRequest{}
	ctx.ShouldBindJSON(&request)

	shouldThisUserBeReplying := f.friendshipService.ShouldThisUserBeReplying(request.FriendshipID, user.ID)
	if !shouldThisUserBeReplying {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "this friendship does not belong to this user"})
		return
	}

	if err := f.friendshipService.ReplyFriendship(request.FriendshipID, request.Accept); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not accept friendship: %v", err.Error())})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// RemoveFriendByFriendshipID	godoc
// @Summary						Deletes a Friendship by id.
// @Param						friendshipId path string true "delete Friendship by id"
// @Description					Deletes a Friendship.
// @Produce						application/json
// @Tags						Friendship
// @Success						204
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)

}

// RemoveFriendByUserID		godoc
// @Summary					Deletes a Friendship by the target userId Id.
// @Param					friendId path string true "Target User ID"
// @Description				Logged user deletes a friendship.
// @Produce					application/json
// @Tags					Friendship
// @Success					204
// @Router					/friendships/friend/{userId} [delete]
func (f *FriendshipController) RemoveFriendByUserID(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	if err := f.friendshipService.RemoveFriend(user.ID, userId); err != nil {
		if errors.Is(err, service.ErrFriendshipNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("friendship not found between users %v and %v", user.ID, userId)})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)

}
