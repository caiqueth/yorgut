package dto

type CreateScrapRequest struct {
	UserAuthorID    uint
	UserRecipientID uint
	Content         string
}

type FindScrapOfUserRecipientResponse struct{}
