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

type CommunityDto struct {
	Name        string
	Description string
	CreatedBy   uint
	PictureUrl  string
}

type FindUserByIdResponse struct {
	UserDto
}

type FindManyUsersResponse struct {
	Users []UserDto `json:"users"`
}

type FindUserCommunitiesResponse struct {
	Communities []CommunityDto `json:"communities"`
}

type FindUserMembershipsResponse struct {
	Communities []CommunityDto `json:"communities"`
}
