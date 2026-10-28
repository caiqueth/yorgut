package service

import (
	"errors"
	"math"
	"sort"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var (
	ErrFriendshipNotFound = newError(ErrNotFound, "Esta amizade não existe.")
	ErrFriendshipExists   = newError(ErrConflict, "Vocês já são amigos ou já há um pedido de amizade entre vocês.")
	ErrFriendshipAnswered = newError(ErrConflict, "Este pedido de amizade já foi aceito.")
	ErrFriendshipNotYours = newError(ErrForbidden, "Esta amizade não é sua.")
	ErrFriendshipSelf     = newError(ErrInvalid, "Você não pode adicionar a si mesmo como amigo.")
)

type FriendshipService struct {
	transactor           repository.Transactor
	friendshipRepository repository.FriendshipRepository
	userRepository       repository.UserRepository
	userReviewRepository repository.UserReviewRepository
	userFanRepository    repository.UserFanRepository
}

func NewFriendshipService(
	transactor repository.Transactor,
	friendshipRepository repository.FriendshipRepository,
	userRepository repository.UserRepository,
	userReviewRepository repository.UserReviewRepository,
	userFanRepository repository.UserFanRepository,
) *FriendshipService {
	return &FriendshipService{
		transactor:           transactor,
		friendshipRepository: friendshipRepository,
		userRepository:       userRepository,
		userReviewRepository: userReviewRepository,
		userFanRepository:    userFanRepository,
	}
}

func (f *FriendshipService) CreateFriendship(userRequestedID uint, userRequesteeID uint) error {
	if userRequestedID == userRequesteeID {
		return ErrFriendshipSelf
	}

	userRequested, err := f.userRepository.FindById(userRequestedID)
	if err != nil {
		return fromRepository(err, "o usuário")
	}

	userRequestee, err := f.userRepository.FindById(userRequesteeID)
	if err != nil {
		return fromRepository(err, "o usuário")
	}

	pending, err := f.friendshipRepository.HasPendingBetweenUsers(userRequestedID, userRequesteeID)
	if err != nil {
		return err
	}
	if pending {
		return ErrFriendshipExists
	}

	friendshipModel := model.NewFriendship(userRequested, userRequestee)

	err = f.friendshipRepository.Create(friendshipModel)
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrFriendshipExists
	}
	return err
}

// ReplyFriendship accepts or refuses a pending request. Only its requestee
// (userId) may answer it.
func (f *FriendshipService) ReplyFriendship(userId uint, friendshipId uint, accepted bool) error {
	friendship, err := f.friendshipRepository.FindByIdWithUsers(friendshipId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrFriendshipNotFound
	}
	if err != nil {
		return err
	}
	if friendship.UserRequesteeID != userId {
		return ErrFriendshipNotYours
	}
	if friendship.Accepted {
		return ErrFriendshipAnswered
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

// deleteFriendship deletes the friendship and recounts both users' friends
// (recounting instead of decrementing is right for pending requests too).
func (f *FriendshipService) deleteFriendship(friendship *model.Friendship) error {
	return f.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Friendship.Delete(friendship.ID); err != nil {
			return err
		}
		return updateFriendCounts(tx, friendship)
	})
}

// areFriends is the rule of every friends-only action (testimonials, reviews,
// fans): the two users are friends, which also rules out being the same user.
func areFriends(friendships repository.FriendshipRepository, userId uint, otherId uint) (bool, error) {
	if userId == otherId {
		return false, nil
	}
	return friendships.UsersAreFriends(userId, otherId)
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
	friendship, err := f.friendshipRepository.FindAcceptedBetweenUsers(userId, friendId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrFriendshipNotFound
	}
	if err != nil {
		return err
	}
	return f.deleteFriendship(friendship)
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

// FindMyFriends returns one page of the user's friends, most recently seen
// first, each with the user's rating of them and fan star, loaded for the
// whole page in two queries (the friends page shows them on every row).
func (f *FriendshipService) FindMyFriends(userId uint, page int, pageSize int) (*transport.FindMyFriendsResponse, error) {
	friends, count, err := f.friendshipRepository.FindFriendsOfUser(userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(friends))
	for i := range friends {
		ids[i] = friends[i].ID
	}
	reviews, err := f.userReviewRepository.FindOfReviewerForUsers(userId, ids)
	if err != nil {
		return nil, err
	}
	fanOf, err := f.userFanRepository.FindIdolIdsOfFan(userId, ids)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindMyFriendsResponse{}
	resp.FromFriends(friends, reviews, fanOf, page, pageSize, count)
	return resp, nil
}

// FindFriendsOfUser returns one page of the user's friends, most recently
// seen first.
func (f *FriendshipService) FindFriendsOfUser(userId uint, page int, pageSize int) (*transport.FindFriendsOfUserResponse, error) {
	friends, count, err := f.friendshipRepository.FindFriendsOfUser(userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindFriendsOfUserResponse{}
	resp.FromFriends(friends, page, pageSize, count)

	return resp, nil
}

// UsersFriendshipStatus tells how the user relates to targetUserId: friends,
// or a request waiting on one of them.
func (f *FriendshipService) UsersFriendshipStatus(userId uint, targetUserId uint) (*transport.AmIFriendOfResponse, error) {
	friendships, err := f.friendshipRepository.FindFriendshipsBetweenUsers(userId, targetUserId)
	if err != nil {
		return nil, err
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
	}, nil
}

// CountIndirectConnections counts the friends of the user's friends who
// aren't their friends already.
func (f *FriendshipService) CountIndirectConnections(userId uint) (uint, error) {
	count, err := f.friendshipRepository.CountFriendsOfFriends(userId)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

// birthdayWindowDays is how far ahead the home page looks for birthdays.
const birthdayWindowDays = 30

// nextBirthday returns the date of the coming birthday of someone born on
// birthDate, today counting as coming, and how many days away it is. Someone
// born on 29 February has it on 1 March in other years.
func nextBirthday(birthDate time.Time, now time.Time) (time.Time, int) {
	today := dayOf(now)
	birth := birthDate.In(time.Local)
	next := time.Date(today.Year(), birth.Month(), birth.Day(), 0, 0, 0, 0, time.Local)
	if next.Before(today) {
		next = time.Date(today.Year()+1, birth.Month(), birth.Day(), 0, 0, 0, 0, time.Local)
	}
	days := int(math.Round(next.Sub(today).Hours() / 24))
	return next, days
}

// FindUpcomingBirthdays lists the user's friends whose birthday is in the
// next 30 days (today included), soonest first. Friends without a birth date
// are left out.
func (f *FriendshipService) FindUpcomingBirthdays(userId uint, now time.Time) (*transport.FindUpcomingBirthdaysResponse, error) {
	friends, err := f.friendshipRepository.FindAllFriendsOfUser(userId)
	if err != nil {
		return nil, err
	}

	resp := &transport.FindUpcomingBirthdaysResponse{Birthdays: []transport.BirthdayItem{}}
	for i := range friends {
		friend := &friends[i]
		if friend.BirthDate == nil {
			continue
		}
		next, days := nextBirthday(*friend.BirthDate, now)
		if days >= birthdayWindowDays {
			continue
		}
		item := transport.BirthdayItem{}
		item.FromUser(friend, next, days)
		resp.Birthdays = append(resp.Birthdays, item)
	}

	sort.Slice(resp.Birthdays, func(i, j int) bool {
		a, b := resp.Birthdays[i], resp.Birthdays[j]
		if a.DaysUntil != b.DaysUntil {
			return a.DaysUntil < b.DaysUntil
		}
		return a.FirstName < b.FirstName
	})
	return resp, nil
}

// HasRatedFriends reports whether the user has rated any of their current
// friends (for the home page's next steps).
func (f *FriendshipService) HasRatedFriends(userId uint) (bool, error) {
	count, err := f.friendshipRepository.CountRatedFriends(userId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
