package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoRepository interface {
	Create(p model.Photo) error
	Delete(id int)
	FindById(id int) (model.Photo, error)
	FindMany() []model.Photo
	FindOfUser(userId int)
	FindOfPhotoAlbum(photoAlbumId int)
}

type PhotoRepositoryImpl struct {
	Repository[model.Photo]
}

func (p *PhotoRepositoryImpl) FindOfPhotoAlbum(photoAlbumId int) {
	panic("unimplemented")
}

func (p *PhotoRepositoryImpl) FindOfUser(userId int) {
	panic("unimplemented")
}

func NewPhotoRepositoryImpl(Db *gorm.DB) PhotoRepository {
	return &PhotoRepositoryImpl{
		Repository: Repository[model.Photo]{Db: Db},
	}
}
