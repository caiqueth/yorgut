package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoRepository struct {
	repository[*model.Photo]
}

func NewPhotoRepository(Db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{
		repository: repository[*model.Photo]{db: Db},
	}
}

// FindOfAlbum returns the album's photos, newest first.
func (p *PhotoRepository) FindOfAlbum(albumId uint) ([]model.Photo, error) {
	var photos []model.Photo
	result := p.db.
		Where("photo_album_id = ?", albumId).
		Order("created_at desc").
		Find(&photos)

	if result.Error != nil {
		return nil, result.Error
	}
	return photos, nil
}

func (p *PhotoRepository) CountOfAlbum(albumId uint) (int64, error) {
	var count int64
	result := p.db.
		Model(&model.Photo{}).
		Where("photo_album_id = ?", albumId).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// FindByIdWithAlbum returns the photo with its album, so callers can check
// who owns it.
func (p *PhotoRepository) FindByIdWithAlbum(photoId uint) (*model.Photo, error) {
	var photo model.Photo
	result := p.db.Preload("PhotoAlbum").First(&photo, photoId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &photo, nil
}
