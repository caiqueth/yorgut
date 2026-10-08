package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type ScrapRepository struct {
	repository[*model.Scrap]
}

func NewScrapRepository(Db *gorm.DB) *ScrapRepository {
	return &ScrapRepository{
		repository: repository[*model.Scrap]{db: Db},
	}
}

func (s *ScrapRepository) FindOfUserRecipient(userId uint) ([]model.Scrap, error) {
	var scraps []model.Scrap
	result := s.db.
		Where("user_recipient_id = ?", userId).
		Preload("UserAuthor").
		Preload("UserAuthor.Profile").
		Order("created_at desc").
		Find(&scraps)

	if result.Error != nil {
		return nil, result.Error
	}
	return scraps, nil
}

func (s *ScrapRepository) CountOfUserRecipient(userId uint) (int64, error) {
	var count int64
	result := s.db.
		Model(&model.Scrap{}).
		Where("user_recipient_id = ?", userId).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
