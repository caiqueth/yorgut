package controller

import "yorgut/internal/service"

type TestimonialController struct {
	TestimonialService service.TestimonialService
}

func NewTestimonialController(service service.TestimonialService) *TestimonialController {
	return &TestimonialController{
		TestimonialService: service,
	}
}

func (t *TestimonialController) Create() {
	panic("unimplemented")
}

func (t *TestimonialController) Delete() {
	panic("unimplemented")
}

func (t *TestimonialController) FindOfUserAuthor() {
	panic("unimplemented")
}

func (t *TestimonialController) FindOfUserRecipient() {
	panic("unimplemented")
}
