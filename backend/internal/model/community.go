package model

import (
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name        string
	Description string
	CreatedByID uint
	CreatedBy   *User `gorm:"foreignKey:CreatedByID;references:ID"`
}

func NewCommunity(name string, description string, createdBy *User) *Community {
	return &Community{
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
	}
}
