package service

import (
	"yorgut/internal/repository"
	"yorgut/internal/storage"
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
	ProfileVisit        *ProfileVisitService
	Poll                *PollService
	Info                *InfoService
}

// store keeps the uploaded pictures (storage.NewMemory() in tests).
func NewServices(r *repository.Repositories, store storage.ObjectStore) *Services {
	return &Services{
		Auth:                NewAuthService(*r.User, *r.UserProfile, store),
		User:                NewUserService(r.Transactor, *r.User, *r.Friendship, store),
		UserProfile:         NewUserProfileService(*r.UserProfile, *r.User, store),
		UserReview:          NewUserReviewService(*r.UserReview, *r.Friendship),
		UserFan:             NewUserFanService(*r.UserFan, *r.Friendship),
		Community:           NewCommunityService(r.Transactor, *r.Community, *r.CommunityProfile, *r.CommunityMembership, *r.CommunityJoinRequest, *r.User, store),
		CommunityMembership: NewCommunityMembershipService(r.Transactor, *r.CommunityMembership, *r.CommunityJoinRequest, *r.CommunityInvitation, *r.Community, *r.User, *r.Friendship),
		Friendship:          NewFriendshipService(r.Transactor, *r.Friendship, *r.User, *r.UserReview, *r.UserFan),
		Photo:               NewPhotoService(*r.Photo, *r.PhotoAlbum, *r.User, *r.Friendship, store),
		Scrap:               NewScrapService(*r.Scrap, *r.User, *r.Friendship),
		Testimonial:         NewTestimonialService(*r.Testimonial, *r.User, *r.Friendship),
		Topic:               NewTopicService(*r.Topic, *r.Community, *r.CommunityMembership),
		ProfileVisit:        NewProfileVisitService(*r.ProfileVisit, *r.User),
		Info:                NewInfoService(),
		Poll:                NewPollService(r.Transactor, *r.Poll, *r.PollOption, *r.PollVote, *r.Community, *r.CommunityMembership),
	}
}
