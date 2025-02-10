package model

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	PhotoAlbumID uuid.UUID
	PhotoAlbum   PhotoAlbum `gorm:"foreignKey:PhotoAlbumID;references:ID"`
	URL          string
	Caption      string
	UploadedAt   time.Time
}
