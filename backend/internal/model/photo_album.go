package model

import "gorm.io/gorm"

type PhotoAlbum struct {
	gorm.Model
	UserID uint
	User   *User `gorm:"foreignKey:UserID;references:ID"`
	Name   string
}

func NewPhotoAlbum(user *User, name string) *PhotoAlbum {
	return &PhotoAlbum{
		UserID: user.ID,
		User:   user,
		Name:   name,
	}
}
