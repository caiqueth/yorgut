package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type TestimonialRepository struct {
	Repository[model.Testimonial]
}

func (t *TestimonialRepository) FindOfUserAuthor(userId int) {
	panic("unimplemented")
}

func (t *TestimonialRepository) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

func NewTestimonialRepository(Db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{
		Repository: Repository[model.Testimonial]{Db: Db},
	}
}
