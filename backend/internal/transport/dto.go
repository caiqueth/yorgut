package transport

import (
	"time"
)

type PagedResponse struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Count    int64 `json:"count"`
}

type Response struct {
	Data any `json:"data,omitempty"`
}

type LoginWithGoogleRequest struct {
	// Credential is the Google ID token returned by Sign in with Google.
	Credential string `json:"credential" binding:"required"`
}

type LoginResponse struct {
	// Token is the Yorgut session token, sent back as "Authorization: Bearer <token>".
	Token string `json:"token"`
}

type FindUserByIdResponse struct {
	ID          uint
	Email       string
	FirstName   string
	LastName    string
	BirthDate   time.Time
	FriendCount uint
}

type CreateCommunityMembershipRequest struct {
	CommunityID uint `json:"communityId" binding:"required"`
}

type CommunityMembershipResponse struct {
	ID uint `json:"id"`
}

type FindCommunityMembershipOfUserResponse struct {
}

// CreateCommunityRequest is sent as multipart/form-data, with an optional
// "picture" file next to these fields.
type CreateCommunityRequest struct {
	Name     string `form:"name" binding:"required,max=100" example:"My Community"`
	Category string `form:"category" binding:"required,max=100" example:"Computadores e Internet"`
}

type CreateCommunityResponse struct {
	ID uint `json:"id"`
}

// UpdateCommunityRequest replaces a community's details; an empty description
// clears it.
type UpdateCommunityRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Category    string `json:"category" binding:"required,max=100"`
	Description string `json:"description" binding:"max=2048"`
}

type FindCommunityByIdResponse struct {
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	CreatedByID    uint      `json:"createdById"`
	OwnerID        uint      `json:"ownerId"`
	OwnerFirstName string    `json:"ownerFirstName"`
	CreatedAt      time.Time `json:"createdAt"`
}

type FindCommunityProfileByIdResponse struct {
	Description    *string `json:"description" extensions:"x-nullable"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
}

type CommunityMemberEntry struct {
	ID             uint   `json:"id"`
	FirstName      string `json:"firstName"`
	ProfilePicture string `json:"profilePicture"`
	FriendCount    uint   `json:"friendCount"`
}

type FindCommunityMembersResponse struct {
	MemberCount uint                   `json:"memberCount"`
	Members     []CommunityMemberEntry `json:"members"`
}

type FindManyCommunitiesResponse struct {
	Communities []CommunityItem `json:"communities"`
}

type FriendshipItem struct {
	UserRequestedID uint      `json:"userRequestedId"`
	UserRequesteeID uint      `json:"userRequesteeId"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CreateFriendshipRequest struct {
	UserRequesteeID uint `json:"userRequesteeId"`
}

type ReplyFriendshipRequest struct {
	FriendshipID uint `json:"friendshipId"`
	Accept       bool
}

type FindFriendshipOfUserResponse struct {
	Friendships []FriendshipItem `json:"friendships"`
}

type PendingFriendship struct {
	ID              uint   `json:"id"`
	UserRequestedID uint   `json:"userRequestedId"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	ProfilePicture  string `json:"profilePicture"`
}

type FindPendingFriendshipOfUserResponse struct {
	Count              uint                `json:"count"`
	PendingFriendships []PendingFriendship `json:"pendingFriendships"`
}

type FindFriendsOfUserResponse struct {
	Friends []FriendItem `json:"friends"`
}

type PhotoItem struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"createdAt"`
}

type FindPhotosOfUserResponse struct {
	Photos []PhotoItem `json:"photos"`
	Count  int         `json:"count"`
	// MaxPhotos is how many photos the personal album can hold.
	MaxPhotos int `json:"maxPhotos"`
	// CanUpload is true when the logged user owns the album.
	CanUpload bool `json:"canUpload"`
}

type CreateScrapRequest struct {
	UserRecipientID uint   `json:"userRecipientId"`
	Content         string `json:"content"`
}

type ScrapItem struct {
	ID                   uint      `json:"id"`
	Content              string    `json:"content"`
	AuthorID             uint      `json:"authorId"`
	AuthorName           string    `json:"authorName"`
	AuthorProfilePicture string    `json:"authorProfilePicture"`
	CreatedAt            time.Time `json:"createdAt"`
}

type FindScrapOfUserRecipientResponse struct {
	Scraps []ScrapItem `json:"scraps"`
}

type CreateTestimonialRequest struct {
	UserRecipientID uint   `json:"userRecipientId" binding:"required"`
	Content         string `json:"content" binding:"required,max=2048"`
}

type TestimonialItem struct {
	ID                   uint      `json:"id"`
	Content              string    `json:"content"`
	AuthorID             uint      `json:"authorId"`
	AuthorName           string    `json:"authorName"`
	AuthorProfilePicture string    `json:"authorProfilePicture"`
	Approved             bool      `json:"approved"`
	CreatedAt            time.Time `json:"createdAt"`
}

type FindTestimonialsOfUserResponse struct {
	// PagedResponse describes the approved testimonials; Count is their total.
	PagedResponse
	// Testimonials holds the approved testimonials, visible to everyone.
	Testimonials []TestimonialItem `json:"testimonials"`
	// Pending is only filled when the recipient is the one looking.
	Pending  []TestimonialItem `json:"pending"`
	CanWrite bool              `json:"canWrite"`
}

type CreateUserRequest struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
}

// UpdateMeRequest replaces the logged user's basic data and profile fields.
// Email is not editable: it identifies the Google account. An empty profile
// field clears it.
type UpdateMeRequest struct {
	FirstName string `json:"firstName" binding:"required,max=50"`
	LastName  string `json:"lastName" binding:"max=50"`
	// BirthDate is a calendar date, "YYYY-MM-DD"; empty keeps the current one.
	BirthDate string `json:"birthDate" binding:"omitempty,datetime=2006-01-02"`

	Description        string `json:"description" binding:"max=2048"`
	RelationshipStatus string `json:"relationshipStatus" binding:"max=2048"`
	Interests          string `json:"interests" binding:"max=2048"`
	About              string `json:"about" binding:"max=2048"`
	Religion           string `json:"religion" binding:"max=2048"`
	Humor              string `json:"humor" binding:"max=2048"`
	Fashion            string `json:"fashion" binding:"max=2048"`
	LivesWith          string `json:"livesWith" binding:"max=2048"`
	BirthCity          string `json:"birthCity" binding:"max=2048"`
	Website            string `json:"website" binding:"max=2048"`
	Passions           string `json:"passions" binding:"max=2048"`
	Sports             string `json:"sports" binding:"max=2048"`
	Hobbies            string `json:"hobbies" binding:"max=2048"`
	Books              string `json:"books" binding:"max=2048"`
	Music              string `json:"music" binding:"max=2048"`
	Movies             string `json:"movies" binding:"max=2048"`
	Cuisine            string `json:"cuisine" binding:"max=2048"`
}

type UserItem struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	BirthDate   time.Time `json:"birthDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	FriendCount uint      `json:"friendCount"`
}

type UserProfileItem struct {
	ProfilePicture     *string `json:"profilePicture" extensions:"x-nullable"`
	Description        *string `json:"description" extensions:"x-nullable"`
	RelationshipStatus *string `json:"relationshipStatus" extensions:"x-nullable"`
	Interests          *string `json:"interests" extensions:"x-nullable"`
	About              *string `json:"about" extensions:"x-nullable"`
	Religion           *string `json:"religion" extensions:"x-nullable"`
	Humor              *string `json:"humor" extensions:"x-nullable"`
	Fashion            *string `json:"fashion" extensions:"x-nullable"`
	LivesWith          *string `json:"livesWith" extensions:"x-nullable"`
	BirthCity          *string `json:"birthCity" extensions:"x-nullable"`
	Website            *string `json:"website" extensions:"x-nullable"`
	Passions           *string `json:"passions" extensions:"x-nullable"`
	Sports             *string `json:"sports" extensions:"x-nullable"`
	Hobbies            *string `json:"hobbies" extensions:"x-nullable"`
	Books              *string `json:"books" extensions:"x-nullable"`
	Music              *string `json:"music" extensions:"x-nullable"`
	Movies             *string `json:"movies" extensions:"x-nullable"`
	Cuisine            *string `json:"cuisine" extensions:"x-nullable"`
}

type UserSearchItem struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"firstName"`
	BirthDate      time.Time `json:"birthDate"`
	ProfilePicture *string   `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint      `json:"friendCount"`
}

type FriendItem struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"firstName"`
	BirthDate      time.Time `json:"birthDate"`
	ProfilePicture string    `json:"profilePicture"`
	FriendCount    uint      `json:"friendCount"`
}

type CommunityItem struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	CreatedBy       uint    `json:"firstName"`
	ProfilePicture  *string `json:"profilePicture" extensions:"x-nullable"`
	MembershipCount uint    `json:"membershipCount"`
	OwnerID         uint    `json:"owner"`
}

type MyOwnedCommunityItem struct {
	ID              uint    `json:"id"`
	MembershipCount uint    `json:"membershipCount"`
	Name            string  `json:"name"`
	ProfilePicture  *string `json:"profilePicture" extensions:"x-nullable"`
}

type FindUserResponse struct {
	UserItem
}

type ProfileViews struct {
	AllTime  uint `json:"allTime"`
	LastWeek uint `json:"lastWeek"`
}

type RecentVisitor struct {
	FirstName string `json:"firstName"`
	UserID    uint   `json:"userId"`
}

type WelcomeResponse struct {
	UserFirstName       string          `json:"userFirstName"`
	UserCreatedAt       time.Time       `json:"userCreatedAt"`
	IndirectConnections uint            `json:"indirectConnections"`
	DirectConnections   uint            `json:"directConnections"`
	NumberOfFans        uint            `json:"numberOfFans"`
	NumberOfScraps      uint            `json:"numberOfScraps"`
	ProfileViews        ProfileViews    `json:"profileViews"`
	RecentVisitors      []RecentVisitor `json:"recentVisitors"`
	TodaysLuck          string          `json:"todaysLuck"`
}

type FindUserWithProfileByIdResponse struct {
	UserItem
	UserProfileItem
}

type FindManyUsersResponse struct {
	Users []UserItem `json:"users"`
}

type SearchUsersResponse struct {
	PagedResponse
	Users []*UserSearchItem `json:"users"`
}

type CommunitySearchItem struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
	MemberCount    uint    `json:"memberCount"`
}

type SearchCommunitiesResponse struct {
	PagedResponse
	Communities []CommunitySearchItem `json:"communities"`
}

type FindMyOwnedCommunitiesResponse struct {
	Communities []MyOwnedCommunityItem `json:"communities"`
}

type FindUserMembershipsResponse struct {
	Communities     []CommunityItem `json:"communities"`
	MembershipCount uint            `json:"membershipCount"`
}

type AmIMemberOfResponse struct {
	AmIMember bool `json:"amIMember"`
}

type AmIFriendOfResponse struct {
	FriendshipId        *uint `json:"friendshipId" extensions:"x-nullable"`
	AmIFriend           bool  `json:"amIFriend"`
	IsPendingTheirReply bool  `json:"isPendingTheirReply"`
	IsPendingMyReply    bool  `json:"isPendingMyReply"`
}

type TopicSummary struct {
	ID                   uint      `json:"id"`
	AuthorProfilePicture string    `json:"authorProfilePicture"`
	Title                string    `json:"title"`
	Replies              int64     `json:"replies"`
	LastReplyPreview     string    `json:"lastReplyPreview"`
	LastReplyDate        time.Time `json:"lastReplyDate"`
}

type TopicsResponse struct {
	Topics []TopicSummary `json:"topics"`
}

type Topic struct {
	ID                   uint      `json:"id"`
	AuthorID             uint      `json:"authorId"`
	AuthorFirstName      string    `json:"authorFirstName"`
	AuthorProfilePicture string    `json:"authorProfilePicture"`
	Title                string    `json:"title"`
	Text                 string    `json:"text"`
	CreatedAt            time.Time `json:"createdAt"`
	Replies              int64     `json:"replies"`
}

type TopicReply struct {
	ID                   uint      `json:"id"`
	AuthorID             uint      `json:"authorId"`
	AuthorFirstName      string    `json:"authorFirstName"`
	AuthorProfilePicture string    `json:"authorProfilePicture"`
	Title                string    `json:"title"`
	Text                 string    `json:"text"`
	CreatedAt            time.Time `json:"createdAt"`
}

type TopicRepliesResponse struct {
	Topic   *Topic       `json:"topic" extensions:"x-nullable"`
	Replies []TopicReply `json:"replies"`
}

type CreateTopicRequest struct {
	CommunityID uint   `json:"communityId" binding:"required"`
	Title       string `json:"title" binding:"required,max=200"`
	Text        string `json:"text" binding:"required,max=2048"`
}

// ReplyTopicRequest has no community: a reply always belongs to the community
// of the topic it answers.
type ReplyTopicRequest struct {
	Title string `json:"title" binding:"max=200"`
	Text  string `json:"text" binding:"required,max=2048"`
}

type UpsertUserReviewRequest struct {
	Cool        *uint8 `json:"cool" binding:"required,min=0,max=3" example:"3"`
	Sexy        *uint8 `json:"sexy" binding:"required,min=0,max=3" example:"2"`
	Trustworthy *uint8 `json:"trustworthy" binding:"required,min=0,max=3" example:"3"`
}

type UserReviewItem struct {
	Cool        uint8 `json:"cool"`
	Sexy        uint8 `json:"sexy"`
	Trustworthy uint8 `json:"trustworthy"`
}

type UserFanSummaryResponse struct {
	FanCount int64 `json:"fanCount"`
	AmIFan   bool  `json:"amIFan"`
	CanFan   bool  `json:"canFan"`
}

type UserReviewSummaryResponse struct {
	Cool        float64         `json:"cool"`
	Sexy        float64         `json:"sexy"`
	Trustworthy float64         `json:"trustworthy"`
	ReviewCount int64           `json:"reviewCount"`
	CanReview   bool            `json:"canReview"`
	MyReview    *UserReviewItem `json:"myReview" extensions:"x-nullable"`
}
