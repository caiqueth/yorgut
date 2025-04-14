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
	Owner       *User
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
	Email       string
	FirstName   string
	LastName    string
	BirthDate   time.Time
	FriendCount uint
	UserProfile *UserProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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

// Constructors
func NewCommunity(
	name string,
	category string,
	createdBy *User,
) *Community {
	return &Community{
		Name:        name,
		CreatedByID: createdBy.ID,
		CreatedBy:   createdBy,
		Category:    category,
		Owner:       createdBy,
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

func NewPhotoAlbum(user *User, name string) *PhotoAlbum {
	return &PhotoAlbum{
		UserID: user.ID,
		User:   user,
		Name:   name,
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

func NewUserProfile(user *User, opts ...func(*UserProfile)) *UserProfile {
	profile := &UserProfile{
		User:   user,
		UserID: user.ID,
	}

	for _, opt := range opts {
		opt(profile)
	}

	return profile
}

func NewCommunityProfile(community *Community, opts ...func(*CommunityProfile)) *CommunityProfile {
	profile := &CommunityProfile{
		Community:   community,
		CommunityID: community.ID,
	}

	for _, opt := range opts {
		opt(profile)
	}

	return profile
}

func AddUserProfilePicture(profile *UserProfile, profilePicturePath string) {
	profile.ProfilePicture = &profilePicturePath
}

func AddCommunityProfilePicture(profile *CommunityProfile, profilePicturePath string) {
	profile.ProfilePicture = &profilePicturePath
}

func NewUser(firstName string, lastName string, email string, birthDate time.Time) *User {
	return &User{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		BirthDate:   birthDate,
		FriendCount: 0,
	}
}
