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

// betweenUsers keeps the friendships (accepted or pending) between two users,
// whichever of them asked.
func betweenUsers(userId1 uint, userId2 uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(
			"((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?))",
			userId1, userId2, userId2, userId1,
		)
	}
}

func (f *FriendshipRepository) FindFriendshipsBetweenUsers(userId1 uint, userId2 uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Scopes(betweenUsers(userId1, userId2)).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

// friendIdsOf is a subquery selecting the ids of the user's accepted friends.
func (f *FriendshipRepository) friendIdsOf(userId uint) *gorm.DB {
	return f.db.
		Model(&model.Friendship{}).
		Select("case when user_requested_id = ? then user_requestee_id else user_requested_id end", userId).
		Where("(user_requested_id = ? or user_requestee_id = ?) and accepted = ?", userId, userId, true)
}

// CountRatedFriends counts the user's current friends that the user has
// rated (any trait).
func (f *FriendshipRepository) CountRatedFriends(userId uint) (int64, error) {
	var count int64
	result := f.db.
		Model(&model.UserReview{}).
		Where("reviewer_id = ? and reviewed_id in (?)", userId, f.friendIdsOf(userId)).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// FindAllFriendsOfUser returns all of the user's accepted friends, with
// profiles. For computations over every friend (birthdays); lists page with
// FindFriendsOfUser.
func (f *FriendshipRepository) FindAllFriendsOfUser(userId uint) ([]model.User, error) {
	friends := []model.User{}
	result := f.db.
		Where("users.id in (?)", f.friendIdsOf(userId)).
		Preload("Profile").
		Find(&friends)
	if result.Error != nil {
		return nil, result.Error
	}
	return friends, nil
}

// CountFriendsOfFriends counts the users who are friends of the user's
// friends, without the user and their own friends.
func (f *FriendshipRepository) CountFriendsOfFriends(userId uint) (int64, error) {
	friendIds := f.friendIdsOf(userId)

	// For each accepted friendship touching one of the friends, the other side.
	// A friendship between two friends yields a friend, which is excluded below.
	friendsOfFriends := f.db.
		Model(&model.Friendship{}).
		Select("case when user_requested_id in (?) then user_requestee_id else user_requested_id end", friendIds).
		Where("accepted = ? and (user_requested_id in (?) or user_requestee_id in (?))", true, friendIds, friendIds)

	var count int64
	result := f.db.
		Model(&model.User{}).
		Where("users.id in (?)", friendsOfFriends).
		Where("users.id not in (?)", friendIds).
		Where("users.id <> ?", userId).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// FindFriendsOfUser returns one page of the user's accepted friends, most
// recently seen first, plus how many friends they have in total.
func (f *FriendshipRepository) FindFriendsOfUser(userId uint, page int, pageSize int) ([]model.User, int64, error) {
	friendIds := f.friendIdsOf(userId)

	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("users.id in (?)", friendIds)
	}

	var count int64
	if err := f.db.Model(&model.User{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	friends := []model.User{}
	result := f.db.
		Scopes(filter, recentlySeenFirst, newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Find(&friends)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return friends, count, nil
}

// FindAcceptedBetweenUsers returns the accepted friendship between two users.
func (f *FriendshipRepository) FindAcceptedBetweenUsers(userId1 uint, userId2 uint) (*model.Friendship, error) {
	var ent model.Friendship
	result := f.db.
		Scopes(betweenUsers(userId1, userId2)).
		Where("accepted = ?", true).
		First(&ent)

	if result.Error != nil {
		return nil, translate(result.Error)
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

// UsersAreFriends tells whether two users have an accepted friendship.
func (f *FriendshipRepository) UsersAreFriends(userId1 uint, userId2 uint) (bool, error) {
	return f.existsBetweenUsers(userId1, userId2, true)
}

// HasPendingBetweenUsers tells whether either user asked the other to be
// friends and is still waiting.
func (f *FriendshipRepository) HasPendingBetweenUsers(userId1 uint, userId2 uint) (bool, error) {
	return f.existsBetweenUsers(userId1, userId2, false)
}

func (f *FriendshipRepository) existsBetweenUsers(userId1 uint, userId2 uint, accepted bool) (bool, error) {
	var count int64
	result := f.db.
		Model(&model.Friendship{}).
		Scopes(betweenUsers(userId1, userId2)).
		Where("accepted = ?", accepted).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// FindByIdWithUsers returns the friendship with both users and their profiles.
func (f *FriendshipRepository) FindByIdWithUsers(friendshipId uint) (*model.Friendship, error) {
	var ent model.Friendship

	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Preload("UserRequested.Profile").
		Preload("UserRequestee.Profile").
		Where("id = ?", friendshipId).
		First(&ent)

	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &ent, nil
}
