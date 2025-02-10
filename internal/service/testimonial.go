package service

import (
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type TestimonialServiceImpl struct {
	TestimonialRepository repository.TestimonialRepository
	Validate              *validator.Validate
}

// Create implements TestimonialService.
func (t *TestimonialServiceImpl) Create(testimonial request.CreateTestimonialRequest) {
	panic("unimplemented")
}

// Delete implements TestimonialService.
func (t *TestimonialServiceImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUserAuthor implements TestimonialService.
func (t *TestimonialServiceImpl) FindOfUserAuthor(userId int) response.FindOfUserResponse {
	panic("unimplemented")
}

// FindOfUserRecipient implements TestimonialService.
func (t *TestimonialServiceImpl) FindOfUserRecipient(userId int) response.FindOfUserResponse {
	panic("unimplemented")
}

func NewTestimonialServiceImpl(testimonialRepository repository.TestimonialRepository, validate *validator.Validate) TestimonialService {
	return &TestimonialServiceImpl{
		TestimonialRepository: testimonialRepository,
		Validate:              validate,
	}
}
