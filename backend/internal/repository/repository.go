package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository[T any] struct {
	db *gorm.DB
}

type paginate struct {
	page     int
	pageSize int
}

func newPaginate(page int, pageSize int) *paginate {
	return &paginate{page: page, pageSize: pageSize}
}

func (p *paginate) paginatedResult(db *gorm.DB) *gorm.DB {
	offset := (p.page - 1) * p.pageSize

	return db.Offset(offset).Limit(p.pageSize)
}

func (r *repository[T]) Create(entity T) error {
	result := r.db.Create(&entity)
	if result.Error != nil {
		return fmt.Errorf("error creating entity: %v", result.Error)
	}
	return nil
}

func (r *repository[T]) CreateMany(entities []T) error {
	result := r.db.CreateInBatches(&entities, 100)
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
	result := r.db.First(&ent, id)
	if result.Error != nil {
		return ent, result.Error
	}
	return ent, nil
}

func (r *repository[T]) FindByIdEager(id uint) (T, error) {
	var ent T
	result := r.db.Preload(clause.Associations).First(&ent, id)
	if result.Error != nil {
		return ent, result.Error
	}
	return ent, nil
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

// Repositories holds one instance of every repository, so the API server and
// the seed command are wired the same way.
type Repositories struct {
	User                *UserRepository
	UserProfile         *UserProfileRepository
	UserReview          *UserReviewRepository
	UserFan             *UserFanRepository
	Community           *CommunityRepository
	CommunityProfile    *CommunityProfileRepository
	CommunityMembership *CommunityMembershipRepository
	Friendship          *FriendshipRepository
	PhotoAlbum          *PhotoAlbumRepository
	Photo               *PhotoRepository
	Scrap               *ScrapRepository
	Testimonial         *TestimonialRepository
	Topic               *TopicRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:                NewUserRepository(db),
		UserProfile:         NewUserProfileRepository(db),
		UserReview:          NewUserReviewRepository(db),
		UserFan:             NewUserFanRepository(db),
		Community:           NewCommunityRepository(db),
		CommunityProfile:    NewCommunityProfileRepository(db),
		CommunityMembership: NewCommunityMembershipRepository(db),
		Friendship:          NewFriendshipRepository(db),
		PhotoAlbum:          NewPhotoAlbumRepository(db),
		Photo:               NewPhotoRepository(db),
		Scrap:               NewScrapRepository(db),
		Testimonial:         NewTestimonialRepository(db),
		Topic:               NewTopicRepository(db),
	}
}
