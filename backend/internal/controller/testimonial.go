package controller

import (
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type TestimonialController struct {
	TestimonialService service.TestimonialService
}

func NewTestimonialController(service service.TestimonialService) *TestimonialController {
	return &TestimonialController{
		TestimonialService: service,
	}
}

// Create				godoc
// @Summary				Create a Testimonial.
// @Description			Creates a Testimonial.
// @Param				Testimonial body request.CreateTestimonialRequest true "Create Testimonial"
// @Produce				application/json
// @Tags				Testimonial
// @Success				204 {object} response.Response{}
// @Router				/testimonial [post]
func (t *TestimonialController) Create(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete 				godoc
// @Summary				Delete Testimonial by id.
// @Param				testimonialId path string true "delete Testimonial by id"
// @Description			Deletes a Testimonial.
// @Produce				application/json
// @Tags				Testimonial
// @Success				204 {object} response.Response{}
// @Router				/testimonial/{testimonialId} [delete]
func (t *TestimonialController) Delete(ctx *gin.Context) {
	panic("unimplemented")
}
