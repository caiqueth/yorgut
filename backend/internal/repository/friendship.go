package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type FriendshipRepository interface {
	Create(f model.Friendship) error
	Delete(id int)
	FindById(id int) (model.Friendship, error)
	FindMany() []model.Friendship
	FindOfUser(userId int)
}

type FriendshipRepositoryImpl struct {
	Repository[model.Friendship]
}

func (f *FriendshipRepositoryImpl) FindOfUser(userId int) {
	panic("unimplemented")
}

func NewFriendshipRepositoryImpl(Db *gorm.DB) FriendshipRepository {
	return &FriendshipRepositoryImpl{
		Repository: Repository[model.Friendship]{Db: Db},
	}
}
