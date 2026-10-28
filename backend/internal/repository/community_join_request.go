package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityJoinRequestRepository struct {
	repository[*model.CommunityJoinRequest]
}

func NewCommunityJoinRequestRepository(Db *gorm.DB) *CommunityJoinRequestRepository {
	return &CommunityJoinRequestRepository{
		repository: repository[*model.CommunityJoinRequest]{db: Db},
	}
}

func (r *CommunityJoinRequestRepository) FindOfUserAndCommunity(communityId uint, userId uint) (*model.CommunityJoinRequest, error) {
	var request model.CommunityJoinRequest
	result := r.db.
		Where("user_id = ? and community_id = ?", userId, communityId).
		First(&request)
	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &request, nil
}

// FindOfCommunity returns the community's pending requests with their users,
// oldest first.
func (r *CommunityJoinRequestRepository) FindOfCommunity(communityId uint) ([]model.CommunityJoinRequest, error) {
	requests := []model.CommunityJoinRequest{}
	result := r.db.
		Preload("User").
		Preload("User.Profile").
		Where("community_id = ?", communityId).
		Order("created_at, id").
		Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

// FindCommunityIdsOfUser returns which of communityIds the user has a pending
// request to join.
func (r *CommunityJoinRequestRepository) FindCommunityIdsOfUser(userId uint, communityIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(communityIds) == 0 {
		return ids, nil
	}
	result := r.db.
		Model(&model.CommunityJoinRequest{}).
		Where("user_id = ? and community_id in ?", userId, communityIds).
		Pluck("community_id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}

// DeleteOfCommunity deletes all of the community's pending requests.
func (r *CommunityJoinRequestRepository) DeleteOfCommunity(communityId uint) error {
	result := r.db.Where("community_id = ?", communityId).Delete(&model.CommunityJoinRequest{})
	return result.Error
}

// DeleteOfUserAndCommunity deletes the user's pending request to join the
// community, if there is one.
func (r *CommunityJoinRequestRepository) DeleteOfUserAndCommunity(communityId uint, userId uint) error {
	result := r.db.Where("community_id = ? and user_id = ?", communityId, userId).Delete(&model.CommunityJoinRequest{})
	return result.Error
}
