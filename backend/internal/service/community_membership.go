package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityMembershipService struct {
	CommunityMembershipRepository repository.CommunityMembershipRepository
	Validate                      *validator.Validate
}

func (c *CommunityMembershipService) Create(community request.CreateCommunityMembershipRequest) {
	panic("unimplemented")
}

func (c *CommunityMembershipService) Delete(id int) {
	panic("unimplemented")
}

func (c *CommunityMembershipService) FindOfUser(userId int) response.FindCommunityMembershipOfUserResponse {
	panic("unimplemented")
}

func NewCommunityMembershipService(
	communityMembershipRepository repository.CommunityMembershipRepository,
	validate *validator.Validate,
) *CommunityMembershipService {
	return &CommunityMembershipService{
		CommunityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}
