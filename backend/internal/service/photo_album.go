package service

import (
	"errors"
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoAlbumService struct {
	PhotoAlbumRepository repository.PhotoAlbumRepository
	UserRepository       repository.UserRepository
	Validate             *validator.Validate
}

func (p *PhotoAlbumService) Create(photoAlbum dto.CreatePhotoAlbumRequest) error {
	owner, err := p.UserRepository.FindById(photoAlbum.UserID)
	if err != nil {
		return errors.New("user not found")
	}
	photoAlbumModel := model.NewPhotoAlbum(owner, photoAlbum.Name)

	err = p.PhotoAlbumRepository.Create(photoAlbumModel)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoAlbumService) Delete(id uint) error {
	err := p.PhotoAlbumRepository.Delete(id)
	if err != nil {
		return err
	}
	return err
}

func (p *PhotoAlbumService) FindOfUser(userId int) (dto.FindPhotoAlmbumOfUserResponse, error) {
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
