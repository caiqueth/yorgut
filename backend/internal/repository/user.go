package repository

import (
	"time"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	repository[*model.User]
}

func NewUserRepository(Db *gorm.DB) *UserRepository {
	return &UserRepository{
		repository: repository[*model.User]{db: Db},
	}
}

func (u *UserRepository) FindByEmail(email string) (*model.User, error) {
	var ent model.User
	result := u.db.Where("email = ?", email).First(&ent)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return &ent, nil
	}
}

func (u *UserRepository) FindByGoogleSub(googleSub string) (*model.User, error) {
	var ent model.User
	result := u.db.Where("google_sub = ?", googleSub).First(&ent)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ent, nil
}

func (u *UserRepository) SetGoogleSub(userId uint, googleSub string) error {
	result := u.db.Model(&model.User{}).
		Where("id = ?", userId).
		Update("google_sub", googleSub)

	return result.Error
}

func (r *UserRepository) FindUserWithProfile(userId uint) (*model.User, error) {
	var ent model.User
	results := r.db.Preload("Profile").Where("id = ?", userId).First(&ent)
	if results.Error != nil {
		return nil, results.Error
	}
	return &ent, nil
}

// UpdateFriendCount recounts the user's accepted friendships and stores the
// result, in a single statement so concurrent changes can't lose updates.
func (r *UserRepository) UpdateFriendCount(userId uint) error {
	friends := r.db.
		Model(&model.Friendship{}).
		Select("count(*)").
		Where("(user_requested_id = ? or user_requestee_id = ?) and accepted = ?", userId, userId, true)

	result := r.db.
		Model(&model.User{}).
		Where("id = ?", userId).
		Update("friend_count", friends)

	return result.Error
}

func (r *UserRepository) SearchUsers(name string, page int, pageSize int) ([]*model.User, int64) {
	var ents []*model.User
	result := r.db.
		Scopes(newPaginate(page, pageSize).paginatedResult).
		Where("lower(concat(first_name, ' ', last_name)) like lower(?)", "%"+name+"%").
		Preload("Profile").
		Order("first_name").
		Find(&ents)

	if result.Error != nil {
		return []*model.User{}, 0
	}

	var count int64
	r.db.Where("lower(concat(first_name, ' ', last_name)) like lower(?)", "%"+name+"%").Count(&count)

	return ents, count
}

// UpdateBasicData sets the user's names and, when given, birth date. A map is
// used so empty values are written instead of skipped.
func (r *UserRepository) UpdateBasicData(userId uint, firstName string, lastName string, birthDate *time.Time) error {
	fields := map[string]any{
		"first_name": firstName,
		"last_name":  lastName,
	}
	if birthDate != nil {
		fields["birth_date"] = *birthDate
	}

	result := r.db.Model(&model.User{}).Where("id = ?", userId).Updates(fields)
	return result.Error
}
