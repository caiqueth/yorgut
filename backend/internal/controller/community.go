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

// CreateCommunities	godoc
// @Symmary				Create Community
// @Decription			Save community in Db.
// @Param				community body request.CreateCommunityRequest true "Create Community"
// @Produce				application/json
// @Tags				community
// @Success				200 {object} response.Response{}
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

	ctx.JSON(http.StatusOK, webResponse)
}

// FindByIdCommunity 	godoc
// @Summary				Get Single community by id.
// @Param				communityId path string true "get community by id"
// @Description			Returns the community matching the requested id
// @Produce				application/json
// @Tags				tags
// @Success				200 {object} response.FindCommunityByIdResponse{}
// @Router				/community/{communityId} [get]
func (controller *CommunityController) FindCommunityById(ctx *gin.Context) {
	communityId := ctx.Param("communityId")
	id, err := strconv.Atoi(communityId)
	if err != nil {
		panic(err)
	}

	community := controller.communityService.FindById(id)

	webResponse := response.Response{
		Data: community,
	}
	ctx.JSON(http.StatusOK, webResponse)

}
