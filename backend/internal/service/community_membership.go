package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var (
	ErrOwnerCannotLeave = newError(ErrForbidden, "the owner cannot leave their own community")
	ErrAlreadyMember    = newError(ErrConflict, "user is already a member of that community")
)

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
		return fromRepository(err, "community")
	}

	user, err := c.userRepository.FindById(userId)
	if err != nil {
		return fromRepository(err, "user")
	}

	membershipExisting, _ := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if membershipExisting != nil {
		return ErrAlreadyMember
	}

	membership := model.NewCommunityMembership(community, user)
	err = c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.CommunityMembership.Create(membership); err != nil {
			return err
		}
		return tx.Community.UpdateMemberCount(communityId)
	})
	if errors.Is(err, repository.ErrDuplicate) {
		// Two joins at the same time: the unique index rejected the second.
		return ErrAlreadyMember
	}
	return err
}

// LeaveCommunity removes the user's membership. The owner can't leave.
func (c *CommunityMembershipService) LeaveCommunity(userId uint, communityId uint) error {
	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return fromRepository(err, "community")
	}

	if community.OwnerID == userId {
		return ErrOwnerCannotLeave
	}

	membership, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return fromRepository(err, "membership")
	}

	return c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.CommunityMembership.Delete(membership.ID); err != nil {
			return err
		}
		return tx.Community.UpdateMemberCount(communityId)
	})
}
