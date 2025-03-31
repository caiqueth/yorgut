package service

import (
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type TestimonialService struct {
	TestimonialRepository repository.TestimonialRepository
	UserRepository        repository.UserRepository
	Validate              *validator.Validate
}

type UserService struct {
	userRepository                repository.UserRepository
	communityRepository           repository.CommunityRepository
	communityMembershipRepository repository.CommunityMembershipRepository
	photoAlbumRepository          repository.PhotoAlbumRepository
	friendshipRepository          repository.FriendshipRepository
	Validate                      *validator.Validate
}

type ScrapService struct {
	scrapRepository repository.ScrapRepository
	userRepository  repository.UserRepository
	Validate        *validator.Validate
}

type PhotoService struct {
	photoRepository      repository.PhotoRepository
	photoAlbumRepository repository.PhotoAlbumRepository
	Validate             *validator.Validate
}

type PhotoAlbumService struct {
	photoAlbumRepository repository.PhotoAlbumRepository
	userRepository       repository.UserRepository
	Validate             *validator.Validate
}

type FriendshipService struct {
	friendshipRepository repository.FriendshipRepository
	userRepository       repository.UserRepository
	Validate             *validator.Validate
}

type CommunityService struct {
	communityRepository repository.CommunityRepository
	userRepository      repository.UserRepository
	Validate            *validator.Validate
}

type CommunityMembershipService struct {
	communityMembershipRepository repository.CommunityMembershipRepository
	communityRepository           repository.CommunityRepository
	userRepository                repository.UserRepository
	Validate                      *validator.Validate
}

func NewCommunityMembershipService(
	communityMembershipRepository repository.CommunityMembershipRepository,
	validate *validator.Validate,
) *CommunityMembershipService {
	return &CommunityMembershipService{
		communityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}

func NewCommunityService(
	communityRepository repository.CommunityRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *CommunityService {
	return &CommunityService{
		communityRepository: communityRepository,
		userRepository:      userRepository,
		Validate:            validate,
	}
}

func NewFriendshipService(
	friendshipRepository repository.FriendshipRepository,
	validate *validator.Validate,
) *FriendshipService {
	return &FriendshipService{
		friendshipRepository: friendshipRepository,
		Validate:             validate,
	}
}

func NewPhotoAlbumService(
	photoAlbumRepository repository.PhotoAlbumRepository,
	validate *validator.Validate,
) *PhotoAlbumService {
	return &PhotoAlbumService{
		photoAlbumRepository: photoAlbumRepository,
		Validate:             validate,
	}
}

func NewPhotoService(photoRepository repository.PhotoRepository, validate *validator.Validate) *PhotoService {
	return &PhotoService{
		photoRepository: photoRepository,
		Validate:        validate,
	}
}

func NewScrapService(scrapRepository repository.ScrapRepository, validate *validator.Validate) *ScrapService {
	return &ScrapService{
		scrapRepository: scrapRepository,
		Validate:        validate,
	}
}

func NewTestimonialService(testimonialRepository repository.TestimonialRepository, validate *validator.Validate) *TestimonialService {
	return &TestimonialService{
		TestimonialRepository: testimonialRepository,
		Validate:              validate,
	}
}

func NewUserService(
	userRepository repository.UserRepository,
	communityRepository repository.CommunityRepository,
	communityMembershipRepository repository.CommunityMembershipRepository,
	photoAlbumRepository repository.PhotoAlbumRepository,
	friendshipRepository repository.FriendshipRepository,
	validate *validator.Validate,
) *UserService {
	return &UserService{
		userRepository:                userRepository,
		communityRepository:           communityRepository,
		communityMembershipRepository: communityMembershipRepository,
		photoAlbumRepository:          photoAlbumRepository,
		friendshipRepository:          friendshipRepository,
		Validate:                      validate,
	}
}
