package controller

import (
	"errors"
	"net/http"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type TestimonialController struct {
	testimonialService service.TestimonialService
}

func NewTestimonialController(service service.TestimonialService) *TestimonialController {
	return &TestimonialController{
		testimonialService: service,
	}
}

// testimonialErrorStatus maps the service's sentinel errors to HTTP statuses.
func testimonialErrorStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrTestimonialNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrTestimonialForbidden), errors.Is(err, service.ErrTestimonialNotFriend):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// FindOfUser			godoc
// @Summary				Gets the testimonials of a user.
// @Param				userId path string true "User ID"
// @Param				page query int false "Page of approved testimonials, starting at 1 (omit for all)"
// @Param				pageSize query int false "Approved testimonials per page (default 5, max 50)"
// @Description			Returns the approved testimonials of a user. When the logged user is the recipient, pending ones are included too.
// @Produce				application/json
// @Tags				Testimonial
// @Success				200 {object} transport.FindTestimonialsOfUserResponse{}
// @Router				/testimonials/user/{userId} [get]
func (t *TestimonialController) FindOfUser(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId := helpers.MustParseUint(ctx.Param("userId"))
	page := helpers.MustParseInt(ctx.Query("page"))
	pageSize := helpers.MustParseInt(ctx.Query("pageSize"))
	if page < 0 {
		page = 0
	}
	if page == 0 {
		pageSize = 0
	} else if pageSize <= 0 {
		pageSize = 5
	} else if pageSize > 50 {
		pageSize = 50
	}

	testimonials, err := t.testimonialService.FindTestimonialsOfUser(viewer.ID, userId, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, testimonials)
}

// Create				godoc
// @Summary				Create a Testimonial.
// @Description			Logged user writes a testimonial for a friend. It stays pending until the recipient approves it.
// @Param				Testimonial body transport.CreateTestimonialRequest true "Create Testimonial"
// @Produce				application/json
// @Tags				Testimonial
// @Success				201 {object} transport.Response{}
// @Router				/testimonials [post]
func (t *TestimonialController) Create(ctx *gin.Context) {
	author := ctx.MustGet("user").(*model.User)

	request := transport.CreateTestimonialRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := t.testimonialService.CreateTestimonial(author.ID, request.UserRecipientID, request.Content); err != nil {
		ctx.JSON(testimonialErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transport.Response{})
}

// Approve				godoc
// @Summary				Approve a Testimonial.
// @Param				testimonialId path string true "Testimonial ID"
// @Description			The recipient approves a pending testimonial, making it public.
// @Produce				application/json
// @Tags				Testimonial
// @Success				204
// @Router				/testimonials/{testimonialId}/approve [post]
func (t *TestimonialController) Approve(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	testimonialId := helpers.MustParseUint(ctx.Param("testimonialId"))

	if err := t.testimonialService.ApproveTestimonial(viewer.ID, testimonialId); err != nil {
		ctx.JSON(testimonialErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// Delete 				godoc
// @Summary				Delete Testimonial by id.
// @Param				testimonialId path string true "Testimonial ID"
// @Description			The recipient rejects or removes a testimonial, or the author deletes one they wrote.
// @Produce				application/json
// @Tags				Testimonial
// @Success				204
// @Router				/testimonials/{testimonialId} [delete]
func (t *TestimonialController) Delete(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	testimonialId := helpers.MustParseUint(ctx.Param("testimonialId"))

	if err := t.testimonialService.DeleteTestimonial(viewer.ID, testimonialId); err != nil {
		ctx.JSON(testimonialErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
