package repository

import (
	"errors"
	"yorgut/internal/dto/request"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityRepository interface {
	Save(c model.Community)
	Update(c model.Community)
	Delete(id int)
	FindById(id int) (model.Community, error)
	FindMany() []model.Community
}

type CommunityRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements CommunityRepository.
func (c *CommunityRepositoryImpl) Delete(id int) {
	var community model.Community
	result := c.Db.Where("id = ?", id).Delete(&community)
	if result.Error != nil {
		panic(result.Error)
	}
}

// FindById implements CommunityRepository.
func (c *CommunityRepositoryImpl) FindById(id int) (model.Community, error) {
	var community model.Community
	result := c.Db.Find(&community, id)
	if result != nil {
		return community, nil
	} else {
		return community, errors.New("community not found")
	}
}

// FindMany implements CommunityRepository.
func (c *CommunityRepositoryImpl) FindMany() []model.Community {
	var communities []model.Community
	results := c.Db.Find(&communities)
	if results.Error != nil {
		panic(results.Error)
	}
	return communities
}

// Save implements CommunityRepository.
func (c *CommunityRepositoryImpl) Save(community model.Community) {
	result := c.Db.Create(&community)
	if result.Error != nil {
		panic(result.Error)
	}
}

// Update implements CommunityRepository.
func (c *CommunityRepositoryImpl) Update(community model.Community) {
	updateCommunity := request.UpdateCommunityRequest{
		Name:        community.Name,
		Description: community.Description,
	}
	result := c.Db.Model(&community).Updates(updateCommunity)
	if result.Error != nil {
		panic(result.Error)
	}
}

func NewCommunityRepositoryImpl(Db *gorm.DB) CommunityRepository {
	return &CommunityRepositoryImpl{Db: Db}
}
