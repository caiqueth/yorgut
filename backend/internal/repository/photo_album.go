package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoAlbumRepository interface {
	Create(pa model.PhotoAlbum) error
	Delete(id int)
	FindById(id int) (model.PhotoAlbum, error)
	FindMany() []model.PhotoAlbum
	FindOfUser(userId int)
}

type PhotoAlbumRepositoryImpl struct {
	Repository[model.PhotoAlbum]
}

func (p *PhotoAlbumRepositoryImpl) FindOfUser(userId int) {
	panic("unimplemented")
}

func NewPhotoAlbumRepositoryImpl(Db *gorm.DB) PhotoAlbumRepository {
	return &PhotoAlbumRepositoryImpl{
		Repository: Repository[model.PhotoAlbum]{Db: Db},
	}
}
