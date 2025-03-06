package dto

type CreateFriendshipRequest struct {
	UserRequested uint
	UserRequestee uint
}

type FindFriendshipOfUserResponse struct{}
