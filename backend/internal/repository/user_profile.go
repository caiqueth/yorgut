package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type UserProfileRepository struct {
	repository[*model.UserProfile]
}

func NewUserProfileRepository(Db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{
		repository: repository[*model.UserProfile]{db: Db},
	}
}

// findOrCreateOfUser returns the user's profile row, creating an empty one if
// the user has none yet.
func (r *UserProfileRepository) findOrCreateOfUser(userId uint) (*model.UserProfile, error) {
	var profile model.UserProfile
	result := r.db.
		Where(model.UserProfile{UserID: userId}).
		FirstOrCreate(&profile)

	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

// UpdateDetailsOfUser overwrites every descriptive field of the user's
// profile with the values in details; nil fields are stored as NULL.
func (r *UserProfileRepository) UpdateDetailsOfUser(userId uint, details *model.UserProfile) error {
	profile, err := r.findOrCreateOfUser(userId)
	if err != nil {
		return err
	}

	result := r.db.
		Model(profile).
		Select(
			"description",
			"relationship_status",
			"interests",
			"about",
			"religion",
			"humor",
			"fashion",
			"lives_with",
			"birth_city",
			"website",
			"passions",
			"sports",
			"hobbies",
			"books",
			"music",
			"movies",
			"cuisine",
		).
		Updates(details)

	return result.Error
}

func (r *UserProfileRepository) UpdatePictureOfUser(userId uint, picturePath string) error {
	profile, err := r.findOrCreateOfUser(userId)
	if err != nil {
		return err
	}

	result := r.db.Model(profile).Update("profile_picture", picturePath)
	return result.Error
}
