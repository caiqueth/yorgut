package transport

import (
	"time"
)

type PagedResponse struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Count    int64 `json:"count"`
}

// ErrorResponse is the body of every error answer (4xx and 5xx).
type ErrorResponse struct {
	Error string `json:"error"`
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

// CreateCommunityRequest is sent as multipart/form-data, with an optional
// "picture" file next to these fields.
type CreateCommunityRequest struct {
	Name     string            `form:"name" binding:"required,max=100" example:"My Community"`
	Category CommunityCategory `form:"category" binding:"required,community_category" example:"computers_internet"`
	// Description is optional; a blank one leaves the community without one.
	Description string `form:"description" binding:"max=2048"`
	// Private communities show their forum only to members, and the owner
	// accepts who joins.
	Private bool `form:"private"`
}

type CreateCommunityResponse struct {
	ID uint `json:"id"`
}

// UpdateCommunityRequest replaces a community's details; an empty description
// clears it.
type UpdateCommunityRequest struct {
	Name        string            `json:"name" binding:"required,max=100"`
	Category    CommunityCategory `json:"category" binding:"required,community_category"`
	Description string            `json:"description" binding:"max=2048"`
	// Making a private community public accepts its pending join requests.
	Private bool `json:"private"`
}

type FindCommunityByIdResponse struct {
	Name           string            `json:"name"`
	Category       CommunityCategory `json:"category"`
	CreatedByID    uint              `json:"createdById"`
	OwnerID        uint              `json:"ownerId"`
	OwnerFirstName string            `json:"ownerFirstName"`
	Private        bool              `json:"private"`
	CreatedAt      time.Time         `json:"createdAt"`
}

// JoinCommunityResponse says whether the user became a member (Pending
// false) or asked a private community's owner to accept them (Pending true).
type JoinCommunityResponse struct {
	Pending bool `json:"pending"`
}

// JoinRequestItem is a user waiting to join a private community.
type JoinRequestItem struct {
	UserID         uint      `json:"userId"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	ProfilePicture *string   `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint      `json:"friendCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

// FindJoinRequestsResponse lists a private community's pending join requests,
// oldest first.
type FindJoinRequestsResponse struct {
	Requests []JoinRequestItem `json:"requests"`
}

type ReplyJoinRequestRequest struct {
	// Accept false (or missing) declines the request.
	Accept bool `json:"accept"`
}

type InviteToCommunityRequest struct {
	// UserID is the friend being invited.
	UserID uint `json:"userId" binding:"required"`
}

// InvitableFriendItem is one of the logged user's friends, with whether they
// can still be invited to the community.
type InvitableFriendItem struct {
	ID             uint    `json:"id"`
	FirstName      string  `json:"firstName"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint    `json:"friendCount"`
	// Member: already in the community. Invited: someone already invited them.
	Member  bool `json:"member"`
	Invited bool `json:"invited"`
}

// FindInvitableFriendsResponse is one page of the logged user's friends, to
// pick whom to invite to a community.
type FindInvitableFriendsResponse struct {
	PagedResponse
	Friends []InvitableFriendItem `json:"friends"`
}

// CommunityInvitationItem is an invitation the logged user received.
type CommunityInvitationItem struct {
	ID                    uint      `json:"id"`
	CommunityID           uint      `json:"communityId"`
	CommunityName         string    `json:"communityName"`
	CommunityPicture      *string   `json:"communityPicture" extensions:"x-nullable"`
	CommunityPrivate      bool      `json:"communityPrivate"`
	InviterID             uint      `json:"inviterId"`
	InviterFirstName      string    `json:"inviterFirstName"`
	InviterProfilePicture *string   `json:"inviterProfilePicture" extensions:"x-nullable"`
	CreatedAt             time.Time `json:"createdAt"`
}

// FindMyCommunityInvitationsResponse lists the logged user's pending
// community invitations, newest first.
type FindMyCommunityInvitationsResponse struct {
	Invitations []CommunityInvitationItem `json:"invitations"`
}

type FindCommunityProfileByIdResponse struct {
	Description    *string `json:"description" extensions:"x-nullable"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
}

type CommunityMemberEntry struct {
	ID             uint    `json:"id"`
	FirstName      string  `json:"firstName"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint    `json:"friendCount"`
}

// FindCommunityMembersResponse is one page of members; Count is how many the
// community has.
type FindCommunityMembersResponse struct {
	PagedResponse
	Members []CommunityMemberEntry `json:"members"`
}

type FindManyCommunitiesResponse struct {
	Communities []CommunityItem `json:"communities"`
}

type CreateFriendshipRequest struct {
	UserRequesteeID uint `json:"userRequesteeId" binding:"required"`
}

type ReplyFriendshipRequest struct {
	FriendshipID uint `json:"friendshipId" binding:"required"`
	// Accept false (or missing) declines the request.
	Accept bool `json:"accept"`
}

type PendingFriendship struct {
	ID              uint    `json:"id"`
	UserRequestedID uint    `json:"userRequestedId"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	ProfilePicture  *string `json:"profilePicture" extensions:"x-nullable"`
}

type FindPendingFriendshipOfUserResponse struct {
	Count              uint                `json:"count"`
	PendingFriendships []PendingFriendship `json:"pendingFriendships"`
}

// FindFriendsOfUserResponse is one page of friends; Count is how many the user
// has.
type FindFriendsOfUserResponse struct {
	PagedResponse
	Friends []FriendItem `json:"friends"`
}

// MyFriendItem is a friend of the logged user, with what the logged user
// gave them: their rating (zeros when not rated yet) and fan star.
type MyFriendItem struct {
	FriendItem
	MyReview UserReviewItem `json:"myReview"`
	AmIFan   bool           `json:"amIFan"`
}

// FindMyFriendsResponse is one page of the logged user's friends; Count is
// how many they have.
type FindMyFriendsResponse struct {
	PagedResponse
	Friends []MyFriendItem `json:"friends"`
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
	UserRecipientID uint   `json:"userRecipientId" binding:"required"`
	Content         string `json:"content" binding:"required,notblank,max=2048"`
}

type ScrapItem struct {
	ID                   uint      `json:"id"`
	Content              string    `json:"content"`
	AuthorID             uint      `json:"authorId"`
	AuthorName           string    `json:"authorName"`
	AuthorProfilePicture *string   `json:"authorProfilePicture" extensions:"x-nullable"`
	CreatedAt            time.Time `json:"createdAt"`
}

// FindScrapOfUserRecipientResponse is one page of a scrapbook; Count is how
// many scraps it has.
type FindScrapOfUserRecipientResponse struct {
	PagedResponse
	Scraps []ScrapItem `json:"scraps"`
}

type CreateTestimonialRequest struct {
	UserRecipientID uint   `json:"userRecipientId" binding:"required"`
	Content         string `json:"content" binding:"required,notblank,max=2048"`
}

type TestimonialItem struct {
	ID                   uint      `json:"id"`
	Content              string    `json:"content"`
	AuthorID             uint      `json:"authorId"`
	AuthorName           string    `json:"authorName"`
	AuthorProfilePicture *string   `json:"authorProfilePicture" extensions:"x-nullable"`
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

// UserItem is what any logged user may see about another user. It has no
// email on purpose: emails are private (only FindMeResponse carries one).
type UserItem struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
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
	ID             uint    `json:"id"`
	FirstName      string  `json:"firstName"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint    `json:"friendCount"`
}

type FriendItem struct {
	ID             uint    `json:"id"`
	FirstName      string  `json:"firstName"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint    `json:"friendCount"`
}

type CommunityItem struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	CreatedBy       uint    `json:"createdBy"`
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

// FindMeResponse is the logged user, with their own email and settings.
type FindMeResponse struct {
	UserItem
	Email   string          `json:"email"`
	Privacy PrivacySettings `json:"privacy"`
	// Theme is the look the user picked for the app.
	Theme string `json:"theme" enums:"2005,2008"`
}

// ThemeKeys are the looks a user may pick for the app; the first one is the
// default. The enums tags on Theme fields list them too, for the client.
var ThemeKeys = []string{"2005", "2008"}

// UpdateMyThemeRequest picks the look of the app for the logged user.
type UpdateMyThemeRequest struct {
	Theme string `json:"theme" binding:"required" enums:"2005,2008"`
}

// UpdateMyPrivacyRequest replaces every privacy setting. Pointers so that
// false passes "required".
type UpdateMyPrivacyRequest struct {
	RecordVisits       *bool `json:"recordVisits" binding:"required"`
	PublicAlbum        *bool `json:"publicAlbum" binding:"required"`
	PublicCommunities  *bool `json:"publicCommunities" binding:"required"`
	PublicScraps       *bool `json:"publicScraps" binding:"required"`
	PublicTestimonials *bool `json:"publicTestimonials" binding:"required"`
	PublicBirthday     *bool `json:"publicBirthday" binding:"required"`
	// Keys of the profile fields non-friends don't see; empty shows them all.
	HiddenProfileFields []string `json:"hiddenProfileFields" binding:"required"`
}

// ProfileViews counts visits to a profile, one per visitor and day.
type ProfileViews struct {
	AllTime uint `json:"allTime"`
	// LastSevenDays includes today.
	LastSevenDays uint `json:"lastSevenDays"`
}

type RecentVisitor struct {
	FirstName string `json:"firstName"`
	UserID    uint   `json:"userId"`
}

// BirthdayItem is a friend whose birthday is coming up.
type BirthdayItem struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"firstName"`
	ProfilePicture *string   `json:"profilePicture" extensions:"x-nullable"`
	BirthDate      time.Time `json:"birthDate"`
	// NextBirthday is the date of the coming birthday (today counts).
	NextBirthday time.Time `json:"nextBirthday"`
	DaysUntil    int       `json:"daysUntil"`
}

// FindUpcomingBirthdaysResponse lists friends with a birthday in the next
// days, soonest first.
type FindUpcomingBirthdaysResponse struct {
	Birthdays []BirthdayItem `json:"birthdays"`
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
	// For the "próximas etapas": whether the user rated any current friend,
	// and how many communities they're in (DirectConnections is their friends).
	HasRatedFriends bool `json:"hasRatedFriends"`
	CommunityCount  uint `json:"communityCount"`
	// RecordVisits false: the user turned off "registrar minhas visitas", so
	// RecentVisitors is empty on purpose.
	RecordVisits bool `json:"recordVisits"`
}

// FindUserWithProfileByIdResponse is a profile as the viewer may see it: for
// someone who isn't the user's friend, the fields and birthday the user keeps
// to friends are left out.
type FindUserWithProfileByIdResponse struct {
	UserItem
	UserProfileItem
	// BirthDate is a calendar date; null when the viewer may not see it or
	// the user never set it.
	BirthDate *time.Time `json:"birthDate" extensions:"x-nullable"`
	// HiddenFields are the profile fields left out for this viewer (by their
	// JSON names, "birthDate" included), so the page leaves their rows out
	// instead of showing "-".
	HiddenFields []string `json:"hiddenFields"`
	// Visibility tells which parts of the profile the viewer may see.
	Visibility ProfileVisibility `json:"visibility"`
}

// ProfileVisibility tells which parts of a profile a viewer may see: all of
// them for the user and their friends, the ones the user opened for others.
type ProfileVisibility struct {
	Album        bool `json:"album"`
	Communities  bool `json:"communities"`
	Scraps       bool `json:"scraps"`
	Testimonials bool `json:"testimonials"`
}

// ProfileFieldKeys are the profile fields a user may keep to friends, by
// their JSON names, in the order the profile shows them.
var ProfileFieldKeys = []string{
	"description", "relationshipStatus", "interests", "about", "religion", "humor",
	"fashion", "livesWith", "birthCity", "website", "passions", "sports",
	"hobbies", "books", "music", "movies", "cuisine",
}

// PrivacySettings is what the user lets non-friends see (friends always see
// everything), plus the "registrar minhas visitas" setting.
type PrivacySettings struct {
	RecordVisits       bool `json:"recordVisits"`
	PublicAlbum        bool `json:"publicAlbum"`
	PublicCommunities  bool `json:"publicCommunities"`
	PublicScraps       bool `json:"publicScraps"`
	PublicTestimonials bool `json:"publicTestimonials"`
	PublicBirthday     bool `json:"publicBirthday"`
	// HiddenProfileFields are keys of ProfileFieldKeys.
	HiddenProfileFields []string `json:"hiddenProfileFields"`
}

type SearchUsersResponse struct {
	PagedResponse
	Users []*UserSearchItem `json:"users"`
}

type CommunitySearchItem struct {
	ID             uint              `json:"id"`
	Name           string            `json:"name"`
	Category       CommunityCategory `json:"category"`
	ProfilePicture *string           `json:"profilePicture" extensions:"x-nullable"`
	MemberCount    uint              `json:"memberCount"`
	Private        bool              `json:"private"`
	// Pending is true when the logged user already asked to join.
	Pending bool `json:"pending"`
}

type SearchCommunitiesResponse struct {
	PagedResponse
	Communities []CommunitySearchItem `json:"communities"`
}

type FindMyOwnedCommunitiesResponse struct {
	PagedResponse
	Communities []MyOwnedCommunityItem `json:"communities"`
}

// FindUserMembershipsResponse is one page of the user's communities; Count is
// how many they are a member of.
type FindUserMembershipsResponse struct {
	PagedResponse
	Communities []CommunityItem `json:"communities"`
}

type AmIMemberOfResponse struct {
	AmIMember bool `json:"amIMember"`
	// Pending is true when the logged user asked to join a private community
	// and the owner hasn't answered yet.
	Pending bool `json:"pending"`
}

type AmIFriendOfResponse struct {
	FriendshipId        *uint `json:"friendshipId" extensions:"x-nullable"`
	AmIFriend           bool  `json:"amIFriend"`
	IsPendingTheirReply bool  `json:"isPendingTheirReply"`
	IsPendingMyReply    bool  `json:"isPendingMyReply"`
}

type TopicSummary struct {
	ID                   uint      `json:"id"`
	AuthorID             uint      `json:"authorId"`
	AuthorProfilePicture *string   `json:"authorProfilePicture" extensions:"x-nullable"`
	Title                string    `json:"title"`
	Replies              int64     `json:"replies"`
	LastReplyPreview     string    `json:"lastReplyPreview"`
	LastReplyDate        time.Time `json:"lastReplyDate"`
}

// TopicsResponse is one page of a community's topics; Count is how many it
// has.
type TopicsResponse struct {
	PagedResponse
	Topics []TopicSummary `json:"topics"`
}

type Topic struct {
	ID                   uint      `json:"id"`
	AuthorID             uint      `json:"authorId"`
	AuthorFirstName      string    `json:"authorFirstName"`
	AuthorProfilePicture *string   `json:"authorProfilePicture" extensions:"x-nullable"`
	Title                string    `json:"title"`
	Text                 string    `json:"text"`
	CreatedAt            time.Time `json:"createdAt"`
	Replies              int64     `json:"replies"`
}

type TopicReply struct {
	ID                   uint      `json:"id"`
	AuthorID             uint      `json:"authorId"`
	AuthorFirstName      string    `json:"authorFirstName"`
	AuthorProfilePicture *string   `json:"authorProfilePicture" extensions:"x-nullable"`
	Title                string    `json:"title"`
	Text                 string    `json:"text"`
	CreatedAt            time.Time `json:"createdAt"`
}

// TopicRepliesResponse is the topic and one page of its replies; Count is how
// many replies it has.
type TopicRepliesResponse struct {
	PagedResponse
	Topic   *Topic       `json:"topic" extensions:"x-nullable"`
	Replies []TopicReply `json:"replies"`
}

type CreateTopicRequest struct {
	CommunityID uint   `json:"communityId" binding:"required"`
	Title       string `json:"title" binding:"required,notblank,max=200"`
	Text        string `json:"text" binding:"required,notblank,max=2048"`
}

// ReplyTopicRequest has no community: a reply always belongs to the community
// of the topic it answers.
type ReplyTopicRequest struct {
	Title string `json:"title" binding:"max=200"`
	Text  string `json:"text" binding:"required,notblank,max=2048"`
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

// FanItem is a user who is a fan of someone.
type FanItem struct {
	ID             uint    `json:"id"`
	FirstName      string  `json:"firstName"`
	ProfilePicture *string `json:"profilePicture" extensions:"x-nullable"`
	FriendCount    uint    `json:"friendCount"`
}

// FindUserFanListResponse is one page of a user's fans, the newest first;
// Count is how many fans they have.
type FindUserFanListResponse struct {
	PagedResponse
	Fans []FanItem `json:"fans"`
}

type UserReviewSummaryResponse struct {
	Cool        float64         `json:"cool"`
	Sexy        float64         `json:"sexy"`
	Trustworthy float64         `json:"trustworthy"`
	ReviewCount int64           `json:"reviewCount"`
	CanReview   bool            `json:"canReview"`
	MyReview    *UserReviewItem `json:"myReview" extensions:"x-nullable"`
}

// CreatePollRequest creates a poll in a community. Options are trimmed and
// must all be different. validate:"optional" only marks a field optional in
// the spec (swag runs with --requiredByDefault); Gin reads binding.
type CreatePollRequest struct {
	Question    string   `json:"question" binding:"required,max=200"`
	Description string   `json:"description" binding:"max=2048" validate:"optional"`
	Options     []string `json:"options" binding:"required,min=2,max=10,dive,required,max=100"`
	// ClosesAt is when voting ends; omitted, the poll stays open until its
	// author or the community's owner closes it.
	ClosesAt *time.Time `json:"closesAt" validate:"optional"`
}

type CreatePollResponse struct {
	ID uint `json:"id"`
}

type VotePollRequest struct {
	OptionID uint `json:"optionId" binding:"required"`
}

// PollSummary is a poll in a community's list.
type PollSummary struct {
	ID              uint       `json:"id"`
	Question        string     `json:"question"`
	AuthorID        uint       `json:"authorId"`
	AuthorFirstName string     `json:"authorFirstName"`
	CreatedAt       time.Time  `json:"createdAt"`
	ClosesAt        *time.Time `json:"closesAt" extensions:"x-nullable"`
	Closed          bool       `json:"closed"`
	VoteCount       int64      `json:"voteCount"`
	// Voted: the logged user voted in it.
	Voted bool `json:"voted"`
}

// PollsResponse is one page of a community's polls, newest first.
type PollsResponse struct {
	PagedResponse
	Polls []PollSummary `json:"polls"`
}

type PollOptionItem struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
	// Votes is null until the logged user voted or the poll closed, so
	// results don't sway the vote.
	Votes *int64 `json:"votes" extensions:"x-nullable"`
}

// PollResponse is a poll with its options, and its results once the logged
// user may see them.
type PollResponse struct {
	ID                   uint       `json:"id"`
	CommunityID          uint       `json:"communityId"`
	Question             string     `json:"question"`
	Description          *string    `json:"description" extensions:"x-nullable"`
	AuthorID             uint       `json:"authorId"`
	AuthorFirstName      string     `json:"authorFirstName"`
	AuthorProfilePicture *string    `json:"authorProfilePicture" extensions:"x-nullable"`
	CreatedAt            time.Time  `json:"createdAt"`
	ClosesAt             *time.Time `json:"closesAt" extensions:"x-nullable"`
	Closed               bool       `json:"closed"`
	VoteCount            int64      `json:"voteCount"`
	// MyOptionID is the logged user's vote, null when they haven't voted.
	MyOptionID *uint `json:"myOptionId" extensions:"x-nullable"`
	// CanVote: a member who hasn't voted, in an open poll.
	CanVote bool `json:"canVote"`
	// CanClose: the author or the community's owner, while the poll is open.
	CanClose bool             `json:"canClose"`
	Options  []PollOptionItem `json:"options"`
}

// InstanceInfoResponse is what this Yorgut instance's administrator set up
// for everyone to see, logged in or not.
type InstanceInfoResponse struct {
	// AdminEmail is how to reach the administrator (ADMIN_EMAIL), null when
	// it isn't set.
	AdminEmail *string `json:"adminEmail" extensions:"x-nullable"`
}
