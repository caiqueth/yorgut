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

// FindOfUserRecipient returns one page of the scraps sent to the user,
// newest first, plus how many they have.
func (s *ScrapRepository) FindOfUserRecipient(userId uint, page int, pageSize int) ([]model.Scrap, int64, error) {
	count, err := s.CountOfUserRecipient(userId)
	if err != nil {
		return nil, 0, err
	}

	scraps := []model.Scrap{}
	result := s.db.
		Scopes(newPaginate(page, pageSize).paginatedResult).
		Where("user_recipient_id = ?", userId).
		Preload("UserAuthor").
		Preload("UserAuthor.Profile").
		Order("created_at desc, id desc").
		Find(&scraps)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return scraps, count, nil
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

// DeleteForever removes the scrap from the database (not a soft delete):
// a deleted scrap is gone.
func (s *ScrapRepository) DeleteForever(scrapId uint) error {
	result := s.db.Unscoped().Where("id = ?", scrapId).Delete(&model.Scrap{})
	return result.Error
}
