package controller

import (
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
// @Param				Friendship body request.CreateFriendshipRequest true "Create Friendship"
// @Produce				application/json
// @Tags				Friendship
// @Success				201 {object} response.Response{}
// @Router				/friendship [post]
func (f *FriendshipController) Create(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete Friendship by id.
// @Param				friendshipId path string true "delete Friendship by id"
// @Description			Deletes a Friendship.
// @Produce				application/json
// @Tags				Friendship
// @Success				204 {object} response.Response{}
// @Router				/friendship/{friendshipId} [delete]
func (f *FriendshipController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}
