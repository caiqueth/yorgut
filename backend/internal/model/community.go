package model

import (
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name        string
	Description string
	CreatedByID uint
	PictureUrl  string
	CreatedBy   *User `gorm:"foreignKey:CreatedByID;references:ID"`
}

func NewCommunity(name string, description string, pictureUrl string, createdBy *User) *Community {
	return &Community{
		Name:        name,
		Description: description,
		PictureUrl:  pictureUrl,
		CreatedBy:   createdBy,
	}
}
