package service

import (
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityMembershipServiceImpl struct {
	CommunityMembershipRepository repository.CommunityMembershipRepository
	Validate                      *validator.Validate
}

// Create implements CommunityMembershipService.
func (c *CommunityMembershipServiceImpl) Create(community request.CreateCommunityMembershipRequest) {
	panic("unimplemented")
}

// Delete implements CommunityMembershipService.
func (c *CommunityMembershipServiceImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUser implements CommunityMembershipService.
func (c *CommunityMembershipServiceImpl) FindOfUser(userId int) response.FindCommunityMembershipOfUserResponse {
	panic("unimplemented")
}

func NewCommunityMembershipServiceImpl(communityMembershipRepository repository.CommunityMembershipRepository, validate *validator.Validate) CommunityMembershipService {
	return &CommunityMembershipServiceImpl{
		CommunityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}
