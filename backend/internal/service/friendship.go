package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var (
	ErrFriendshipNotFound = newError(ErrNotFound, "friendship not found")
	ErrFriendshipExists   = newError(ErrConflict, "there is already a friendship or request between these users")
)

type FriendshipService struct {
	transactor           repository.Transactor
	friendshipRepository repository.FriendshipRepository
	userRepository       repository.UserRepository
	Validate             *validator.Validate
}

func NewFriendshipService(
	transactor repository.Transactor,
	friendshipRepository repository.FriendshipRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *FriendshipService {
	return &FriendshipService{
		transactor:           transactor,
		friendshipRepository: friendshipRepository,
		userRepository:       userRepository,
		Validate:             validate,
	}
}

func (f *FriendshipService) CreateFriendship(userRequestedID uint, userRequesteeID uint) error {
	userRequested, err := f.userRepository.FindById(userRequestedID)
	if err != nil {
		return fromRepository(err, "user that requested the friendship")
	}

	userRequestee, err := f.userRepository.FindById(userRequesteeID)
	if err != nil {
		return fromRepository(err, "target user (requestee) of the friendship")
	}

	if f.friendshipRepository.UsersHavePendingFriendshipInvites(userRequestedID, userRequesteeID) {
		return ErrFriendshipExists
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.friendshipRepository.Create(friendshipModel)
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrFriendshipExists
	}
	return err
}

func (f *FriendshipService) ReplyFriendship(friendshipId uint, accepted bool) error {
	friendship, err := f.friendshipRepository.FindPendingFriendshipWithUsersById(friendshipId)
	if err != nil {
		return fromRepository(err, "friendship")
	}

	if !accepted {
		// A pending request doesn't count as a friend, so no counts change.
		return f.friendshipRepository.Delete(friendshipId)
	}

	return f.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Friendship.AcceptFriendship(friendshipId); err != nil {
			return err
		}
		return updateFriendCounts(tx, friendship)
	})
}

func (f *FriendshipService) ShouldThisUserBeReplying(friendshipId uint, userId uint) bool {
	friendship, err := f.friendshipRepository.FindById(friendshipId)
	if err != nil {
		return false
	}

	if friendship.UserRequesteeID == userId {
		return true
	}

	return false
}

func (f *FriendshipService) ShouldThisUserBeDeleting(friendshipId uint, userId uint) bool {
	friendship, err := f.friendshipRepository.FindById(friendshipId)
	if err != nil {
		return false
	}

	if friendship.UserRequesteeID == userId || friendship.UserRequestedID == userId {
		return true
	}

	return false
}

func (f *FriendshipService) DeleteFriendship(friendshipId uint) error {
	friendship, err := f.friendshipRepository.FindById(friendshipId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrFriendshipNotFound
	}
	if err != nil {
		return err
	}

	// Also used to cancel pending requests: recounting (instead of
	// decrementing) keeps the counts right in both cases.
	return f.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Friendship.Delete(friendship.ID); err != nil {
			return err
		}
		return updateFriendCounts(tx, friendship)
	})
}

// updateFriendCounts recounts the friends of both users of a friendship.
func updateFriendCounts(tx *repository.Repositories, friendship *model.Friendship) error {
	for _, userId := range []uint{friendship.UserRequestedID, friendship.UserRequesteeID} {
		if err := tx.User.UpdateFriendCount(userId); err != nil {
			return err
		}
	}
	return nil
}

// RemoveFriend deletes the accepted friendship between two users, returning
// ErrFriendshipNotFound when they are not friends.
func (f *FriendshipService) RemoveFriend(userId uint, friendId uint) error {
	friendship, err := f.friendshipRepository.FindActiveAcceptedFriends(userId, friendId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrFriendshipNotFound
	}
	if err != nil {
		return err
	}
	return f.DeleteFriendship(friendship.ID)
}

func (f *FriendshipService) FindPendingFriendshipsOfUser(userId uint) (*transport.FindPendingFriendshipOfUserResponse, error) {
	invites, err := f.friendshipRepository.FindPendingFriendshipsOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindPendingFriendshipOfUserResponse{}
	resp.FromUsers(invites)

	return resp, nil
}

func (f *FriendshipService) FindFriendsOfUser(userId uint) (*transport.FindFriendsOfUserResponse, error) {
	friends, err := f.friendshipRepository.FindFriendsOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindFriendsOfUserResponse{}
	resp.FromFriends(friends)

	return resp, nil
}

func (f *FriendshipService) UsersFriendshipStatus(userId uint, targetUserId uint) *transport.AmIFriendOfResponse {
	if userId == targetUserId {
		return nil
	}

	friendships, err := f.friendshipRepository.FindFriendshipsBetweenUsers(userId, targetUserId)
	if err != nil {
		return nil
	}

	var friendshipId *uint
	usersAreFriends := false
	isPendingMyReply := false
	isPendingTheirReply := false

	for _, fsp := range friendships {
		friendshipId = &fsp.ID
		if fsp.Accepted {
			usersAreFriends = true
		} else {
			if fsp.UserRequestedID == userId {
				isPendingTheirReply = true
			}
			if fsp.UserRequestedID == targetUserId {
				isPendingMyReply = true
			}
		}
	}

	return &transport.AmIFriendOfResponse{
		FriendshipId:        friendshipId,
		AmIFriend:           usersAreFriends,
		IsPendingMyReply:    isPendingMyReply,
		IsPendingTheirReply: isPendingTheirReply,
	}

}
