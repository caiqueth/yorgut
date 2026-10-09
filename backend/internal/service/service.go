package service

import (
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

// Services holds one instance of every service, so the API server and the
// seed command are wired the same way.
type Services struct {
	Auth                *AuthService
	User                *UserService
	UserProfile         *UserProfileService
	UserReview          *UserReviewService
	UserFan             *UserFanService
	Community           *CommunityService
	CommunityMembership *CommunityMembershipService
	Friendship          *FriendshipService
	Photo               *PhotoService
	Scrap               *ScrapService
	Testimonial         *TestimonialService
	Topic               *TopicService
}

func NewServices(r *repository.Repositories, validate *validator.Validate) *Services {
	return &Services{
		Auth:                NewAuthService(*r.User, validate),
		User:                NewUserService(*r.User, validate),
		UserProfile:         NewUserProfileService(*r.UserProfile, *r.User, validate),
		UserReview:          NewUserReviewService(*r.UserReview, *r.Friendship, validate),
		UserFan:             NewUserFanService(*r.UserFan, *r.Friendship, validate),
		Community:           NewCommunityService(*r.Community, *r.CommunityProfile, *r.CommunityMembership, *r.User, validate),
		CommunityMembership: NewCommunityMembershipService(*r.CommunityMembership, *r.Community, *r.User, validate),
		Friendship:          NewFriendshipService(*r.Friendship, *r.User, validate),
		Photo:               NewPhotoService(*r.Photo, *r.PhotoAlbum, validate),
		Scrap:               NewScrapService(*r.Scrap, *r.User, validate),
		Testimonial:         NewTestimonialService(*r.Testimonial, *r.User, *r.Friendship, validate),
		Topic:               NewTopicService(*r.Topic, *r.CommunityMembership, validate),
	}
}
