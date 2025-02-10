package response

import "time"

type FindUserByIdResponse struct {
	Id        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birthDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
