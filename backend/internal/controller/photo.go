package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type PhotoController struct {
	photoService service.PhotoService
}

func NewPhotoController(service service.PhotoService) *PhotoController {
	return &PhotoController{
		photoService: service,
	}
}

// FindOfUser			godoc
// @ID					findUserPhotos
// @Summary				Lists the photos of a user's personal album.
// @Param				userId path int true "User ID"
// @Param				as query string false "\"stranger\": as someone who isn't a friend sees it (the profile preview)"
// @Description			Returns the album's photos (newest first), how many it holds and its limit.
// @Produce				application/json
// @Tags				Photo
// @Success				200 {object} transport.FindPhotosOfUserResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/photos/user/{userId} [get]
func (p *PhotoController) FindOfUser(ctx *gin.Context) {
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	photos, err := p.photoService.FindPhotosOfUser(viewerOf(ctx), userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, photos)
}

// Create				godoc
// @ID					uploadPhoto
// @Summary				Uploads a photo to the logged user's album.
// @Description			Adds a PNG or JPEG (max 5 MB and 50 megapixels; stored shrunk to 1024px and 1 MB) to the logged user's personal album, which holds up to 12 photos (409 when full).
// @Param				picture formData file true "Photo"
// @Param				caption formData string false "Caption (max 200 characters)"
// @Accept				multipart/form-data
// @Produce				application/json
// @Tags				Photo
// @Success				201 {object} transport.Response{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/photos [post]
func (p *PhotoController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	picture, ok := readPicture(ctx, false)
	if !ok {
		return
	}

	if err := p.photoService.UploadPhoto(ctx.Request.Context(), user.ID, picture, ctx.PostForm("caption")); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}

// Delete 				godoc
// @ID					deletePhoto
// @Summary				Deletes a photo from the logged user's album.
// @Param				photoId path int true "Photo ID"
// @Description			Removes one of the logged user's photos.
// @Produce				application/json
// @Tags				Photo
// @Success				204
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/photos/{photoId} [delete]
func (p *PhotoController) Delete(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	photoId, ok := idParam(ctx, "photoId")
	if !ok {
		return
	}

	if err := p.photoService.DeletePhoto(ctx.Request.Context(), user.ID, photoId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
