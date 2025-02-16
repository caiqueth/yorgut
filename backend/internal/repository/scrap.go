package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type ScrapRepository interface {
	Create(s model.Scrap) error
	Update(s model.Scrap)
	Delete(id int)
	FindById(id int) (model.Scrap, error)
	FindMany() []model.Scrap
	FindOfUserRecipient(userId int)
}

type ScrapRepositoryImpl struct {
	Repository[model.Scrap]
}

// FindOfUserRecipient implements ScrapRepository.
func (s *ScrapRepositoryImpl) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

func NewScrapRepositoryImpl(Db *gorm.DB) ScrapRepository {
	return &ScrapRepositoryImpl{
		Repository: Repository[model.Scrap]{Db: Db},
	}
}
