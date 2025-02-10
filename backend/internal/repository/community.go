package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements CommunityRepository.
func (c *CommunityRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindById implements CommunityRepository.
func (*CommunityRepositoryImpl) FindById(id int) (c model.Community, err error) {
	panic("unimplemented")
}

// FindMany implements CommunityRepository.
func (c *CommunityRepositoryImpl) FindMany() []model.Community {
	panic("unimplemented")
}

// Save implements CommunityRepository.
func (*CommunityRepositoryImpl) Save(c model.Community) {
	panic("unimplemented")
}

// Update implements CommunityRepository.
func (*CommunityRepositoryImpl) Update(c model.Community) {
	panic("unimplemented")
}

func NewCommunityRepositoryImpl(Db *gorm.DB) CommunityRepository {
	return &CommunityRepositoryImpl{Db: Db}
}
