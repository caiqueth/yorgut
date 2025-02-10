package model

import (
	"time"

	"github.com/google/uuid"
)

type CommunityMembership struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CommunityID uuid.UUID
	Community   Community `gorm:"foreignKey:CommunityID;references:ID"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID;references:ID"`
	JoinedAt    time.Time
}
