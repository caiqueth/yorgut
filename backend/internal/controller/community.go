package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
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
// @Param				community body request.CreateCommunityRequest true "Create Community"
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} response.Response{}
// @Router				/community [post]
func (controller *CommunityController) Create(ctx *gin.Context) {
	createCommunityRequest := request.CreateCommunityRequest{}

	err := ctx.ShouldBindJSON(&createCommunityRequest)
	if err != nil {
		panic(err)
	}

	controller.communityService.Create(createCommunityRequest)
	webResponse := response.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusNoContent, webResponse)
}

// FindById				godoc
// @Summary				Updates a community.
// @Description			Updates properties of a single community
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} response.Response{}
// @Router				/community [put]
func (controller *CommunityController) Update(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete community by id.
// @Param				communityId path string true "delete community by id"
// @Description			Deletes a community
// @Produce				application/json
// @Tags				Community
// @Success				204 {object} response.Response{}
// @Router				/community/{communityId} [delete]
func (controller *CommunityController) Delete(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		panic(err)
	}

	err = controller.communityService.Delete(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, response.Response{})

}

// FindById				godoc
// @Summary				Get Single community by id.
// @Param				communityId path string true "get community by id"
// @Description			Returns the community matching the requested id
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} response.FindCommunityByIdResponse{}
// @Router				/community/{communityId} [get]
func (controller *CommunityController) FindById(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.ParseUint(communityId, 0, 0)
	if err != nil {
		panic(err)
	}

	community, err := controller.communityService.FindById(uint(id))

	webResponse := response.Response{
		Data: community,
	}
	ctx.JSON(http.StatusOK, webResponse)

}

// FindMany				godoc
// @Summary				Get Multiple communities.
// @Description			Returns an array of communities
// @Produce				application/json
// @Tags				Community
// @Success				200 {object} response.FindManyCommunitiesResponse{}
// @Router				/community [get]
func (controller *CommunityController) FindMany(ctx *gin.Context) {
	panic("unimplemented")
}
