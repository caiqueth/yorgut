package response

import "time"

type UserDto struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FindUserByIdResponse struct {
	UserDto
}

type FindManyUsersResponse struct {
	Users []UserDto `json:"users"`
}
