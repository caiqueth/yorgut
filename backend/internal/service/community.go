package service

import (
	"errors"
	"strings"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var (
	ErrCommunityNotFound  = newError(ErrNotFound, "community not found")
	ErrCommunityForbidden = newError(ErrForbidden, "only the owner can edit this community")
)

type CommunityService struct {
	transactor                 repository.Transactor
	communityRepository        repository.CommunityRepository
	communityProfileRepository repository.CommunityProfileRepository
	membershipRepository       repository.CommunityMembershipRepository
	userRepository             repository.UserRepository
	Validate                   *validator.Validate
}

func NewCommunityService(
	transactor repository.Transactor,
	communityRepository repository.CommunityRepository,
	communityProfileRepository repository.CommunityProfileRepository,
	membershipRepository repository.CommunityMembershipRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *CommunityService {
	return &CommunityService{
		transactor:                 transactor,
		communityRepository:        communityRepository,
		communityProfileRepository: communityProfileRepository,
		membershipRepository:       membershipRepository,
		userRepository:             userRepository,
		Validate:                   validate,
	}
}

// CreateCommunity creates the community together with its profile, using
// picture (optional; same limits as every upload, see picture.go) as its
// image, and makes the creator its first member. It returns the new
// community's id.
func (c *CommunityService) CreateCommunity(creatorId uint, name string, category string, picture []byte) (uint, error) {
	creator, err := c.userRepository.FindById(creatorId)
	if err != nil {
		return 0, fromRepository(err, "community creator")
	}

	profile := &model.CommunityProfile{}
	if len(picture) > 0 {
		if err := validatePicture(picture); err != nil {
			return 0, err
		}
		storedFileName, err := helpers.UploadPicture(picture)
		if err != nil {
			return 0, err
		}
		profile.ProfilePicture = storedFileName
	}

	communityModel := model.NewCommunity(strings.TrimSpace(name), strings.TrimSpace(category), creator)
	communityModel.Profile = profile
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
		return 0, fromRepository(err, "community")
	}
	return communityModel.ID, nil
}

func (c *CommunityService) DeleteCommunity(id uint) error {
	err := c.communityRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) FindCommunityById(id uint) (*transport.FindCommunityByIdResponse, error) {
	community, err := c.communityRepository.FindByIdEager(id)
	if err != nil {
		return nil, fromRepository(err, "community")
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

func (c *CommunityService) FindManyCommunities() (*transport.FindManyCommunitiesResponse, error) {
	communities, err := c.communityRepository.FindCommunitiesWithProfiles()
	if err != nil {
		return nil, err
	}

	resp := transport.FindManyCommunitiesResponse{}
	resp.FromCommunities(communities)

	return &resp, nil
}

func (c *CommunityService) FindCommunityMembers(communityId uint) (*transport.FindCommunityMembersResponse, error) {
	memberships, err := c.membershipRepository.FindMembersOfCommunity(communityId)
	if err != nil {
		return nil, err
	}

	community, err := c.communityRepository.FindById(communityId)
	if err != nil {
		return nil, fromRepository(err, "community")
	}

	response := transport.FindCommunityMembersResponse{}
	response.FromMemberships(community, memberships)

	return &response, nil
}

func (c *CommunityService) FindCommunitiesCreatedByUser(userId uint) (*transport.FindMyOwnedCommunitiesResponse, error) {
	communities, err := c.communityRepository.FindCreatedByUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindMyOwnedCommunitiesResponse{}
	resp.FromCommunities(communities)

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

	return resp, nil
}

// findOwnedCommunity loads the community and checks that userId owns it.
func (c *CommunityService) findOwnedCommunity(userId uint, communityId uint) error {
	community, err := c.communityRepository.FindById(communityId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrCommunityNotFound
	}
	if err != nil {
		return err
	}
	if community.OwnerID != userId {
		return ErrCommunityForbidden
	}
	return nil
}

// UpdateCommunity replaces the community's name, category and description.
// Only the owner may; a blank description clears it.
func (c *CommunityService) UpdateCommunity(userId uint, communityId uint, name string, category string, description string) error {
	if err := c.findOwnedCommunity(userId, communityId); err != nil {
		return err
	}

	return c.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Community.UpdateDetails(communityId, strings.TrimSpace(name), strings.TrimSpace(category)); err != nil {
			return err
		}
		return tx.CommunityProfile.UpdateDescriptionOfCommunity(communityId, nilIfBlank(description))
	})
}

// UpdateCommunityPicture replaces the community's picture (limits in
// picture.go). Only the owner may.
func (c *CommunityService) UpdateCommunityPicture(userId uint, communityId uint, picture []byte) error {
	if err := c.findOwnedCommunity(userId, communityId); err != nil {
		return err
	}

	if err := validatePicture(picture); err != nil {
		return err
	}

	storedFileName, err := helpers.UploadPicture(picture)
	if err != nil {
		return err
	}

	return c.communityProfileRepository.UpdatePictureOfCommunity(communityId, *storedFileName)
}
