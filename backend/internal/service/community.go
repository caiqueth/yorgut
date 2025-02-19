package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityService struct {
	CommunityRepository repository.CommunityRepository
	Validate            *validator.Validate
}

func (c *CommunityService) Create(community request.CreateCommunityRequest) error {
	panic("unimplemented")
}

func (c *CommunityService) Delete(communityId int) {
	panic("unimplemented")
}

func (c *CommunityService) FindById(communityId int) response.FindCommunityByIdResponse {
	panic("unimplemented")
}

func (c *CommunityService) FindMany() response.Response {
	panic("unimplemented")
}

func (c *CommunityService) Update(user request.UpdateCommunityRequest) response.Response {
	panic("unimplemented")
}

func NewCommunityService(communityRepository repository.CommunityRepository, validate *validator.Validate) *CommunityService {
	return &CommunityService{
		CommunityRepository: communityRepository,
		Validate:            validate,
	}
}
