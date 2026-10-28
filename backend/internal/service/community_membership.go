package service

import (
	"errors"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var (
	ErrOwnerCannotLeave      = newError(ErrForbidden, "O dono não pode sair da própria comunidade.")
	ErrAlreadyMember         = newError(ErrConflict, "Você já participa desta comunidade.")
	ErrAlreadyRequested      = newError(ErrConflict, "Você já pediu para participar desta comunidade.")
	ErrJoinRequestNotFound   = newError(ErrNotFound, "Este pedido para participar não existe.")
	ErrJoinRequestsForbidden = newError(ErrForbidden, "Só o dono vê e responde os pedidos para participar.")
	ErrInviterNotMember      = newError(ErrForbidden, "Só quem participa da comunidade pode convidar amigos.")
	ErrInviteeNotFriend      = newError(ErrForbidden, "Você só pode convidar seus amigos.")
	ErrInviteeAlreadyMember  = newError(ErrConflict, "Essa pessoa já participa desta comunidade.")
	ErrAlreadyInvited        = newError(ErrConflict, "Essa pessoa já foi convidada para esta comunidade.")
	ErrInvitationNotFound    = newError(ErrNotFound, "Este convite não existe.")
)

// CommunityMembershipService handles joining and leaving communities,
// including the join requests of private communities and the invitations
// members send their friends.
type CommunityMembershipService struct {
	transactor                     repository.Transactor
	communityMembershipRepository  repository.CommunityMembershipRepository
	communityJoinRequestRepository repository.CommunityJoinRequestRepository
	communityInvitationRepository  repository.CommunityInvitationRepository
	communityRepository            repository.CommunityRepository
	userRepository                 repository.UserRepository
	friendshipRepository           repository.FriendshipRepository
}

func NewCommunityMembershipService(
	transactor repository.Transactor,
	communityMembershipRepository repository.CommunityMembershipRepository,
	communityJoinRequestRepository repository.CommunityJoinRequestRepository,
	communityInvitationRepository repository.CommunityInvitationRepository,
	communityRepository repository.CommunityRepository,
	userRepository repository.UserRepository,
	friendshipRepository repository.FriendshipRepository,
) *CommunityMembershipService {
	return &CommunityMembershipService{
		transactor:                     transactor,
		communityMembershipRepository:  communityMembershipRepository,
		communityJoinRequestRepository: communityJoinRequestRepository,
		communityInvitationRepository:  communityInvitationRepository,
		communityRepository:            communityRepository,
		userRepository:                 userRepository,
		friendshipRepository:           friendshipRepository,
	}
}

// FindUserMemberships returns one page of the communities the user is a
// member of, by name or, with random, shuffled (for the profile box). The
// shuffle stays the same all day, so reloading the box doesn't change it.
func (c *CommunityMembershipService) FindUserMemberships(viewerId uint, userId uint, page int, pageSize int, random bool) (*transport.FindUserMembershipsResponse, error) {
	if err := requireAccess(c.userRepository, c.friendshipRepository, viewerId, userId, sectionCommunities); err != nil {
		return nil, err
	}

	var seed *int64
	if random {
		daily := dailySeed(userId, time.Now())
		seed = &daily
	}
	communities, count, err := c.communityMembershipRepository.FindCommunitiesOfUser(userId, page, pageSize, seed)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindUserMembershipsResponse{}
	resp.FromCommunities(communities, page, pageSize, count)

	return resp, nil
}

// FindMembershipStatus says whether the user is a member of the community,
// or is waiting for its owner to accept them.
func (c *CommunityMembershipService) FindMembershipStatus(communityId uint, userId uint) (*transport.AmIMemberOfResponse, error) {
	resp := &transport.AmIMemberOfResponse{}

	_, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err == nil {
		resp.AmIMember = true
		return resp, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	_, err = c.communityJoinRequestRepository.FindOfUserAndCommunity(communityId, userId)
	if err == nil {
		resp.Pending = true
		return resp, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	return resp, nil
}

// JoinCommunity makes the user a member of a public community. For a
// private one it asks the owner instead, and returns pending true.
func (c *CommunityMembershipService) JoinCommunity(userId uint, communityId uint) (pending bool, err error) {
	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return false, fromRepository(err, "a comunidade")
	}

	user, err := c.userRepository.FindById(userId)
	if err != nil {
		return false, fromRepository(err, "o usuário")
	}

	membershipExisting, _ := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if membershipExisting != nil {
		return false, ErrAlreadyMember
	}

	// The owner never needs to ask (only the seed hits this: creating a
	// community already makes the owner a member).
	if community.Private && community.OwnerID != userId {
		return true, c.requestToJoin(communityId, userId)
	}

	return false, c.addMember(community.ID, user.ID)
}

// requestToJoin records the user's request to join a private community.
func (c *CommunityMembershipService) requestToJoin(communityId uint, userId uint) error {
	requestExisting, _ := c.communityJoinRequestRepository.FindOfUserAndCommunity(communityId, userId)
	if requestExisting != nil {
		return ErrAlreadyRequested
	}

	err := c.communityJoinRequestRepository.Create(model.NewCommunityJoinRequest(communityId, userId))
	if errors.Is(err, repository.ErrDuplicate) {
		// Two requests at the same time: the unique index rejected the second.
		return ErrAlreadyRequested
	}
	return err
}

func (c *CommunityMembershipService) addMember(communityId uint, userId uint) error {
	err := c.transactor.Transaction(func(tx *repository.Repositories) error {
		return joinInTransaction(tx, communityId, userId)
	})
	if errors.Is(err, repository.ErrDuplicate) {
		// Two joins at the same time: the unique index rejected the second.
		return ErrAlreadyMember
	}
	return err
}

// joinInTransaction makes the user a member of the community. Joining
// answers whatever was still pending for them there, so it also drops their
// invitation and their join request.
func joinInTransaction(tx *repository.Repositories, communityId uint, userId uint) error {
	if err := tx.CommunityInvitation.DeleteOfInviteeAndCommunity(communityId, userId); err != nil {
		return err
	}
	if err := tx.CommunityJoinRequest.DeleteOfUserAndCommunity(communityId, userId); err != nil {
		return err
	}
	if err := tx.CommunityMembership.Create(&model.CommunityMembership{CommunityID: communityId, UserID: userId}); err != nil {
		return err
	}
	return tx.Community.UpdateMemberCount(communityId)
}

// LeaveCommunity removes the user's membership. The owner can't leave.
func (c *CommunityMembershipService) LeaveCommunity(userId uint, communityId uint) error {
	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return fromRepository(err, "a comunidade")
	}

	if community.OwnerID == userId {
		return ErrOwnerCannotLeave
	}

	membership, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if err != nil {
		return fromRepository(err, "a participação")
	}

	return c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.CommunityMembership.Delete(membership.ID); err != nil {
			return err
		}
		return tx.Community.UpdateMemberCount(communityId)
	})
}

// CancelJoinRequest withdraws the user's pending request to join a community.
func (c *CommunityMembershipService) CancelJoinRequest(userId uint, communityId uint) error {
	request, err := c.communityJoinRequestRepository.FindOfUserAndCommunity(communityId, userId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrJoinRequestNotFound
	}
	if err != nil {
		return err
	}
	return c.communityJoinRequestRepository.Delete(request.ID)
}

// findCommunityOwnedBy loads the community and checks that ownerId owns it.
func (c *CommunityMembershipService) findCommunityOwnedBy(ownerId uint, communityId uint) (*model.Community, error) {
	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return nil, fromRepository(err, "a comunidade")
	}
	if community.OwnerID != ownerId {
		return nil, ErrJoinRequestsForbidden
	}
	return community, nil
}

// FindJoinRequests lists the users waiting to join the community, oldest
// first. Only the owner may see them.
func (c *CommunityMembershipService) FindJoinRequests(ownerId uint, communityId uint) (*transport.FindJoinRequestsResponse, error) {
	if _, err := c.findCommunityOwnedBy(ownerId, communityId); err != nil {
		return nil, err
	}

	requests, err := c.communityJoinRequestRepository.FindOfCommunity(communityId)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindJoinRequestsResponse{}
	resp.FromJoinRequests(requests)
	return resp, nil
}

// ReplyJoinRequest answers userId's request to join the community: accepting
// makes them a member, declining just drops the request. Only the owner may.
func (c *CommunityMembershipService) ReplyJoinRequest(ownerId uint, communityId uint, userId uint, accept bool) error {
	community, err := c.findCommunityOwnedBy(ownerId, communityId)
	if err != nil {
		return err
	}

	request, err := c.communityJoinRequestRepository.FindOfUserAndCommunity(communityId, userId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrJoinRequestNotFound
	}
	if err != nil {
		return err
	}

	if !accept {
		return c.communityJoinRequestRepository.Delete(request.ID)
	}

	return c.addMember(community.ID, userId)
}

// CountUserMemberships counts the communities the user is a member of.
func (c *CommunityMembershipService) CountUserMemberships(userId uint) (uint, error) {
	count, err := c.communityMembershipRepository.CountOfUser(userId)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

// requireMember checks that the community exists and the user is a member.
func (c *CommunityMembershipService) requireMember(communityId uint, userId uint) error {
	if _, err := c.communityRepository.FindById(communityId); err != nil {
		return fromRepository(err, "a comunidade")
	}
	_, err := c.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrInviterNotMember
	}
	return err
}

// FindInvitableFriends returns one page of the user's friends, saying which
// are already members of the community or already invited to it. Only
// members may invite, so only they may see it.
func (c *CommunityMembershipService) FindInvitableFriends(userId uint, communityId uint, page int, pageSize int) (*transport.FindInvitableFriendsResponse, error) {
	if err := c.requireMember(communityId, userId); err != nil {
		return nil, err
	}

	friends, count, err := c.friendshipRepository.FindFriendsOfUser(userId, page, pageSize)
	if err != nil {
		return nil, err
	}
	friendIds := make([]uint, len(friends))
	for i, f := range friends {
		friendIds[i] = f.ID
	}
	memberIds, err := c.communityMembershipRepository.FindMemberIdsOfCommunity(communityId, friendIds)
	if err != nil {
		return nil, err
	}
	invitedIds, err := c.communityInvitationRepository.FindInviteeIdsOfCommunity(communityId, friendIds)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindInvitableFriendsResponse{}
	resp.FromFriends(friends, memberIds, invitedIds, page, pageSize, count)
	return resp, nil
}

// InviteToCommunity invites a friend of the inviter to join the community.
// Only members may invite, only their friends, and only someone who is
// neither a member nor already invited (by anyone).
func (c *CommunityMembershipService) InviteToCommunity(inviterId uint, communityId uint, inviteeId uint) error {
	if err := c.requireMember(communityId, inviterId); err != nil {
		return err
	}

	friends, err := c.friendshipRepository.UsersAreFriends(inviterId, inviteeId)
	if err != nil {
		return err
	}
	if !friends {
		return ErrInviteeNotFriend
	}

	_, err = c.communityMembershipRepository.FindOfUserAndCommunity(communityId, inviteeId)
	if err == nil {
		return ErrInviteeAlreadyMember
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	err = c.communityInvitationRepository.Create(model.NewCommunityInvitation(communityId, inviterId, inviteeId))
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrAlreadyInvited
	}
	return err
}

// FindMyInvitations lists the user's pending community invitations, newest
// first.
func (c *CommunityMembershipService) FindMyInvitations(userId uint) (*transport.FindMyCommunityInvitationsResponse, error) {
	invitations, err := c.communityInvitationRepository.FindOfInvitee(userId)
	if err != nil {
		return nil, err
	}
	resp := &transport.FindMyCommunityInvitationsResponse{}
	resp.FromInvitations(invitations)
	return resp, nil
}

// findInvitationOf loads the invitation, which must be addressed to userId.
// Someone else's invitation is reported as missing, not forbidden, so ids
// don't reveal who was invited where.
func (c *CommunityMembershipService) findInvitationOf(userId uint, invitationId uint) (*model.CommunityInvitation, error) {
	invitation, err := c.communityInvitationRepository.FindByIdWithCommunity(invitationId)
	if errors.Is(err, repository.ErrNotFound) || (err == nil && invitation.InviteeID != userId) {
		return nil, ErrInvitationNotFound
	}
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

// AcceptInvitation makes the user a member of the community they were
// invited to. In a private community only the owner's invitation lets them
// straight in; another member's becomes a request for the owner to accept,
// and pending is true.
func (c *CommunityMembershipService) AcceptInvitation(userId uint, invitationId uint) (pending bool, err error) {
	invitation, err := c.findInvitationOf(userId, invitationId)
	if err != nil {
		return false, err
	}
	community := invitation.Community

	_, err = c.communityMembershipRepository.FindOfUserAndCommunity(community.ID, userId)
	if err == nil {
		// They joined some other way meanwhile: the invitation is moot.
		return false, c.communityInvitationRepository.Delete(invitation.ID)
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return false, err
	}

	if community.Private && invitation.InviterID != community.OwnerID {
		err = c.transactor.Transaction(func(tx *repository.Repositories) error {
			if err := tx.CommunityInvitation.Delete(invitation.ID); err != nil {
				return err
			}
			_, err := tx.CommunityJoinRequest.FindOfUserAndCommunity(community.ID, userId)
			if err == nil {
				// They had already asked: that request stands.
				return nil
			}
			if !errors.Is(err, repository.ErrNotFound) {
				return err
			}
			return tx.CommunityJoinRequest.Create(model.NewCommunityJoinRequest(community.ID, userId))
		})
		if errors.Is(err, repository.ErrDuplicate) {
			return true, ErrAlreadyRequested
		}
		return true, err
	}

	return false, c.addMember(community.ID, userId)
}

// DeclineInvitation drops the user's invitation.
func (c *CommunityMembershipService) DeclineInvitation(userId uint, invitationId uint) error {
	invitation, err := c.findInvitationOf(userId, invitationId)
	if err != nil {
		return err
	}
	return c.communityInvitationRepository.Delete(invitation.ID)
}
