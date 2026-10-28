package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
)

var (
	ErrCommunityNotFound  = newError(ErrNotFound, "A comunidade não existe.")
	ErrCommunityForbidden = newError(ErrForbidden, "Só o dono pode editar esta comunidade.")
)

type CommunityService struct {
	transactor                 repository.Transactor
	communityRepository        repository.CommunityRepository
	communityProfileRepository repository.CommunityProfileRepository
	membershipRepository       repository.CommunityMembershipRepository
	joinRequestRepository      repository.CommunityJoinRequestRepository
	userRepository             repository.UserRepository
	store                      storage.ObjectStore
}

func NewCommunityService(
	transactor repository.Transactor,
	communityRepository repository.CommunityRepository,
	communityProfileRepository repository.CommunityProfileRepository,
	membershipRepository repository.CommunityMembershipRepository,
	joinRequestRepository repository.CommunityJoinRequestRepository,
	userRepository repository.UserRepository,
	store storage.ObjectStore,
) *CommunityService {
	return &CommunityService{
		transactor:                 transactor,
		communityRepository:        communityRepository,
		communityProfileRepository: communityProfileRepository,
		membershipRepository:       membershipRepository,
		joinRequestRepository:      joinRequestRepository,
		userRepository:             userRepository,
		store:                      store,
	}
}

// CreateCommunity creates the community along with its profile (a blank
// description leaves it without one), using picture (optional; same limits
// as every upload, see picture.go) as its image, and makes the creator its
// first member. It returns the new community's id.
func (c *CommunityService) CreateCommunity(ctx context.Context, creatorId uint, name string, category string, description string, private bool, picture []byte) (uint, error) {
	creator, err := c.userRepository.FindById(creatorId)
	if err != nil {
		return 0, fromRepository(err, "o usuário")
	}

	profile := &model.CommunityProfile{Description: nilIfBlank(description)}
	if len(picture) > 0 {
		picture, err = preparePicture(picture)
		if err != nil {
			return 0, err
		}
		key, err := storage.PutImage(ctx, c.store, picture)
		if err != nil {
			return 0, err
		}
		profile.ProfilePicture = &key
	}

	communityModel := model.NewCommunity(strings.TrimSpace(name), strings.TrimSpace(category), creator)
	communityModel.Profile = profile
	communityModel.Private = private
	err = c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Community.Create(communityModel); err != nil {
			return err
		}
		// The creator owns the community and is its first member.
		if err := tx.CommunityMembership.Create(model.NewCommunityMembership(communityModel, creator)); err != nil {
			return err
		}
		return tx.Community.UpdateMemberCount(communityModel.ID)
	})
	if err != nil {
		// Nothing points to the picture anymore.
		deletePicture(ctx, c.store, profile.ProfilePicture)
		return 0, fromRepository(err, "a comunidade")
	}
	return communityModel.ID, nil
}

func (c *CommunityService) FindCommunityById(id uint) (*transport.FindCommunityByIdResponse, error) {
	community, err := c.communityRepository.FindByIdEager(id)
	if err != nil {
		return nil, fromRepository(err, "a comunidade")
	}

	response := transport.FindCommunityByIdResponse{}
	response.FromCommunity(community)

	return &response, nil
}

func (c *CommunityService) FindCommunityProfileById(communityId uint) (*transport.FindCommunityProfileByIdResponse, error) {
	response := transport.FindCommunityProfileByIdResponse{}

	profile, err := c.communityProfileRepository.FindOfCommunity(communityId)
	// Communities may have no profile row yet; that's an empty profile, not an error.
	if errors.Is(err, repository.ErrNotFound) {
		return &response, nil
	}
	if err != nil {
		return nil, err
	}

	response.FromCommunityProfile(profile)

	return &response, nil
}

// FindCommunityMembers returns one page of the community's members, most
// recently seen first.
func (c *CommunityService) FindCommunityMembers(communityId uint, page int, pageSize int) (*transport.FindCommunityMembersResponse, error) {
	if _, err := c.communityRepository.FindById(communityId); err != nil {
		return nil, fromRepository(err, "a comunidade")
	}

	members, count, err := c.membershipRepository.FindMembersOfCommunity(communityId, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := transport.FindCommunityMembersResponse{}
	response.FromMembers(members, page, pageSize, count)

	return &response, nil
}

// relatedCommunitiesShown is how many communities the "comunidades
// relacionadas" box shows.
const relatedCommunitiesShown = 9

// FindRelatedCommunities returns a few random communities other than this one,
// the same ones all day. Communities don't record related ones yet, so any
// community will do.
func (c *CommunityService) FindRelatedCommunities(communityId uint) (*transport.FindManyCommunitiesResponse, error) {
	if _, err := c.communityRepository.FindById(communityId); err != nil {
		return nil, fromRepository(err, "a comunidade")
	}

	communities, err := c.communityRepository.FindShuffledExcept(communityId, relatedCommunitiesShown, dailySeed(communityId, time.Now()))
	if err != nil {
		return nil, err
	}

	resp := transport.FindManyCommunitiesResponse{}
	resp.FromCommunities(communities)

	return &resp, nil
}

// FindCommunitiesCreatedByUser returns one page of the communities the user
// created, by name.
func (c *CommunityService) FindCommunitiesCreatedByUser(userId uint, page int, pageSize int) (*transport.FindMyOwnedCommunitiesResponse, error) {
	communities, count, err := c.communityRepository.FindCreatedByUser(userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindMyOwnedCommunitiesResponse{}
	resp.FromCommunities(communities, page, pageSize, count)

	return resp, nil
}

// SearchCommunitiesToDiscover lists communities the user hasn't joined yet,
// optionally filtered by name.
func (c *CommunityService) SearchCommunitiesToDiscover(userId uint, name string, page int, pageSize int) (*transport.SearchCommunitiesResponse, error) {
	communities, count, err := c.communityRepository.SearchNotMemberOf(userId, name, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &transport.SearchCommunitiesResponse{}
	resp.FromCommunities(communities, page, pageSize, count)

	// Private communities the user already asked to join.
	ids := make([]uint, len(communities))
	for i := range communities {
		ids[i] = communities[i].ID
	}
	pendingIds, err := c.joinRequestRepository.FindCommunityIdsOfUser(userId, ids)
	if err != nil {
		return nil, err
	}
	pending := make(map[uint]bool, len(pendingIds))
	for _, id := range pendingIds {
		pending[id] = true
	}
	for i := range resp.Communities {
		resp.Communities[i].Pending = pending[resp.Communities[i].ID]
	}

	return resp, nil
}

// findOwnedCommunity loads the community and checks that userId owns it.
func (c *CommunityService) findOwnedCommunity(userId uint, communityId uint) (*model.Community, error) {
	community, err := c.communityRepository.FindById(communityId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrCommunityNotFound
	}
	if err != nil {
		return nil, err
	}
	if community.OwnerID != userId {
		return nil, ErrCommunityForbidden
	}
	return community, nil
}

// UpdateCommunity replaces the community's name, category, description and
// privacy. Only the owner may; a blank description clears it. Making a
// private community public accepts everyone who asked to join.
func (c *CommunityService) UpdateCommunity(userId uint, communityId uint, name string, category string, description string, private bool) error {
	community, err := c.findOwnedCommunity(userId, communityId)
	if err != nil {
		return err
	}

	return c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Community.UpdateDetails(communityId, strings.TrimSpace(name), strings.TrimSpace(category), private); err != nil {
			return err
		}
		if err := tx.CommunityProfile.UpdateDescriptionOfCommunity(communityId, nilIfBlank(description)); err != nil {
			return err
		}
		if community.Private && !private {
			return acceptAllJoinRequests(tx, communityId)
		}
		return nil
	})
}

// acceptAllJoinRequests turns every pending request to join the community
// into a membership. Runs inside the caller's transaction.
func acceptAllJoinRequests(tx *repository.Repositories, communityId uint) error {
	requests, err := tx.CommunityJoinRequest.FindOfCommunity(communityId)
	if err != nil {
		return err
	}
	if len(requests) == 0 {
		return nil
	}
	for _, request := range requests {
		// Join refuses members and accepting deletes the request, so the
		// user can't already be a member here.
		membership := &model.CommunityMembership{CommunityID: communityId, UserID: request.UserID}
		if err := tx.CommunityMembership.Create(membership); err != nil {
			return err
		}
	}
	if err := tx.CommunityJoinRequest.DeleteOfCommunity(communityId); err != nil {
		return err
	}
	return tx.Community.UpdateMemberCount(communityId)
}

// UpdateCommunityPicture replaces the community's picture (limits in
// picture.go). Only the owner may.
func (c *CommunityService) UpdateCommunityPicture(ctx context.Context, userId uint, communityId uint, picture []byte) error {
	if _, err := c.findOwnedCommunity(userId, communityId); err != nil {
		return err
	}

	picture, err := preparePicture(picture)
	if err != nil {
		return err
	}

	key, err := storage.PutImage(ctx, c.store, picture)
	if err != nil {
		return err
	}

	previous, err := c.communityProfileRepository.UpdatePictureOfCommunity(communityId, key)
	if err != nil {
		deletePicture(ctx, c.store, &key)
		return err
	}
	deletePicture(ctx, c.store, previous)
	return nil
}
