package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
)

// Picture limits (size, format, resolution) live in picture.go.
const (
	MaxPhotosPerAlbum = 12
	MaxCaptionLength  = 200
)

var (
	ErrPhotoCaptionTooLong = newError(ErrInvalid, fmt.Sprintf("A legenda pode ter no máximo %d caracteres.", MaxCaptionLength))
	ErrAlbumFull           = newError(ErrConflict, fmt.Sprintf("O álbum já tem %d fotos. Exclua uma para adicionar outra.", MaxPhotosPerAlbum))
	ErrPhotoNotFound       = newError(ErrNotFound, "A foto não existe.")
	ErrPhotoForbidden      = newError(ErrForbidden, "Esta foto não é sua.")
)

type PhotoService struct {
	photoRepository      repository.PhotoRepository
	photoAlbumRepository repository.PhotoAlbumRepository
	userRepository       repository.UserRepository
	friendshipRepository repository.FriendshipRepository
	store                storage.ObjectStore
}

func NewPhotoService(
	photoRepository repository.PhotoRepository,
	photoAlbumRepository repository.PhotoAlbumRepository,
	userRepository repository.UserRepository,
	friendshipRepository repository.FriendshipRepository,
	store storage.ObjectStore,
) *PhotoService {
	return &PhotoService{
		photoRepository:      photoRepository,
		photoAlbumRepository: photoAlbumRepository,
		userRepository:       userRepository,
		friendshipRepository: friendshipRepository,
		store:                store,
	}
}

// FindPhotosOfUser lists the user's personal album; canUpload tells whether
// the viewer owns it.
func (p *PhotoService) FindPhotosOfUser(viewerId uint, userId uint) (*transport.FindPhotosOfUserResponse, error) {
	if err := requireAccess(p.userRepository, p.friendshipRepository, viewerId, userId, sectionAlbum); err != nil {
		return nil, err
	}

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
func (p *PhotoService) UploadPhoto(ctx context.Context, userId uint, picture []byte, caption string) error {
	caption = strings.TrimSpace(caption)
	if len([]rune(caption)) > MaxCaptionLength {
		return ErrPhotoCaptionTooLong
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

	// After the cheap checks: resizing a big photo takes a moment.
	picture, err = preparePicture(picture)
	if err != nil {
		return err
	}

	key, err := storage.PutImage(ctx, p.store, picture)
	if err != nil {
		return err
	}

	if err := p.photoRepository.Create(model.NewPhoto(album, key, caption)); err != nil {
		deletePicture(ctx, p.store, &key)
		return err
	}
	return nil
}

// DeletePhoto removes a photo and its file from the user's album; only the
// owner may.
func (p *PhotoService) DeletePhoto(ctx context.Context, userId uint, photoId uint) error {
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

	if err := p.photoRepository.Delete(photoId); err != nil {
		return err
	}
	deletePicture(ctx, p.store, &photo.URL)
	return nil
}
