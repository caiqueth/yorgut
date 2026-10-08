package service

import (
	"strings"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserProfileService struct {
	userProfileRepository repository.UserProfileRepository
	userRepository        repository.UserRepository
	Validate              *validator.Validate
}

func NewUserProfileService(
	userProfileRepository repository.UserProfileRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *UserProfileService {
	return &UserProfileService{
		userProfileRepository: userProfileRepository,
		userRepository:        userRepository,
		Validate:              validate,
	}
}

// UpdateUserProfileInput holds every descriptive profile field. A blank value
// clears the field.
type UpdateUserProfileInput struct {
	Description        string
	RelationshipStatus string
	Interests          string
	About              string
	Religion           string
	Humor              string
	Fashion            string
	LivesWith          string
	BirthCity          string
	Website            string
	Passions           string
	Sports             string
	Hobbies            string
	Books              string
	Music              string
	Movies             string
	Cuisine            string
}

func nilIfBlank(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func (u *UserProfileService) UpdateProfile(userId uint, input UpdateUserProfileInput) error {
	details := &model.UserProfile{
		Description:        nilIfBlank(input.Description),
		RelationshipStatus: nilIfBlank(input.RelationshipStatus),
		Interests:          nilIfBlank(input.Interests),
		About:              nilIfBlank(input.About),
		Religion:           nilIfBlank(input.Religion),
		Humor:              nilIfBlank(input.Humor),
		Fashion:            nilIfBlank(input.Fashion),
		LivesWith:          nilIfBlank(input.LivesWith),
		BirthCity:          nilIfBlank(input.BirthCity),
		Website:            nilIfBlank(input.Website),
		Passions:           nilIfBlank(input.Passions),
		Sports:             nilIfBlank(input.Sports),
		Hobbies:            nilIfBlank(input.Hobbies),
		Books:              nilIfBlank(input.Books),
		Music:              nilIfBlank(input.Music),
		Movies:             nilIfBlank(input.Movies),
		Cuisine:            nilIfBlank(input.Cuisine),
	}

	return u.userProfileRepository.UpdateDetailsOfUser(userId, details)
}

// UpdateProfilePicture follows the same limits as album photos (picture.go).
func (u *UserProfileService) UpdateProfilePicture(userId uint, picture []byte) error {
	if err := validatePicture(picture); err != nil {
		return err
	}

	storedFileName, err := helpers.UploadPicture(picture)
	if err != nil {
		return err
	}

	return u.userProfileRepository.UpdatePictureOfUser(userId, *storedFileName)
}
