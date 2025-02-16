package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityRepository interface {
	Create(c model.Community) error
	Update(c model.Community)
	Delete(id int)
	FindById(id int) (model.Community, error)
	FindMany() []model.Community
}

type CommunityRepositoryImpl struct {
	Repository[model.Community]
}

func NewCommunityRepositoryImpl(Db *gorm.DB) CommunityRepository {
	return &CommunityRepositoryImpl{
		Repository: Repository[model.Community]{Db: Db},
	}
}
