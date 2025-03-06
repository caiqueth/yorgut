package service

import (
	"errors"
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type PhotoAlbumService struct {
	photoAlbumRepository repository.PhotoAlbumRepository
	UserRepository       repository.UserRepository
	Validate             *validator.Validate
}

func (p *PhotoAlbumService) Create(photoAlbum dto.CreatePhotoAlbumRequest) error {
	owner, err := p.UserRepository.FindById(photoAlbum.UserID)
	if err != nil {
		return errors.New("user not found")
	}
	photoAlbumModel := model.NewPhotoAlbum(owner, photoAlbum.Name)

	err = p.photoAlbumRepository.Create(photoAlbumModel)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoAlbumService) Delete(id uint) error {
	err := p.photoAlbumRepository.Delete(id)
	if err != nil {
		return err
	}
	return err
}

func (p *PhotoAlbumService) FindOfUser(userId int) (*dto.FindPhotoAlbumOfUserResponse, error) {
	photoAlbums, err := p.photoAlbumRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPhotoAlbumOfUserResponse{}
	resp.FromPhotoAlbums(photoAlbums)

	return resp, nil
}

func NewPhotoAlbumService(
	photoAlbumRepository repository.PhotoAlbumRepository,
	validate *validator.Validate,
) *PhotoAlbumService {
	return &PhotoAlbumService{
		photoAlbumRepository: photoAlbumRepository,
		Validate:             validate,
	}
}
