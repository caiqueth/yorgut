package model

import (
	"time"

	"gorm.io/gorm"
)

// Structs
type Community struct {
	gorm.Model
	Name        string
	CreatedByID uint
	CreatedBy   *User `gorm:"foreignKey:CreatedByID;references:ID"`
	Category    string
	OwnerID     uint
	Owner       *User `gorm:"foreignKey:OwnerID;references:ID"`
	MemberCount uint
	Profile     *CommunityProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CommunityProfile struct {
	gorm.Model
	Description    *string
	ProfilePicture *string
	CommunityID    uint       `gorm:"uniqueIndex"`
	Community      *Community `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CommunityMembership struct {
	gorm.Model
	CommunityID uint
	Community   *Community `gorm:"foreignKey:CommunityID;references:ID"`
	UserID      uint
	User        *User `gorm:"foreignKey:UserID;references:ID"`
}

type Friendship struct {
	gorm.Model
	UserRequestedID uint
	UserRequested   *User `gorm:"foreignKey:UserRequestedID;references:ID"`
	UserRequesteeID uint
	UserRequestee   *User `gorm:"foreignKey:UserRequesteeID;references:ID"`
	Accepted        bool
	AcceptedAt      *time.Time
}

type PhotoAlbum struct {
	gorm.Model
	UserID uint
	User   *User `gorm:"foreignKey:UserID;references:ID"`
	Name   string
}

type Photo struct {
	gorm.Model
	PhotoAlbumID uint
	PhotoAlbum   *PhotoAlbum `gorm:"foreignKey:PhotoAlbumID;references:ID"`
	URL          string
	Caption      string
}

type Scrap struct {
	gorm.Model
	UserAuthorID    uint
	UserAuthor      *User `gorm:"foreignKey:UserAuthorID;references:ID"`
	UserRecipientID uint
	UserRecipient   *User `gorm:"foreignKey:UserRecipientID;references:ID"`
	Content         string
}

type Testimonial struct {
	gorm.Model
	UserAuthorID    uint
	UserAuthor      *User `gorm:"foreignKey:UserAuthorID;references:ID"`
	UserRecipientID uint
	UserRecipient   *User `gorm:"foreignKey:UserRecipientID;references:ID"`
	Content         string
	Approved        bool
}

type User struct {
	gorm.Model
	// GoogleSub is Google's stable id for the account (the ID token's "sub").
	// Nil for seeded users until they first log in.
	GoogleSub   *string `gorm:"uniqueIndex"`
	Email       string  `gorm:"uniqueIndex"`
	FirstName   string
	LastName    string
	BirthDate   time.Time
	FriendCount uint
	Profile     *UserProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type UserReview struct {
	gorm.Model
	ReviewerID  uint  `gorm:"uniqueIndex:idx_user_review_pair"`
	Reviewer    *User `gorm:"foreignKey:ReviewerID;references:ID"`
	ReviewedID  uint  `gorm:"uniqueIndex:idx_user_review_pair"`
	Reviewed    *User `gorm:"foreignKey:ReviewedID;references:ID"`
	Cool        uint8
	Sexy        uint8
	Trustworthy uint8
}

// UserFan marks FanID as a fan of IdolID (the star on Orkut's friend list).
type UserFan struct {
	gorm.Model
	FanID  uint  `gorm:"uniqueIndex:idx_user_fan_pair"`
	Fan    *User `gorm:"foreignKey:FanID;references:ID"`
	IdolID uint  `gorm:"uniqueIndex:idx_user_fan_pair"`
	Idol   *User `gorm:"foreignKey:IdolID;references:ID"`
}

type UserProfile struct {
	gorm.Model
	UserID             uint  `gorm:"uniqueIndex"`
	User               *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProfilePicture     *string
	Description        *string
	RelationshipStatus *string
	Interests          *string
	About              *string
	Religion           *string
	Humor              *string
	Fashion            *string
	LivesWith          *string
	BirthCity          *string
	Website            *string
	Passions           *string
	Sports             *string
	Hobbies            *string
	Books              *string
	Music              *string
	Movies             *string
	Cuisine            *string
}

type Topic struct {
	gorm.Model
	CommunityID   uint
	Community     *Community
	AuthorID      uint
	Author        *User
	ParentTopicID *uint
	ParentTopic   *Topic
	Title         string
	Text          string
}

// Constructors
func NewCommunity(
	name string,
	category string,
	createdBy *User,
) *Community {
	return &Community{
		Name:        name,
		Category:    category,
		CreatedByID: createdBy.ID,
		OwnerID:     createdBy.ID,
		MemberCount: 0,
	}
}

func NewCommunityMembership(community *Community, user *User) *CommunityMembership {
	return &CommunityMembership{
		CommunityID: community.ID,
		Community:   community,
		UserID:      user.ID,
		User:        user,
	}
}

func NewFriendship(userRequested *User, userRequestee *User) *Friendship {
	return &Friendship{
		UserRequestedID: userRequested.ID,
		UserRequested:   userRequested,
		UserRequesteeID: userRequestee.ID,
		UserRequestee:   userRequestee,
		Accepted:        false,
		AcceptedAt:      nil,
	}
}

func NewPhoto(photoAlbum *PhotoAlbum, url string, caption string) *Photo {
	return &Photo{
		PhotoAlbumID: photoAlbum.ID,
		PhotoAlbum:   photoAlbum,
		URL:          url,
		Caption:      caption,
	}
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

func NewUserReview(reviewerID uint, reviewedID uint, cool uint8, sexy uint8, trustworthy uint8) *UserReview {
	return &UserReview{
		ReviewerID:  reviewerID,
		ReviewedID:  reviewedID,
		Cool:        cool,
		Sexy:        sexy,
		Trustworthy: trustworthy,
	}
}

func NewUserFan(fanID uint, idolID uint) *UserFan {
	return &UserFan{
		FanID:  fanID,
		IdolID: idolID,
	}
}

func NewUserProfile(user *User) *UserProfile {
	profile := &UserProfile{
		User:   user,
		UserID: user.ID,
	}
	return profile
}

func AddUserProfilePicture(profile *UserProfile, profilePicturePath string) {
	profile.ProfilePicture = &profilePicturePath
}

func NewUser(firstName string, lastName string, email string, birthDate time.Time, profilePicture *string) *User {
	return &User{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		BirthDate:   birthDate,
		FriendCount: 0,
		Profile: &UserProfile{
			ProfilePicture: profilePicture,
		},
	}
}

func NewTopic(authorId uint, communityId uint, parentTopicId *uint, title string, text string) *Topic {
	topic := Topic{
		CommunityID:   communityId,
		AuthorID:      authorId,
		ParentTopicID: parentTopicId,
		Title:         title,
		Text:          text,
	}

	return &topic
}
