package service

import (
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

type UserReviewService struct {
	userReviewRepository repository.UserReviewRepository
	friendshipRepository repository.FriendshipRepository
	Validate             *validator.Validate
}

func NewUserReviewService(
	userReviewRepository repository.UserReviewRepository,
	friendshipRepository repository.FriendshipRepository,
	validate *validator.Validate,
) *UserReviewService {
	return &UserReviewService{
		userReviewRepository: userReviewRepository,
		friendshipRepository: friendshipRepository,
		Validate:             validate,
	}
}

func (u *UserReviewService) CanReview(reviewerId uint, reviewedId uint) bool {
	if reviewerId == reviewedId {
		return false
	}
	return u.friendshipRepository.UsersAreFriends(reviewerId, reviewedId)
}

func (u *UserReviewService) UpsertReview(reviewerId uint, reviewedId uint, cool uint8, sexy uint8, trustworthy uint8) error {
	review := model.NewUserReview(reviewerId, reviewedId, cool, sexy, trustworthy)
	return u.userReviewRepository.Upsert(review)
}

func (u *UserReviewService) FindSummaryOfUser(viewerId uint, reviewedId uint) (*transport.UserReviewSummaryResponse, error) {
	averages, err := u.userReviewRepository.AveragesOfUser(reviewedId)
	if err != nil {
		return nil, err
	}

	response := transport.UserReviewSummaryResponse{
		Cool:        averages.Cool,
		Sexy:        averages.Sexy,
		Trustworthy: averages.Trustworthy,
		ReviewCount: averages.ReviewCount,
		CanReview:   u.CanReview(viewerId, reviewedId),
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
