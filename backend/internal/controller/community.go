package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/transport"

	"yorgut/internal/service"

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
// @Summary				Create Community
// @Description			Creates a community owned by the logged user, optionally with a picture (PNG or JPEG, max 1 MB and 1024x1024).
// @Param				name formData string true "Community name"
// @Param				category formData string true "Community category"
// @Param				picture formData file false "Community picture"
// @Accept				multipart/form-data
// @Produce				application/json
// @Tags				Community
// @Success				201 {object} transport.Response{data=transport.CreateCommunityResponse}
// @Router				/communities [post]
func (c *CommunityController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.CreateCommunityRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var picture []byte
	file, err := ctx.FormFile("picture")
	switch {
	case errors.Is(err, http.ErrMissingFile):
		// the picture is optional
	case err != nil:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get picture file"})
		return
	case file.Size > service.MaxPictureBytes:
		// Refuse oversized files before reading them; the service checks again.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrPictureTooLarge.Error()})
		return
	default:
		picture, err = helpers.FileToBytes(file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read picture file"})
			return
		}
	}

	communityId, err := c.communityService.CreateCommunity(user.ID, request.Name, request.Category, picture)
	if errors.Is(err, service.ErrInvalidPicture) ||
		errors.Is(err, service.ErrPictureTooLarge) ||
		errors.Is(err, service.ErrPictureTooBig) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{
		Data: transport.CreateCommunityResponse{ID: communityId},
	})
}

// communityErrorStatus maps the service's sentinel errors to HTTP statuses.
func communityErrorStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrCommunityNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrCommunityForbidden):
		return http.StatusForbidden
	case errors.Is(err, service.ErrInvalidPicture),
		errors.Is(err, service.ErrPictureTooLarge),
		errors.Is(err, service.ErrPictureTooBig):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Update				godoc
// @Summary				Updates a community.
// @Param				communityId path string true "Community ID"
// @Param				Community body transport.UpdateCommunityRequest true "Community details"
// @Description			Replaces the community's name, category and description. Only the owner may (403 otherwise).
// @Produce				application/json
// @Tags				Community
// @Success				204
// @Router				/communities/{communityId} [put]
func (c *CommunityController) Update(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId := helpers.MustParseUint(ctx.Param("communityId"))

	request := transport.UpdateCommunityRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.communityService.UpdateCommunity(user.ID, communityId, request.Name, request.Category, request.Description)
	if err != nil {
		ctx.JSON(communityErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// UpdatePicture		godoc
// @Summary				Replaces a community's picture.
// @Param				communityId path string true "Community ID"
// @Param				picture formData file true "Community picture"
// @Description			Uploads a PNG or JPEG (max 1 MB and 1024x1024) as the community's picture. Only the owner may (403 otherwise).
// @Accept				multipart/form-data
// @Produce				application/json
// @Tags				Community
// @Success				204
// @Router				/communities/{communityId}/picture [put]
func (c *CommunityController) UpdatePicture(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId := helpers.MustParseUint(ctx.Param("communityId"))

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get picture file"})
		return
	}
	// Refuse oversized files before reading them; the service checks again.
	if file.Size > service.MaxPictureBytes {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrPictureTooLarge.Error()})
		return
	}

	picture, err := helpers.FileToBytes(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read picture file"})
		return
	}

	if err := c.communityService.UpdateCommunityPicture(user.ID, communityId, picture); err != nil {
		ctx.JSON(communityErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// Delete 				godoc
// @Summary				Delete community by id.
// @Param				communityId path string true "Community ID"
// @Description			Deletes a community
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} transport.Response{}
// @Router				/communities/{communityId} [delete]
func (c *CommunityController) Delete(ctx *gin.Context) {
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	if err := c.communityService.DeleteCommunity(uint(communityId)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, transport.Response{})

}

// FindById				godoc
// @Summary				Get Single community by id.
// @Param				communityId path string true "Community ID"
// @Description			Returns the community matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.FindCommunityByIdResponse{}
// @Router				/communities/{communityId} [get]
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
// @Success				200 {object} transport.FindCommunityProfileByIdResponse{}
// @Router				/communities/{communityId}/profile [get]
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
// @Success				200 {object} transport.FindCommunityMembersResponse{}
// @Router				/communities/{communityId}/members [get]
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
// @Success				200 {object} transport.FindManyCommunitiesResponse{}
// @Router				/communities [get]
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
// @Success				200 {object} transport.FindManyCommunitiesResponse{}
// @Router				/communities/{communityId}/related [get]
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

// FindHotTopics		godoc
// @Summary				Gets a list of most recently replied topics.
// @Description			Returns a list of topics
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.TopicsResponse{}
// @Router				/communities/{communityId}/hot-topics [get]
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
// @Success				200 {object} transport.TopicsResponse{}
// @Router				/communities/{communityId}/topics [get]
func (c *CommunityController) FindAllTopics(ctx *gin.Context) {
	communityId := helpers.MustParseUint(ctx.Param("communityId"))
	topics, err := c.topicService.FindTopics(communityId, 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, topics)

}

// Search				godoc
// @Summary				Searches communities the logged user hasn't joined.
// @Param				name query string false "Part of the community name (omit to list all)"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Communities per page (default 10, max 50)"
// @Description			Returns a page of communities the logged user is not a member of, biggest first.
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} transport.SearchCommunitiesResponse{}
// @Router				/communities/search [get]
func (c *CommunityController) Search(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	name := ctx.Query("name")
	page := helpers.MustParseInt(ctx.Query("page"))
	pageSize := helpers.MustParseInt(ctx.Query("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 50 {
		pageSize = 50
	}

	communities, err := c.communityService.SearchCommunitiesToDiscover(user.ID, name, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}
