package service

import (
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	userRepository                repository.UserRepository
	communityRepository           repository.CommunityRepository
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

func NewUserService(
	userRepository repository.UserRepository,
	communityRepository repository.CommunityRepository,
	communityMembershipRepository repository.CommunityMembershipRepository,
	validate *validator.Validate,
) *UserService {
	return &UserService{
		userRepository:                userRepository,
		communityRepository:           communityRepository,
		communityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}
