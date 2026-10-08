package controller

import (
	"net/http"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

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
// @Param				Scrap body transport.CreateScrapRequest true "Create Scrap"
// @Produce				application/json
// @Tags				Scrap
// @Success				201 {object} transport.Response{}
// @Router				/scraps [post]
func (s *ScrapController) Create(ctx *gin.Context) {
	author := ctx.MustGet("user").(*model.User)
	// WARN: Only accepting AuthorID while we cannot extract it from JWT
	req := transport.CreateScrapRequest{}

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.scrapService.CreateScrap(author.ID, req.UserRecipientID, req.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := transport.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// Delete 				godoc
// @Summary				Delete Scrap by id.
// @Param				scrapId path string true "Scrap ID"
// @Description			Deletes a Scrap.
// @Produce				application/json
// @Tags				Scrap
// @Success				204 {object} transport.Response{}
// @Router				/scraps/{scrapId} [delete]
func (s *ScrapController) Delete(ctx *gin.Context) {
	notImplemented(ctx)
}

// FindOfUser			godoc
// @Summary				Gets a list of scraps of a user.
// @Param				userId path string true "User ID"
// @Description			Returns a list of scraps
// @Produce				application/json
// @Tags				Scrap
// @Success				200 {object} transport.FindScrapOfUserRecipientResponse{}
// @Router				/scraps/user/{userId} [get]
func (s *ScrapController) FindOfUser(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	scraps, err := s.scrapService.FindScrapOfUserRecipient(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, scraps)

}
