package controller

import (
	"fmt"
	"net/http"
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type FriendshipController struct {
	FriendshipService service.FriendshipService
}

func NewFriendshipController(service service.FriendshipService) *FriendshipController {
	return &FriendshipController{
		FriendshipService: service,
	}
}

// Create				godoc
// @Summary				Create a Friendship.
// @Description			Creates a Friendship.
// @Param				Friendship body dto.CreateFriendshipRequest true "Create Friendship"
// @Produce				application/json
// @Tags				Friendship
// @Success				201 {object} dto.Response{}
// @Router				/friendship [post]
func (f *FriendshipController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := dto.CreateFriendshipRequest{}
	ctx.ShouldBindJSON(&request)

	if err := f.FriendshipService.CreateFriendship(user.ID, request.UserRequestee); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("could not create friendship: %v", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Data: "",
	})

}

// ReplyFriendship		godoc
// @Summary				Reply a Friendship request.
// @Description			Repliess a Friendship request.
// @Param				FriendshipRequest body dto.ReplyFriendshipRequest true "Accept Friendship"
// @Produce				application/json
// @Tags				Friendship
// @Success				201 {object} dto.Response{}
// @Router				/friendship/accept [post]
func (f *FriendshipController) ReplyFriendship(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := dto.ReplyFriendshipRequest{}
	ctx.ShouldBindJSON(&request)

	shouldThisUserBeReplying := f.FriendshipService.ShouldThisUserBeReplying(request.FriendshipID, user.ID)
	if !shouldThisUserBeReplying {
		ctx.JSON(http.StatusForbidden, "this friendship does not belong to this user")
		return
	}

	if err := f.FriendshipService.ReplyFriendship(request.FriendshipID, request.Accept); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("could not accept friendship: %v", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Data: "",
	})
}

// Delete 				godoc
// @Summary				Delete Friendship by id.
// @Param				friendshipId path string true "delete Friendship by id"
// @Description			Deletes a Friendship.
// @Produce				application/json
// @Tags				Friendship
// @Success				204 {object} dto.Response{}
// @Router				/friendship/{friendshipId} [delete]
func (f *FriendshipController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}
