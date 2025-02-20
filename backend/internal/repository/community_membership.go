package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityMembershipRepository struct {
	Repository[*model.CommunityMembership]
}

func (c *CommunityMembershipRepository) FindOfUser(id int) {
	panic("unimplemented")
}

func NewCommunityMembershipRepository(Db *gorm.DB) *CommunityMembershipRepository {
	return &CommunityMembershipRepository{
		Repository: Repository[*model.CommunityMembership]{Db: Db},
	}
}
