package model

import (
	"time"

	"github.com/google/uuid"
)

type Friendship struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserRequestedID uint
	UserRequested   User `gorm:"foreignKey:UserRequestedID;references:ID"`
	UserRequesteeID uint
	UserRequestee   User `gorm:"foreignKey:UserRequesteeID;references:ID"`
	IsAccepted      bool
	AcceptedAt      time.Time
	CreatedAt       time.Time
}
