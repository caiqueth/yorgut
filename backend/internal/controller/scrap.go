package controller

import (
	"net/http"
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
// @ID					createScrap
// @Summary				Create a Scrap.
// @Description			Creates a Scrap.
// @Param				Scrap body transport.CreateScrapRequest true "Create Scrap"
// @Accept				application/json
// @Produce				application/json
// @Tags				Scrap
// @Success				201 {object} transport.Response{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/scraps [post]
func (s *ScrapController) Create(ctx *gin.Context) {
	author := ctx.MustGet("user").(*model.User)

	req := transport.CreateScrapRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := s.scrapService.CreateScrap(author.ID, req.UserRecipientID, req.Content); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}

// Delete 				godoc
// @ID					deleteScrap
// @Summary				Deletes a scrap.
// @Param				scrapId path int true "Scrap ID"
// @Description			Deletes the scrap for good. Only the owner of the scrapbook it is in and its author may (403 otherwise).
// @Produce				application/json
// @Tags				Scrap
// @Success				204
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/scraps/{scrapId} [delete]
func (s *ScrapController) Delete(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	scrapId, ok := idParam(ctx, "scrapId")
	if !ok {
		return
	}

	if err := s.scrapService.DeleteScrap(user.ID, scrapId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// FindOfUser			godoc
// @ID					findUserScraps
// @Summary				Gets a page of a user's scraps.
// @Param				userId path int true "User ID"
// @Param				as query string false "\"stranger\": as someone who isn't a friend sees it (the profile preview)"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns a page of the scraps sent to the user, newest first
// @Produce				application/json
// @Tags				Scrap
// @Success				200 {object} transport.FindScrapOfUserRecipientResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/scraps/user/{userId} [get]
func (s *ScrapController) FindOfUser(ctx *gin.Context) {
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)
	scraps, err := s.scrapService.FindScrapOfUserRecipient(viewerOf(ctx), userId, page, pageSize)

	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, scraps)
}
