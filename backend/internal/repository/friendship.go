package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type FriendshipRepository struct {
	Repository[*model.Friendship]
}

func (f *FriendshipRepository) FindOfUser(userId int) {
	panic("unimplemented")
}

func NewFriendshipRepository(Db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		Repository: Repository[*model.Friendship]{Db: Db},
	}
}
