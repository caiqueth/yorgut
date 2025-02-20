package request

type CreateCommunityRequest struct {
	Name        string
	Description string
	CreatedByID uint
}

type UpdateCommunityRequest struct {
	Name        string
	Description string
}
