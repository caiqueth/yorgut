package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoRepository struct {
	Repository[model.Photo]
}

func (p *PhotoRepository) FindOfPhotoAlbum(photoAlbumId int) {
	panic("unimplemented")
}

func (p *PhotoRepository) FindOfUser(userId int) {
	panic("unimplemented")
}

func NewPhotoRepository(Db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{
		Repository: Repository[model.Photo]{Db: Db},
	}
}
