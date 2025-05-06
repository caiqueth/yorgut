package repository

import (
	"fmt"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
)

// Membership
func (c *CommunityMembershipRepository) FindOfUser(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.db.Preload("Community").Preload("Community.Profile").Where("user_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func (c *CommunityMembershipRepository) FindOfUserAndCommunity(communityId uint, userId uint) (*model.CommunityMembership, error) {
	var membership model.CommunityMembership
	result := c.db.
		Where("user_id = ? and community_id = ?", userId, communityId).
		First(&membership)
	if result.Error != nil {
		return nil, result.Error
	}
	return &membership, nil
}

func (c *CommunityMembershipRepository) FindOfCommunity(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.db.Preload("Community").Where("community_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func (c *CommunityMembershipRepository) FindMembersOfCommunity(communityId uint) ([]model.CommunityMembership, error) {
	var members []model.CommunityMembership
	result := c.db.Preload("User").Preload("User.Profile").Where("community_id = ?", communityId).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// PhotoAlbum
func (p *PhotoAlbumRepository) FindOfUser(userId uint) ([]model.PhotoAlbum, error) {
	var ents []model.PhotoAlbum
	result := p.db.Where("user_id = ?", userId).Find(&ents)
	if result.Error != nil {
		return nil, result.Error
	}
	return ents, nil
}

// Community
func (r *CommunityRepository) FindCommunitiesWithProfiles() ([]model.Community, error) {
	var ents []model.Community
	results := r.db.Preload("CommunityProfile").Find(&ents)
	if results.Error != nil {
		return nil, results.Error
	}
	return ents, nil
}

func (r *CommunityRepository) IncrementCommunityMembers(community *model.Community, increment int) error {
	community.MemberCount = helpers.SafeAdd(community.MemberCount, increment)
	result := r.db.Model(community).Updates(community)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// User
func (u *UserRepository) FindByEmail(email string) (*model.User, error) {
	var ent model.User
	result := u.db.Where("email = ?", email).First(&ent)
	if result != nil {
		return &ent, nil
	} else {
		return &ent, fmt.Errorf("user with e-mail %v not found", email)
	}
}

func (u *UserRepository) FindUserOwnedCommunities(userId uint) ([]model.Community, error) {
	var communities []model.Community

	result := u.db.Where("created_by_id = ?", userId).Find(&communities)

	if result.Error != nil {
		return nil, result.Error
	}

	return communities, nil
}

func (r *UserRepository) FindUserWithProfile(userId uint) (*model.User, error) {
	var ent model.User
	results := r.db.Preload("Profile").Where("id = ?", userId).First(&ent)
	if results.Error != nil {
		return nil, results.Error
	}
	return &ent, nil
}

func (r *UserRepository) IncrementFriendCount(user *model.User, increment int) error {
	user.FriendCount = helpers.SafeAdd(user.FriendCount, increment)
	result := r.db.Model(user).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Friendship
func (f *FriendshipRepository) FindPendingFriendshipsOfUser(userId uint, accepted bool) ([]model.User, error) {
	var ents []model.Friendship
	invites := []model.User{}
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequested.Profile").
		Where("(user_requestee_id = ?) and accepted = ?", userId, accepted).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, fsp := range ents {
		invites = append(invites, *fsp.UserRequested)
	}
	return invites, nil
}

func (f *FriendshipRepository) FindFriendsOfUser(userId uint) ([]model.User, error) {
	var ents []model.Friendship
	friends := []model.User{}
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Preload("UserRequested.Profile").
		Preload("UserRequestee.Profile").
		Where("(user_requested_id = ? or user_requestee_id = ?) and accepted = ?", userId, userId, true).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, fsp := range ents {
		if fsp.UserRequestedID != userId {
			friends = append(friends, *fsp.UserRequested)
		}
		if fsp.UserRequesteeID != userId {
			friends = append(friends, *fsp.UserRequestee)
		}
	}

	return friends, nil
}

func (f *FriendshipRepository) UsersAreFriends(user1 uint, user2 uint) bool {
	var count int64
	f.db.Where(
		"(user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)",
		user1, user2, user2, user1,
	).Count(&count)

	return count > 0
}

func (f *FriendshipRepository) FindPendingFriendshipWithUsersById(friendshipId uint) (*model.Friendship, error) {
	var ent model.Friendship

	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Preload("UserRequested.Profile").
		Preload("UserRequestee.Profile").
		Where("id = ?", friendshipId).
		First(&ent)

	if result.Error != nil {
		return nil, result.Error
	}
	return &ent, nil
}

// Scrap
func (s *ScrapRepository) FindOfUserRecipient(userId uint) ([]model.Scrap, error) {
	var scraps []model.Scrap
	result := s.db.
		Where("user_recipient_id = ?", userId).
		Preload("UserAuthor").
		Preload("UserAuthor.Profile").
		Order("created_at desc").
		Find(&scraps)

	if result.Error != nil {
		return nil, result.Error
	}
	return scraps, nil
}
