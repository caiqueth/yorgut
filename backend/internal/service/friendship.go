package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var ErrFriendshipNotFound = errors.New("friendship not found")

type FriendshipService struct {
	friendshipRepository repository.FriendshipRepository
	userRepository       repository.UserRepository
	Validate             *validator.Validate
}

func NewFriendshipService(
	friendshipRepository repository.FriendshipRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *FriendshipService {
	return &FriendshipService{
		friendshipRepository: friendshipRepository,
		userRepository:       userRepository,
		Validate:             validate,
	}
}

func (f *FriendshipService) CreateFriendship(userRequestedID uint, userRequesteeID uint) error {
	userRequested, err := f.userRepository.FindById(userRequestedID)
	if err != nil {
		return errors.New("could not find user that requested the friendship")
	}

	userRequestee, err := f.userRepository.FindById(userRequesteeID)
	if err != nil {
		return errors.New("could not find target user (requestee) of the friendship")
	}

	if f.friendshipRepository.UsersHavePendingFriendshipInvites(userRequestedID, userRequesteeID) {
		return errors.New("there is already a friendship or request between these users")
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.friendshipRepository.Create(friendshipModel)
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendshipService) ReplyFriendship(friendshipId uint, accepted bool) error {
	friendship, err := f.friendshipRepository.FindPendingFriendshipWithUsersById(friendshipId)
	if err != nil {
		return err
	}

	if !accepted {
		if err := f.friendshipRepository.Delete(friendshipId); err != nil {
			return err
		}
		return nil
	}

	if err := f.friendshipRepository.AcceptFriendship(friendshipId); err != nil {
		return err
	}

	for _, userId := range []uint{friendship.UserRequestedID, friendship.UserRequesteeID} {
		count := f.friendshipRepository.CountFriends(userId)
		if err := f.userRepository.UpdateFriendCount(userId, uint(count)); err != nil {
			return err
		}
	}

	return nil
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
	friendship, err := f.friendshipRepository.FindByIdWithUsers(friendshipId)
	if err != nil {
		return errors.New("friendship not found")
	}

	if err := f.friendshipRepository.Delete(friendship.ID); err != nil {
		return err
	}

	f.userRepository.IncrementFriendCount(friendship.UserRequested, -1)
	f.userRepository.IncrementFriendCount(friendship.UserRequestee, -1)

	return nil
}

// RemoveFriend deletes the accepted friendship between two users, returning
// ErrFriendshipNotFound when they are not friends.
func (f *FriendshipService) RemoveFriend(userId uint, friendId uint) error {
	friendship, err := f.friendshipRepository.FindActiveAcceptedFriends(userId, friendId)
	if err != nil {
		return ErrFriendshipNotFound
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
