package dto

import (
	"time"
	"yorgut/internal/model"
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

// Mappers
func (r *FindUserByIdResponse) FromUser(user *model.User) {
	r.ID = user.ID
	r.Email = user.Email
	r.FirstName = user.FirstName
	r.LastName = user.LastName
	r.BirthDate = user.BirthDate
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
}

func (r *FindUserCommunitiesResponse) FromCommunities(communities []model.Community) {
	r.Communities = make([]CommunityDto, len(communities))
	for i, c := range communities {
		r.Communities[i] = CommunityDto{
			Name:        c.Name,
			Description: c.Description,
			CreatedBy:   c.CreatedByID,
			PictureUrl:  c.PictureUrl,
		}
	}
}

func (r *FindUserMembershipsResponse) FromMemberships(memberships []model.CommunityMembership) {
	r.Communities = make([]CommunityDto, len(memberships))
	for i, m := range memberships {
		r.Communities[i] = CommunityDto{
			Name:        m.Community.Name,
			Description: m.Community.Description,
			CreatedBy:   m.Community.CreatedByID,
			PictureUrl:  m.Community.PictureUrl,
		}
	}
}
