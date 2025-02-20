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
}

func NewUser(name string, email string, birthDate time.Time) *User {
	return &User{
		Name:      name,
		Email:     email,
		BirthDate: birthDate,
	}
}
