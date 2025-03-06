package dto

type CreateCommunityRequest struct {
	Name        string
	Description string
	PictureUrl  string
	CreatedByID uint
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
