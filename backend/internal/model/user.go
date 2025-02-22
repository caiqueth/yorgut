package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string
	FirstName string
	LastName  string
	BirthDate time.Time
}

func NewUser(firstName string, lastName string, email string, birthDate time.Time) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		BirthDate: birthDate,
	}
}
