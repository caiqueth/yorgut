package request

type CreateTestimonialRequest struct {
	UserAuthorID    uint
	UserRecipientID uint
	Content         string
}
