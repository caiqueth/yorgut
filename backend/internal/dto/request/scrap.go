package request

type CreateScrapRequest struct {
	UserAuthorID    uint
	UserRecipientID uint
	Content         string
}
