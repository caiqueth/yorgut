package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoAlbumRepository struct {
	repository[*model.PhotoAlbum]
}

func NewPhotoAlbumRepository(Db *gorm.DB) *PhotoAlbumRepository {
	return &PhotoAlbumRepository{
		repository: repository[*model.PhotoAlbum]{db: Db},
	}
}

// FindOrCreatePersonalOfUser returns the user's personal album (their first
// one), creating it on first use.
func (p *PhotoAlbumRepository) FindOrCreatePersonalOfUser(userId uint) (*model.PhotoAlbum, error) {
	var album model.PhotoAlbum
	result := p.db.
		Where(model.PhotoAlbum{UserID: userId}).
		Attrs(model.PhotoAlbum{Name: "meu álbum"}).
		Order("id").
		FirstOrCreate(&album)

	if result.Error != nil {
		return nil, result.Error
	}
	return &album, nil
}
