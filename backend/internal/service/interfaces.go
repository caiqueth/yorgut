package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
)

type CommunityService interface {
	Create(community request.CreateCommunityRequest) response.Response
	Update(user request.UpdateCommunityRequest) response.Response
	Delete(communityId int)
	FindById(communityId int) response.FindCommunityByIdResponse
	FindMany() response.Response
}

type CommunityMembershipService interface {
	Create(community request.CreateCommunityMembershipRequest)
	Delete(id int)
	FindOfUser(userId int) response.FindCommunityMembershipOfUserResponse
}

type FriendshipService interface {
	Create(friendship request.CreateFriendshipRequest)
	Delete(id int)
	FindOfUser(userId int) response.FindFriendshipOfUserResponse
}

type PhotoAlbumService interface {
	Create(photoAlbum request.CreatePhotoAlbumRequest)
	Delete(id int)
	FindOfUser(userId int) response.FindPhotoAlmbumOfUserResponse
}

type PhotoService interface {
	Create(photo request.CreatePhotoRequest)
	Delete(id int)
	FindOfUser(userId int) response.FindPhotoOfUserResponse
	FindOfAlbum(photoAlbumId int) response.FindPhotoOfAlbumResponse
}

type ScrapService interface {
	Create(scrap request.CreateScrapRequest)
	Delete(id int)
	FindOfUserRecipient(userId int) response.FindScrapOfUserRecipientResponse
}

type TestimonialService interface {
	Create(testimonial request.CreateTestimonialRequest)
	Delete(id int)
	FindOfUserRecipient(userId int) response.FindOfUserResponse
	FindOfUserAuthor(userId int) response.FindOfUserResponse
}

type UserService interface {
	Create(user request.CreateUserRequest) response.Response
	Update(user request.UpdateUserRequest) response.Response
	Delete(userId int)
	FindById(userId int) response.FindUserByIdResponse
	FindMany() response.Response
}
