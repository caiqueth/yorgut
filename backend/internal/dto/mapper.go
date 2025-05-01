package dto

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

func (r *FindUserCommunitiesResponse) FromCommunities(communities []model.Community) {
	r.Communities = make([]CommunityDto, len(communities))
	for i, c := range communities {
		r.Communities[i] = CommunityDto{
			Name:      c.Name,
			CreatedBy: c.CreatedByID,
		}
	}
}

func (r *FindUserMembershipsResponse) FromMemberships(memberships []model.CommunityMembership) {
	r.Communities = make([]CommunityDto, len(memberships))
	for i, m := range memberships {
		r.Communities[i] = CommunityDto{
			ID:              m.Community.ID,
			Name:            m.Community.Name,
			CreatedBy:       m.Community.CreatedByID,
			ProfilePicture:  m.Community.Profile.ProfilePicture,
			MembershipCount: m.Community.MemberCount,
		}
	}
}

func (r *FindPhotoAlbumOfUserResponse) FromPhotoAlbums(photoAlbums []model.PhotoAlbum) {
	r.PhotoAlbums = make([]PhotoAlbumDto, len(photoAlbums))

	for i, pa := range photoAlbums {
		r.PhotoAlbums[i] = PhotoAlbumDto{
			Name:   pa.Name,
			UserID: pa.UserID,
		}
	}
}

func (r *FindFriendshipOfUserResponse) FromFriendships(friendships []model.Friendship) {
	r.Friendships = make([]FriendshipDto, len(friendships))

	for i, f := range friendships {
		r.Friendships[i] = FriendshipDto{
			UserRequestedID: f.UserRequestedID,
			UserRequesteeID: f.UserRequesteeID,
			CreatedAt:       f.CreatedAt,
		}
	}
}

func (r *FindFriendsOfUserResponse) FromFriends(friends []model.User) {
	r.Friends = make([]FriendDto, len(friends))

	for i, f := range friends {
		r.Friends[i] = FriendDto{
			ID:             f.ID,
			FirstName:      f.FirstName,
			BirthDate:      f.BirthDate,
			ProfilePicture: *f.Profile.ProfilePicture,
			FriendCount:    f.FriendCount,
		}
	}
}

func (r *FindPendingFriendshipOfUserResponse) FromUsers(users []model.User) {
	r.PendingFriendships = make([]PendingFriendship, len(users))

	for i, f := range users {
		r.PendingFriendships[i] = PendingFriendship{
			ID:             f.ID,
			FirstName:      f.FirstName,
			LastName:       f.LastName,
			ProfilePicture: *f.Profile.ProfilePicture,
		}
	}
}
