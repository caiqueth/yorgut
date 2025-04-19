package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type CommunityMembershipController struct {
	CommunityMembershipService service.CommunityMembershipService
}

func NewCommunityMembershipController(service service.CommunityMembershipService) *CommunityMembershipController {
	return &CommunityMembershipController{
		CommunityMembershipService: service,
	}
}

// Create				godoc
// @Summary				Create Community Membership
// @Decription			Save membership in Db.
// @Param				community body dto.CreateCommunityMembershipRequest true "Create Community Membership"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				201 {object} dto.Response{}
// @Router				/community_membership [post]
func (cm *CommunityMembershipController) Create(ctx *gin.Context) {
	createMembershipRequest := dto.CreateCommunityMembershipRequest{}

	err := ctx.ShouldBindJSON(&createMembershipRequest)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = cm.CommunityMembershipService.CreateCommunityMembership(createMembershipRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// Create				godoc
// @Summary				Delete Community Membership
// @Decription			Deletes the uint so that an user is not a member of the community anymore
// @Param				userId path uint true "id of the user"
// @Param				communityId path string true "id of the community"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				204 {object} dto.Response{}
// @Router				/community_membership/{userId}/{communityId} [delete]
func (cm *CommunityMembershipController) Delete(ctx *gin.Context) {
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

	communityMembershipId, err := cm.CommunityMembershipService.FindCommunityMembership(uint(communityId), uint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cm.CommunityMembershipService.DeleteCommunityMembership(*communityMembershipId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}
