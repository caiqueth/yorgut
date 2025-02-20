package model

import (
	"gorm.io/gorm"
)

type Scrap struct {
	gorm.Model
	UserAuthorID    uint
	UserAuthor      *User `gorm:"foreignKey:UserAuthorID;references:ID"`
	UserRecipientID uint
	UserRecipient   *User `gorm:"foreignKey:UserRecipientID;references:ID"`
	Content         string
}

func NewScrap(author *User, recipient *User, content string) *Scrap {
	return &Scrap{
		UserAuthor:      author,
		UserAuthorID:    author.ID,
		UserRecipient:   recipient,
		UserRecipientID: recipient.ID,
		Content:         content,
	}
}
