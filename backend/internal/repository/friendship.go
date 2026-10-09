package repository

import (
	"time"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type FriendshipRepository struct {
	repository[*model.Friendship]
}

func NewFriendshipRepository(Db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		repository: repository[*model.Friendship]{db: Db},
	}
}

func (f *FriendshipRepository) FindByIdWithUsers(id uint) (*model.Friendship, error) {
	var ent model.Friendship
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		First(&ent, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &ent, nil
}

func (f *FriendshipRepository) FindPendingFriendshipsOfUser(userId uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequested.Profile").
		Where("(user_requestee_id = ?) and accepted = false", userId).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

func (f *FriendshipRepository) FindFriendshipsOfUser(userId uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Where("(user_requestee_id = ? or user_requested_id = ?)", userId, userId).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

func (f *FriendshipRepository) FindFriendshipsBetweenUsers(userId1 uint, userId2 uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Where("(user_requestee_id = ? and user_requested_id = ?) or (user_requestee_id = ? and user_requested_id = ?)", userId1, userId2, userId2, userId1).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

func (f *FriendshipRepository) FindFriendsOfUser(userId uint) ([]model.User, error) {
	var ents []model.Friendship
	friends := []model.User{}
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Preload("UserRequested.Profile").
		Preload("UserRequestee.Profile").
		Where("(user_requested_id = ? or user_requestee_id = ?) and accepted = true", userId, userId).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, fsp := range ents {
		if fsp.UserRequestedID != userId {
			friends = append(friends, *fsp.UserRequested)
		}
		if fsp.UserRequesteeID != userId {
			friends = append(friends, *fsp.UserRequestee)
		}
	}

	return friends, nil
}

func (f *FriendshipRepository) FindActiveAcceptedFriends(userId1 uint, userId2 uint) (*model.Friendship, error) {
	var ent model.Friendship
	result := f.db.
		Where("((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)) and accepted = true", userId1, userId2, userId2, userId1).
		First(&ent)

	if result.Error != nil {
		return nil, result.Error
	}

	return &ent, nil
}

func (f *FriendshipRepository) AcceptFriendship(friendshipId uint) error {
	now := time.Now()
	result := f.db.
		Model(&model.Friendship{}).
		Where("id = ?", friendshipId).
		Updates(&model.Friendship{Accepted: true, AcceptedAt: &now})

	return result.Error
}

func (f *FriendshipRepository) UsersAreFriends(user1 uint, user2 uint) bool {
	var count int64
	f.db.Model(&model.Friendship{}).Where(
		"((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)) and accepted = true",
		user1, user2, user2, user1,
	).Count(&count)

	return count > 0
}

func (f *FriendshipRepository) UsersHavePendingFriendshipInvites(user1 uint, user2 uint) bool {
	var count int64
	f.db.Model(&model.Friendship{}).Where(
		"((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)) and accepted = false",
		user1, user2, user2, user1,
	).Count(&count)

	return count > 0
}

func (f *FriendshipRepository) FindPendingFriendshipWithUsersById(friendshipId uint) (*model.Friendship, error) {
	var ent model.Friendship

	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Preload("UserRequested.Profile").
		Preload("UserRequestee.Profile").
		Where("id = ?", friendshipId).
		First(&ent)

	if result.Error != nil {
		return nil, result.Error
	}
	return &ent, nil
}
