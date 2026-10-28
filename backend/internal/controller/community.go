package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	communityService service.CommunityService
	topicService     service.TopicService
}

func NewCommunityController(
	communityService service.CommunityService,
	topicService service.TopicService,
) *CommunityController {
	return &CommunityController{
		communityService: communityService,
		topicService:     topicService,
	}
}

// Create				godoc
// @ID					createCommunity
// @Summary				Create Community
// @Description			Creates a community owned by the logged user, optionally with a picture (PNG or JPEG, max 5 MB and 50 megapixels; stored shrunk to 1024px and 1 MB).
// @Param				name formData string true "Community name"
// @Param				category formData string true "Community category, one of transport.CommunityCategory (e.g. computers_internet)"
// @Param				description formData string false "Community description (max 2048 characters)"
// @Param				private formData bool false "Private: only members read the forum, and the owner accepts who joins"
// @Param				picture formData file false "Community picture"
// @Accept				multipart/form-data
// @Produce				application/json
// @Tags				Community
// @Success				201 {object} transport.Response{data=transport.CreateCommunityResponse}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities [post]
func (c *CommunityController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.CreateCommunityRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	picture, ok := readPicture(ctx, true)
	if !ok {
		return
	}

	communityId, err := c.communityService.CreateCommunity(ctx.Request.Context(), user.ID, request.Name, string(request.Category), request.Description, request.Private, picture)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{
		Data: transport.CreateCommunityResponse{ID: communityId},
	})
}

// Update				godoc
// @ID					updateCommunity
// @Summary				Updates a community.
// @Param				communityId path int true "Community ID"
// @Param				Community body transport.UpdateCommunityRequest true "Community details"
// @Description			Replaces the community's name, category, description and privacy. Only the owner may (403 otherwise). Making a private community public accepts its pending join requests.
// @Accept				application/json
// @Produce				application/json
// @Tags				Community
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId} [put]
func (c *CommunityController) Update(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	request := transport.UpdateCommunityRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	err := c.communityService.UpdateCommunity(user.ID, communityId, request.Name, string(request.Category), request.Description, request.Private)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// UpdatePicture		godoc
// @ID					updateCommunityPicture
// @Summary				Replaces a community's picture.
// @Param				communityId path int true "Community ID"
// @Param				picture formData file true "Community picture"
// @Description			Uploads a PNG or JPEG (max 5 MB and 50 megapixels; stored shrunk to 1024px and 1 MB) as the community's picture. Only the owner may (403 otherwise).
// @Accept				multipart/form-data
// @Produce				application/json
// @Tags				Community
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/picture [put]
func (c *CommunityController) UpdatePicture(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	picture, ok := readPicture(ctx, false)
	if !ok {
		return
	}

	if err := c.communityService.UpdateCommunityPicture(ctx.Request.Context(), user.ID, communityId, picture); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// FindById				godoc
// @ID					findCommunity
// @Summary				Get Single community by id.
// @Param				communityId path int true "Community ID"
// @Description			Returns the community matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindCommunityByIdResponse{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId} [get]
func (c *CommunityController) FindById(ctx *gin.Context) {
	id, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	community, err := c.communityService.FindCommunityById(id)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, community)
}

// FindProfileById		godoc
// @ID					findCommunityProfile
// @Summary				Get Single community profile by its id.
// @Param				communityId path int true "Community ID"
// @Description			Returns the community profile matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindCommunityProfileByIdResponse{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/profile [get]
func (c *CommunityController) FindProfileById(ctx *gin.Context) {
	id, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	profile, err := c.communityService.FindCommunityProfileById(id)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// FindMembers			godoc
// @ID					findCommunityMembers
// @Summary				Get a page of community members.
// @Param				communityId path int true "Community ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns a page of the community's members, most recently seen first
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindCommunityMembersResponse{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/members [get]
func (c *CommunityController) FindMembers(ctx *gin.Context) {
	id, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	page, pageSize := pageQuery(ctx)
	members, err := c.communityService.FindCommunityMembers(id, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, members)
}

// FindRelated			godoc
// @ID					findRelatedCommunities
// @Summary				Get related communities.
// @Param				communityId path int true "Community ID"
// @Description			Returns 9 random communities other than this one, the same ones all day (communities don't record related ones yet)
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindManyCommunitiesResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/related [get]
func (c *CommunityController) FindRelated(ctx *gin.Context) {
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	communities, err := c.communityService.FindRelatedCommunities(communityId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, communities)
}

// FindHotTopics		godoc
// @ID					findCommunityHotTopics
// @Summary				Gets the 5 newest topics of a community.
// @Param				communityId path int true "Community ID"
// @Description			Returns a list of topics. In a private community, only members may read (403 otherwise).
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.TopicsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/hot-topics [get]
func (c *CommunityController) FindHotTopics(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}
	topics, err := c.topicService.FindTopics(user.ID, communityId, 1, 5)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, topics)
}

// FindAllTopics		godoc
// @ID					findCommunityTopics
// @Summary				Gets a page of a community's topics.
// @Param				communityId path int true "Community ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns a page of topics, newest first. In a private community, only members may read (403 otherwise).
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.TopicsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/{communityId}/topics [get]
func (c *CommunityController) FindAllTopics(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)
	topics, err := c.topicService.FindTopics(user.ID, communityId, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, topics)
}

// Search				godoc
// @ID					searchCommunities
// @Summary				Searches communities the logged user hasn't joined.
// @Param				name query string false "Part of the community name (omit to list all)"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Communities per page (default 20, max 50)"
// @Description			Returns a page of communities the logged user is not a member of, biggest first.
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.SearchCommunitiesResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/communities/search [get]
func (c *CommunityController) Search(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	name := ctx.Query("name")
	page, pageSize := pageQuery(ctx)

	communities, err := c.communityService.SearchCommunitiesToDiscover(user.ID, name, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, communities)
}
