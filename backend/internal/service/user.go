package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	userRepository                repository.UserRepository
	communityMembershipRepository repository.CommunityMembershipRepository
	Validate                      *validator.Validate
}

func (u *UserService) Delete(id uint) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindMany() (*response.Response, error) {
	users, err := u.userRepository.FindMany()
	if err != nil {
		return nil, err
	}

	userDtos := make([]response.UserDto, len(users))

	for i, user := range users {
		userDtos[i] = response.UserDto{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			BirthDate: user.BirthDate,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return &response.Response{
		Data: map[string][]response.UserDto{"Users": userDtos},
	}, nil
}

func (u *UserService) Create(user request.CreateUserRequest) error {
	err := u.userRepository.Create(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Update(user request.UpdateUserRequest) error {
	err := u.userRepository.Update(model.NewUser(user.FirstName, user.LastName, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindById(id uint) (*response.FindUserByIdResponse, error) {
	userData, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	userResponse := response.FindUserByIdResponse{
		UserDto: response.UserDto{
			ID:        userData.ID,
			Email:     userData.Email,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			BirthDate: userData.BirthDate,
			CreatedAt: userData.CreatedAt,
			UpdatedAt: userData.UpdatedAt,
		}}

	return &userResponse, err

}

func (u *UserService) FindUserOwnedCommunities(userId uint) (*response.FindUserCommunitiesResponse, error) {
	communities, err := u.userRepository.FindUserOwnedCommunities(userId)
	if err != nil {
		return nil, err
	}
	r := make([]response.CommunityDto, len(communities))

	for _, c := range communities {
		r = append(r, response.CommunityDto{
			Name:        c.Name,
			Description: c.Description,
			CreatedBy:   c.CreatedByID,
			PictureUrl:  c.PictureUrl,
		})
	}
	return &response.FindUserCommunitiesResponse{
		Communities: r,
	}, nil
}

func (u *UserService) FindUserMemberships(userId uint) (*response.FindUserMembershipsResponse, error) {
	memberships, err := u.communityMembershipRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}
	r := make([]response.CommunityDto, len(memberships))
	for _, m := range memberships {
		r = append(r, response.CommunityDto{
			Name:        m.Community.Name,
			Description: m.Community.Description,
			CreatedBy:   m.Community.CreatedByID,
			PictureUrl:  m.Community.PictureUrl,
		})
	}
	return &response.FindUserMembershipsResponse{
		Communities: r,
	}, nil
}

func NewUserService(
	userRepository repository.UserRepository,
	communityMembershipRepository repository.CommunityMembershipRepository,
	validate *validator.Validate,
) *UserService {
	return &UserService{
		userRepository:                userRepository,
		communityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}
