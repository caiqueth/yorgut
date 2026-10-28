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
	// Private communities only show their forum to members, and the owner
	// accepts who joins (CommunityJoinRequest).
	Private bool              `gorm:"not null;default:false"`
	Profile *CommunityProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CommunityProfile struct {
	gorm.Model
	Description    *string
	ProfilePicture *string
	CommunityID    uint       `gorm:"uniqueIndex"`
	Community      *Community `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// CommunityMembership rows are unique per user and community among the live
// (not soft-deleted) ones, so leaving and rejoining keeps the old row.
type CommunityMembership struct {
	gorm.Model
	CommunityID uint       `gorm:"uniqueIndex:idx_community_memberships_active_pair,priority:2,where:deleted_at IS NULL"`
	Community   *Community `gorm:"foreignKey:CommunityID;references:ID"`
	UserID      uint       `gorm:"uniqueIndex:idx_community_memberships_active_pair,priority:1,where:deleted_at IS NULL"`
	User        *User      `gorm:"foreignKey:UserID;references:ID"`
}

// CommunityJoinRequest is a user's request to join a private community,
// waiting for the owner's answer. Accepting or declining deletes it. Unique
// per user and community among the live (not soft-deleted) rows.
type CommunityJoinRequest struct {
	gorm.Model
	CommunityID uint       `gorm:"uniqueIndex:idx_community_join_requests_active_pair,priority:2,where:deleted_at IS NULL"`
	Community   *Community `gorm:"foreignKey:CommunityID;references:ID"`
	UserID      uint       `gorm:"uniqueIndex:idx_community_join_requests_active_pair,priority:1,where:deleted_at IS NULL"`
	User        *User      `gorm:"foreignKey:UserID;references:ID"`
}

// CommunityInvitation is a member's invitation for a friend to join a
// community, waiting for the friend's answer. Accepting or declining deletes
// it. Unique per invitee and community among the live (not soft-deleted)
// rows: a second member can't invite someone who is already invited.
type CommunityInvitation struct {
	gorm.Model
	CommunityID uint       `gorm:"uniqueIndex:idx_community_invitations_active_pair,priority:2,where:deleted_at IS NULL"`
	Community   *Community `gorm:"foreignKey:CommunityID;references:ID"`
	InviterID   uint
	Inviter     *User `gorm:"foreignKey:InviterID;references:ID"`
	InviteeID   uint  `gorm:"uniqueIndex:idx_community_invitations_active_pair,priority:1,where:deleted_at IS NULL"`
	Invitee     *User `gorm:"foreignKey:InviteeID;references:ID"`
}

// ProfileVisit records that a user opened another user's profile: one row
// per visitor, visited user and day (Day is local midnight), whose UpdatedAt
// is the latest visit that day.
type ProfileVisit struct {
	gorm.Model
	VisitedID uint      `gorm:"uniqueIndex:idx_profile_visits_daily,priority:1"`
	Visited   *User     `gorm:"foreignKey:VisitedID;references:ID"`
	VisitorID uint      `gorm:"uniqueIndex:idx_profile_visits_daily,priority:2"`
	Visitor   *User     `gorm:"foreignKey:VisitorID;references:ID"`
	Day       time.Time `gorm:"uniqueIndex:idx_profile_visits_daily,priority:3"`
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
	// Nil for seeded users until they first log in. Both have a size because
	// MySQL can't index a TEXT column.
	GoogleSub *string `gorm:"size:255;uniqueIndex"`
	Email     string  `gorm:"size:255;uniqueIndex"`
	FirstName string
	LastName  string
	// BirthDate is nil until the user sets it (Google login doesn't give it).
	BirthDate   *time.Time
	FriendCount uint
	// LastSeenAt is when the user last used the app, refreshed at most every
	// few minutes. Nil until their first request.
	LastSeenAt *time.Time
	// RecordVisits is the "registrar minhas visitas" setting: when false, the
	// user's visits to profiles aren't recorded or listed, and they don't see
	// their own visitors either.
	RecordVisits bool `gorm:"not null;default:true"`
	// What people who aren't the user's friends may see of their profile
	// (friends always see everything). Open by default, as profiles always
	// were; the birthday is friends-only until the user opens it.
	PublicAlbum        bool `gorm:"not null;default:true"`
	PublicCommunities  bool `gorm:"not null;default:true"`
	PublicScraps       bool `gorm:"not null;default:true"`
	PublicTestimonials bool `gorm:"not null;default:true"`
	PublicBirthday     bool `gorm:"not null;default:false"`
	// HiddenProfileFields lists the profile fields non-friends don't see,
	// comma-separated, by their API names ("religion,birthCity").
	HiddenProfileFields string `gorm:"not null;default:''"`
	// Theme is the look the user picked for the app, one of
	// transport.ThemeKeys. Saved here so it follows them across browsers.
	Theme   string       `gorm:"size:32;not null;default:'2005'"`
	Profile *UserProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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

// Poll is a question for a community's members to vote on (Orkut's
// "enquetes"). It closes at ClosesAt; nil means it stays open until its
// author or the community's owner closes it, which sets ClosesAt to now.
type Poll struct {
	gorm.Model
	CommunityID uint
	Community   *Community `gorm:"foreignKey:CommunityID;references:ID"`
	AuthorID    uint
	Author      *User `gorm:"foreignKey:AuthorID;references:ID"`
	Question    string
	Description *string
	ClosesAt    *time.Time
}

// PollOption is one of a poll's answers, shown in id order.
type PollOption struct {
	gorm.Model
	PollID uint
	Poll   *Poll `gorm:"foreignKey:PollID;references:ID"`
	Text   string
}

// PollVote is a user's vote in a poll: one per user and poll, never changed.
type PollVote struct {
	gorm.Model
	PollID       uint  `gorm:"uniqueIndex:idx_poll_vote_pair"`
	Poll         *Poll `gorm:"foreignKey:PollID;references:ID"`
	PollOptionID uint
	PollOption   *PollOption `gorm:"foreignKey:PollOptionID;references:ID"`
	UserID       uint        `gorm:"uniqueIndex:idx_poll_vote_pair"`
	User         *User       `gorm:"foreignKey:UserID;references:ID"`
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

func NewCommunityJoinRequest(communityId uint, userId uint) *CommunityJoinRequest {
	return &CommunityJoinRequest{
		CommunityID: communityId,
		UserID:      userId,
	}
}

func NewCommunityInvitation(communityId uint, inviterId uint, inviteeId uint) *CommunityInvitation {
	return &CommunityInvitation{
		CommunityID: communityId,
		InviterID:   inviterId,
		InviteeID:   inviteeId,
	}
}

func NewPoll(communityId uint, authorId uint, question string, description *string, closesAt *time.Time) *Poll {
	return &Poll{
		CommunityID: communityId,
		AuthorID:    authorId,
		Question:    question,
		Description: description,
		ClosesAt:    closesAt,
	}
}

func NewPollOption(pollId uint, text string) *PollOption {
	return &PollOption{PollID: pollId, Text: text}
}

func NewPollVote(pollId uint, pollOptionId uint, userId uint) *PollVote {
	return &PollVote{PollID: pollId, PollOptionID: pollOptionId, UserID: userId}
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

func NewUser(firstName string, lastName string, email string, birthDate *time.Time, profilePicture *string) *User {
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
