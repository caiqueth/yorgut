package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type TestimonialService struct {
	TestimonialRepository repository.TestimonialRepository
	Validate              *validator.Validate
}

func (t *TestimonialService) Create(testimonial request.CreateTestimonialRequest) {
	panic("unimplemented")
}

func (t *TestimonialService) Delete(id int) {
	panic("unimplemented")
}

func (t *TestimonialService) FindOfUserAuthor(userId int) response.FindOfUserResponse {
	panic("unimplemented")
}

func (t *TestimonialService) FindOfUserRecipient(userId int) response.FindOfUserResponse {
	panic("unimplemented")
}

func NewTestimonialService(testimonialRepository repository.TestimonialRepository, validate *validator.Validate) *TestimonialService {
	return &TestimonialService{
		TestimonialRepository: testimonialRepository,
		Validate:              validate,
	}
}
