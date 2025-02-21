package controller

import "yorgut/internal/service"

type PhotoAlbumController struct {
	PhotoAlbumService service.PhotoAlbumService
}

func NewPhotoAlbumController(service service.PhotoAlbumService) *PhotoAlbumController {
	return &PhotoAlbumController{
		PhotoAlbumService: service,
	}
}

func (pa *PhotoAlbumController) Create() {
	panic("unimplemented")
}

func (pa *PhotoAlbumController) Delete() {
	panic("unimplemented")
}

func (pa *PhotoAlbumController) FindOfUser() {
	panic("unimplemented")
}
