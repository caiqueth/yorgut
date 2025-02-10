package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityMembershipRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements CommunityMembershipRepository.
func (c *CommunityMembershipRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUser implements CommunityMembershipRepository.
func (c *CommunityMembershipRepositoryImpl) FindOfUser(id int) {
	panic("unimplemented")
}

// Save implements CommunityMembershipRepository.
func (c *CommunityMembershipRepositoryImpl) Save(cm model.CommunityMembership) {
	panic("unimplemented")
}

func NewCommunityMembershipRepositoryImpl(Db *gorm.DB) CommunityMembershipRepository {
	return &CommunityMembershipRepositoryImpl{Db: Db}
}
