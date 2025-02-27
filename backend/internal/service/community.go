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
	CommunityRepository repository.CommunityRepository
	UserRepository      repository.UserRepository
	Validate            *validator.Validate
}

func (c *CommunityService) Create(community request.CreateCommunityRequest) error {
	owner, err := c.UserRepository.FindById(community.CreatedByID)
	if err != nil {
		return errors.New("could not find specified community owner")
	}

	communityModel := model.NewCommunity(community.Name, community.Description, community.PictureUrl, owner)
	err = c.CommunityRepository.Create(communityModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) Delete(id uint) error {
	err := c.CommunityRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) FindById(id uint) (response.FindCommunityByIdResponse, error) {
	community, err := c.CommunityRepository.FindById(id)
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

func (c *CommunityService) FindMany() (response.Response, error) {
	panic("unimplemented")
}

func (c *CommunityService) Update(user request.UpdateCommunityRequest) (response.Response, error) {
	panic("unimplemented")
}

func NewCommunityService(communityRepository repository.CommunityRepository, validate *validator.Validate) *CommunityService {
	return &CommunityService{
		CommunityRepository: communityRepository,
		Validate:            validate,
	}
}
