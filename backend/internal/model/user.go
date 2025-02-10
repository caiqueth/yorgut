package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string
	Name      string
	BirthDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
