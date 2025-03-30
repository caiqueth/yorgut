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

type UserProfile struct {
	gorm.Model
	UserID             uint
	User               *User `gorm:"foreignKey:UserID;references:ID"`
	ProfilePicture     *string
	Description        *string
	RelationshipStatus *string
	Interests          *string
	About              *string
	Religion           *string
	Humor              *string
	Fashion            *string
	LivesWith          *string
	BirthCity          *string
	Website            *string
	Passions           *string
	Sports             *string
	Hobbies            *string
	Books              *string
	Music              *string
	Movies             *string
	Cuisine            *string
}

func NewUser(firstName string, lastName string, email string, birthDate time.Time) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		BirthDate: birthDate,
	}
}

func NewUserProfile(user *User, opts ...func(*UserProfile)) *UserProfile {
	profile := &UserProfile{
		User:   user,
		UserID: user.ID,
	}

	for _, opt := range opts {
		opt(profile)
	}

	return profile
}
