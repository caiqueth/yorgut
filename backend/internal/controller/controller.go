package controller

import (
	"errors"
	"net/http"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RespondError answers a failed service call. It is the only place that
// decides which HTTP status a service error becomes:
//
//	ErrInvalidSession, ErrInvalidGoogleCredential -> 401
//	ErrNotFound -> 404, ErrForbidden -> 403, ErrConflict -> 409, ErrInvalid -> 400
//
// Anything else is a bug or an infrastructure failure: it is logged, and the
// client only gets a generic 500 so internals don't leak. Handlers must
// return right after calling it.
func RespondError(ctx *gin.Context, err error) {
	var status int
	switch {
	case errors.Is(err, service.ErrInvalidSession),
		errors.Is(err, service.ErrInvalidGoogleCredential):
		status = http.StatusUnauthorized
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, service.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, service.ErrInvalid):
		status = http.StatusBadRequest
	default:
		log.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("path", ctx.FullPath()).
			Msg("request failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(status, gin.H{"error": err.Error()})
}

// notImplemented answers routes whose handler has not been written yet.
func notImplemented(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
