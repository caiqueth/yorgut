package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoRepository interface {
	Save(p model.Photo)
	Delete(id int)
	FindOfUser(userId int)
	FindOfPhotoAlbum(photoAlbumId int)
}

type PhotoRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements PhotoRepository.
func (p *PhotoRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfPhotoAlbum implements PhotoRepository.
func (p *PhotoRepositoryImpl) FindOfPhotoAlbum(photoAlbumId int) {
	panic("unimplemented")
}

// FindOfUser implements PhotoRepository.
func (p *PhotoRepositoryImpl) FindOfUser(userId int) {
	panic("unimplemented")
}

// Save implements PhotoRepository.
func (*PhotoRepositoryImpl) Save(p model.Photo) {
	panic("unimplemented")
}

func NewPhotoRepositoryImpl(Db *gorm.DB) PhotoRepository {
	return &PhotoRepositoryImpl{Db: Db}
}
