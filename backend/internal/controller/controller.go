package controller

import (
	"errors"
	"net/http"
	"strconv"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

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
	// Users read the message (the frontend shows it): only the user-facing
	// part, never what wrapped it on the way up.
	message := service.UserMessage(err)
	switch {
	case errors.Is(err, service.ErrInvalidSession):
		status, message = http.StatusUnauthorized, service.ErrInvalidSession.Error()
	case errors.Is(err, service.ErrInvalidGoogleCredential):
		// Wrapped with Google's reason, which goes to the log only.
		log.Info().Err(err).Msg("google login refused")
		status, message = http.StatusUnauthorized, service.ErrInvalidGoogleCredential.Error()
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
		ctx.JSON(http.StatusInternalServerError, transport.ErrorResponse{Error: "Algo deu errado do nosso lado. Tente de novo mais tarde."})
		return
	}

	ctx.JSON(status, transport.ErrorResponse{Error: message})
}

// Page sizes of the paginated lists (friends, members, communities).
const (
	defaultPageSize = 20
	maxPageSize     = 50
)

// respondBadRequest answers a request the handler couldn't read (a malformed
// body, path param or upload), before any service is called. message is
// shown to the user.
func respondBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Error: message})
}

// respondInvalidRequest answers a body that failed binding or validation.
// The validator's text is technical and in English, so it only goes to the
// log; the frontend checks the same rules before sending.
func respondInvalidRequest(ctx *gin.Context, err error) {
	log.Debug().Err(err).Str("path", ctx.FullPath()).Msg("invalid request")
	respondBadRequest(ctx, "Os dados enviados são inválidos.")
}

// viewerOf is whom a profile read is for: the logged user, or, with
// ?as=stranger, "anyone who isn't a friend" (viewer 0), which shows users
// their own profile as non-friends see it.
func viewerOf(ctx *gin.Context) uint {
	if ctx.Query("as") == "stranger" {
		return 0
	}
	return ctx.MustGet("user").(*model.User).ID
}

// idParam reads a numeric id from the path (":userId", ":communityId", ...).
// On a malformed or zero id it answers 400 and returns false; the handler
// must return right away.
func idParam(ctx *gin.Context, name string) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param(name), 10, 0)
	if err != nil || id == 0 {
		respondBadRequest(ctx, "Endereço inválido.")
		return 0, false
	}
	return uint(id), true
}

// readPicture reads the "picture" file of a multipart upload, refusing
// oversized files before reading them (the service validates the picture
// again). When optional, a missing file gives nil. On failure it answers 400
// and returns false; the handler must return right away.
func readPicture(ctx *gin.Context, optional bool) ([]byte, bool) {
	file, err := ctx.FormFile("picture")
	switch {
	case optional && errors.Is(err, http.ErrMissingFile):
		return nil, true
	case err != nil:
		respondBadRequest(ctx, "A foto não chegou. Tente enviar de novo.")
		return nil, false
	case file.Size > service.MaxUploadBytes:
		respondBadRequest(ctx, service.ErrPictureTooLarge.Error())
		return nil, false
	}

	picture, err := helpers.FileToBytes(file)
	if err != nil {
		respondBadRequest(ctx, "Não foi possível ler a foto. Tente enviar de novo.")
		return nil, false
	}
	return picture, true
}

// queryInt reads an integer query param; missing or malformed gives 0.
func queryInt(ctx *gin.Context, name string) int {
	value, _ := strconv.Atoi(ctx.Query(name))
	return value
}

// pageQuery reads the page and pageSize query params: page defaults to 1, and
// pageSize to defaultPageSize, capped at maxPageSize.
func pageQuery(ctx *gin.Context) (page int, pageSize int) {
	page = queryInt(ctx, "page")
	pageSize = queryInt(ctx, "pageSize")
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

// notImplemented answers routes whose handler has not been written yet.
func notImplemented(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, transport.ErrorResponse{Error: "Ainda não está pronto."})
}
