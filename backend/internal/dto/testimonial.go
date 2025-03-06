package dto

type CreateTestimonialRequest struct {
	UserAuthorID    uint
	UserRecipientID uint
	Content         string
}

type FindTestimonialOfUserResponse struct{}
