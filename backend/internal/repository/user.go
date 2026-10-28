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
		return nil, translate(result.Error)
	} else {
		return &ent, nil
	}
}

func (u *UserRepository) FindByGoogleSub(googleSub string) (*model.User, error) {
	var ent model.User
	result := u.db.Where("google_sub = ?", googleSub).First(&ent)
	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &ent, nil
}

func (u *UserRepository) SetGoogleSub(userId uint, googleSub string) error {
	result := u.db.Model(&model.User{}).
		Where("id = ?", userId).
		Update("google_sub", googleSub)

	return translate(result.Error)
}

func (r *UserRepository) FindUserWithProfile(userId uint) (*model.User, error) {
	var ent model.User
	results := r.db.Preload("Profile").Where("id = ?", userId).First(&ent)
	if results.Error != nil {
		return nil, translate(results.Error)
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

// TouchLastSeen sets the user's last_seen_at to now, unless it is already
// after staleBefore, so most requests don't write. UpdateColumn leaves
// updated_at alone: being seen doesn't change the user.
func (r *UserRepository) TouchLastSeen(userId uint, now time.Time, staleBefore time.Time) error {
	result := r.db.
		Model(&model.User{}).
		Where("id = ? and (last_seen_at is null or last_seen_at < ?)", userId, staleBefore).
		UpdateColumn("last_seen_at", now)

	return result.Error
}

// UpdatePrivacy stores the user's privacy settings (see model.User). A map,
// so false values are written instead of skipped.
func (r *UserRepository) UpdatePrivacy(userId uint, settings *model.User) error {
	result := r.db.Model(&model.User{}).Where("id = ?", userId).Updates(map[string]any{
		"record_visits":         settings.RecordVisits,
		"public_album":          settings.PublicAlbum,
		"public_communities":    settings.PublicCommunities,
		"public_scraps":         settings.PublicScraps,
		"public_testimonials":   settings.PublicTestimonials,
		"public_birthday":       settings.PublicBirthday,
		"hidden_profile_fields": settings.HiddenProfileFields,
	})
	return result.Error
}

// UpdateTheme stores the look the user picked for the app.
func (r *UserRepository) UpdateTheme(userId uint, theme string) error {
	result := r.db.Model(&model.User{}).Where("id = ?", userId).Update("theme", theme)
	return result.Error
}

// RecordsVisits reads the user's "registrar minhas visitas" setting.
func (r *UserRepository) RecordsVisits(userId uint) (bool, error) {
	var recordVisits bool
	result := r.db.Model(&model.User{}).Where("id = ?", userId).Pluck("record_visits", &recordVisits)
	if result.Error != nil {
		return false, result.Error
	}
	return recordVisits, nil
}

// SearchUsers returns one page of the users whose full name contains name
// (case insensitive), by first name, plus how many match.
func (r *UserRepository) SearchUsers(name string, page int, pageSize int) ([]*model.User, int64, error) {
	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("lower(concat(first_name, ' ', last_name)) like lower(?)", "%"+name+"%")
	}

	var count int64
	if err := r.db.Model(&model.User{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var ents []*model.User
	result := r.db.
		Scopes(filter, newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Order("first_name, id").
		Find(&ents)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return ents, count, nil
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

// DeleteForever removes the user and every row about them from the database:
// what they wrote (scraps, testimonials, topics with all their replies, polls),
// what they did (friendships, fans, ratings, visits, votes, memberships,
// requests, invitations), their photos and profile. Communities they own must
// have been handed over or deleted first. Picture files are the caller's to
// delete, and recounting other users' and communities' counters too.
func (r *UserRepository) DeleteForever(userId uint) error {
	var topicIds []uint
	if err := r.db.Unscoped().Model(&model.Topic{}).Where("author_id = ?", userId).Pluck("id", &topicIds).Error; err != nil {
		return err
	}
	if err := deleteTopicsForever(r.db, topicIds); err != nil {
		return err
	}
	if err := deletePollsForever(r.db, "author_id = ?", userId); err != nil {
		return err
	}

	albumIds := r.db.Unscoped().Model(&model.PhotoAlbum{}).Select("id").Where("user_id = ?", userId)
	steps := []struct {
		value any
		query string
		args  []any
	}{
		{&model.PollVote{}, "user_id = ?", []any{userId}},
		{&model.Scrap{}, "user_author_id = ? or user_recipient_id = ?", []any{userId, userId}},
		{&model.Testimonial{}, "user_author_id = ? or user_recipient_id = ?", []any{userId, userId}},
		{&model.Friendship{}, "user_requested_id = ? or user_requestee_id = ?", []any{userId, userId}},
		{&model.UserFan{}, "fan_id = ? or idol_id = ?", []any{userId, userId}},
		{&model.UserReview{}, "reviewer_id = ? or reviewed_id = ?", []any{userId, userId}},
		{&model.ProfileVisit{}, "visitor_id = ? or visited_id = ?", []any{userId, userId}},
		{&model.CommunityInvitation{}, "inviter_id = ? or invitee_id = ?", []any{userId, userId}},
		{&model.CommunityJoinRequest{}, "user_id = ?", []any{userId}},
		{&model.CommunityMembership{}, "user_id = ?", []any{userId}},
		{&model.Photo{}, "photo_album_id in (?)", []any{albumIds}},
		{&model.PhotoAlbum{}, "user_id = ?", []any{userId}},
		{&model.UserProfile{}, "user_id = ?", []any{userId}},
		{&model.User{}, "id = ?", []any{userId}},
	}
	for _, step := range steps {
		if err := deleteForever(r.db, step.value, step.query, step.args...); err != nil {
			return err
		}
	}
	return nil
}
