package service

import (
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var ErrReviewNotFriend = newError(ErrForbidden, "Só amigos podem avaliar.")

type UserReviewService struct {
	userReviewRepository repository.UserReviewRepository
	friendshipRepository repository.FriendshipRepository
}

func NewUserReviewService(
	userReviewRepository repository.UserReviewRepository,
	friendshipRepository repository.FriendshipRepository,
) *UserReviewService {
	return &UserReviewService{
		userReviewRepository: userReviewRepository,
		friendshipRepository: friendshipRepository,
	}
}

// UpsertReview saves the reviewer's ratings of a friend, replacing earlier
// ones. Only friends review each other, never themselves.
func (u *UserReviewService) UpsertReview(reviewerId uint, reviewedId uint, cool uint8, sexy uint8, trustworthy uint8) error {
	friends, err := areFriends(u.friendshipRepository, reviewerId, reviewedId)
	if err != nil {
		return err
	}
	if !friends {
		return ErrReviewNotFriend
	}
	review := model.NewUserReview(reviewerId, reviewedId, cool, sexy, trustworthy)
	return u.userReviewRepository.Upsert(review)
}

func (u *UserReviewService) FindSummaryOfUser(viewerId uint, reviewedId uint) (*transport.UserReviewSummaryResponse, error) {
	averages, err := u.userReviewRepository.AveragesOfUser(reviewedId)
	if err != nil {
		return nil, err
	}

	canReview, err := areFriends(u.friendshipRepository, viewerId, reviewedId)
	if err != nil {
		return nil, err
	}

	response := transport.UserReviewSummaryResponse{
		Cool:        averages.Cool,
		Sexy:        averages.Sexy,
		Trustworthy: averages.Trustworthy,
		ReviewCount: averages.ReviewCount,
		CanReview:   canReview,
	}

	myReview, err := u.userReviewRepository.FindOfReviewerAndReviewed(viewerId, reviewedId)
	if err != nil {
		return nil, err
	}

	if myReview != nil {
		response.MyReview = &transport.UserReviewItem{}
		response.MyReview.FromUserReview(myReview)
	}

	return &response, nil
}
