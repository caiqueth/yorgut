package response

type FindCommunityByIdResponse struct {
	Name        string
	Description string
	CreatedBy   uint
	PictureUrl  string
}

type FindManyCommunitiesResponse struct{}
