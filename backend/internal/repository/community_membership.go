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

// FindCommunitiesOfUser returns one page of the communities the user is a
// member of, plus how many there are in total. They come by name, or shuffled
// by shuffleSeed when it isn't nil (see shuffledBy).
func (c *CommunityMembershipRepository) FindCommunitiesOfUser(userId uint, page int, pageSize int, shuffleSeed *int64) ([]model.Community, int64, error) {
	communityIds := c.db.
		Model(&model.CommunityMembership{}).
		Select("community_id").
		Where("user_id = ?", userId)

	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("communities.id in (?)", communityIds)
	}

	var count int64
	if err := c.db.Model(&model.Community{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	order := func(db *gorm.DB) *gorm.DB { return db.Order("communities.name, communities.id") }
	if shuffleSeed != nil {
		order = shuffledBy("communities", *shuffleSeed)
	}

	communities := []model.Community{}
	result := c.db.
		Scopes(filter, order, newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Find(&communities)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return communities, count, nil
}

func (c *CommunityMembershipRepository) CountOfUser(userId uint) (int64, error) {
	var count int64
	result := c.db.Model(&model.CommunityMembership{}).Where("user_id = ?", userId).Count(&count)
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

// FindMembersOfCommunity returns one page of the community's members, most
// recently seen first, plus how many members it has in total.
func (c *CommunityMembershipRepository) FindMembersOfCommunity(communityId uint, page int, pageSize int) ([]model.User, int64, error) {
	memberIds := c.db.
		Model(&model.CommunityMembership{}).
		Select("user_id").
		Where("community_id = ?", communityId)

	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("users.id in (?)", memberIds)
	}

	var count int64
	if err := c.db.Model(&model.User{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	members := []model.User{}
	result := c.db.
		Scopes(filter, recentlySeenFirst, newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Find(&members)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return members, count, nil
}

// FindMemberIdsOfCommunity returns which of userIds are members of the
// community.
func (c *CommunityMembershipRepository) FindMemberIdsOfCommunity(communityId uint, userIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(userIds) == 0 {
		return ids, nil
	}
	result := c.db.
		Model(&model.CommunityMembership{}).
		Where("community_id = ? and user_id in ?", communityId, userIds).
		Pluck("user_id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}

// FindOldestMemberExcept returns the id of the community's longest-standing
// member other than userId, or ErrNotFound when there is none.
func (c *CommunityMembershipRepository) FindOldestMemberExcept(communityId uint, userId uint) (uint, error) {
	var membership model.CommunityMembership
	result := c.db.
		Where("community_id = ? and user_id <> ?", communityId, userId).
		Order("created_at, id").
		First(&membership)
	if result.Error != nil {
		return 0, translate(result.Error)
	}
	return membership.UserID, nil
}

// FindCommunityIdsOfUser returns the ids of every community the user is a
// member of.
func (c *CommunityMembershipRepository) FindCommunityIdsOfUser(userId uint) ([]uint, error) {
	ids := []uint{}
	result := c.db.Model(&model.CommunityMembership{}).Where("user_id = ?", userId).Pluck("community_id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}
