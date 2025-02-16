package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityMembershipRepository interface {
	Create(cm model.CommunityMembership) error
	Delete(id int)
	FindById(id int) (model.CommunityMembership, error)
	FindMany() []model.CommunityMembership
	FindOfUser(id int)
}

type CommunityMembershipRepositoryImpl struct {
	Repository[model.CommunityMembership]
}

func (c *CommunityMembershipRepositoryImpl) FindOfUser(id int) {
	panic("unimplemented")
}

func NewCommunityMembershipRepositoryImpl(Db *gorm.DB) CommunityMembershipRepository {
	return &CommunityMembershipRepositoryImpl{
		Repository: Repository[model.CommunityMembership]{Db: Db},
	}
}
