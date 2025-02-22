package controller

import (
	"net/http"
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
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
// @Param				community body request.CreateCommunityMembershipRequest true "Create Community Membership"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				204 {object} response.Response{}
// @Router				/community_membership [post]
func (cm *CommunityMembershipController) Create(ctx *gin.Context) {
	createMembershipRequest := request.CreateCommunityMembershipRequest{}

	err := ctx.ShouldBindJSON(&createMembershipRequest)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = cm.CommunityMembershipService.Create(createMembershipRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := response.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusNoContent, webResponse)

}

// Create				godoc
// @Summary				Delete Community Membership
// @Decription			Deletes the uint so that an user is not a member of the community anymore
// @Param				userId path uint true "id of the user"
// @Param				communityId path string true "id of the community"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				204 {object} response.Response{}
// @Router				/community_membership/{userId}/{communityId} [delete]
func (cm *CommunityMembershipController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}

// Create				godoc
// @Summary				Lists communities that a user is member of
// @Decription			Lists all communities that a given user is member
// @Param				userId path uint true "id of the user"
// @Produce				application/json
// @Tags				CommunityMembership
// @Success				204 {object} response.Response{}
// @Router				/community_membership/{userId} [get]
func (cm *CommunityMembershipController) FindOfUser(ctx *gin.Context) {
	panic("unimplemented")
}
