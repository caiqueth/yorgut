package controller

import (
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type PhotoAlbumController struct {
	PhotoAlbumService service.PhotoAlbumService
}

func NewPhotoAlbumController(service service.PhotoAlbumService) *PhotoAlbumController {
	return &PhotoAlbumController{
		PhotoAlbumService: service,
	}
}

// Create				godoc
// @Summary				Create a PhotoAlbum.
// @Description			Creates a PhotoAlbum.
// @Param				PhotoAlbum body request.CreatePhotoAlbumRequest true "Create PhotoAlbum"
// @Produce				application/json
// @Tags				PhotoAlbum
// @Success				201 {object} response.Response{}
// @Router				/photoAlbum [post]
func (pa *PhotoAlbumController) Create(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete PhotoAlbum by id.
// @Param				photoAlbumId path string true "delete PhotoAlbum by id"
// @Description			Deletes a PhotoAlbum.
// @Produce				application/json
// @Tags				PhotoAlbum
// @Success				204 {object} response.Response{}
// @Router				/photoAlbum/{photoAlbumId} [delete]
func (pa *PhotoAlbumController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}

// Update				godoc
// @Summary				Updates a PhotoAlbum.
// @Description			Updates properties of a single PhotoAlbum
// @Produce				application/json
// @Tags				PhotoAlbum
// @Success				204 {object} response.Response{}
// @Router				/photoAlbum [put]
func (pa *PhotoAlbumController) Update(ctx *gin.Context) {
	panic("unimplemented")
}

// FindMany				godoc
// @Summary				Get a list of PhotoAlbums.
// @Description			Returns an array of PhotoAlbums.
// @Produce				application/json
// @Tags				PhotoAlbum
// @Success				200 {object} response.Response{}
// @Router				/photoAlbum [get]
func (pa *PhotoAlbumController) FindPhotos(ctx *gin.Context) {
	panic("unimplemented")
}
