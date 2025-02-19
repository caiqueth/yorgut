package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoService struct {
	PhotoRepository repository.PhotoRepository
	Validate        *validator.Validate
}

func (p *PhotoService) Create(photo request.CreatePhotoRequest) {
	panic("unimplemented")
}

func (p *PhotoService) Delete(id int) {
	panic("unimplemented")
}

func (p *PhotoService) FindOfAlbum(photoAlbumId int) response.FindPhotoOfAlbumResponse {
	panic("unimplemented")
}

func (p *PhotoService) FindOfUser(userId int) response.FindPhotoOfUserResponse {
	panic("unimplemented")
}

func NewPhotoService(photoRepository repository.PhotoRepository, validate *validator.Validate) *PhotoService {
	return &PhotoService{
		PhotoRepository: photoRepository,
		Validate:        validate,
	}
}
