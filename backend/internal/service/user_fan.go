package service

import (
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var ErrFanNotFriend = newError(ErrForbidden, "Só amigos podem virar fãs.")

type UserFanService struct {
	userFanRepository    repository.UserFanRepository
	friendshipRepository repository.FriendshipRepository
}

func NewUserFanService(
	userFanRepository repository.UserFanRepository,
	friendshipRepository repository.FriendshipRepository,
) *UserFanService {
	return &UserFanService{
		userFanRepository:    userFanRepository,
		friendshipRepository: friendshipRepository,
	}
}

// BecomeFan follows the review rule: only friends, never yourself.
func (u *UserFanService) BecomeFan(fanId uint, idolId uint) error {
	friends, err := areFriends(u.friendshipRepository, fanId, idolId)
	if err != nil {
		return err
	}
	if !friends {
		return ErrFanNotFriend
	}
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

	canFan, err := areFriends(u.friendshipRepository, viewerId, idolId)
	if err != nil {
		return nil, err
	}

	return &transport.UserFanSummaryResponse{
		FanCount: count,
		AmIFan:   amIFan,
		CanFan:   canFan,
	}, nil
}

// FindFansOfUser returns one page of the user's fans, the newest first.
func (u *UserFanService) FindFansOfUser(idolId uint, page int, pageSize int) (*transport.FindUserFanListResponse, error) {
	fans, count, err := u.userFanRepository.FindFansOfIdol(idolId, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindUserFanListResponse{}
	resp.FromFans(fans, page, pageSize, count)
	return resp, nil
}
