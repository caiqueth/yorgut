package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityRepository struct {
	repository[*model.Community]
}

func NewCommunityRepository(Db *gorm.DB) *CommunityRepository {
	return &CommunityRepository{
		repository: repository[*model.Community]{db: Db},
	}
}

// FindShuffledExcept returns up to limit communities other than communityId,
// with profiles, shuffled by seed (see shuffledBy).
func (r *CommunityRepository) FindShuffledExcept(communityId uint, limit int, seed int64) ([]model.Community, error) {
	communities := []model.Community{}
	result := r.db.
		Scopes(shuffledBy("communities", seed)).
		Preload("Profile").
		Where("communities.id <> ?", communityId).
		Limit(limit).
		Find(&communities)
	if result.Error != nil {
		return nil, result.Error
	}
	return communities, nil
}

// UpdateMemberCount recounts the community's members and stores the result,
// in a single statement so concurrent joins and leaves can't lose updates.
func (r *CommunityRepository) UpdateMemberCount(communityId uint) error {
	members := r.db.
		Model(&model.CommunityMembership{}).
		Select("count(*)").
		Where("community_id = ?", communityId)

	result := r.db.
		Model(&model.Community{}).
		Where("id = ?", communityId).
		Update("member_count", members)

	return result.Error
}

// FindCreatedByUser returns one page of the communities the user created, by
// name, plus how many there are in total.
func (r *CommunityRepository) FindCreatedByUser(userId uint, page int, pageSize int) ([]model.Community, int64, error) {
	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("created_by_id = ?", userId)
	}

	var count int64
	if err := r.db.Model(&model.Community{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	communities := []model.Community{}
	result := r.db.
		Scopes(filter, newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Order("name, id").
		Find(&communities)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return communities, count, nil
}

// SearchNotMemberOf finds communities whose name contains name (case
// insensitive) and that userId hasn't joined, biggest first. It returns one
// page plus the total count of matches.
func (r *CommunityRepository) SearchNotMemberOf(userId uint, name string, page int, pageSize int) ([]model.Community, int64, error) {
	memberships := r.db.
		Model(&model.CommunityMembership{}).
		Select("community_id").
		Where("user_id = ?", userId)

	filter := func(db *gorm.DB) *gorm.DB {
		return db.
			Where("lower(name) like lower(?)", "%"+name+"%").
			Where("id not in (?)", memberships)
	}

	var count int64
	countResult := r.db.Model(&model.Community{}).Scopes(filter).Count(&count)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	var communities []model.Community
	result := r.db.
		Scopes(filter, newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Order("member_count desc, name").
		Find(&communities)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return communities, count, nil
}

// UpdateDetails sets the community's name, category and privacy.
func (r *CommunityRepository) UpdateDetails(communityId uint, name string, category string, private bool) error {
	result := r.db.
		Model(&model.Community{}).
		Where("id = ?", communityId).
		Updates(map[string]any{"name": name, "category": category, "private": private})

	return result.Error
}

// FindOwnedBy returns every community the user owns.
func (r *CommunityRepository) FindOwnedBy(userId uint) ([]model.Community, error) {
	communities := []model.Community{}
	result := r.db.Where("owner_id = ?", userId).Order("id").Find(&communities)
	if result.Error != nil {
		return nil, result.Error
	}
	return communities, nil
}

// TransferOwnership makes newOwnerId the community's owner. It also becomes
// its creator, the one "minhas comunidades" lists it for.
func (r *CommunityRepository) TransferOwnership(communityId uint, newOwnerId uint) error {
	return r.db.
		Model(&model.Community{}).
		Where("id = ?", communityId).
		Updates(map[string]any{"owner_id": newOwnerId, "created_by_id": newOwnerId}).
		Error
}

// DeleteForever removes the community and everything in it (forum, polls,
// members, requests, invitations, profile) from the database. The picture
// file is the caller's to delete.
func (r *CommunityRepository) DeleteForever(communityId uint) error {
	if err := deletePollsForever(r.db, "community_id = ?", communityId); err != nil {
		return err
	}
	var topicIds []uint
	if err := r.db.Unscoped().Model(&model.Topic{}).Where("community_id = ?", communityId).Pluck("id", &topicIds).Error; err != nil {
		return err
	}
	if err := deleteTopicsForever(r.db, topicIds); err != nil {
		return err
	}
	for _, value := range []any{
		&model.CommunityMembership{},
		&model.CommunityJoinRequest{},
		&model.CommunityInvitation{},
		&model.CommunityProfile{},
	} {
		if err := deleteForever(r.db, value, "community_id = ?", communityId); err != nil {
			return err
		}
	}
	return deleteForever(r.db, &model.Community{}, "id = ?", communityId)
}
