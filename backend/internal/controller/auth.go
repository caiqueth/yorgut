package controller

import (
	"net/http"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// LoginWithGoogle		godoc
// @ID					loginWithGoogle
// @Summary				Logs in with a Google account.
// @Description			Validates the Google ID token from Sign in with Google, creates the user on their first login, and returns a Yorgut session token. Send it as "Authorization: Bearer <token>" on every other request.
// @Accept				application/json
// @Produce				application/json
// @Tags				Auth
// @Param				request body transport.LoginWithGoogleRequest true "Google credential"
// @Success				200 {object} transport.LoginResponse{}
// @Failure				400 {object} transport.ErrorResponse
// @Failure				401 {object} transport.ErrorResponse
// @Router				/auth/google [post]
func (a *AuthController) LoginWithGoogle(ctx *gin.Context) {
	var req transport.LoginWithGoogleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	resp, err := a.authService.LoginWithGoogle(ctx.Request.Context(), req.Credential)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
