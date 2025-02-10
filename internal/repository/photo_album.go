package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PhotoAlbumRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements PhotoAlbumRepository.
func (p *PhotoAlbumRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUser implements PhotoAlbumRepository.
func (p *PhotoAlbumRepositoryImpl) FindOfUser(userId int) {
	panic("unimplemented")
}

// Save implements PhotoAlbumRepository.
func (p *PhotoAlbumRepositoryImpl) Save(pa model.PhotoAlbum) {
	panic("unimplemented")
}

func NewPhotoAlbumRepositoryImpl(Db *gorm.DB) PhotoAlbumRepository {
	return &PhotoAlbumRepositoryImpl{Db: Db}
}
