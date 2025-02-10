package model

import (
	"time"

	"github.com/google/uuid"
)

type Community struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	CreateByID  uint
	CreatedBy   User `gorm:"foreignKey:CreateByID;references:ID"`
	CreatedAt   time.Time
}
