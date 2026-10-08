package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type CommunityProfileRepository struct {
	repository[*model.CommunityProfile]
}

func NewCommunityProfileRepository(Db *gorm.DB) *CommunityProfileRepository {
	return &CommunityProfileRepository{
		repository: repository[*model.CommunityProfile]{db: Db},
	}
}

func (r *CommunityProfileRepository) FindOfCommunity(communityId uint) (*model.CommunityProfile, error) {
	var profile model.CommunityProfile
	result := r.db.Where("community_id = ?", communityId).First(&profile)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

// findOrCreateOfCommunity returns the community's profile row, creating an
// empty one if the community has none yet.
func (r *CommunityProfileRepository) findOrCreateOfCommunity(communityId uint) (*model.CommunityProfile, error) {
	var profile model.CommunityProfile
	result := r.db.
		Where(model.CommunityProfile{CommunityID: communityId}).
		FirstOrCreate(&profile)

	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

// UpdateDescriptionOfCommunity sets the description; nil stores NULL.
func (r *CommunityProfileRepository) UpdateDescriptionOfCommunity(communityId uint, description *string) error {
	profile, err := r.findOrCreateOfCommunity(communityId)
	if err != nil {
		return err
	}

	result := r.db.Model(profile).Update("description", description)
	return result.Error
}

func (r *CommunityProfileRepository) UpdatePictureOfCommunity(communityId uint, picturePath string) error {
	profile, err := r.findOrCreateOfCommunity(communityId)
	if err != nil {
		return err
	}

	result := r.db.Model(profile).Update("profile_picture", picturePath)
	return result.Error
}
