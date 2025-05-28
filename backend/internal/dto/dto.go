package dto

import (
	"mime/multipart"
	"time"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
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
	CommunityID uint `json:"communityId"`
	UserID      uint `json:"userId"`
}

type CommunityMembershipResponse struct {
	ID uint `json:"id"`
}

type FindCommunityMembershipOfUserResponse struct {
}

type CreateCommunityRequest struct {
	Name        string `form:"name" binding:"required" example:"My Community"`
	CreatedByID uint   `form:"createdById" binding:"required" example:"1"`
	Category    string `form:"category" binding:"required" example:"Computadores e Internet"`
}

type CreateCommunityProfileRequest struct {
	CommunityID uint                  `json:"communityId"`
	Description string                `form:"description" binding:"required" example:"This is a test community"`
	Picture     *multipart.FileHeader `form:"picture" swaggerignore:"true"`
}

type UpdateCommunityRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FindCommunityByIdResponse struct {
	Name           string    `json:"name"`
	CreatedByID    uint      `json:"createdById"`
	OwnerID        uint      `json:"ownerId"`
	OwnerFirstName string    `json:"ownerFirstName"`
	CreatedAt      time.Time `json:"createdAt"`
}

type FindCommunityProfileByIdResponse struct {
	Description    *string `json:"description"`
	ProfilePicture *string `json:"profilePicture"`
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
	Communities []CommunityDto `json:"communities"`
}

type FriendshipDto struct {
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
	Friendships []FriendshipDto `json:"friendships"`
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
	Friends []FriendDto `json:"friends"`
}

type PhotoAlbumDto struct {
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
}

type CreatePhotoAlbumRequest struct {
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
}

type FindPhotoAlbumOfUserResponse struct {
	PhotoAlbums []PhotoAlbumDto `json:"photoAlbums"`
}

type CreatePhotoRequest struct {
	PhotoAlbumID uint   `json:"photoAlbumId"`
	URL          string `json:"uRL"`
	Caption      string `json:"caption"`
}

type FindPhotoOfUserResponse struct{}

type FindPhotoOfAlbumResponse struct{}

type CreateScrapRequest struct {
	UserRecipientID uint   `json:"userRecipientId"`
	Content         string `json:"content"`
}

type ScrapDto struct {
	ID                   uint      `json:"id"`
	Content              string    `json:"content"`
	AuthorID             uint      `json:"authorId"`
	AuthorName           string    `json:"authorName"`
	AuthorProfilePicture string    `json:"authorProfilePicture"`
	CreatedAt            time.Time `json:"createdAt"`
}

type FindScrapOfUserRecipientResponse struct {
	Scraps []ScrapDto `json:"scraps"`
}

type CreateTestimonialRequest struct {
	UserAuthorID    uint   `json:"userAuthorId"`
	UserRecipientID uint   `json:"userRecipientId"`
	Content         string `json:"content"`
}

type FindTestimonialOfUserResponse struct{}

type CreateUserRequest struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
}

type UpdateUserRequest struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
}

type UserDto struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	BirthDate   time.Time `json:"birthDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	FriendCount uint      `json:"friendCount"`
}

type UserProfileDto struct {
	ProfilePicture     *string `json:"profilePicture"`
	Description        *string `json:"description"`
	RelationshipStatus *string `json:"relationshipStatus"`
	Interests          *string `json:"interests"`
	About              *string `json:"about"`
	Religion           *string `json:"religion"`
	Humor              *string `json:"humor"`
	Fashion            *string `json:"fashion"`
	LivesWith          *string `json:"livesWith"`
	BirthCity          *string `json:"birthCity"`
	Website            *string `json:"website"`
	Passions           *string `json:"passions"`
	Sports             *string `json:"sports"`
	Hobbies            *string `json:"hobbies"`
	Books              *string `json:"books"`
	Music              *string `json:"music"`
	Movies             *string `json:"movies"`
	Cuisine            *string `json:"cuisine"`
}

type FriendDto struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"firstName"`
	BirthDate      time.Time `json:"birthDate"`
	ProfilePicture string    `json:"profilePicture"`
	FriendCount    uint      `json:"friendCount"`
}

type CommunityDto struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	CreatedBy       uint    `json:"firstName"`
	ProfilePicture  *string `json:"profilePicture"`
	MembershipCount uint    `json:"membershipCount"`
	OwnerID         uint    `json:"owner"`
}

type FindUserResponse struct {
	UserDto
}

type FindUserWithProfileByIdResponse struct {
	UserDto
	UserProfileDto
}

type FindManyUsersResponse struct {
	Users []UserDto `json:"users"`
}

type FindUserCommunitiesResponse struct {
	Communities []CommunityDto `json:"communities"`
}

type FindUserMembershipsResponse struct {
	Communities     []CommunityDto `json:"communities"`
	MembershipCount uint           `json:"membershipCount"`
}

type AmIMemberOfResponse struct {
	AmIMember bool `json:"amIMember"`
}

type AmIFriendOfResponse struct {
	FriendshipId        *uint `json:"friendshipId"`
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
	Topic   *Topic       `json:"topic"`
	Replies []TopicReply `json:"replies"`
}

type CreateTopicRequest struct {
	CommunityID uint   `json:"communityId"`
	Title       string `json:"title"`
	Text        string `json:"text"`
}

type ReplyTopicRequest struct {
	CommunityID uint   `json:"communityId"`
	Title       string `json:"title"`
	Text        string `json:"text"`
}
