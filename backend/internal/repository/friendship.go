package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type FriendshipRepository interface {
	Save(f model.Friendship)
	Delete(id int)
	FindOfUser(userId int)
}

type FriendshipRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements FriendshipRepository.
func (f *FriendshipRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUser implements FriendshipRepository.
func (f *FriendshipRepositoryImpl) FindOfUser(userId int) {
	panic("unimplemented")
}

// Save implements FriendshipRepository.
func (*FriendshipRepositoryImpl) Save(f model.Friendship) {
	panic("unimplemented")
}

func NewFriendshipRepositoryImpl(Db *gorm.DB) FriendshipRepository {
	return &FriendshipRepositoryImpl{Db: Db}
}
