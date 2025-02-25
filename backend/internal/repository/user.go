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

	err := u.Db.Joins("JOIN communities ON communities.created_by_id = user.id").
		Where("communities.created_by_id = ?", userId).
		Find(&communities).
		Error

	if err != nil {
		return nil, err
	}

	return communities, nil
}
