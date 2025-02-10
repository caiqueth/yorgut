package service

import (
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type FriendshipServiceImpl struct {
	FriendshipRepository repository.FriendshipRepository
	Validate             *validator.Validate
}

// Create implements FriendshipService.
func (f *FriendshipServiceImpl) Create(friendship request.CreateFriendshipRequest) {
	panic("unimplemented")
}

// Delete implements FriendshipService.
func (f *FriendshipServiceImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUser implements FriendshipService.
func (f *FriendshipServiceImpl) FindOfUser(userId int) response.FindFriendshipOfUserResponse {
	panic("unimplemented")
}

func NewFriendshipServiceImpl(friendshipRepository repository.FriendshipRepository, validate *validator.Validate) FriendshipService {
	return &FriendshipServiceImpl{
		FriendshipRepository: friendshipRepository,
		Validate:             validate,
	}
}
