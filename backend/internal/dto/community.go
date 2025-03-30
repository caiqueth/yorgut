package dto

import "mime/multipart"

type CreateCommunityRequest struct {
	Name        string                `form:"name" binding:"required" example:"My Community"`
	Description string                `form:"description" binding:"required" example:"This is a test community"`
	Picture     *multipart.FileHeader `form:"picture" swaggerignore:"true"`
	CreatedByID uint                  `form:"createdById" binding:"required" example:"1"`
}

type UpdateCommunityRequest struct {
	Name        string
	Description string
}

type FindCommunityByIdResponse struct {
	Name        string
	Description string
	CreatedBy   uint
	PictureUrl  string
}

type FindManyCommunitiesResponse struct{}
