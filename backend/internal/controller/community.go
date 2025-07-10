package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"

	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	communityService        service.CommunityService
	communityProfileService service.CommunityProfileService
	topicService            service.TopicService
}

func NewCommunityController(
	communityService service.CommunityService,
	communityProfileService service.CommunityProfileService,
	topicService service.TopicService,
) *CommunityController {
	return &CommunityController{
		communityService:        communityService,
		communityProfileService: communityProfileService,
		topicService:            topicService,
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
	user := ctx.MustGet("user").(*model.User)

	request := dto.CreateCommunityRequest{}
	ctx.ShouldBindJSON(&request)

	request.CreatedByID = user.ID

	err := c.communityService.CreateCommunity(request)
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
// @Param				communityId path string true "Community ID"
// @Description			Deletes a community
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} dto.Response{}
// @Router				/community/{communityId} [delete]
func (c *CommunityController) Delete(ctx *gin.Context) {
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	if err := c.communityService.DeleteCommunity(uint(communityId)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.Response{})

}

// FindById				godoc
// @Summary				Get Single community by id.
// @Param				communityId path string true "Community ID"
// @Description			Returns the community matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.FindCommunityByIdResponse{}
// @Router				/community/{communityId} [get]
func (c *CommunityController) FindById(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	community, err := c.communityService.FindCommunityById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, community)

}

// FindProfileById		godoc
// @Summary				Get Single community profile by its id.
// @Param				communityId path string true "Community ID"
// @Description			Returns the community profile matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.FindCommunityProfileByIdResponse{}
// @Router				/community/{communityId}/profile [get]
func (c *CommunityController) FindProfileById(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := c.communityService.FindCommunityProfileById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)

}

// FindMembers			godoc
// @Summary				Get community members.
// @Param				communityId path string true "Community ID"
// @Description			Returns a list of the community members matching the requested id and a few stats
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.FindCommunityMembersResponse{}
// @Router				/community/{communityId}/members [get]
func (c *CommunityController) FindMembers(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	members, err := c.communityService.FindCommunityMembers(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, members)

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

// FindRelated			godoc
// @Summary				Get Related communities.
// @Description			Returns an array of communities
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.FindManyCommunitiesResponse{}
// @Router				/community/{communityId}/related [get]
func (c *CommunityController) FindRelated(ctx *gin.Context) {
	// For now, I don't have a behavior for this, so we just list all communities
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	fmt.Println(communityId)

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
// @Router				/community/profile [post]
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

// FindHotTopics		godoc
// @Summary				Gets a list of most recently replied topics.
// @Description			Returns a list of topics
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.TopicsResponse{}
// @Router				/community/{communityId}/hot_topics [get]
func (c *CommunityController) FindHotTopics(ctx *gin.Context) {
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	topics, err := c.topicService.FindTopics(communityId, 5)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, topics)

}

// FindAllTopics		godoc
// @Summary				Gets a list of most recently replied topics.
// @Description			Returns a list of topics
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} dto.TopicsResponse{}
// @Router				/community/{communityId}/topics [get]
func (c *CommunityController) FindAllTopics(ctx *gin.Context) {
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	topics, err := c.topicService.FindTopics(communityId, 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, topics)

}
