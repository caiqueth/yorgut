package request

import (
	"time"
)

type CreateUserRequest struct {
	Name      string
	Email     string
	BirthDate time.Time
}

type UpdateUserRequest struct {
	Name      string
	Email     string
	BirthDate time.Time
}
