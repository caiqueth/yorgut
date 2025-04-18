package service

import (
	"errors"
	"fmt"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
)

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

func (t *TestimonialService) FindTestimonialOfAuthor(userId uint) dto.FindTestimonialOfUserResponse {
	panic("unimplemented")
}

func (t *TestimonialService) FindTestimonialOfRecipient(userId uint) dto.FindTestimonialOfUserResponse {
	panic("unimplemented")
}

// User
func (u *UserService) DeleteUser(id uint) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindManyUsers() (*dto.Response, error) {
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

	return &dto.Response{
		Data: map[string][]dto.UserDto{"users": userDtos},
	}, nil
}

func (u *UserService) CreateUser(user dto.CreateUserRequest) error {
	err := u.userRepository.Create(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUser(user dto.UpdateUserRequest) error {
	err := u.userRepository.Update(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindUserById(id uint) (*dto.FindUserByIdResponse, error) {
	userData, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	resp := dto.FindUserByIdResponse{}
	resp.FromUser(userData)

	return &resp, nil

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

func (u *UserService) FindUserOwnedCommunities(userId uint) (*dto.FindUserCommunitiesResponse, error) {
	communities, err := u.userRepository.FindUserOwnedCommunities(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindUserCommunitiesResponse{}
	resp.FromCommunities(communities)

	return resp, nil

}

func (u *UserService) FindUserPhotoAlbums(userId int) (*dto.FindPhotoAlbumOfUserResponse, error) {
	photoAlbums, err := u.photoAlbumRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPhotoAlbumOfUserResponse{}
	resp.FromPhotoAlbums(photoAlbums)

	return resp, nil
}

func (u *UserService) FindPendingFriendshipsOfUser(userId uint) (*dto.FindPendingFriendshipOfUserResponse, error) {
	invites, err := u.friendshipRepository.FindPendingFriendshipsOfUser(userId, false)
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

func (u *UserService) JoinCommunity(userId uint, communityId uint) error {
	community, err := u.communityRepository.FindById(communityId)
	if err != nil {
		return err
	}

	user, err := u.userRepository.FindById(userId)
	if err != nil {
		return err
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
func (s *ScrapService) CreateScrap(scrap dto.CreateScrapRequest) error {
	author, err := s.userRepository.FindById(scrap.UserAuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	recipient, err := s.userRepository.FindById(scrap.UserRecipientID)
	if err != nil {
		return errors.New("recipient not found")
	}

	scrapModel := model.NewScrap(author, recipient, scrap.Content)

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

func (s *ScrapService) FindScrapOfUserRecipient(userId int) (dto.FindScrapOfUserRecipientResponse, error) {
	panic("unimplemented")
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

func (p *PhotoService) FindPhotoOfAlbum(photoAlbumId uint) dto.FindPhotoOfAlbumResponse {
	panic("unimplemented")
}

func (p *PhotoService) FindPhotoOfUser(userId uint) dto.FindPhotoOfUserResponse {
	panic("unimplemented")
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

func (p *PhotoAlbumService) FindPhotoAlbumOfUser(userId int) (*dto.FindPhotoAlbumOfUserResponse, error) {
	photoAlbums, err := p.photoAlbumRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPhotoAlbumOfUserResponse{}
	resp.FromPhotoAlbums(photoAlbums)

	return resp, nil
}

// Friendship
func (f *FriendshipService) CreateFriendship(friendship dto.CreateFriendshipRequest) error {
	userRequested, err := f.userRepository.FindById(friendship.UserRequested)
	if err != nil {
		return errors.New("could not find user that requested the friendship")
	}

	userRequestee, err := f.userRepository.FindById(friendship.UserRequestee)
	if err != nil {
		return errors.New("could not find user that requestee the friendship")
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.friendshipRepository.Create(friendshipModel)
	if err != nil {
		return err
	}
	return nil
}

// Friendship
func (f *FriendshipService) ReplyFriendship(friendshipId uint, accepted bool) error {
	friendship, err := f.friendshipRepository.FindPendingFriendshipWithUsersById(friendshipId)
	if err != nil {
		return err
	}

	if !accepted {
		f.friendshipRepository.Delete(friendshipId)
		return nil
	}

	friendship.Accepted = accepted
	if err := f.friendshipRepository.Update(friendship); err != nil {
		return err
	}

	if accepted {
		f.userRepository.IncrementFriendCount(friendship.UserRequested, 1)
		f.userRepository.IncrementFriendCount(friendship.UserRequestee, 1)
	}

	return nil
}

func (f *FriendshipService) DeleteFriendship(id uint) error {
	err := f.friendshipRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) FindFriendshipsOfUser(userId uint) (dto.FindFriendshipOfUserResponse, error) {
	panic("unimplemented")
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

func (c *CommunityService) FindCommunityById(id uint) (dto.FindCommunityByIdResponse, error) {
	community, err := c.communityRepository.FindById(id)
	if err != nil {
		return dto.FindCommunityByIdResponse{}, errors.New("could not find specified community id")
	}

	return dto.FindCommunityByIdResponse{
		Name:      community.Name,
		CreatedBy: community.CreatedByID,
	}, nil
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

func (c *CommunityService) FindManyCommunities() (*dto.Response, error) {
	communities, err := c.communityRepository.FindCommunitiesWithProfiles()
	if err != nil {
		return nil, err
	}

	communityDtos := make([]dto.CommunityDto, len(communities))

	for i, community := range communities {
		communityDtos[i] = dto.CommunityDto{
			Name:           community.Name,
			CreatedBy:      community.CreatedByID,
			ProfilePicture: community.Profile.ProfilePicture,
		}
	}

	return &dto.Response{
		Data: map[string][]dto.CommunityDto{"communities": communityDtos},
	}, nil
}

func (c *CommunityService) UpdateCommunity(user dto.UpdateCommunityRequest) (dto.Response, error) {
	panic("unimplemented")
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

	storedFileName, err := helpers.UploadCommunityPicture(picture)
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

func (c *CommunityMembershipService) FindCommunityMembershipsOfUser(userId uint) (dto.FindCommunityMembershipOfUserResponse, error) {
	panic("unimplemented")
}
