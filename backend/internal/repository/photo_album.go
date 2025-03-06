package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoAlbumRepository struct {
	Repository[*model.PhotoAlbum]
}

func (p *PhotoAlbumRepository) FindOfUser(userId int) ([]model.PhotoAlbum, error) {
	var ents []model.PhotoAlbum
	result := p.Db.Where("user_id = ?", userId).Find(&ents)
	if result.Error != nil {
		return nil, result.Error
	}
	return ents, nil
}

func NewPhotoAlbumRepository(Db *gorm.DB) *PhotoAlbumRepository {
	return &PhotoAlbumRepository{
		Repository: Repository[*model.PhotoAlbum]{Db: Db},
	}
}
