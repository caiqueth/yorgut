package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u model.User) error
	Update(u model.User)
	Delete(id int)
	FindById(id int) (model.User, error)
	FindMany() []model.User
}

type UserRepositoryImpl struct {
	Repository[model.User]
}

func NewUserRepositoryImpl(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		Repository: Repository[model.User]{Db: Db},
	}
}
