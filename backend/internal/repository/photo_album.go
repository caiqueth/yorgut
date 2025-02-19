package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoAlbumRepository struct {
	Repository[model.PhotoAlbum]
}

func (p *PhotoAlbumRepository) FindOfUser(userId int) {
	panic("unimplemented")
}

func NewPhotoAlbumRepository(Db *gorm.DB) *PhotoAlbumRepository {
	return &PhotoAlbumRepository{
		Repository: Repository[model.PhotoAlbum]{Db: Db},
	}
}
