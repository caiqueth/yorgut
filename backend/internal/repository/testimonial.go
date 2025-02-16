package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type TestimonialRepository interface {
	Create(t model.Testimonial) error
	Delete(id int)
	FindById(id int) (model.Testimonial, error)
	FindMany() []model.Testimonial
	FindOfUserRecipient(userId int)
	FindOfUserAuthor(userId int)
}

type TestimonialRepositoryImpl struct {
	Repository[model.Testimonial]
}

func (t *TestimonialRepositoryImpl) FindOfUserAuthor(userId int) {
	panic("unimplemented")
}

func (t *TestimonialRepositoryImpl) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

func NewTestimonialRepositoryImpl(Db *gorm.DB) TestimonialRepository {
	return &TestimonialRepositoryImpl{
		Repository: Repository[model.Testimonial]{Db: Db},
	}
}
