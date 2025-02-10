package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type TestimonialRepositoryImpl struct {
	Db *gorm.DB
}

// Delete implements TestimonialRepository.
func (t *TestimonialRepositoryImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUserAuthor implements TestimonialRepository.
func (t *TestimonialRepositoryImpl) FindOfUserAuthor(userId int) {
	panic("unimplemented")
}

// FindOfUserRecipient implements TestimonialRepository.
func (t *TestimonialRepositoryImpl) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

// Save implements TestimonialRepository.
func (*TestimonialRepositoryImpl) Save(t model.Testimonial) {
	panic("unimplemented")
}

func NewTestimonialRepositoryImpl(Db *gorm.DB) TestimonialRepository {
	return &TestimonialRepositoryImpl{Db: Db}
}
