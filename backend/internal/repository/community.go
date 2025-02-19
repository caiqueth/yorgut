package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityRepository struct {
	Repository[model.Community]
}

func NewCommunityRepository(Db *gorm.DB) *CommunityRepository {
	return &CommunityRepository{
		Repository: Repository[model.Community]{Db: Db},
	}
}
