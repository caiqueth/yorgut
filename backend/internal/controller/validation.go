package controller

import (
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

// Custom binding tags for transport requests, registered on Gin's validator
// before any request is bound:
//
//	community_category: one of transport.CommunityCategories.
//	notblank: not empty and not only whitespace (posts, titles).
func init() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("gin's validator is not go-playground/validator")
	}
	err := validate.RegisterValidation("community_category", func(field validator.FieldLevel) bool {
		return transport.CommunityCategory(field.Field().String()).IsValid()
	})
	if err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("notblank", validators.NotBlank); err != nil {
		panic(err)
	}
}
