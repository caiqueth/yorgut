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

func (r *CommunityRepository) FindCommunitiesWithProfiles() ([]model.Community, error) {
	var ents []model.Community
	results := r.db.Preload("Profile").Find(&ents)
	if results.Error != nil {
		return nil, results.Error
	}
	return ents, nil
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

func (r *CommunityRepository) FindCreatedByUser(userId uint) ([]model.Community, error) {
	var communities []model.Community

	result := r.db.Preload("Profile").Where("created_by_id = ?", userId).Find(&communities)

	if result.Error != nil {
		return nil, result.Error
	}

	return communities, nil
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

// UpdateDetails sets the community's name and category.
func (r *CommunityRepository) UpdateDetails(communityId uint, name string, category string) error {
	result := r.db.
		Model(&model.Community{}).
		Where("id = ?", communityId).
		Updates(map[string]any{"name": name, "category": category})

	return result.Error
}
