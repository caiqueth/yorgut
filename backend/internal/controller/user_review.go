package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type UserReviewController struct {
	userReviewService service.UserReviewService
}

func NewUserReviewController(service service.UserReviewService) *UserReviewController {
	return &UserReviewController{
		userReviewService: service,
	}
}

// FindSummary			godoc
// @ID					findUserReviews
// @Summary				Gets the review summary of a user.
// @Param				userId path int true "User ID"
// @Description			Returns the average cool, sexy and trustworthy ratings of a user, plus the logged user's own review.
// @Produce				application/json
// @Tags				UserReview
// @Success				200 {object} transport.UserReviewSummaryResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/reviews [get]
func (u *UserReviewController) FindSummary(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	summary, err := u.userReviewService.FindSummaryOfUser(viewer.ID, userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// Upsert				godoc
// @ID					upsertUserReview
// @Summary				Reviews a user.
// @Param				userId path int true "User ID"
// @Param				UserReview body transport.UpsertUserReviewRequest true "Review (each trait from 0 to 3)"
// @Accept				application/json
// @Description			Logged user creates or updates their review of a friend.
// @Produce				application/json
// @Tags				UserReview
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/review [put]
func (u *UserReviewController) Upsert(ctx *gin.Context) {
	reviewer := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	request := transport.UpsertUserReviewRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := u.userReviewService.UpsertReview(reviewer.ID, userId, *request.Cool, *request.Sexy, *request.Trustworthy); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
