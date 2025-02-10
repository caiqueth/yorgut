package service

import (
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoServiceImpl struct {
	PhotoRepository repository.PhotoRepository
	Validate        *validator.Validate
}

// Create implements PhotoService.
func (p *PhotoServiceImpl) Create(photo request.CreatePhotoRequest) {
	panic("unimplemented")
}

// Delete implements PhotoService.
func (p *PhotoServiceImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfAlbum implements PhotoService.
func (p *PhotoServiceImpl) FindOfAlbum(photoAlbumId int) response.FindPhotoOfAlbumResponse {
	panic("unimplemented")
}

// FindOfUser implements PhotoService.
func (p *PhotoServiceImpl) FindOfUser(userId int) response.FindPhotoOfUserResponse {
	panic("unimplemented")
}

func NewPhotoServiceImpl(photoRepository repository.PhotoRepository, validate *validator.Validate) PhotoService {
	return &PhotoServiceImpl{
		PhotoRepository: photoRepository,
		Validate:        validate,
	}
}
