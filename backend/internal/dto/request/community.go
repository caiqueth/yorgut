package request

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
