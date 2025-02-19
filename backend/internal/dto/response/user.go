package response

import "time"

type UserDto struct {
	Id        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
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
