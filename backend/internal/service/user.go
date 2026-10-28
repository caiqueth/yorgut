package service

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
)

var ErrFirstNameRequired = newError(ErrInvalid, "O nome é obrigatório.")

type UserService struct {
	transactor           repository.Transactor
	userRepository       repository.UserRepository
	friendshipRepository repository.FriendshipRepository
	store                storage.ObjectStore
}

func NewUserService(
	transactor repository.Transactor,
	userRepository repository.UserRepository,
	friendshipRepository repository.FriendshipRepository,
	store storage.ObjectStore,
) *UserService {
	return &UserService{
		transactor:           transactor,
		userRepository:       userRepository,
		friendshipRepository: friendshipRepository,
		store:                store,
	}
}

type UpdateUserInput struct {
	FirstName string
	LastName  string
	// BirthDate nil keeps the current one.
	BirthDate *time.Time
}

func (u *UserService) UpdateUser(userId uint, input UpdateUserInput) error {
	firstName := strings.TrimSpace(input.FirstName)
	if firstName == "" {
		return ErrFirstNameRequired
	}

	return u.userRepository.UpdateBasicData(userId, firstName, strings.TrimSpace(input.LastName), input.BirthDate)
}

func (u *UserService) FindUserById(id uint) (*transport.FindUserResponse, error) {
	userData, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, fromRepository(err, "o usuário")
	}

	resp := transport.FindUserResponse{}
	resp.FromUser(userData)

	return &resp, nil

}

func (u *UserService) SearchUsers(name string, page int, pageSize int) (*transport.SearchUsersResponse, error) {
	users, count, err := u.userRepository.SearchUsers(name, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &transport.SearchUsersResponse{}
	resp.FromUsers(users, page, pageSize, count)

	return resp, nil
}

// FindUserWithProfileById returns the profile as viewerId may see it (see
// privacy.go): non-friends don't get the fields and birthday kept to friends.
func (u *UserService) FindUserWithProfileById(viewerId uint, id uint) (*transport.FindUserWithProfileByIdResponse, error) {
	userData, err := u.userRepository.FindUserWithProfile(id)
	if err != nil {
		return nil, fromRepository(err, "o usuário")
	}

	friend := viewerId != 0 && viewerId == id
	if !friend {
		friend, err = areFriends(u.friendshipRepository, viewerId, id)
		if err != nil {
			return nil, err
		}
	}

	resp := transport.FindUserWithProfileByIdResponse{}
	resp.FromUserFor(userData, friend)

	return &resp, nil
}

var ErrUnknownProfileField = newError(ErrInvalid, "Campo do perfil desconhecido.")

// PrivacyInput holds every privacy setting (see model.User).
type PrivacyInput struct {
	RecordVisits        bool
	PublicAlbum         bool
	PublicCommunities   bool
	PublicScraps        bool
	PublicTestimonials  bool
	PublicBirthday      bool
	HiddenProfileFields []string
}

// UpdatePrivacy replaces the user's privacy settings.
func (u *UserService) UpdatePrivacy(userId uint, input PrivacyInput) error {
	known := map[string]bool{}
	for _, key := range transport.ProfileFieldKeys {
		known[key] = true
	}
	// Kept in ProfileFieldKeys order, once each.
	hidden := map[string]bool{}
	for _, key := range input.HiddenProfileFields {
		if !known[key] {
			return ErrUnknownProfileField
		}
		hidden[key] = true
	}
	keys := []string{}
	for _, key := range transport.ProfileFieldKeys {
		if hidden[key] {
			keys = append(keys, key)
		}
	}

	return u.userRepository.UpdatePrivacy(userId, &model.User{
		RecordVisits:        input.RecordVisits,
		PublicAlbum:         input.PublicAlbum,
		PublicCommunities:   input.PublicCommunities,
		PublicScraps:        input.PublicScraps,
		PublicTestimonials:  input.PublicTestimonials,
		PublicBirthday:      input.PublicBirthday,
		HiddenProfileFields: strings.Join(keys, ","),
	})
}

var ErrUnknownTheme = newError(ErrInvalid, "Tema desconhecido.")

// UpdateTheme changes the look of the app for the user.
func (u *UserService) UpdateTheme(userId uint, theme string) error {
	if !slices.Contains(transport.ThemeKeys, theme) {
		return ErrUnknownTheme
	}
	return u.userRepository.UpdateTheme(userId, theme)
}

// DeleteAccount deletes the user and everything they put in Yorgut, for good
// (UserRepository.DeleteForever). Each community they own goes to its
// longest-standing other member; one nobody else is in is deleted with
// everything in it. Their friends' and communities' counts are recounted.
func (u *UserService) DeleteAccount(ctx context.Context, userId uint) error {
	// The files of the pictures that go, deleted once the rows are gone.
	var pictures []string
	err := u.transactor.Transaction(func(tx *repository.Repositories) error {
		pictures = nil

		user, err := tx.User.FindUserWithProfile(userId)
		if err != nil {
			return fromRepository(err, "o usuário")
		}
		if user.Profile != nil && user.Profile.ProfilePicture != nil {
			pictures = append(pictures, *user.Profile.ProfilePicture)
		}
		photos, err := tx.Photo.FindOfUser(userId)
		if err != nil {
			return err
		}
		for _, photo := range photos {
			pictures = append(pictures, photo.URL)
		}

		owned, err := tx.Community.FindOwnedBy(userId)
		if err != nil {
			return err
		}
		deleted := map[uint]bool{}
		for _, community := range owned {
			heir, err := tx.CommunityMembership.FindOldestMemberExcept(community.ID, userId)
			if err == nil {
				if err := tx.Community.TransferOwnership(community.ID, heir); err != nil {
					return err
				}
				continue
			}
			if !errors.Is(err, repository.ErrNotFound) {
				return err
			}
			profile, err := tx.CommunityProfile.FindOfCommunity(community.ID)
			if err == nil && profile.ProfilePicture != nil {
				pictures = append(pictures, *profile.ProfilePicture)
			} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
				return err
			}
			if err := tx.Community.DeleteForever(community.ID); err != nil {
				return err
			}
			deleted[community.ID] = true
		}

		// Read before the rows that say so are gone.
		friends, err := tx.Friendship.FindAllFriendsOfUser(userId)
		if err != nil {
			return err
		}
		communityIds, err := tx.CommunityMembership.FindCommunityIdsOfUser(userId)
		if err != nil {
			return err
		}

		if err := tx.User.DeleteForever(userId); err != nil {
			return err
		}

		for _, friend := range friends {
			if err := tx.User.UpdateFriendCount(friend.ID); err != nil {
				return err
			}
		}
		for _, communityId := range communityIds {
			if deleted[communityId] {
				continue
			}
			if err := tx.Community.UpdateMemberCount(communityId); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for i := range pictures {
		deletePicture(ctx, u.store, &pictures[i])
	}
	return nil
}
