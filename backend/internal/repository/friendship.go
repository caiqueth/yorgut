package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type FriendshipRepository struct {
	Repository[*model.Friendship]
}

func (f *FriendshipRepository) FindOfUser(userId int) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.Where("user_requested_id = ? or user_requestee_id = ?", userId, userId).Find(&ents)
	if result.Error != nil {
		return nil, result.Error
	}
	return ents, nil
}

func (f *FriendshipRepository) AreFriends(user1 uint, user2 uint) bool {
	var count int64
	f.db.Where(
		"(user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)",
		user1, user2, user2, user1,
	).Count(&count)

	return count > 0
}

func NewFriendshipRepository(Db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		Repository: Repository[*model.Friendship]{db: Db},
	}
}
