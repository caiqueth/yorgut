package controller

import (
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type PhotoController struct {
	PhotoService service.PhotoService
}

func NewPhotoController(service service.PhotoService) *PhotoController {
	return &PhotoController{
		PhotoService: service,
	}
}

// Create				godoc
// @Summary				Create a Photo.
// @Description			Creates a Photo.
// @Param				Photo body request.CreatePhotoRequest true "Create Photo"
// @Produce				application/json
// @Tags				Photo
// @Success				204 {object} response.Response{}
// @Router				/photo [post]
func (p *PhotoController) Create(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete Photo by id.
// @Param				photoId path string true "delete Photo by id"
// @Description			Deletes a Photo.
// @Produce				application/json
// @Tags				Photo
// @Success				204 {object} response.Response{}
// @Router				/photo/{photoId} [delete]
func (p *PhotoController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}
