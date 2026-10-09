package controller

import (
	"net/http"
	"yorgut/internal/helpers"
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
// @Description			Returns the album's photos (newest first), how many it holds and its limit.
// @Produce				application/json
// @Tags				Photo
// @Success				200 {object} transport.FindPhotosOfUserResponse{}
// @Failure			401
// @Security			BearerAuth
// @Router				/photos/user/{userId} [get]
func (p *PhotoController) FindOfUser(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))

	photos, err := p.photoService.FindPhotosOfUser(viewer.ID, userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, photos)
}

// Create				godoc
// @ID					uploadPhoto
// @Summary				Uploads a photo to the logged user's album.
// @Description			Adds a PNG or JPEG (max 1 MB and 1024x1024) to the logged user's personal album, which holds up to 15 photos (409 when full).
// @Param				picture formData file true "Photo"
// @Param				caption formData string false "Caption (max 200 characters)"
// @Accept				multipart/form-data
// @Produce				application/json
// @Tags				Photo
// @Success				201 {object} transport.Response{}
// @Failure			400
// @Failure			401
// @Failure			409
// @Security			BearerAuth
// @Router				/photos [post]
func (p *PhotoController) Create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get picture file"})
		return
	}
	// Refuse oversized files before reading them; the service checks again.
	if file.Size > service.MaxPictureBytes {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrPictureTooLarge.Error()})
		return
	}

	picture, err := helpers.FileToBytes(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read picture file"})
		return
	}

	if err := p.photoService.UploadPhoto(user.ID, picture, ctx.PostForm("caption")); err != nil {
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
// @Failure			401
// @Failure			403
// @Failure			404
// @Security			BearerAuth
// @Router				/photos/{photoId} [delete]
func (p *PhotoController) Delete(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	photoId := helpers.MustParseUint(ctx.Param("photoId"))

	if err := p.photoService.DeletePhoto(user.ID, photoId); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
