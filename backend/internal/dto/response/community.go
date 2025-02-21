package response

type FindCommunityByIdResponse struct {
	Name        string
	Description string
	CreatedBy   uint
}

type FindManyCommunitiesResponse struct{}
