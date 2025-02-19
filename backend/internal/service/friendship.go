package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type FriendshipService struct {
	FriendshipRepository repository.FriendshipRepository
	Validate             *validator.Validate
}

func (f *FriendshipService) Create(friendship request.CreateFriendshipRequest) {
	panic("unimplemented")
}

func (f *FriendshipService) Delete(id int) {
	panic("unimplemented")
}

func (f *FriendshipService) FindOfUser(userId int) response.FindFriendshipOfUserResponse {
	panic("unimplemented")
}

func NewFriendshipService(
	friendshipRepository repository.FriendshipRepository,
	validate *validator.Validate,
) *FriendshipService {
	return &FriendshipService{
		FriendshipRepository: friendshipRepository,
		Validate:             validate,
	}
}
