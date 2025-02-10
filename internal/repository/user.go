package repository

import (
	"errors"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func (u *UserRepositoryImpl) Delete(userId int) {
	var user model.User
	result := u.Db.Where("id = ?", userId).Delete(&user)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (u *UserRepositoryImpl) FindById(userId int) (model.User, error) {
	var user model.User
	result := u.Db.Find(&user, userId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}

func (u *UserRepositoryImpl) FindMany() []model.User {
	var users []model.User
	results := u.Db.Find(&users)
	if results.Error != nil {
		panic(results.Error)
	}
	return users
}

func (u *UserRepositoryImpl) Save(user model.User) {
	result := u.Db.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (u *UserRepositoryImpl) Update(user model.User) {
	panic("unimplemented")
}

func NewUserRepositoryImpl(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: Db}
}
