package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto"

	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	communityService service.CommunityService
}

func NewCommunityController(service service.CommunityService) *CommunityController {
	return &CommunityController{
		communityService: service,
	}
}

// Create				godoc
// @Summary				Create Community
// @Decription			Save community in Db.
// @Param				community body dto.CreateCommunityRequest true "Create Community"
// @Produce				application/json
// @Tags				Community
// @Success				201 {object} dto.Response{}
// @Router				/community [post]
func (c *CommunityController) Create(ctx *gin.Context) {
	createCommunityRequest := dto.CreateCommunityRequest{}

	err := ctx.ShouldBindJSON(&createCommunityRequest)
	if err != nil {
		panic(err)
	}

	err = c.communityService.Create(createCommunityRequest)
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

	err = c.communityService.Delete(uint(id))
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

	community, err := c.communityService.FindById(uint(id))
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
	communities, err := c.communityService.FindMany()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}
