package repository

import (
	"yorgut/internal/model"
)

type CommunityRepository interface {
	Save(c model.Community)
	Update(c model.Community)
	Delete(id int)
	FindById(id int) (c model.Community, err error)
	FindMany() []model.Community
}

type CommunityMembershipRepository interface {
	Save(cm model.CommunityMembership)
	Delete(id int)
	FindOfUser(id int)
}

type FriendshipRepository interface {
	Save(f model.Friendship)
	Delete(id int)
	FindOfUser(userId int)
}

type PhotoAlbumRepository interface {
	Save(pa model.PhotoAlbum)
	Delete(id int)
	FindOfUser(userId int)
}

type PhotoRepository interface {
	Save(p model.Photo)
	Delete(id int)
	FindOfUser(userId int)
	FindOfPhotoAlbum(photoAlbumId int)
}

type ScrapRepository interface {
	Save(s model.Scrap)
	Delete(id int)
	FindOfUserRecipient(userId int)
}

type TestimonialRepository interface {
	Save(t model.Testimonial)
	Delete(id int)
	FindOfUserRecipient(userId int)
	FindOfUserAuthor(userId int)
}

type UserRepository interface {
	Save(u model.User)
	Update(u model.User)
	Delete(id int)
	FindById(id int) (u model.User, err error)
	FindMany() []model.User
}
