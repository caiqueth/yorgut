package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var ErrOwnerCannotLeave = errors.New("the owner cannot leave their own community")

type CommunityMembershipService struct {
	transactor                    repository.Transactor
	communityMembershipRepository repository.CommunityMembershipRepository
	communityRepository           repository.CommunityRepository
	userRepository                repository.UserRepository
	Validate                      *validator.Validate
}

func NewCommunityMembershipService(
	transactor repository.Transactor,
	communityMembershipRepository repository.CommunityMembershipRepository,
	communityRepository repository.CommunityRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *CommunityMembershipService {
	return &CommunityMembershipService{
		transactor:                    transactor,
		communityMembershipRepository: communityMembershipRepository,
		communityRepository:           communityRepository,
		userRepository:                userRepository,
		Validate:                      validate,
	}
}

func (c *CommunityMembershipService) FindUserMemberships(userId uint) (*transport.FindUserMembershipsResponse, error) {
	memberships, err := c.communityMembershipRepository.FindOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindUserMembershipsResponse{}
	resp.FromMemberships(memberships)
	resp.MembershipCount = c.CountUserMemberships(userId)

	return resp, nil
}

func (c *CommunityMembershipService) IsUserMemberOf(communityId uint, userId uint) bool {
	membership, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return false
	}

	if membership == nil {
		return false
	}

	return true
}

func (c *CommunityMembershipService) CountUserMemberships(userId uint) uint {
	count, err := c.communityMembershipRepository.CountOfUser(userId)
	if err != nil {
		return 0
	}
	return uint(count)
}

func (c *CommunityMembershipService) JoinCommunity(userId uint, communityId uint) error {
	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return err
	}

	user, err := c.userRepository.FindById(userId)
	if err != nil {
		return err
	}

	membershipExisting, _ := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if membershipExisting != nil {
		return errors.New("user is already a member of that community")
	}

	membership := model.NewCommunityMembership(community, user)
	return c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.CommunityMembership.Create(membership); err != nil {
			return err
		}
		return tx.Community.UpdateMemberCount(communityId)
	})
}

// LeaveCommunity removes the user's membership. The owner can't leave.
func (c *CommunityMembershipService) LeaveCommunity(userId uint, communityId uint) error {
	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return err
	}

	if community.OwnerID == userId {
		return ErrOwnerCannotLeave
	}

	membership, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return err
	}

	return c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.CommunityMembership.Delete(membership.ID); err != nil {
			return err
		}
		return tx.Community.UpdateMemberCount(communityId)
	})
}
