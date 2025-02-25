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
	results := c.Db.Find(&memberships)
	if results.Error != nil {
		return nil, results.Error
	}
	return memberships, nil
}

func NewCommunityMembershipRepository(Db *gorm.DB) *CommunityMembershipRepository {
	return &CommunityMembershipRepository{
		Repository: Repository[*model.CommunityMembership]{Db: Db},
	}
}
