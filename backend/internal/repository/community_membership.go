package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityMembershipRepository struct {
	repository[*model.CommunityMembership]
}

func NewCommunityMembershipRepository(Db *gorm.DB) *CommunityMembershipRepository {
	return &CommunityMembershipRepository{
		repository: repository[*model.CommunityMembership]{db: Db},
	}
}

func (c *CommunityMembershipRepository) FindOfUser(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.db.Preload("Community").Preload("Community.Profile").Where("user_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func (c *CommunityMembershipRepository) CountOfUser(id uint) (int64, error) {
	var count int64
	result := c.db.Model(&model.CommunityMembership{}).Where("user_id = ?", id).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (c *CommunityMembershipRepository) FindOfUserAndCommunity(communityId uint, userId uint) (*model.CommunityMembership, error) {
	var membership model.CommunityMembership
	result := c.db.
		Where("user_id = ? and community_id = ?", userId, communityId).
		First(&membership)
	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &membership, nil
}

func (c *CommunityMembershipRepository) FindOfCommunity(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.db.Preload("Community").Where("community_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func (c *CommunityMembershipRepository) FindMembersOfCommunity(communityId uint) ([]model.CommunityMembership, error) {
	var members []model.CommunityMembership
	result := c.db.Preload("User").Preload("User.Profile").Where("community_id = ?", communityId).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}
