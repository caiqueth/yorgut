package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type ScrapRepository struct {
	Repository[model.Scrap]
}

func (s *ScrapRepository) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

func NewScrapRepository(Db *gorm.DB) *ScrapRepository {
	return &ScrapRepository{
		Repository: Repository[model.Scrap]{Db: Db},
	}
}
