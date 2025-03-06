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
	result := f.Db.Where("user_requested_id = ? or user_requestee_id = ?", userId, userId).Find(&ents)
	if result.Error != nil {
		return nil, result.Error
	}
	return ents, nil
}

func NewFriendshipRepository(Db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		Repository: Repository[*model.Friendship]{Db: Db},
	}
}
