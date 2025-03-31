package repository

import "yorgut/internal/model"

func (u *UserRepository) FindUserOwnedCommunities(userId uint) ([]model.Community, error) {
	var communities []model.Community

	result := u.db.Where("created_by_id = ?", userId).Find(&communities)

	if result.Error != nil {
		return nil, result.Error
	}

	return communities, nil
}

func (c *CommunityMembershipRepository) FindOfUser(id uint) ([]model.CommunityMembership, error) {
	var memberships []model.CommunityMembership
	result := c.db.Preload("Community").Where("user_id = ?", id).Find(&memberships)
	if result.Error != nil {
		return nil, result.Error
	}
	return memberships, nil
}

func (c *CommunityMembershipRepository) FindOfUserAndCommunity(communityId uint, userId uint) (*model.CommunityMembership, error) {
	var membership model.CommunityMembership
	result := c.db.Where("user_id = ? and communityId = ?", userId, communityId).First(&membership)
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

func (f *FriendshipRepository) FindOfUser(userId int) ([]model.Friendship, error) {
	var ents []model.Friendship
	result := f.db.Where("user_requested_id = ? or user_requestee_id = ?", userId, userId).Find(&ents)
	if result.Error != nil {
		return nil, result.Error
	}
	return ents, nil
}

func (f *FriendshipRepository) UsersAreFriends(user1 uint, user2 uint) bool {
	var count int64
	f.db.Where(
		"(user_requested_id = ? and user_requestee_id = ?) or (user_requested_id = ? and user_requestee_id = ?)",
		user1, user2, user2, user1,
	).Count(&count)

	return count > 0
}

func (p *PhotoAlbumRepository) FindOfUser(userId int) ([]model.PhotoAlbum, error) {
	var ents []model.PhotoAlbum
	result := p.db.Where("user_id = ?", userId).Find(&ents)
	if result.Error != nil {
		return nil, result.Error
	}
	return ents, nil
}

func (s *ScrapRepository) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}

func (p *PhotoRepository) FindOfPhotoAlbum(photoAlbumId int) {
	panic("unimplemented")
}

func (p *PhotoRepository) FindOfUser(userId int) {
	panic("unimplemented")
}

func (t *TestimonialRepository) FindOfUserAuthor(userId int) {
	panic("unimplemented")
}

func (t *TestimonialRepository) FindOfUserRecipient(userId int) {
	panic("unimplemented")
}
