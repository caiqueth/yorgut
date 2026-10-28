package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type TestimonialRepository struct {
	repository[*model.Testimonial]
}

func NewTestimonialRepository(Db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{
		repository: repository[*model.Testimonial]{db: Db},
	}
}

// FindOfUserRecipient returns one page of a user's testimonials (newest first)
// plus the total count. A page of 0 returns all of them.
func (t *TestimonialRepository) FindOfUserRecipient(userId uint, approved bool, page int, pageSize int) ([]model.Testimonial, int64, error) {
	var count int64
	countResult := t.db.
		Model(&model.Testimonial{}).
		Where("user_recipient_id = ? and approved = ?", userId, approved).
		Count(&count)

	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	query := t.db
	if page > 0 {
		query = query.Scopes(newPaginate(page, pageSize).paginatedResult)
	}

	var testimonials []model.Testimonial
	result := query.
		Where("user_recipient_id = ? and approved = ?", userId, approved).
		Preload("UserAuthor").
		Preload("UserAuthor.Profile").
		Order("created_at desc").
		Find(&testimonials)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return testimonials, count, nil
}

func (t *TestimonialRepository) Approve(testimonialId uint) error {
	result := t.db.
		Model(&model.Testimonial{}).
		Where("id = ?", testimonialId).
		Update("approved", true)

	return result.Error
}
