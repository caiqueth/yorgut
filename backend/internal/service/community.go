package service

import (
	"errors"
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityService struct {
	communityRepository repository.CommunityRepository
	userRepository      repository.UserRepository
	Validate            *validator.Validate
}

func (c *CommunityService) Create(community request.CreateCommunityRequest) error {
	owner, err := c.userRepository.FindById(community.CreatedByID)
	if err != nil {
		return errors.New("could not find specified community owner")
	}

	communityModel := model.NewCommunity(community.Name, community.Description, community.PictureUrl, owner)
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

func (c *CommunityService) FindById(id uint) (response.FindCommunityByIdResponse, error) {
	community, err := c.communityRepository.FindById(id)
	if err != nil {
		return response.FindCommunityByIdResponse{}, errors.New("could not find specified community id")
	}

	r := response.FindCommunityByIdResponse{
		Name:        community.Name,
		Description: community.Description,
		CreatedBy:   community.CreatedByID,
	}
	return r, nil
}

func (c *CommunityService) FindMany() (*response.Response, error) {
	communities, err := c.communityRepository.FindMany()
	if err != nil {
		return nil, err
	}

	communityDtos := make([]response.CommunityDto, len(communities))

	for i, community := range communities {
		communityDtos[i] = response.CommunityDto{
			Name:        community.Name,
			Description: community.Description,
			CreatedBy:   community.CreatedByID,
			PictureUrl:  community.PictureUrl,
		}
	}

	return &response.Response{
		Data: map[string][]response.CommunityDto{"communities": communityDtos},
	}, nil
}

func (c *CommunityService) Update(user request.UpdateCommunityRequest) (response.Response, error) {
	panic("unimplemented")
}

func NewCommunityService(
	communityRepository repository.CommunityRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *CommunityService {
	return &CommunityService{
		communityRepository: communityRepository,
		userRepository:      userRepository,
		Validate:            validate,
	}
}
