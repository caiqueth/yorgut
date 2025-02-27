package controller

import (
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type ScrapController struct {
	ScrapService service.ScrapService
}

func NewScrapController(service service.ScrapService) *ScrapController {
	return &ScrapController{
		ScrapService: service,
	}
}

// Create				godoc
// @Summary				Create a Scrap.
// @Description			Creates a Scrap.
// @Param				Scrap body request.CreateScrapRequest true "Create Scrap"
// @Produce				application/json
// @Tags				Scrap
// @Success				201 {object} response.Response{}
// @Router				/scrap [post]
func (s *ScrapController) Create(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete Scrap by id.
// @Param				scrapId path string true "delete Scrap by id"
// @Description			Deletes a Scrap.
// @Produce				application/json
// @Tags				Scrap
// @Success				204 {object} response.Response{}
// @Router				/scrap/{scrapId} [delete]
func (s *ScrapController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}
