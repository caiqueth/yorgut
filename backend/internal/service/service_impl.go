package service

import (
	"errors"
	"fmt"
	"path"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
)

func (t *TestimonialService) Create(testimonial dto.CreateTestimonialRequest) error {
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

func (t *TestimonialService) Delete(id uint) error {
	err := t.TestimonialRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TestimonialService) FindOfUserAuthor(userId uint) dto.FindTestimonialOfUserResponse {
	panic("unimplemented")
}

func (t *TestimonialService) FindOfUserRecipient(userId uint) dto.FindTestimonialOfUserResponse {
	panic("unimplemented")
}

func (u *UserService) Delete(id uint) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindMany() (*dto.Response, error) {
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

func (u *UserService) Create(user dto.CreateUserRequest) error {
	err := u.userRepository.Create(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Update(user dto.UpdateUserRequest) error {
	err := u.userRepository.Update(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindById(id uint) (*dto.FindUserByIdResponse, error) {
	userData, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	resp := dto.FindUserByIdResponse{}
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

func (u *UserService) FindFriendshipsOfUser(userId int) (*dto.FindFriendshipOfUserResponse, error) {
	friendships, err := u.friendshipRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindFriendshipOfUserResponse{}
	resp.FromFriendships(friendships)

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

	m := model.NewCommunityMembership(community, user)

	u.communityMembershipRepository.Create(m)
	return nil
}

func (s *ScrapService) Create(scrap dto.CreateScrapRequest) error {
	author, err := s.UserRepository.FindById(scrap.UserAuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	recipient, err := s.UserRepository.FindById(scrap.UserRecipientID)
	if err != nil {
		return errors.New("recipient not found")
	}

	scrapModel := model.NewScrap(author, recipient, scrap.Content)

	err = s.ScrapRepository.Create(scrapModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScrapService) Delete(id uint) error {
	err := s.ScrapRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScrapService) FindOfUserRecipient(userId int) (dto.FindScrapOfUserRecipientResponse, error) {
	panic("unimplemented")
}

func (p *PhotoService) Create(photo dto.CreatePhotoRequest) error {
	photoAlbum, err := p.PhotoAlbumRepository.FindById(photo.PhotoAlbumID)
	if err != nil {
		return errors.New("photo album not found")
	}

	photoModel := model.NewPhoto(photoAlbum, photo.URL, photo.Caption)

	err = p.PhotoRepository.Create(photoModel)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoService) Delete(id uint) error {
	err := p.PhotoRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PhotoService) FindOfAlbum(photoAlbumId uint) dto.FindPhotoOfAlbumResponse {
	panic("unimplemented")
}

func (p *PhotoService) FindOfUser(userId uint) dto.FindPhotoOfUserResponse {
	panic("unimplemented")
}

func (p *PhotoAlbumService) Create(photoAlbum dto.CreatePhotoAlbumRequest) error {
	owner, err := p.UserRepository.FindById(photoAlbum.UserID)
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

func (p *PhotoAlbumService) Delete(id uint) error {
	err := p.photoAlbumRepository.Delete(id)
	if err != nil {
		return err
	}
	return err
}

func (p *PhotoAlbumService) FindOfUser(userId int) (*dto.FindPhotoAlbumOfUserResponse, error) {
	photoAlbums, err := p.photoAlbumRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.FindPhotoAlbumOfUserResponse{}
	resp.FromPhotoAlbums(photoAlbums)

	return resp, nil
}

func (f *FriendshipService) Create(friendship dto.CreateFriendshipRequest) error {
	userRequested, err := f.UserRepository.FindById(friendship.UserRequested)
	if err != nil {
		return errors.New("could not find user that requested the friendship")
	}

	userRequestee, err := f.UserRepository.FindById(friendship.UserRequestee)
	if err != nil {
		return errors.New("could not find user that requestee the friendship")
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.FriendshipRepository.Create(friendshipModel)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) Delete(id uint) error {
	err := f.FriendshipRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) FindOfUser(userId uint) (dto.FindFriendshipOfUserResponse, error) {
	panic("unimplemented")
}

func (c *CommunityService) Create(request dto.CreateCommunityRequest) error {
	owner, err := c.userRepository.FindById(request.CreatedByID)
	if err != nil {
		return errors.New("could not find specified community owner")
	}

	picture, err := helpers.FileToBytes(request.Picture)
	if err != nil {
		return err
	}

	imageHasValidExtension := helpers.CheckValidImage(picture)
	if !imageHasValidExtension {
		return errors.New("invalid image format")
	}

	storedFileName, err := helpers.WriteFileToLocalStorage(path.Join(helpers.GetUploadsPath(), "community_pictures"), picture)
	if err != nil {
		return err
	}

	communityModel := model.NewCommunity(request.Name, request.Description, storedFileName, owner)
	err = c.communityRepository.Create(communityModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) Delete(id uint) error {
	err := c.communityRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) FindById(id uint) (dto.FindCommunityByIdResponse, error) {
	community, err := c.communityRepository.FindById(id)
	if err != nil {
		return dto.FindCommunityByIdResponse{}, errors.New("could not find specified community id")
	}

	r := dto.FindCommunityByIdResponse{
		Name:        community.Name,
		Description: community.Description,
		CreatedBy:   community.CreatedByID,
	}
	return r, nil
}

func (c *CommunityService) FindMany() (*dto.Response, error) {
	communities, err := c.communityRepository.FindMany()
	if err != nil {
		return nil, err
	}

	communityDtos := make([]dto.CommunityDto, len(communities))

	for i, community := range communities {
		communityDtos[i] = dto.CommunityDto{
			Name:        community.Name,
			Description: community.Description,
			CreatedBy:   community.CreatedByID,
			PictureUrl:  community.PictureUrl,
		}
	}

	return &dto.Response{
		Data: map[string][]dto.CommunityDto{"communities": communityDtos},
	}, nil
}

func (c *CommunityService) Update(user dto.UpdateCommunityRequest) (dto.Response, error) {
	panic("unimplemented")
}

func (c *CommunityMembershipService) Create(communityMembership dto.CreateCommunityMembershipRequest) error {
	community, err := c.CommunityRepository.FindById(communityMembership.CommunityID)
	if err != nil {
		return errors.New("could not find specified community")
	}

	user, err := c.UserRepository.FindById(communityMembership.UserID)
	if err != nil {
		return errors.New("could not find community owner")
	}

	communityMembershipModel := model.NewCommunityMembership(community, user)
	err = c.CommunityMembershipRepository.Create(communityMembershipModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityMembershipService) Delete(id uint) error {
	err := c.CommunityMembershipRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityMembershipService) FindOfUser(userId uint) (dto.FindCommunityMembershipOfUserResponse, error) {
	panic("unimplemented")
}
