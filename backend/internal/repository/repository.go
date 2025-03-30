package repository

import (
	"errors"
	"fmt"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type repository[T any] struct {
	db *gorm.DB
}

func (r *repository[T]) Create(entity T) error {
	result := r.db.Create(&entity)
	if result.Error != nil {
		return fmt.Errorf("error creating entity: %v", result.Error)
	}
	return nil
}

func (r *repository[T]) Delete(id uint) error {
	var ent T
	result := r.db.Where("id = ?", id).Delete(&ent)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository[T]) FindById(id uint) (T, error) {
	var ent T
	result := r.db.Find(&ent, id)
	if result != nil {
		return ent, nil
	} else {
		return ent, errors.New("entity not found")
	}
}

func (r *repository[T]) FindMany() ([]T, error) {
	var ents []T
	results := r.db.Find(&ents)
	if results.Error != nil {
		return nil, results.Error
	}
	return ents, nil
}

func (r *repository[T]) Update(ent T) error {
	updateEntity := ent
	result := r.db.Model(&ent).Updates(updateEntity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type UserProfileRepository struct {
	repository[*model.UserProfile]
}

type UserRepository struct {
	repository[*model.User]
}

type CommunityMembershipRepository struct {
	repository[*model.CommunityMembership]
}

type FriendshipRepository struct {
	repository[*model.Friendship]
}

type PhotoRepository struct {
	repository[*model.Photo]
}

type CommunityRepository struct {
	repository[*model.Community]
}

type PhotoAlbumRepository struct {
	repository[*model.PhotoAlbum]
}

type ScrapRepository struct {
	repository[*model.Scrap]
}

type TestimonialRepository struct {
	repository[*model.Testimonial]
}

func NewTestimonialRepository(Db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{
		repository: repository[*model.Testimonial]{db: Db},
	}
}
func NewScrapRepository(Db *gorm.DB) *ScrapRepository {
	return &ScrapRepository{
		repository: repository[*model.Scrap]{db: Db},
	}
}

func NewPhotoAlbumRepository(Db *gorm.DB) *PhotoAlbumRepository {
	return &PhotoAlbumRepository{
		repository: repository[*model.PhotoAlbum]{db: Db},
	}
}

func NewCommunityRepository(Db *gorm.DB) *CommunityRepository {
	return &CommunityRepository{
		repository: repository[*model.Community]{db: Db},
	}
}

func NewPhotoRepository(Db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{
		repository: repository[*model.Photo]{db: Db},
	}
}

func NewFriendshipRepository(Db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		repository: repository[*model.Friendship]{db: Db},
	}
}

func NewUserProfileRepository(Db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{
		repository: repository[*model.UserProfile]{db: Db},
	}
}

func NewUserRepository(Db *gorm.DB) *UserRepository {
	return &UserRepository{
		repository: repository[*model.User]{db: Db},
	}
}

func NewCommunityMembershipRepository(Db *gorm.DB) *CommunityMembershipRepository {
	return &CommunityMembershipRepository{
		repository: repository[*model.CommunityMembership]{db: Db},
	}
}
