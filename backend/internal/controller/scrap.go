package controller

import (
	"net/http"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type ScrapController struct {
	scrapService service.ScrapService
}

func NewScrapController(service service.ScrapService) *ScrapController {
	return &ScrapController{
		scrapService: service,
	}
}

// Create				godoc
// @Summary				Create a Scrap.
// @Description			Creates a Scrap.
// @Param				Scrap body dto.CreateScrapRequest true "Create Scrap"
// @Produce				application/json
// @Tags				Scrap
// @Success				201 {object} dto.Response{}
// @Router				/scrap [post]
func (s *ScrapController) Create(ctx *gin.Context) {

	// WARN: Only accepting AuthorID while we cannot extract it from JWT
	req := dto.CreateScrapRequest{}

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	s.scrapService.CreateScrap(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// Delete 				godoc
// @Summary				Delete Scrap by id.
// @Param				scrapId path string true "delete Scrap by id"
// @Description			Deletes a Scrap.
// @Produce				application/json
// @Tags				Scrap
// @Success				204 {object} dto.Response{}
// @Router				/scrap/{scrapId} [delete]
func (s *ScrapController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}

// FindMany				godoc
// @Summary				Gets a list of scraps of a user.
// @Description			Returns a list of scraps
// @Produce				application/json
// @Tags				User
// @Success				200 {object} dto.FindScrapOfUserRecipientResponse{}
// @Router				/scraps/{userId} [get]
func (s *ScrapController) FindMany(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	scraps, err := s.scrapService.FindScrapOfUserRecipient(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, scraps)

}
