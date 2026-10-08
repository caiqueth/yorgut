package controller

import (
	"net/http"
	"yorgut/internal/helpers"
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
// @Summary				Gets the review summary of a user.
// @Param				userId path string true "User ID"
// @Description			Returns the average cool, sexy and trustworthy ratings of a user, plus the logged user's own review.
// @Produce				application/json
// @Tags				UserReview
// @Success				200 {object} transport.UserReviewSummaryResponse{}
// @Router				/users/{userId}/reviews [get]
func (u *UserReviewController) FindSummary(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	summary, err := u.userReviewService.FindSummaryOfUser(viewer.ID, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// Upsert				godoc
// @Summary				Reviews a user.
// @Param				userId path string true "User ID"
// @Param				UserReview body transport.UpsertUserReviewRequest true "Review (each trait from 0 to 3)"
// @Description			Logged user creates or updates their review of a friend.
// @Produce				application/json
// @Tags				UserReview
// @Success				204
// @Router				/users/{userId}/review [put]
func (u *UserReviewController) Upsert(ctx *gin.Context) {
	reviewer := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	request := transport.UpsertUserReviewRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !u.userReviewService.CanReview(reviewer.ID, userId) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only friends can review this user"})
		return
	}

	if err := u.userReviewService.UpsertReview(reviewer.ID, userId, *request.Cool, *request.Sexy, *request.Trustworthy); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
