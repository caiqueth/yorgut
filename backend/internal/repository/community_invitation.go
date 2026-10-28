package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityInvitationRepository struct {
	repository[*model.CommunityInvitation]
}

func NewCommunityInvitationRepository(Db *gorm.DB) *CommunityInvitationRepository {
	return &CommunityInvitationRepository{
		repository: repository[*model.CommunityInvitation]{db: Db},
	}
}

// FindByIdWithCommunity returns the invitation with its community.
func (r *CommunityInvitationRepository) FindByIdWithCommunity(id uint) (*model.CommunityInvitation, error) {
	var invitation model.CommunityInvitation
	result := r.db.Preload("Community").First(&invitation, id)
	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &invitation, nil
}

// FindOfInvitee returns the user's pending invitations with their
// communities and inviters, newest first.
func (r *CommunityInvitationRepository) FindOfInvitee(inviteeId uint) ([]model.CommunityInvitation, error) {
	invitations := []model.CommunityInvitation{}
	result := r.db.
		Preload("Community").
		Preload("Community.Profile").
		Preload("Inviter").
		Preload("Inviter.Profile").
		Where("invitee_id = ?", inviteeId).
		Order("created_at desc, id desc").
		Find(&invitations)
	if result.Error != nil {
		return nil, result.Error
	}
	return invitations, nil
}

// FindInviteeIdsOfCommunity returns which of userIds have a pending
// invitation to the community.
func (r *CommunityInvitationRepository) FindInviteeIdsOfCommunity(communityId uint, userIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(userIds) == 0 {
		return ids, nil
	}
	result := r.db.
		Model(&model.CommunityInvitation{}).
		Where("community_id = ? and invitee_id in ?", communityId, userIds).
		Pluck("invitee_id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}

// DeleteOfInviteeAndCommunity deletes the user's pending invitation to the
// community, if there is one.
func (r *CommunityInvitationRepository) DeleteOfInviteeAndCommunity(communityId uint, inviteeId uint) error {
	result := r.db.Where("community_id = ? and invitee_id = ?", communityId, inviteeId).Delete(&model.CommunityInvitation{})
	return result.Error
}
