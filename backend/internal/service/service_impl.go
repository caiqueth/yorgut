package service

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

var userIDCache = sync.Map{}

// Testimonial
func (t *TestimonialService) CreateTestimonial(testimonial dto.CreateTestimonialRequest) error {
	author, err := t.UserRepository.FindById(testimonial.UserAuthorID)
	if err != nil {
		return fmt.Errorf("author not found %v", testimonial.UserAuthorID)
	}

	recipient, err := t.UserRepository.FindById(testimonial.UserRecipientID)
	if err != nil {
		return fmt.Errorf("recipient not found %v", testimonial.UserRecipientID)
	}

	testimonialModel := model.NewTestimonial(author, recipient, testimonial.Content)
	err = t.TestimonialRepository.Create(testimonialModel)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestimonialService) DeleteTestimonial(id uint) error {
	err := t.TestimonialRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// User
func (u *UserService) DeleteUser(id uint) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindManyUsers() (*dto.FindManyUsersResponse, error) {
	users, err := u.userRepository.FindMany()
	if err != nil {
		return nil, err
	}

	userDtos := make([]dto.UserDto, len(users))

	for i, user := range users {
		userDtos[i] = dto.UserDto{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			BirthDate: user.BirthDate,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return &dto.FindManyUsersResponse{
		Users: userDtos,
	}, nil
}

func (u *UserService) CreateUser(firstName string, lastName string, email string, birthDate time.Time, profilePicture *string) error {
	err := u.userRepository.Create(model.NewUser(firstName, lastName, email, birthDate, profilePicture))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUser(user dto.UpdateUserRequest) error {
	return errors.New("not implemented")
	// err := u.userRepository.Update(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (u *UserService) FindUserById(id uint) (*dto.FindUserResponse, error) {
	userData, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	resp := dto.FindUserResponse{}
	resp.FromUser(userData)

	return &resp, nil

}

func (u *UserService) FindUserByEmail(email string) (*model.User, error) {
	return u.userRepository.FindByEmail(email)
}

func (u *UserService) FindUserByToken(token *model.Token) (*model.User, error) {
	return u.userRepository.FindByEmail(token.Email)
}

func (u *UserService) GetOrCreateByToken(token *model.Token) (*model.User, error) {
	user, err := u.FindUserByToken(token)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		createErr := u.CreateUser(token.Given_name, token.Family_name, token.Email, time.Time{}, &token.Picture)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create user: %w", createErr)
		}

		user, err = u.FindUserByToken(token)
		if err != nil {
			return nil, fmt.Errorf("user created but not found: %w", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("could not find user: %w", err)
}

func (u *UserService) FindUserIdFromEmail(email string) (uint, error) {
	if val, ok := userIDCache.Load(email); ok {
		return val.(uint), nil
	}
	return 0, fmt.Errorf("could not find userID for %v", email)
}

func (u *UserService) FindUserWithProfileById(id uint) (*dto.FindUserWithProfileByIdResponse, error) {
	userData, err := u.userRepository.FindUserWithProfile(id)
	if err != nil {
		return nil, err
	}

	resp := dto.FindUserWithProfileByIdResponse{}
	resp.FromUser(userData)

	return &resp, nil

}

func (u *UserService) FindUserOwnedCommunities(userId uint) (*dto.FindMyOwnedCommunitiesResponse, error) {
	communities, err := u.userRepository.FindUserOwnedCommunities(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindMyOwnedCommunitiesResponse{}
	resp.FromCommunities(communities)

	return resp, nil

}

func (u *UserService) FindUserPhotoAlbums(userId uint) (*dto.FindPhotoAlbumOfUserResponse, error) {
	photoAlbums, err := u.photoAlbumRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPhotoAlbumOfUserResponse{}
	resp.FromPhotoAlbums(photoAlbums)

	return resp, nil
}

func (u *UserService) FindPendingFriendshipsOfUser(userId uint) (*dto.FindPendingFriendshipOfUserResponse, error) {
	invites, err := u.friendshipRepository.FindPendingFriendshipsOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPendingFriendshipOfUserResponse{}
	resp.FromUsers(invites)

	return resp, nil
}

func (u *UserService) FindFriendsOfUser(userId uint) (*dto.FindFriendsOfUserResponse, error) {
	friends, err := u.friendshipRepository.FindFriendsOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindFriendsOfUserResponse{}
	resp.FromFriends(friends)

	return resp, nil
}

func (u *UserService) FindUserMemberships(userId uint) (*dto.FindUserMembershipsResponse, error) {
	memberships, err := u.communityMembershipRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindUserMembershipsResponse{}
	resp.FromMemberships(memberships)

	return resp, nil
}

func (u *UserService) IsUserMemberOf(communityId uint, userId uint) bool {
	membership, err := u.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return false
	}

	if membership == nil {
		return false
	}

	return true
}

func (u *UserService) UsersAreFriends(userId uint, targetUserId uint) bool {
	return u.friendshipRepository.UsersAreFriends(userId, targetUserId)
}

func (u *UserService) UsersFriendshipStatus(userId uint, targetUserId uint) *dto.AmIFriendOfResponse {
	if userId == targetUserId {
		return nil
	}

	friendships, err := u.friendshipRepository.FindFriendshipsBetweenUsers(userId, targetUserId)
	if err != nil {
		return nil
	}

	var friendshipId *uint
	usersAreFriends := false
	isPendingMyReply := false
	isPendingTheirReply := false

	for _, fsp := range friendships {
		friendshipId = &fsp.ID
		if fsp.Accepted {
			usersAreFriends = true
		} else {
			if fsp.UserRequestedID == userId {
				isPendingTheirReply = true
			}
			if fsp.UserRequestedID == targetUserId {
				isPendingMyReply = true
			}
		}
	}

	return &dto.AmIFriendOfResponse{
		FriendshipId:        friendshipId,
		AmIFriend:           usersAreFriends,
		IsPendingMyReply:    isPendingMyReply,
		IsPendingTheirReply: isPendingTheirReply,
	}

}

func (u *UserService) CountUserMemberships(userId uint) uint {
	count, err := u.communityMembershipRepository.CountOfUser(userId)
	if err != nil {
		return 0
	}
	return uint(count)
}

func (u *UserService) JoinCommunity(userId uint, communityId uint) error {
	community, err := u.communityRepository.FindById(communityId)
	if err != nil {
		return err
	}

	user, err := u.userRepository.FindById(userId)
	if err != nil {
		return err
	}

	membershipExisting, _ := u.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if membershipExisting != nil {
		return errors.New("user is already a member of that community")
	}

	membership := model.NewCommunityMembership(community, user)
	if err := u.communityMembershipRepository.Create(membership); err != nil {
		return err
	}

	if err := u.communityRepository.IncrementCommunityMembers(community, 1); err != nil {
		return err
	}

	return nil
}

func (u *UserService) LeaveCommunity(userId uint, communityId uint) error {
	membership, err := u.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return err
	}

	community, err := u.communityRepository.FindById(communityId)
	if err != nil {
		return err
	}

	if err := u.communityMembershipRepository.Delete(membership.ID); err != nil {
		return err
	}

	if err := u.communityRepository.IncrementCommunityMembers(community, -1); err != nil {
		return err
	}

	return nil
}

// Scrap
func (s *ScrapService) CreateScrap(userAuthorID uint, userRecipientID uint, content string) error {
	author, err := s.userRepository.FindById(userAuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	recipient, err := s.userRepository.FindById(userRecipientID)
	if err != nil {
		return errors.New("recipient not found")
	}

	scrapModel := model.NewScrap(author, recipient, content)

	if err := s.scrapRepository.Create(scrapModel); err != nil {
		return err
	}

	return nil
}

func (s *ScrapService) DeleteScrap(id uint) error {
	err := s.scrapRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScrapService) FindScrapOfUserRecipient(userId uint) (*dto.FindScrapOfUserRecipientResponse, error) {
	scraps, err := s.scrapRepository.FindOfUserRecipient(userId)
	if err != nil {
		return nil, err
	}

	response := dto.FindScrapOfUserRecipientResponse{
		Scraps: []dto.ScrapDto{},
	}

	for _, scrap := range scraps {
		response.Scraps = append(response.Scraps, dto.ScrapDto{
			ID:                   scrap.ID,
			Content:              scrap.Content,
			AuthorID:             scrap.UserAuthorID,
			AuthorName:           scrap.UserAuthor.FirstName,
			AuthorProfilePicture: *scrap.UserAuthor.Profile.ProfilePicture,
			CreatedAt:            scrap.CreatedAt,
		})
	}

	return &response, nil
}

// Photo
func (p *PhotoService) CreatePhoto(photo dto.CreatePhotoRequest) error {
	photoAlbum, err := p.photoAlbumRepository.FindById(photo.PhotoAlbumID)
	if err != nil {
		return errors.New("photo album not found")
	}

	photoModel := model.NewPhoto(photoAlbum, photo.URL, photo.Caption)

	err = p.photoRepository.Create(photoModel)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoService) DeletePhoto(id uint) error {
	err := p.photoRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// PhotoAlbum
func (p *PhotoAlbumService) CreatePhotoAlbum(photoAlbum dto.CreatePhotoAlbumRequest) error {
	owner, err := p.userRepository.FindById(photoAlbum.UserID)
	if err != nil {
		return errors.New("user not found")
	}
	photoAlbumModel := model.NewPhotoAlbum(owner, photoAlbum.Name)

	err = p.photoAlbumRepository.Create(photoAlbumModel)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoAlbumService) DeletePhotoAlbum(id uint) error {
	err := p.photoAlbumRepository.Delete(id)
	if err != nil {
		return err
	}
	return err
}

func (p *PhotoAlbumService) FindPhotoAlbumOfUser(userId uint) (*dto.FindPhotoAlbumOfUserResponse, error) {
	photoAlbums, err := p.photoAlbumRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPhotoAlbumOfUserResponse{}
	resp.FromPhotoAlbums(photoAlbums)

	return resp, nil
}

// Friendship
func (f *FriendshipService) CreateFriendship(userRequestedID uint, userRequesteeID uint) error {
	userRequested, err := f.userRepository.FindById(userRequestedID)
	if err != nil {
		return errors.New("could not find user that requested the friendship")
	}

	userRequestee, err := f.userRepository.FindById(userRequesteeID)
	if err != nil {
		return errors.New("could not find target user (requestee) of the friendship")
	}

	if f.friendshipRepository.UsersHavePendingFriendshipInvites(userRequestedID, userRequesteeID) {
		return errors.New("there is already a friendship or request between these users")
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.friendshipRepository.Create(friendshipModel)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) ReplyFriendship(friendshipId uint, accepted bool) error {
	friendship, err := f.friendshipRepository.FindPendingFriendshipWithUsersById(friendshipId)
	if err != nil {
		return err
	}

	if !accepted {
		f.friendshipRepository.Delete(friendshipId)
		return nil
	}

	now := time.Now()
	friendship.Accepted = accepted
	friendship.AcceptedAt = &now

	if err := f.friendshipRepository.Update(friendship); err != nil {
		return err
	}

	if accepted {
		f.userRepository.IncrementFriendCount(friendship.UserRequested, 1)
		f.userRepository.IncrementFriendCount(friendship.UserRequestee, 1)
	}

	return nil
}

func (f *FriendshipService) ShouldThisUserBeReplying(friendshipId uint, userId uint) bool {
	friendship, err := f.friendshipRepository.FindById(friendshipId)
	if err != nil {
		return false
	}

	if friendship.UserRequesteeID == userId {
		return true
	}

	return false
}

func (f *FriendshipService) ShouldThisUserBeDeleting(friendshipId uint, userId uint) bool {
	friendship, err := f.friendshipRepository.FindById(friendshipId)
	if err != nil {
		return false
	}

	if friendship.UserRequesteeID == userId || friendship.UserRequestedID == userId {
		return true
	}

	return false
}

func (f *FriendshipService) DeleteFriendship(friendshipId uint) error {
	friendship, err := f.friendshipRepository.FindByIdWithUsers(friendshipId)
	if err != nil {
		return errors.New("friendship not found")
	}

	if err := f.friendshipRepository.Delete(friendship.ID); err != nil {
		return err
	}

	f.userRepository.IncrementFriendCount(friendship.UserRequested, -1)
	f.userRepository.IncrementFriendCount(friendship.UserRequestee, -1)

	return nil
}

func (f *FriendshipService) FindFriendshipByUsersId(userId1 uint, userId2 uint) *model.Friendship {
	friendship, err := f.friendshipRepository.FindActiveAcceptedFriends(userId1, userId2)
	if err != nil {
		return nil
	}
	return friendship
}

// Community
func (c *CommunityService) CreateCommunity(request dto.CreateCommunityRequest) error {
	creator, err := c.userRepository.FindById(request.CreatedByID)
	if err != nil {
		return errors.New("could not find specified community creator")
	}

	communityModel := model.NewCommunity(request.Name, request.Category, creator)
	err = c.communityRepository.Create(communityModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) DeleteCommunity(id uint) error {
	err := c.communityRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) FindCommunityById(id uint) (*dto.FindCommunityByIdResponse, error) {
	response := dto.FindCommunityByIdResponse{}
	community, err := c.communityRepository.FindByIdEager(id)
	if err != nil {
		return &response, errors.New("could not find specified community id")
	}

	response.Name = community.Name
	response.CreatedByID = community.CreatedByID
	response.OwnerFirstName = community.Owner.FirstName
	response.OwnerID = community.Owner.ID
	response.CreatedAt = community.CreatedAt

	return &response, nil
}

func (c *CommunityService) FindCommunityProfileById(id uint) (dto.FindCommunityProfileByIdResponse, error) {
	profile, err := c.communityProfileRepository.FindById(id)
	if err != nil {
		return dto.FindCommunityProfileByIdResponse{}, errors.New("could not find specified community profile id")
	}

	return dto.FindCommunityProfileByIdResponse{
		Description:    profile.Description,
		ProfilePicture: profile.ProfilePicture,
	}, nil
}

func (c *CommunityService) FindManyCommunities() (*dto.FindManyCommunitiesResponse, error) {
	communities, err := c.communityRepository.FindCommunitiesWithProfiles()
	if err != nil {
		return nil, err
	}

	resp := dto.FindManyCommunitiesResponse{}
	resp.FromCommunities(communities)

	return &resp, nil
}

func (c *CommunityService) FindCommunityMembers(communityId uint) (*dto.FindCommunityMembersResponse, error) {
	memberships, err := c.membershipRepository.FindMembersOfCommunity(communityId)
	if err != nil {
		return nil, err
	}

	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return nil, err
	}

	members := []dto.CommunityMemberEntry{}
	for _, membership := range memberships {
		members = append(members, dto.CommunityMemberEntry{
			ID:             membership.UserID,
			FirstName:      membership.User.FirstName,
			ProfilePicture: *membership.User.Profile.ProfilePicture,
			FriendCount:    membership.User.FriendCount,
		})
	}

	return &dto.FindCommunityMembersResponse{
		MemberCount: community.MemberCount,
		Members:     members,
	}, nil
}

// Community Profile
func (c *CommunityProfileService) CreateCommunityProfile(request dto.CreateCommunityProfileRequest) error {
	community, err := c.communityRepository.FindById(request.CommunityID)
	if err != nil {
		return errors.New("could not find specified community")
	}

	picture, err := helpers.FileToBytes(request.Picture)
	if err != nil {
		return err
	}

	storedFileName, err := helpers.UploadPicture(picture)
	if err != nil {
		return err
	}

	communityProfileModel := model.NewCommunityProfile(community)
	model.AddCommunityProfilePicture(communityProfileModel, *storedFileName)

	err = c.communityProfileRepository.Create(communityProfileModel)
	if err != nil {
		return err
	}
	return nil
}

// CommunityMembership
func (c *CommunityMembershipService) CreateCommunityMembership(communityMembership dto.CreateCommunityMembershipRequest) error {
	community, err := c.communityRepository.FindById(communityMembership.CommunityID)
	if err != nil {
		return errors.New("could not find specified community")
	}

	user, err := c.userRepository.FindById(communityMembership.UserID)
	if err != nil {
		return errors.New("could not find community owner")
	}

	communityMembershipModel := model.NewCommunityMembership(community, user)
	err = c.communityMembershipRepository.Create(communityMembershipModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityMembershipService) DeleteCommunityMembership(id uint) error {
	err := c.communityMembershipRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityMembershipService) FindCommunityMembership(communityId uint, userId uint) (*uint, error) {
	membership, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return nil, err
	}
	return &membership.ID, nil
}

func (ts *TopicService) FindTopics(communityId uint, limit int) (*dto.TopicsResponse, error) {
	response := dto.TopicsResponse{}
	response.Topics = []dto.TopicSummary{}

	topics, err := ts.topicRepository.FindMostRecentOfCommunity(communityId, limit)

	for _, t := range topics {
		count, replies, err := ts.topicRepository.FindRepliesOfTopic(t.ID)
		if err != nil {
			return nil, err
		}

		topicSummary := dto.TopicSummary{
			ID:                   t.ID,
			AuthorProfilePicture: *t.Author.Profile.ProfilePicture,
			Title:                t.Title,
			Replies:              count,
		}

		if count != 0 {
			topicSummary.LastReplyPreview = replies[len(replies)-1].Text
			topicSummary.LastReplyDate = replies[len(replies)-1].CreatedAt
		}

		response.Topics = append(response.Topics, topicSummary)
	}

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (ts *TopicService) FindTopicReplies(topicId uint) (*dto.TopicRepliesResponse, error) {
	response := dto.TopicRepliesResponse{}
	response.Replies = []dto.TopicReply{}

	topic, err := ts.topicRepository.FindByIdWithAuthorProfile(topicId)
	if err != nil {
		return nil, err
	}

	count, replies, err := ts.topicRepository.FindRepliesOfTopic(topicId)
	if err != nil {
		return nil, err
	}

	response.Topic = &dto.Topic{
		ID:                   topic.ID,
		AuthorID:             topic.AuthorID,
		AuthorFirstName:      topic.Author.FirstName,
		AuthorProfilePicture: *topic.Author.Profile.ProfilePicture,
		Title:                topic.Title,
		Text:                 topic.Text,
		CreatedAt:            topic.CreatedAt,
		Replies:              count,
	}

	for _, r := range replies {
		response.Replies = append(response.Replies, dto.TopicReply{
			ID:                   r.ID,
			AuthorID:             r.AuthorID,
			AuthorFirstName:      r.Author.FirstName,
			AuthorProfilePicture: *r.Author.Profile.ProfilePicture,
			Title:                r.Title,
			Text:                 r.Text,
			CreatedAt:            r.CreatedAt,
		})
	}

	return &response, nil
}

func (ts *TopicService) CreateTopic(authorId uint, communityId uint, title string, text string) error {
	if err := ts.topicRepository.Create(model.NewTopic(authorId, communityId, nil, title, text)); err != nil {
		return err
	}
	return nil
}

func (ts *TopicService) ReplyTopic(authorId uint, communityId uint, parentTopicId uint, title string, text string) error {
	parentTopic, err := ts.topicRepository.FindById(parentTopicId)
	if err != nil {
		return fmt.Errorf("could not find topic with id %v", parentTopicId)
	}

	if err := ts.topicRepository.Create(model.NewTopic(authorId, communityId, &parentTopic.ID, title, text)); err != nil {
		return err
	}

	return nil
}
