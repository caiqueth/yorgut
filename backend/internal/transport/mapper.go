package transport

import "yorgut/internal/model"

// Mappers
func (r *FindUserResponse) FromUser(user *model.User) {
	r.ID = user.ID
	r.Email = user.Email
	r.FirstName = user.FirstName
	r.LastName = user.LastName
	r.BirthDate = user.BirthDate
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
	r.FriendCount = user.FriendCount
}

func (r *FindUserWithProfileByIdResponse) FromUser(user *model.User) {
	r.ID = user.ID
	r.Email = user.Email
	r.FirstName = user.FirstName
	r.LastName = user.LastName
	r.BirthDate = user.BirthDate
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
	r.ProfilePicture = user.Profile.ProfilePicture
	r.Description = user.Profile.Description
	r.RelationshipStatus = user.Profile.RelationshipStatus
	r.Interests = user.Profile.Interests
	r.About = user.Profile.About
	r.Religion = user.Profile.Religion
	r.Humor = user.Profile.Humor
	r.Fashion = user.Profile.Fashion
	r.LivesWith = user.Profile.LivesWith
	r.BirthCity = user.Profile.BirthCity
	r.Website = user.Profile.Website
	r.Passions = user.Profile.Passions
	r.Sports = user.Profile.Sports
	r.Hobbies = user.Profile.Hobbies
	r.Books = user.Profile.Books
	r.Music = user.Profile.Music
	r.Movies = user.Profile.Movies
	r.Cuisine = user.Profile.Cuisine
}

func (r *FindMyOwnedCommunitiesResponse) FromCommunities(communities []model.Community) {
	r.Communities = make([]MyOwnedCommunityItem, len(communities))
	for i, c := range communities {
		r.Communities[i] = MyOwnedCommunityItem{
			ID:              c.ID,
			Name:            c.Name,
			MembershipCount: c.MemberCount,
		}
		if c.Profile != nil {
			r.Communities[i].ProfilePicture = c.Profile.ProfilePicture
		}
	}
}

func (r *FindUserMembershipsResponse) FromMemberships(memberships []model.CommunityMembership) {
	r.Communities = make([]CommunityItem, len(memberships))
	r.MembershipCount = 0
	for i, m := range memberships {
		r.Communities[i] = CommunityItem{
			ID:              m.Community.ID,
			Name:            m.Community.Name,
			CreatedBy:       m.Community.CreatedByID,
			ProfilePicture:  m.Community.Profile.ProfilePicture,
			MembershipCount: m.Community.MemberCount,
		}
	}
}

func (r *FindFriendshipOfUserResponse) FromFriendships(friendships []model.Friendship) {
	r.Friendships = make([]FriendshipItem, len(friendships))

	for i, f := range friendships {
		r.Friendships[i] = FriendshipItem{
			UserRequestedID: f.UserRequestedID,
			UserRequesteeID: f.UserRequesteeID,
			CreatedAt:       f.CreatedAt,
		}
	}
}

func (r *FindFriendsOfUserResponse) FromFriends(friends []model.User) {
	r.Friends = make([]FriendItem, len(friends))

	for i, f := range friends {
		r.Friends[i] = FriendItem{
			ID:             f.ID,
			FirstName:      f.FirstName,
			BirthDate:      f.BirthDate,
			ProfilePicture: *f.Profile.ProfilePicture,
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
			ProfilePicture:  *f.UserRequested.Profile.ProfilePicture,
		}
	}
}

func (r *FindManyCommunitiesResponse) FromCommunities(communities []model.Community) {
	r.Communities = make([]CommunityItem, len(communities))
	for i, c := range communities {
		r.Communities[i] = CommunityItem{
			ID:             c.ID,
			Name:           c.Name,
			CreatedBy:      c.CreatedByID,
			ProfilePicture: c.Profile.ProfilePicture,
			OwnerID:        c.OwnerID,
		}
	}
}

func (r *SearchUsersResponse) FromUsers(users []*model.User, page int, pageSize int, count int64) {
	r.PagedResponse = PagedResponse{
		Page:     page,
		PageSize: pageSize,
		Count:    count,
	}

	r.Users = make([]*UserSearchItem, len(users))
	for i, u := range users {
		r.Users[i] = &UserSearchItem{
			ID:             u.ID,
			FirstName:      u.FirstName,
			BirthDate:      u.BirthDate,
			ProfilePicture: u.Profile.ProfilePicture,
			FriendCount:    u.FriendCount,
		}
	}
}

func (d *UserItem) FromUser(user *model.User) {
	d.ID = user.ID
	d.Email = user.Email
	d.FirstName = user.FirstName
	d.LastName = user.LastName
	d.BirthDate = user.BirthDate
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt
	d.FriendCount = user.FriendCount
}

func (r *FindManyUsersResponse) FromUsers(users []*model.User) {
	r.Users = make([]UserItem, len(users))
	for i, u := range users {
		r.Users[i].FromUser(u)
	}
}

func (d *ScrapItem) FromScrap(scrap *model.Scrap) {
	d.ID = scrap.ID
	d.Content = scrap.Content
	d.AuthorID = scrap.UserAuthorID
	d.AuthorName = scrap.UserAuthor.FirstName
	d.AuthorProfilePicture = *scrap.UserAuthor.Profile.ProfilePicture
	d.CreatedAt = scrap.CreatedAt
}

func (r *FindScrapOfUserRecipientResponse) FromScraps(scraps []model.Scrap) {
	r.Scraps = make([]ScrapItem, len(scraps))
	for i := range scraps {
		r.Scraps[i].FromScrap(&scraps[i])
	}
}

func (r *FindCommunityByIdResponse) FromCommunity(community *model.Community) {
	r.Name = community.Name
	r.Category = community.Category
	r.CreatedByID = community.CreatedByID
	r.OwnerFirstName = community.Owner.FirstName
	r.OwnerID = community.Owner.ID
	r.CreatedAt = community.CreatedAt
}

func (r *FindCommunityProfileByIdResponse) FromCommunityProfile(profile *model.CommunityProfile) {
	r.Description = profile.Description
	r.ProfilePicture = profile.ProfilePicture
}

func (r *FindCommunityMembersResponse) FromMemberships(community *model.Community, memberships []model.CommunityMembership) {
	r.MemberCount = community.MemberCount
	r.Members = make([]CommunityMemberEntry, len(memberships))
	for i, m := range memberships {
		r.Members[i] = CommunityMemberEntry{
			ID:             m.UserID,
			FirstName:      m.User.FirstName,
			ProfilePicture: *m.User.Profile.ProfilePicture,
			FriendCount:    m.User.FriendCount,
		}
	}
}

func (d *TopicSummary) FromTopic(topic *model.Topic, replyCount int64, replies []model.Topic) {
	d.ID = topic.ID
	d.AuthorProfilePicture = *topic.Author.Profile.ProfilePicture
	d.Title = topic.Title
	d.Replies = replyCount
	if len(replies) > 0 {
		lastReply := replies[len(replies)-1]
		d.LastReplyPreview = lastReply.Text
		d.LastReplyDate = lastReply.CreatedAt
	}
}

func (d *Topic) FromTopic(topic *model.Topic, replyCount int64) {
	d.ID = topic.ID
	d.AuthorID = topic.AuthorID
	d.AuthorFirstName = topic.Author.FirstName
	d.AuthorProfilePicture = *topic.Author.Profile.ProfilePicture
	d.Title = topic.Title
	d.Text = topic.Text
	d.CreatedAt = topic.CreatedAt
	d.Replies = replyCount
}

func (d *TopicReply) FromTopic(reply *model.Topic) {
	d.ID = reply.ID
	d.AuthorID = reply.AuthorID
	d.AuthorFirstName = reply.Author.FirstName
	d.AuthorProfilePicture = *reply.Author.Profile.ProfilePicture
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
	d.AuthorProfilePicture = *testimonial.UserAuthor.Profile.ProfilePicture
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
	d.Category = community.Category
	d.MemberCount = community.MemberCount
	if community.Profile != nil {
		d.ProfilePicture = community.Profile.ProfilePicture
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
	d.URL = photo.URL
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
