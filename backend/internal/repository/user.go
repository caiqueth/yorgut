package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[*model.User]
}

func NewUserRepository(Db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: Repository[*model.User]{Db: Db},
	}
}

func (u *UserRepository) FindUserOwnedCommunities(userId uint) ([]model.Community, error) {
	var communities []model.Community

	result := u.Db.Where("created_by_id = ?", userId).Find(&communities)

	if result.Error != nil {
		return nil, result.Error
	}

	return communities, nil
}
