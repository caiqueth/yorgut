package model

import (
	"time"

	"github.com/google/uuid"
)

type PhotoAlbum struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    int
	User      User `gorm:"foreignKey:UserID;references:ID"`
	Name      string
	CreatedAt time.Time
}
