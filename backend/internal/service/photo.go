package service

import (
	"errors"
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoService struct {
	PhotoRepository      repository.PhotoRepository
	PhotoAlbumRepository repository.PhotoAlbumRepository
	Validate             *validator.Validate
}

func (p *PhotoService) Create(photo dto.CreatePhotoRequest) error {
	photoAlbum, err := p.PhotoAlbumRepository.FindById(photo.PhotoAlbumID)
	if err != nil {
		return errors.New("photo album not found")
	}

	photoModel := model.NewPhoto(photoAlbum, photo.URL, photo.Caption)

	err = p.PhotoRepository.Create(photoModel)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoService) Delete(id uint) error {
	err := p.PhotoRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoService) FindOfAlbum(photoAlbumId uint) dto.FindPhotoOfAlbumResponse {
	panic("unimplemented")
}

func (p *PhotoService) FindOfUser(userId uint) dto.FindPhotoOfUserResponse {
	panic("unimplemented")
}

func NewPhotoService(photoRepository repository.PhotoRepository, validate *validator.Validate) *PhotoService {
	return &PhotoService{
		PhotoRepository: photoRepository,
		Validate:        validate,
	}
}
