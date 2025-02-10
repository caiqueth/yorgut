package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type ScrapRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements ScrapRepository.
func (s *ScrapRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUserRecipient implements ScrapRepository.
func (s *ScrapRepositoryImpl) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

// Save implements ScrapRepository.
func (*ScrapRepositoryImpl) Save(s model.Scrap) {
	panic("unimplemented")
}

func NewScrapRepositoryImpl(Db *gorm.DB) ScrapRepository {
	return &ScrapRepositoryImpl{Db: Db}
}
