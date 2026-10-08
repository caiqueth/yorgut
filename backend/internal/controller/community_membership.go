package controller

import (
	"errors"
	"net/http"
	"strconv"
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

// Create				godoc
// @Summary				Create Community Membership
// @Decription			Save membership in Db.
// @Param				community body transport.CreateCommunityMembershipRequest true "Create Community Membership"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				201 {object} transport.Response{}
// @Router				/community-memberships [post]
func (cm *CommunityMembershipController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	createMembershipRequest := transport.CreateCommunityMembershipRequest{}

	err := ctx.ShouldBindJSON(&createMembershipRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cm.communityMembershipService.JoinCommunity(user.ID, createMembershipRequest.CommunityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := transport.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// Delete				godoc
// @Summary				Delete Community Membership
// @Description			The logged user leaves a community (same as /me/leave). userId must be the logged user (403 otherwise), and the owner can't leave (403).
// @Param				userId path uint true "id of the user"
// @Param				communityId path string true "id of the community"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				204
// @Router				/community-memberships/{userId}/{communityId} [delete]
func (cm *CommunityMembershipController) Delete(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	communityId, err := strconv.Atoi(ctx.Param("communityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uint(userId) != user.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "users can only remove their own memberships"})
		return
	}

	err = cm.communityMembershipService.LeaveCommunity(user.ID, uint(communityId))
	if errors.Is(err, service.ErrOwnerCannotLeave) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
