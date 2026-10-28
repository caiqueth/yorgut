package controller

import (
	"net/http"
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
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/fans [get]
func (u *UserFanController) FindSummary(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	summary, err := u.userFanService.FindSummaryOfUser(viewer.ID, userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// FindList			godoc
// @ID					findUserFanList
// @Summary				Gets a page of a user's fans.
// @Param				userId path int true "User ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description			Returns a page of the users who are fans of the user, the newest first
// @Produce				application/json
// @Tags				UserFan
// @Success				200 {object} transport.FindUserFanListResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/fans/list [get]
func (u *UserFanController) FindList(ctx *gin.Context) {
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)

	fans, err := u.userFanService.FindFansOfUser(userId, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, fans)
}

// BecomeFan			godoc
// @ID					becomeFan
// @Summary				Becomes a fan of a user.
// @Param				userId path int true "User ID"
// @Description			Logged user becomes a fan of a friend. Idempotent.
// @Produce				application/json
// @Tags				UserFan
// @Success				204
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/fan [put]
func (u *UserFanController) BecomeFan(ctx *gin.Context) {
	fan := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
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
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/fan [delete]
func (u *UserFanController) StopBeingFan(ctx *gin.Context) {
	fan := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	if err := u.userFanService.StopBeingFan(fan.ID, userId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
