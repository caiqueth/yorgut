package model

import (
	"time"

	"github.com/google/uuid"
)

type Scrap struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserAuthorID    int
	UserAuthor      User `gorm:"foreignKey:UserAuthorID;references:ID"`
	UserRecipientID int
	UserRecipient   User `gorm:"foreignKey:UserRecipientID;references:ID"`
	Content         string
	CreatedAt       time.Time
}
