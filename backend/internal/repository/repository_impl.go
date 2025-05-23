package repository

import (
	"errors"
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

func (c *CommunityMembershipRepository) CountOfUser(id uint) (int64, error) {
	var count int64
	result := c.db.Model(&model.CommunityMembership{}).Where("user_id = ?", id).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
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
	results := r.db.Preload("Profile").Find(&ents)
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

	// Had to do it this way, otherwise gorm would ignore 0 values
	result := r.db.Model(user).Select("FriendCount").Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Friendship
func (f *FriendshipRepository) FindByIdWithUsers(id uint) (*model.Friendship, error) {
	var ent model.Friendship
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Find(&ent, id)

	if result != nil {
		return &ent, nil
	} else {
		return nil, errors.New("friendship not found")
	}
}

func (f *FriendshipRepository) FindPendingFriendshipsOfUser(userId uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequested.Profile").
		Where("(user_requestee_id = ?) and accepted = false", userId).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

func (f *FriendshipRepository) FindFriendshipsOfUser(userId uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Where("(user_requestee_id = ? or user_requested_id = ?)", userId, userId).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

func (f *FriendshipRepository) FindFriendshipsBetweenUsers(userId1 uint, userId2 uint) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.
		Where("(user_requestee_id = ? and user_requested_id = ?) or (user_requestee_id = ? and user_requested_id = ?)", userId1, userId2, userId2, userId1).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	return ents, nil
}

func (f *FriendshipRepository) FindFriendsOfUser(userId uint) ([]model.User, error) {
	var ents []model.Friendship
	friends := []model.User{}
	result := f.db.
		Preload("UserRequested").
		Preload("UserRequestee").
		Preload("UserRequested.Profile").
		Preload("UserRequestee.Profile").
		Where("(user_requested_id = ? or user_requestee_id = ?) and accepted = true", userId, userId).
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

func (f *FriendshipRepository) FindActiveAcceptedFriends(userId1 uint, userId2 uint) (*model.Friendship, error) {
	var ent model.Friendship
	result := f.db.
		Where("((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)) and accepted = true", userId1, userId2, userId2, userId1).
		First(&ent)

	if result.Error != nil {
		return nil, result.Error
	}

	return &ent, nil
}

func (f *FriendshipRepository) UsersAreFriends(user1 uint, user2 uint) bool {
	var count int64
	f.db.Model(&model.Friendship{}).Where(
		"((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)) and accepted = true",
		user1, user2, user2, user1,
	).Count(&count)

	return count > 0
}

func (f *FriendshipRepository) UsersHavePendingFriendshipInvites(user1 uint, user2 uint) bool {
	var count int64
	f.db.Model(&model.Friendship{}).Where(
		"((user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)) and accepted = false",
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

// Topic
func (t *TopicRepository) FindMostRecentOfCommunity(communityId uint, limit int) ([]model.Topic, error) {
	var topics []model.Topic
	result := t.db.
		Where("community_id = ? and parent_topic_id is null", communityId).
		Preload("Author").
		Preload("Author.Profile").
		Order("created_at desc").
		Limit(limit).
		Find(&topics)

	if result.Error != nil {
		return nil, result.Error
	}

	return topics, nil
}

func (t *TopicRepository) FindRepliesOfTopic(topicId uint) (int64, []model.Topic, error) {
	var replies []model.Topic
	var count int64
	resultCount := t.db.Model(&model.Topic{}).Where("parent_topic_id = ?", topicId).Count(&count)

	if resultCount.Error != nil {
		return 0, nil, resultCount.Error
	}

	resultReplies := t.db.
		Where("parent_topic_id = ?", topicId).
		Preload("Author").
		Preload("Author.Profile").
		Order("created_at asc").
		Find(&replies)

	if resultReplies.Error != nil {
		return 0, nil, resultReplies.Error
	}

	return count, replies, nil
}

func (t *TopicRepository) FindByIdWithAuthorProfile(topicId uint) (*model.Topic, error) {
	var topic model.Topic
	result := t.db.
		Preload("Author").
		Preload("Author.Profile").
		Find(&topic, topicId)

	if result != nil {
		return &topic, nil
	} else {
		return &topic, errors.New("topic not found")
	}

}
