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
