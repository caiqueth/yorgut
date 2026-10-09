package service

import (
	"errors"
	"fmt"
	"strings"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

// Picture limits (size, format, resolution) live in picture.go.
const (
	MaxPhotosPerAlbum = 15
	MaxCaptionLength  = 200
)

var (
	ErrPhotoCaptionTooLong = newError(ErrInvalid, fmt.Sprintf("the caption must be at most %d characters", MaxCaptionLength))
	ErrAlbumFull           = newError(ErrConflict, fmt.Sprintf("the album already has %d photos", MaxPhotosPerAlbum))
	ErrPhotoNotFound       = newError(ErrNotFound, "photo not found")
	ErrPhotoForbidden      = newError(ErrForbidden, "this photo does not belong to this user")
)

type PhotoService struct {
	photoRepository      repository.PhotoRepository
	photoAlbumRepository repository.PhotoAlbumRepository
	Validate             *validator.Validate
}

func NewPhotoService(
	photoRepository repository.PhotoRepository,
	photoAlbumRepository repository.PhotoAlbumRepository,
	validate *validator.Validate,
) *PhotoService {
	return &PhotoService{
		photoRepository:      photoRepository,
		photoAlbumRepository: photoAlbumRepository,
		Validate:             validate,
	}
}

// FindPhotosOfUser lists the user's personal album; canUpload tells whether
// the viewer owns it.
func (p *PhotoService) FindPhotosOfUser(viewerId uint, userId uint) (*transport.FindPhotosOfUserResponse, error) {
	album, err := p.photoAlbumRepository.FindOrCreatePersonalOfUser(userId)
	if err != nil {
		return nil, err
	}

	photos, err := p.photoRepository.FindOfAlbum(album.ID)
	if err != nil {
		return nil, err
	}

	response := transport.FindPhotosOfUserResponse{}
	response.FromPhotos(photos, MaxPhotosPerAlbum, viewerId == userId)

	return &response, nil
}

// UploadPhoto adds a photo to the user's personal album.
func (p *PhotoService) UploadPhoto(userId uint, picture []byte, caption string) error {
	caption = strings.TrimSpace(caption)
	if len([]rune(caption)) > MaxCaptionLength {
		return ErrPhotoCaptionTooLong
	}

	if err := validatePicture(picture); err != nil {
		return err
	}

	album, err := p.photoAlbumRepository.FindOrCreatePersonalOfUser(userId)
	if err != nil {
		return err
	}

	count, err := p.photoRepository.CountOfAlbum(album.ID)
	if err != nil {
		return err
	}
	if count >= MaxPhotosPerAlbum {
		return ErrAlbumFull
	}

	storedFileName, err := helpers.UploadPicture(picture)
	if err != nil {
		return err
	}

	return p.photoRepository.Create(model.NewPhoto(album, *storedFileName, caption))
}

// DeletePhoto removes a photo from the user's album; only the owner may.
func (p *PhotoService) DeletePhoto(userId uint, photoId uint) error {
	photo, err := p.photoRepository.FindByIdWithAlbum(photoId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrPhotoNotFound
	}
	if err != nil {
		return err
	}

	if photo.PhotoAlbum == nil || photo.PhotoAlbum.UserID != userId {
		return ErrPhotoForbidden
	}

	return p.photoRepository.Delete(photoId)
}
