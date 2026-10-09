package controller

import (
	"net/http"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type UserFanController struct {
	userFanService service.UserFanService
}

func NewUserFanController(service service.UserFanService) *UserFanController {
	return &UserFanController{
		userFanService: service,
	}
}

// FindSummary			godoc
// @ID					findUserFans
// @Summary				Gets the fan summary of a user.
// @Param				userId path int true "User ID"
// @Description			Returns how many fans a user has, and whether the logged user is (or can become) one of them.
// @Produce				application/json
// @Tags				UserFan
// @Success				200 {object} transport.UserFanSummaryResponse{}
// @Failure			401
// @Security			BearerAuth
// @Router				/users/{userId}/fans [get]
func (u *UserFanController) FindSummary(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	summary, err := u.userFanService.FindSummaryOfUser(viewer.ID, userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// BecomeFan			godoc
// @ID					becomeFan
// @Summary				Becomes a fan of a user.
// @Param				userId path int true "User ID"
// @Description			Logged user becomes a fan of a friend. Idempotent.
// @Produce				application/json
// @Tags				UserFan
// @Success				204
// @Failure			401
// @Failure			403
// @Security			BearerAuth
// @Router				/users/{userId}/fan [put]
func (u *UserFanController) BecomeFan(ctx *gin.Context) {
	fan := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	if !u.userFanService.CanBecomeFan(fan.ID, userId) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only friends can become fans of this user"})
		return
	}

	if err := u.userFanService.BecomeFan(fan.ID, userId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// StopBeingFan			godoc
// @ID					stopBeingFan
// @Summary				Stops being a fan of a user.
// @Param				userId path int true "User ID"
// @Description			Logged user removes their star from a user. Idempotent.
// @Produce				application/json
// @Tags				UserFan
// @Success				204
// @Failure			401
// @Security			BearerAuth
// @Router				/users/{userId}/fan [delete]
func (u *UserFanController) StopBeingFan(ctx *gin.Context) {
	fan := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	if err := u.userFanService.StopBeingFan(fan.ID, userId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
