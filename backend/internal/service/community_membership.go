package service

import (
	"errors"
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityMembershipService struct {
	CommunityMembershipRepository repository.CommunityMembershipRepository
	CommunityRepository           repository.CommunityRepository
	UserRepository                repository.UserRepository
	Validate                      *validator.Validate
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

func NewCommunityMembershipService(
	communityMembershipRepository repository.CommunityMembershipRepository,
	validate *validator.Validate,
) *CommunityMembershipService {
	return &CommunityMembershipService{
		CommunityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}
