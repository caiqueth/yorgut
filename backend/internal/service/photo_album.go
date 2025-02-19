package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoAlbumService struct {
	PhotoAlbumRepository repository.PhotoAlbumRepository
	Validate             *validator.Validate
}

func (p *PhotoAlbumService) Create(photoAlbum request.CreatePhotoAlbumRequest) {
	panic("unimplemented")
}

func (p *PhotoAlbumService) Delete(id int) {
	panic("unimplemented")
}

func (p *PhotoAlbumService) FindOfUser(userId int) response.FindPhotoAlmbumOfUserResponse {
	panic("unimplemented")
}

func NewPhotoAlbumService(
	photoAlbumRepository repository.PhotoAlbumRepository,
	validate *validator.Validate,
) *PhotoAlbumService {
	return &PhotoAlbumService{
		PhotoAlbumRepository: photoAlbumRepository,
		Validate:             validate,
	}
}
