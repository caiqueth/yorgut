package service

import (
	"fmt"
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type TestimonialService struct {
	TestimonialRepository repository.TestimonialRepository
	UserRepository        repository.UserRepository
	Validate              *validator.Validate
}

func (t *TestimonialService) Create(testimonial request.CreateTestimonialRequest) error {
	author, err := t.UserRepository.FindById(testimonial.UserAuthorID)
	if err != nil {
		return fmt.Errorf("author not found %v", testimonial.UserAuthorID)
	}

	recipient, err := t.UserRepository.FindById(testimonial.UserRecipientID)
	if err != nil {
		return fmt.Errorf("recipient not found %v", testimonial.UserRecipientID)
	}

	testimonialModel := model.NewTestimonial(author, recipient, testimonial.Content)
	err = t.TestimonialRepository.Create(testimonialModel)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestimonialService) Delete(id uint) error {
	err := t.TestimonialRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TestimonialService) FindOfUserAuthor(userId uint) response.FindTestimonialOfUserResponse {
	panic("unimplemented")
}

func (t *TestimonialService) FindOfUserRecipient(userId uint) response.FindTestimonialOfUserResponse {
	panic("unimplemented")
}

func NewTestimonialService(testimonialRepository repository.TestimonialRepository, validate *validator.Validate) *TestimonialService {
	return &TestimonialService{
		TestimonialRepository: testimonialRepository,
		Validate:              validate,
	}
}
