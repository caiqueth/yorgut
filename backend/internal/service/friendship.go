package service

import (
	"errors"
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type FriendshipService struct {
	FriendshipRepository repository.FriendshipRepository
	UserRepository       repository.UserRepository
	Validate             *validator.Validate
}

func (f *FriendshipService) Create(friendship request.CreateFriendshipRequest) error {
	userRequested, err := f.UserRepository.FindById(friendship.UserRequested)
	if err != nil {
		return errors.New("could not find user that requested the friendship")
	}

	userRequestee, err := f.UserRepository.FindById(friendship.UserRequestee)
	if err != nil {
		return errors.New("could not find user that requestee the friendship")
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.FriendshipRepository.Create(friendshipModel)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) Delete(id uint) error {
	err := f.FriendshipRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) FindOfUser(userId uint) (response.FindFriendshipOfUserResponse, error) {
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
