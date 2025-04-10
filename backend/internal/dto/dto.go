package dto

import (
	"mime/multipart"
	"time"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type CreateCommunityMembershipRequest struct {
	CommunityID uint `json:"communityId"`
	UserID      uint `json:"userId"`
}

type CommunityMembershipResponse struct {
	ID uint `json:"id"`
}

type FindCommunityMembershipOfUserResponse struct {
}

type CreateCommunityRequest struct {
	Name        string `form:"name" binding:"required" example:"My Community"`
	CreatedByID uint   `form:"createdById" binding:"required" example:"1"`
}

type CreateCommunityProfileRequest struct {
	CommunityID uint                  `json:"communityId"`
	Description string                `form:"description" binding:"required" example:"This is a test community"`
	Picture     *multipart.FileHeader `form:"picture" swaggerignore:"true"`
}

type UpdateCommunityRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FindCommunityByIdResponse struct {
	Name           string  `json:"name"`
	CreatedBy      uint    `json:"createdBy"`
	ProfilePicture *string `json:"profilePicture"`
}

type FindManyCommunitiesResponse struct{}

type FriendshipDto struct {
	UserRequestedID uint      `json:"userRequestedId"`
	UserRequesteeID uint      `json:"userRequesteeId"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CreateFriendshipRequest struct {
	UserRequested uint `json:"userRequested"`
	UserRequestee uint `json:"userRequestee"`
}

type FindFriendshipOfUserResponse struct {
	Friendships []FriendshipDto `json:"friendships"`
}

type FindFriendsOfUserResponse struct {
	Friends []FriendDto `json:"friends"`
}

type PhotoAlbumDto struct {
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
}

type CreatePhotoAlbumRequest struct {
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
}

type FindPhotoAlbumOfUserResponse struct {
	PhotoAlbums []PhotoAlbumDto `json:"photoAlbums"`
}

type CreatePhotoRequest struct {
	PhotoAlbumID uint   `json:"photoAlbumId"`
	URL          string `json:"uRL"`
	Caption      string `json:"caption"`
}

type FindPhotoOfUserResponse struct{}

type FindPhotoOfAlbumResponse struct{}

type CreateScrapRequest struct {
	UserAuthorID    uint   `json:"userAuthorId"`
	UserRecipientID uint   `json:"userRecipientId"`
	Content         string `json:"content"`
}

type FindScrapOfUserRecipientResponse struct{}

type CreateTestimonialRequest struct {
	UserAuthorID    uint   `json:"userAuthorId"`
	UserRecipientID uint   `json:"userRecipientId"`
	Content         string `json:"content"`
}

type FindTestimonialOfUserResponse struct{}

type CreateUserRequest struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
}

type UpdateUserRequest struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
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

type FriendDto struct {
	ID             uint
	FirstName      string    `json:"firstName"`
	BirthDate      time.Time `json:"birthDate"`
	ProfilePicture string    `json:"profilePicture"`
}

type CommunityDto struct {
	Name           string  `json:"name"`
	CreatedBy      uint    `json:"firstName"`
	ProfilePicture *string `json:"profilePicture"`
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
