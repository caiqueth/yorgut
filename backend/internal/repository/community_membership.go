package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityMembershipRepository struct {
	Repository[*model.CommunityMembership]
}

func (c *CommunityMembershipRepository) FindOfUser(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.Db.Preload("Community").Where("user_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func (c *CommunityMembershipRepository) FindOfCommunity(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.Db.Preload("Community").Where("community_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func NewCommunityMembershipRepository(Db *gorm.DB) *CommunityMembershipRepository {
	return &CommunityMembershipRepository{
		Repository: Repository[*model.CommunityMembership]{Db: Db},
	}
}
