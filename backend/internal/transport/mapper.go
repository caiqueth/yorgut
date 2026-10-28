package transport

import (
	"slices"
	"strings"
	"time"
	"yorgut/internal/model"
)

// mediaURL turns a stored object key into the URL clients load it from. The
// API and the seed set it once at startup (SetMediaURL); until then keys are
// returned as they are, which is what the service tests compare against.
var mediaURL = func(key string) string { return key }

// SetMediaURL is called once at startup, before serving requests, usually
// with the object store's URL method.
func SetMediaURL(resolve func(key string) string) {
	mediaURL = resolve
}

// pictureURL resolves a stored picture. Google profile pictures are stored as
// absolute URLs and pass through unchanged.
func pictureURL(picture *string) *string {
	if picture == nil {
		return nil
	}
	if strings.HasPrefix(*picture, "http://") || strings.HasPrefix(*picture, "https://") {
		return picture
	}
	url := mediaURL(*picture)
	return &url
}

func newPagedResponse(page int, pageSize int, count int64) PagedResponse {
	return PagedResponse{Page: page, PageSize: pageSize, Count: count}
}

// Mappers
func (r *FindUserResponse) FromUser(user *model.User) {
	r.UserItem.FromUser(user)
}

func (r *FindMeResponse) FromUser(user *model.User) {
	r.UserItem.FromUser(user)
	r.Email = user.Email
	r.Privacy = PrivacySettings{
		RecordVisits:        user.RecordVisits,
		PublicAlbum:         user.PublicAlbum,
		PublicCommunities:   user.PublicCommunities,
		PublicScraps:        user.PublicScraps,
		PublicTestimonials:  user.PublicTestimonials,
		PublicBirthday:      user.PublicBirthday,
		HiddenProfileFields: SplitProfileFields(user.HiddenProfileFields),
	}
	r.Theme = user.Theme
}

// SplitProfileFields reads model.User.HiddenProfileFields.
func SplitProfileFields(fields string) []string {
	keys := []string{}
	for _, key := range strings.Split(fields, ",") {
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

// FromUserFor fills the profile as viewer sees it. friend is true for the
// user themselves and their friends, who see everything; anyone else sees
// only what the user opened (see model.User.Public*).
func (r *FindUserWithProfileByIdResponse) FromUserFor(user *model.User, friend bool) {
	r.FromUser(user)
	r.HiddenFields = []string{}
	if friend {
		r.Visibility = ProfileVisibility{Album: true, Communities: true, Scraps: true, Testimonials: true}
		return
	}
	r.Visibility = ProfileVisibility{
		Album:        user.PublicAlbum,
		Communities:  user.PublicCommunities,
		Scraps:       user.PublicScraps,
		Testimonials: user.PublicTestimonials,
	}
	if !user.PublicBirthday {
		r.BirthDate = nil
		r.HiddenFields = append(r.HiddenFields, "birthDate")
	}
	fields := r.profileFields()
	for _, key := range SplitProfileFields(user.HiddenProfileFields) {
		if field, ok := fields[key]; ok {
			*field = nil
			r.HiddenFields = append(r.HiddenFields, key)
		}
	}
}

// profileFields maps each of ProfileFieldKeys to its field.
func (r *FindUserWithProfileByIdResponse) profileFields() map[string]**string {
	return map[string]**string{
		"description":        &r.Description,
		"relationshipStatus": &r.RelationshipStatus,
		"interests":          &r.Interests,
		"about":              &r.About,
		"religion":           &r.Religion,
		"humor":              &r.Humor,
		"fashion":            &r.Fashion,
		"livesWith":          &r.LivesWith,
		"birthCity":          &r.BirthCity,
		"website":            &r.Website,
		"passions":           &r.Passions,
		"sports":             &r.Sports,
		"hobbies":            &r.Hobbies,
		"books":              &r.Books,
		"music":              &r.Music,
		"movies":             &r.Movies,
		"cuisine":            &r.Cuisine,
	}
}

// FromUser fills the whole profile, as the user themselves sees it.
func (r *FindUserWithProfileByIdResponse) FromUser(user *model.User) {
	r.UserItem.FromUser(user)
	r.HiddenFields = []string{}
	r.Visibility = ProfileVisibility{Album: true, Communities: true, Scraps: true, Testimonials: true}
	// Nil for users created at first login, until they set it.
	if user.BirthDate != nil {
		birthDate := *user.BirthDate
		r.BirthDate = &birthDate
	}
	r.ProfilePicture = profilePictureOf(user)
	if p := user.Profile; p != nil {
		r.Description = p.Description
		r.RelationshipStatus = p.RelationshipStatus
		r.Interests = p.Interests
		r.About = p.About
		r.Religion = p.Religion
		r.Humor = p.Humor
		r.Fashion = p.Fashion
		r.LivesWith = p.LivesWith
		r.BirthCity = p.BirthCity
		r.Website = p.Website
		r.Passions = p.Passions
		r.Sports = p.Sports
		r.Hobbies = p.Hobbies
		r.Books = p.Books
		r.Music = p.Music
		r.Movies = p.Movies
		r.Cuisine = p.Cuisine
	}
}

func (r *FindMyOwnedCommunitiesResponse) FromCommunities(communities []model.Community, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Communities = make([]MyOwnedCommunityItem, len(communities))
	for i, c := range communities {
		r.Communities[i] = MyOwnedCommunityItem{
			ID:              c.ID,
			Name:            c.Name,
			MembershipCount: c.MemberCount,
		}
		if c.Profile != nil {
			r.Communities[i].ProfilePicture = pictureURL(c.Profile.ProfilePicture)
		}
	}
}

func (r *FindUserMembershipsResponse) FromCommunities(communities []model.Community, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Communities = make([]CommunityItem, len(communities))
	for i := range communities {
		r.Communities[i] = communityItemOf(&communities[i])
	}
}

func (r *FindFriendsOfUserResponse) FromFriends(friends []model.User, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Friends = make([]FriendItem, len(friends))

	for i, f := range friends {
		r.Friends[i] = FriendItem{
			ID:             f.ID,
			FirstName:      f.FirstName,
			ProfilePicture: profilePictureOf(&f),
			FriendCount:    f.FriendCount,
		}
	}
}

func (r *FindPendingFriendshipOfUserResponse) FromUsers(users []model.Friendship) {
	r.Count = 0
	r.PendingFriendships = make([]PendingFriendship, len(users))

	for i, f := range users {
		r.Count += 1
		r.PendingFriendships[i] = PendingFriendship{
			ID:              f.ID,
			UserRequestedID: f.UserRequested.ID,
			FirstName:       f.UserRequested.FirstName,
			LastName:        f.UserRequested.LastName,
			ProfilePicture:  profilePictureOf(f.UserRequested),
		}
	}
}

func (r *FindManyCommunitiesResponse) FromCommunities(communities []model.Community) {
	r.Communities = make([]CommunityItem, len(communities))
	for i := range communities {
		r.Communities[i] = communityItemOf(&communities[i])
	}
}

func communityItemOf(c *model.Community) CommunityItem {
	return CommunityItem{
		ID:              c.ID,
		Name:            c.Name,
		CreatedBy:       c.CreatedByID,
		ProfilePicture:  communityPictureOf(c),
		MembershipCount: c.MemberCount,
		OwnerID:         c.OwnerID,
	}
}

func (r *SearchUsersResponse) FromUsers(users []*model.User, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Users = make([]*UserSearchItem, len(users))
	for i, u := range users {
		r.Users[i] = &UserSearchItem{
			ID:             u.ID,
			FirstName:      u.FirstName,
			ProfilePicture: profilePictureOf(u),
			FriendCount:    u.FriendCount,
		}
	}
}

func (d *UserItem) FromUser(user *model.User) {
	d.ID = user.ID
	d.FirstName = user.FirstName
	d.LastName = user.LastName
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt
	d.FriendCount = user.FriendCount
}

func (d *ScrapItem) FromScrap(scrap *model.Scrap) {
	d.ID = scrap.ID
	d.Content = scrap.Content
	d.AuthorID = scrap.UserAuthorID
	d.AuthorName = scrap.UserAuthor.FirstName
	d.AuthorProfilePicture = profilePictureOf(scrap.UserAuthor)
	d.CreatedAt = scrap.CreatedAt
}

func (r *FindScrapOfUserRecipientResponse) FromScraps(scraps []model.Scrap, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Scraps = make([]ScrapItem, len(scraps))
	for i := range scraps {
		r.Scraps[i].FromScrap(&scraps[i])
	}
}

func (r *FindCommunityByIdResponse) FromCommunity(community *model.Community) {
	r.Name = community.Name
	r.Category = CommunityCategory(community.Category)
	r.CreatedByID = community.CreatedByID
	r.OwnerFirstName = community.Owner.FirstName
	r.OwnerID = community.Owner.ID
	r.Private = community.Private
	r.CreatedAt = community.CreatedAt
}

// FromFriends maps a page of friends; memberIds and invitedIds are the ones
// among them already in the community or already invited to it.
func (r *FindInvitableFriendsResponse) FromFriends(friends []model.User, memberIds []uint, invitedIds []uint, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Friends = make([]InvitableFriendItem, len(friends))
	for i := range friends {
		f := &friends[i]
		r.Friends[i] = InvitableFriendItem{
			ID:             f.ID,
			FirstName:      f.FirstName,
			ProfilePicture: profilePictureOf(f),
			FriendCount:    f.FriendCount,
			Member:         slices.Contains(memberIds, f.ID),
			Invited:        slices.Contains(invitedIds, f.ID),
		}
	}
}

func (r *FindMyCommunityInvitationsResponse) FromInvitations(invitations []model.CommunityInvitation) {
	r.Invitations = make([]CommunityInvitationItem, 0, len(invitations))
	for _, inv := range invitations {
		if inv.Community == nil || inv.Inviter == nil {
			// The community or the inviter was deleted after inviting.
			continue
		}
		r.Invitations = append(r.Invitations, CommunityInvitationItem{
			ID:                    inv.ID,
			CommunityID:           inv.CommunityID,
			CommunityName:         inv.Community.Name,
			CommunityPicture:      communityPictureOf(inv.Community),
			CommunityPrivate:      inv.Community.Private,
			InviterID:             inv.InviterID,
			InviterFirstName:      inv.Inviter.FirstName,
			InviterProfilePicture: profilePictureOf(inv.Inviter),
			CreatedAt:             inv.CreatedAt,
		})
	}
}

func (r *FindJoinRequestsResponse) FromJoinRequests(requests []model.CommunityJoinRequest) {
	r.Requests = make([]JoinRequestItem, 0, len(requests))
	for _, req := range requests {
		if req.User == nil {
			// The user was deleted after asking.
			continue
		}
		r.Requests = append(r.Requests, JoinRequestItem{
			UserID:         req.UserID,
			FirstName:      req.User.FirstName,
			LastName:       req.User.LastName,
			ProfilePicture: profilePictureOf(req.User),
			FriendCount:    req.User.FriendCount,
			CreatedAt:      req.CreatedAt,
		})
	}
}

func (r *FindCommunityProfileByIdResponse) FromCommunityProfile(profile *model.CommunityProfile) {
	r.Description = profile.Description
	r.ProfilePicture = pictureURL(profile.ProfilePicture)
}

func (r *FindCommunityMembersResponse) FromMembers(members []model.User, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Members = make([]CommunityMemberEntry, len(members))
	for i := range members {
		m := &members[i]
		r.Members[i] = CommunityMemberEntry{
			ID:             m.ID,
			FirstName:      m.FirstName,
			ProfilePicture: profilePictureOf(m),
			FriendCount:    m.FriendCount,
		}
	}
}

// FromTopic fills the summary; lastReply is nil for topics without replies.
func (d *TopicSummary) FromTopic(topic *model.Topic, replyCount int64, lastReply *model.Topic) {
	d.ID = topic.ID
	d.AuthorID = topic.AuthorID
	d.AuthorProfilePicture = profilePictureOf(topic.Author)
	d.Title = topic.Title
	d.Replies = replyCount
	if lastReply != nil {
		d.LastReplyPreview = lastReply.Text
		d.LastReplyDate = lastReply.CreatedAt
	}
}

func (d *Topic) FromTopic(topic *model.Topic, replyCount int64) {
	d.ID = topic.ID
	d.AuthorID = topic.AuthorID
	d.AuthorFirstName = topic.Author.FirstName
	d.AuthorProfilePicture = profilePictureOf(topic.Author)
	d.Title = topic.Title
	d.Text = topic.Text
	d.CreatedAt = topic.CreatedAt
	d.Replies = replyCount
}

func (d *TopicReply) FromTopic(reply *model.Topic) {
	d.ID = reply.ID
	d.AuthorID = reply.AuthorID
	d.AuthorFirstName = reply.Author.FirstName
	d.AuthorProfilePicture = profilePictureOf(reply.Author)
	d.Title = reply.Title
	d.Text = reply.Text
	d.CreatedAt = reply.CreatedAt
}

func (d *UserReviewItem) FromUserReview(review *model.UserReview) {
	d.Cool = review.Cool
	d.Sexy = review.Sexy
	d.Trustworthy = review.Trustworthy
}

func (d *TestimonialItem) FromTestimonial(testimonial *model.Testimonial) {
	d.ID = testimonial.ID
	d.Content = testimonial.Content
	d.AuthorID = testimonial.UserAuthorID
	d.AuthorName = testimonial.UserAuthor.FirstName
	d.AuthorProfilePicture = profilePictureOf(testimonial.UserAuthor)
	d.Approved = testimonial.Approved
	d.CreatedAt = testimonial.CreatedAt
}

func testimonialItems(testimonials []model.Testimonial) []TestimonialItem {
	items := make([]TestimonialItem, len(testimonials))
	for i := range testimonials {
		items[i].FromTestimonial(&testimonials[i])
	}
	return items
}

func (r *FindTestimonialsOfUserResponse) FromTestimonials(approved []model.Testimonial, pending []model.Testimonial, canWrite bool, page int, pageSize int, count int64) {
	r.PagedResponse = PagedResponse{
		Page:     page,
		PageSize: pageSize,
		Count:    count,
	}
	r.Testimonials = testimonialItems(approved)
	r.Pending = testimonialItems(pending)
	r.CanWrite = canWrite
}

func (d *CommunitySearchItem) FromCommunity(community *model.Community) {
	d.ID = community.ID
	d.Name = community.Name
	d.Category = CommunityCategory(community.Category)
	d.MemberCount = community.MemberCount
	d.Private = community.Private
	if community.Profile != nil {
		d.ProfilePicture = pictureURL(community.Profile.ProfilePicture)
	}
}

func (r *SearchCommunitiesResponse) FromCommunities(communities []model.Community, page int, pageSize int, count int64) {
	r.PagedResponse = PagedResponse{
		Page:     page,
		PageSize: pageSize,
		Count:    count,
	}
	r.Communities = make([]CommunitySearchItem, len(communities))
	for i := range communities {
		r.Communities[i].FromCommunity(&communities[i])
	}
}

func (d *PhotoItem) FromPhoto(photo *model.Photo) {
	d.ID = photo.ID
	d.URL = mediaURL(photo.URL)
	d.Caption = photo.Caption
	d.CreatedAt = photo.CreatedAt
}

func (r *FindPhotosOfUserResponse) FromPhotos(photos []model.Photo, maxPhotos int, canUpload bool) {
	r.Photos = make([]PhotoItem, len(photos))
	for i := range photos {
		r.Photos[i].FromPhoto(&photos[i])
	}
	r.Count = len(photos)
	r.MaxPhotos = maxPhotos
	r.CanUpload = canUpload
}

// profilePictureOf returns the URL of the user's profile picture, or nil when the user,
// its profile or the picture is missing.
func profilePictureOf(u *model.User) *string {
	if u == nil || u.Profile == nil {
		return nil
	}
	return pictureURL(u.Profile.ProfilePicture)
}

// communityPictureOf is the community counterpart of profilePictureOf.
func communityPictureOf(c *model.Community) *string {
	if c == nil || c.Profile == nil {
		return nil
	}
	return pictureURL(c.Profile.ProfilePicture)
}

func RecentVisitorsFromUsers(users []model.User) []RecentVisitor {
	visitors := make([]RecentVisitor, len(users))
	for i, u := range users {
		visitors[i] = RecentVisitor{UserID: u.ID, FirstName: u.FirstName}
	}
	return visitors
}

func (d *BirthdayItem) FromUser(user *model.User, nextBirthday time.Time, daysUntil int) {
	d.ID = user.ID
	d.FirstName = user.FirstName
	d.ProfilePicture = profilePictureOf(user)
	d.BirthDate = *user.BirthDate
	d.NextBirthday = nextBirthday
	d.DaysUntil = daysUntil
}

func (r *FindUserFanListResponse) FromFans(fans []model.User, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Fans = make([]FanItem, len(fans))
	for i := range fans {
		f := &fans[i]
		r.Fans[i] = FanItem{
			ID:             f.ID,
			FirstName:      f.FirstName,
			ProfilePicture: profilePictureOf(f),
			FriendCount:    f.FriendCount,
		}
	}
}

// FromFriends fills the page; reviews are the logged user's reviews of these
// friends and fanOf the ids of the friends they are a fan of.
func (r *FindMyFriendsResponse) FromFriends(friends []model.User, reviews []model.UserReview, fanOf []uint, page int, pageSize int, count int64) {
	reviewOf := make(map[uint]model.UserReview, len(reviews))
	for _, review := range reviews {
		reviewOf[review.ReviewedID] = review
	}
	isFan := make(map[uint]bool, len(fanOf))
	for _, id := range fanOf {
		isFan[id] = true
	}

	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Friends = make([]MyFriendItem, len(friends))
	for i := range friends {
		f := &friends[i]
		review := reviewOf[f.ID]
		r.Friends[i] = MyFriendItem{
			FriendItem: FriendItem{
				ID:             f.ID,
				FirstName:      f.FirstName,
				ProfilePicture: profilePictureOf(f),
				FriendCount:    f.FriendCount,
			},
			MyReview: UserReviewItem{Cool: review.Cool, Sexy: review.Sexy, Trustworthy: review.Trustworthy},
			AmIFan:   isFan[f.ID],
		}
	}
}

// FromPolls maps a page of polls. voteCounts holds each poll's votes,
// votedIds the polls the logged user voted in; closed says whether a poll is
// closed (the service decides, as it knows the time).
func (r *PollsResponse) FromPolls(polls []model.Poll, voteCounts map[uint]int64, votedIds []uint, closed func(*model.Poll) bool, page int, pageSize int, count int64) {
	r.PagedResponse = newPagedResponse(page, pageSize, count)
	r.Polls = make([]PollSummary, len(polls))
	for i := range polls {
		p := &polls[i]
		r.Polls[i] = PollSummary{
			ID:        p.ID,
			Question:  p.Question,
			AuthorID:  p.AuthorID,
			CreatedAt: p.CreatedAt,
			ClosesAt:  p.ClosesAt,
			Closed:    closed(p),
			VoteCount: voteCounts[p.ID],
			Voted:     slices.Contains(votedIds, p.ID),
		}
		if p.Author != nil {
			r.Polls[i].AuthorFirstName = p.Author.FirstName
		}
	}
}

// FromPoll maps a poll and its options. optionVotes holds each option's
// votes and is only read when showResults is true; otherwise every option's
// Votes stays null.
func (r *PollResponse) FromPoll(poll *model.Poll, options []model.PollOption, optionVotes map[uint]int64, showResults bool) {
	r.ID = poll.ID
	r.CommunityID = poll.CommunityID
	r.Question = poll.Question
	r.Description = poll.Description
	r.AuthorID = poll.AuthorID
	if poll.Author != nil {
		r.AuthorFirstName = poll.Author.FirstName
	}
	r.AuthorProfilePicture = profilePictureOf(poll.Author)
	r.CreatedAt = poll.CreatedAt
	r.ClosesAt = poll.ClosesAt
	r.VoteCount = 0
	r.Options = make([]PollOptionItem, len(options))
	for i, o := range options {
		r.VoteCount += optionVotes[o.ID]
		r.Options[i] = PollOptionItem{ID: o.ID, Text: o.Text}
		if showResults {
			votes := optionVotes[o.ID]
			r.Options[i].Votes = &votes
		}
	}
}
