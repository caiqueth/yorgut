package service

import (
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

type UserFanService struct {
	userFanRepository    repository.UserFanRepository
	friendshipRepository repository.FriendshipRepository
	Validate             *validator.Validate
}

func NewUserFanService(
	userFanRepository repository.UserFanRepository,
	friendshipRepository repository.FriendshipRepository,
	validate *validator.Validate,
) *UserFanService {
	return &UserFanService{
		userFanRepository:    userFanRepository,
		friendshipRepository: friendshipRepository,
		Validate:             validate,
	}
}

// CanBecomeFan follows the review rule: only friends, never yourself.
func (u *UserFanService) CanBecomeFan(fanId uint, idolId uint) bool {
	if fanId == idolId {
		return false
	}
	return u.friendshipRepository.UsersAreFriends(fanId, idolId)
}

func (u *UserFanService) BecomeFan(fanId uint, idolId uint) error {
	return u.userFanRepository.Add(model.NewUserFan(fanId, idolId))
}

// StopBeingFan is always allowed, so a star can be removed after unfriending.
func (u *UserFanService) StopBeingFan(fanId uint, idolId uint) error {
	return u.userFanRepository.Remove(fanId, idolId)
}

func (u *UserFanService) CountFansOfUser(idolId uint) (uint, error) {
	count, err := u.userFanRepository.CountOfIdol(idolId)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (u *UserFanService) FindSummaryOfUser(viewerId uint, idolId uint) (*transport.UserFanSummaryResponse, error) {
	count, err := u.userFanRepository.CountOfIdol(idolId)
	if err != nil {
		return nil, err
	}

	amIFan, err := u.userFanRepository.IsFanOf(viewerId, idolId)
	if err != nil {
		return nil, err
	}

	return &transport.UserFanSummaryResponse{
		FanCount: count,
		AmIFan:   amIFan,
		CanFan:   u.CanBecomeFan(viewerId, idolId),
	}, nil
}
