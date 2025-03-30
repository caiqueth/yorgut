package dto

import (
	"mime/multipart"
	"time"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type CreateCommunityMembershipRequest struct {
	CommunityID uint
	UserID      uint
}

type FindCommunityMembershipOfUserResponse struct {
}

type CreateCommunityRequest struct {
	Name        string                `form:"name" binding:"required" example:"My Community"`
	Description string                `form:"description" binding:"required" example:"This is a test community"`
	Picture     *multipart.FileHeader `form:"picture" swaggerignore:"true"`
	CreatedByID uint                  `form:"createdById" binding:"required" example:"1"`
}

type UpdateCommunityRequest struct {
	Name        string
	Description string
}

type FindCommunityByIdResponse struct {
	Name        string
	Description string
	CreatedBy   uint
	PictureUrl  string
}

type FindManyCommunitiesResponse struct{}

type FriendshipDto struct {
	UserRequestedID uint
	UserRequesteeID uint
	CreatedAt       time.Time
}

type CreateFriendshipRequest struct {
	UserRequested uint
	UserRequestee uint
}

type FindFriendshipOfUserResponse struct {
	Friendships []FriendshipDto `json:"friendships"`
}
type PhotoAlbumDto struct {
	UserID uint
	Name   string
}

type CreatePhotoAlbumRequest struct {
	UserID uint
	Name   string
}

type FindPhotoAlbumOfUserResponse struct {
	PhotoAlbums []PhotoAlbumDto `json:"photoAlbums"`
}

type CreatePhotoRequest struct {
	PhotoAlbumID uint
	URL          string
	Caption      string
}

type FindPhotoOfUserResponse struct{}

type FindPhotoOfAlbumResponse struct{}

type CreateScrapRequest struct {
	UserAuthorID    uint
	UserRecipientID uint
	Content         string
}

type FindScrapOfUserRecipientResponse struct{}

type CreateTestimonialRequest struct {
	UserAuthorID    uint
	UserRecipientID uint
	Content         string
}

type FindTestimonialOfUserResponse struct{}

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
