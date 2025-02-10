package service

import (
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoAlbumServiceImpl struct {
	PhotoAlbumRepository repository.PhotoAlbumRepository
	Validate             *validator.Validate
}

// Create implements PhotoAlbumService.
func (p *PhotoAlbumServiceImpl) Create(photoAlbum request.CreatePhotoAlbumRequest) {
	panic("unimplemented")
}

// Delete implements PhotoAlbumService.
func (p *PhotoAlbumServiceImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUser implements PhotoAlbumService.
func (p *PhotoAlbumServiceImpl) FindOfUser(userId int) response.FindPhotoAlmbumOfUserResponse {
	panic("unimplemented")
}

func NewPhotoAlbumServiceImpl(photoAlbumRepository repository.PhotoAlbumRepository, validate *validator.Validate) PhotoAlbumService {
	return &PhotoAlbumServiceImpl{
		PhotoAlbumRepository: photoAlbumRepository,
		Validate:             validate,
	}
}
