package model

import (
	"gorm.io/gorm"
)

type Testimonial struct {
	gorm.Model
	UserAuthorID    uint
	UserAuthor      *User `gorm:"foreignKey:UserAuthorID;references:ID"`
	UserRecipientID uint
	UserRecipient   *User `gorm:"foreignKey:UserRecipientID;references:ID"`
	Content         string
	Approved        bool
}

func NewTestimonial(author *User, recipient *User, content string) *Testimonial {
	return &Testimonial{
		UserAuthorID:    author.ID,
		UserAuthor:      author,
		UserRecipientID: recipient.ID,
		UserRecipient:   recipient,
		Content:         content,
		Approved:        false,
	}
}
