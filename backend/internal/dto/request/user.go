package request

import (
	"time"
)

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	BirthDate time.Time
}

type UpdateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	BirthDate time.Time
}
