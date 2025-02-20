package model

import (
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	PhotoAlbumID uint
	PhotoAlbum   *PhotoAlbum `gorm:"foreignKey:PhotoAlbumID;references:ID"`
	URL          string
	Caption      string
}

func NewPhoto(photoAlbum *PhotoAlbum, url string, caption string) *Photo {
	return &Photo{
		PhotoAlbumID: photoAlbum.ID,
		PhotoAlbum:   photoAlbum,
		URL:          url,
		Caption:      caption,
	}
}
