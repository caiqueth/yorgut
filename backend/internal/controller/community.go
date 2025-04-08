package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"

	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	communityService        service.CommunityService
	communityProfileService service.CommunityProfileService
}

func NewCommunityController(
	communityService service.CommunityService,
	communityProfileService service.CommunityProfileService,
) *CommunityController {
	return &CommunityController{
		communityService:        communityService,
		communityProfileService: communityProfileService,
	}
}

// Create				godoc
// @Summary				Create Community
// @Decription			Save community in Db.
// @Param				name formData string true "Community Name"
// @Param				description formData string true "Community Description"
// @Param				createdById formData int true "ID of Creator"
// @Param				picture formData file true "Community Picture"
// @Produce				application/json
// @Tags				Community
// @Success				201 {object} dto.Response{}
// @Router				/community [post]
func (c *CommunityController) Create(ctx *gin.Context) {
	createCommunityRequest := dto.CreateCommunityRequest{
		Name:        ctx.PostForm("name"),
		CreatedByID: uint(helpers.MustParseUint(ctx.PostForm("createdById"))),
	}

	err := c.communityService.CreateCommunity(createCommunityRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)
}

// Update				godoc
// @Summary				Updates a community.
// @Description			Updates properties of a single community
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} dto.Response{}
// @Router				/community [put]
func (c *CommunityController) Update(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete community by id.
// @Param				communityId path string true "delete community by id"
// @Description			Deletes a community
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} dto.Response{}
// @Router				/community/{communityId} [delete]
func (c *CommunityController) Delete(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		panic(err)
	}

	err = c.communityService.DeleteCommunity(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.Response{})

}

// FindById				godoc
// @Summary				Get Single community by id.
// @Param				communityId path string true "get community by id"
// @Description			Returns the community matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.FindCommunityByIdResponse{}
// @Router				/community/{communityId} [get]
func (c *CommunityController) FindById(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		panic(err)
	}

	community, err := c.communityService.FindCommunityById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: community,
	}
	ctx.JSON(http.StatusOK, webResponse)

}

// FindMany				godoc
// @Summary				Get Multiple communities.
// @Description			Returns an array of communities
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.FindManyCommunitiesResponse{}
// @Router				/community [get]
func (c *CommunityController) FindMany(ctx *gin.Context) {
	communities, err := c.communityService.FindManyCommunities()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}

// CreateCommunityProfile				godoc
// @Summary				Create Community Profile
// @Decription			Save community profile in Db.
// @Param				communityId formData string true "Community Name"
// @Param				description formData string true "Community Description"
// @Param				picture formData file true "Community Picture"
// @Produce				application/json
// @Tags				Community
// @Success				201 {object} dto.Response{}
// @Router				/community_profile [post]
func (c *CommunityController) CreateCommunityProfile(ctx *gin.Context) {
	request := dto.CreateCommunityProfileRequest{
		CommunityID: helpers.MustParseUint(ctx.PostForm("communityId")),
		Description: ctx.PostForm("description"),
	}

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get picture file"})
		return
	}

	request.Picture = file

	err = c.communityProfileService.CreateCommunityProfile(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)
}
