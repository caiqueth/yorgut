package controller

import "yorgut/internal/service"

type PhotoController struct {
	PhotoService service.PhotoService
}

func NewPhotoController(service service.PhotoService) *PhotoController {
	return &PhotoController{
		PhotoService: service,
	}
}

func (p *PhotoController) Create() {
	panic("unimplemented")
}

func (p *PhotoController) Delete() {
	panic("unimplemented")
}

func (p *PhotoController) FindOfAlbum() {
	panic("unimplemented")
}

func (p *PhotoController) FindOfUser() {
	panic("unimplemented")
}
