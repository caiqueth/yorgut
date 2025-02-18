package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityServiceImpl struct {
	CommunityRepository repository.CommunityRepository
	Validate            *validator.Validate
}

// Create implements CommunityService.
func (c *CommunityServiceImpl) Create(community request.CreateCommunityRequest) error {
	panic("unimplemented")
}

// Delete implements CommunityService.
func (c *CommunityServiceImpl) Delete(communityId int) {
	panic("unimplemented")
}

// FindById implements CommunityService.
func (c *CommunityServiceImpl) FindById(communityId int) response.FindCommunityByIdResponse {
	panic("unimplemented")
}

// FindMany implements CommunityService.
func (c *CommunityServiceImpl) FindMany() response.Response {
	panic("unimplemented")
}

// Update implements CommunityService.
func (c *CommunityServiceImpl) Update(user request.UpdateCommunityRequest) response.Response {
	panic("unimplemented")
}

func NewCommunityServiceImpl(communityRepository repository.CommunityRepository, validate *validator.Validate) CommunityService {
	return &CommunityServiceImpl{
		CommunityRepository: communityRepository,
		Validate:            validate,
	}
}
