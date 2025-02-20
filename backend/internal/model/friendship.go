package model

import (
	"time"

	"gorm.io/gorm"
)

type Friendship struct {
	gorm.Model
	UserRequestedID uint
	UserRequested   *User `gorm:"foreignKey:UserRequestedID;references:ID"`
	UserRequesteeID uint
	UserRequestee   *User `gorm:"foreignKey:UserRequesteeID;references:ID"`
	Accepted        bool
	AcceptedAt      *time.Time
}

func NewFriendship(userRequested *User, userRequestee *User) *Friendship {
	return &Friendship{
		UserRequestedID: userRequested.ID,
		UserRequested:   userRequested,
		UserRequesteeID: userRequestee.ID,
		UserRequestee:   userRequestee,
		Accepted:        false,
		AcceptedAt:      nil,
	}
}
