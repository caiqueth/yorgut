package dto

import "yorgut/internal/model"

// Mappers
func (r *FindUserByIdResponse) FromUser(user *model.User) {
	r.ID = user.ID
	r.Email = user.Email
	r.FirstName = user.FirstName
	r.LastName = user.LastName
	r.BirthDate = user.BirthDate
	r.CreatedAt = user.CreatedAt
	r.UpdatedAt = user.UpdatedAt
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
			Name:           m.Community.Name,
			CreatedBy:      m.Community.CreatedByID,
			ProfilePicture: m.Community.Profile.ProfilePicture,
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
		}
	}
}
