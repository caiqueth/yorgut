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

// recentlySeenFirst orders users by when they last used the app, most recent
// first and never-seen last, then by name. The "is null" term works the same
// in MySQL and SQLite, unlike NULLS LAST (which MySQL lacks).
func recentlySeenFirst(db *gorm.DB) *gorm.DB {
	return db.Order("users.last_seen_at is null, users.last_seen_at desc, users.first_name, users.id")
}

// shuffledBy orders rows of table in a pseudo-random order fixed by seed: the
// same seed gives the same order, so reloading doesn't reshuffle. A
// multiplicative hash of the id, which works the same in MySQL and SQLite
// (they have no common seeded random function). Keep seed small (it is added
// to the id before multiplying, and the product must fit in 64 bits).
func shuffledBy(table string, seed int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderBy{Expression: clause.Expr{
			SQL:                "((" + table + ".id + ?) * 2654435761) % 4294967291, " + table + ".id",
			Vars:               []any{seed},
			WithoutParentheses: true,
		}})
	}
}

func (r *repository[T]) Create(entity T) error {
	result := r.db.Create(&entity)
	if result.Error != nil {
		return fmt.Errorf("error creating entity: %w", translate(result.Error))
	}
	return nil
}

func (r *repository[T]) CreateMany(entities []T) error {
	result := r.db.CreateInBatches(&entities, 100)
	if result.Error != nil {
		return fmt.Errorf("error creating entity: %w", translate(result.Error))
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
		return ent, translate(result.Error)
	}
	return ent, nil
}

func (r *repository[T]) FindByIdEager(id uint) (T, error) {
	var ent T
	result := r.db.Preload(clause.Associations).First(&ent, id)
	if result.Error != nil {
		return ent, translate(result.Error)
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

// Transactor runs work that writes more than one row in a single database
// transaction.
type Transactor struct {
	db *gorm.DB
}

// Transaction runs fn with repositories bound to a new transaction. If fn
// returns an error (or panics), everything it wrote is rolled back.
// Repositories outside tx don't see its writes until it commits.
func (t Transactor) Transaction(fn func(tx *Repositories) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewRepositories(tx))
	})
}

// Repositories holds one instance of every repository, so the API server and
// the seed command are wired the same way.
type Repositories struct {
	Transactor           Transactor
	User                 *UserRepository
	UserProfile          *UserProfileRepository
	UserReview           *UserReviewRepository
	UserFan              *UserFanRepository
	Community            *CommunityRepository
	CommunityProfile     *CommunityProfileRepository
	CommunityMembership  *CommunityMembershipRepository
	CommunityJoinRequest *CommunityJoinRequestRepository
	CommunityInvitation  *CommunityInvitationRepository
	Friendship           *FriendshipRepository
	PhotoAlbum           *PhotoAlbumRepository
	Photo                *PhotoRepository
	Scrap                *ScrapRepository
	Testimonial          *TestimonialRepository
	Topic                *TopicRepository
	ProfileVisit         *ProfileVisitRepository
	Poll                 *PollRepository
	PollOption           *PollOptionRepository
	PollVote             *PollVoteRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Transactor:           Transactor{db: db},
		User:                 NewUserRepository(db),
		UserProfile:          NewUserProfileRepository(db),
		UserReview:           NewUserReviewRepository(db),
		UserFan:              NewUserFanRepository(db),
		Community:            NewCommunityRepository(db),
		CommunityProfile:     NewCommunityProfileRepository(db),
		CommunityMembership:  NewCommunityMembershipRepository(db),
		CommunityJoinRequest: NewCommunityJoinRequestRepository(db),
		CommunityInvitation:  NewCommunityInvitationRepository(db),
		Friendship:           NewFriendshipRepository(db),
		PhotoAlbum:           NewPhotoAlbumRepository(db),
		Photo:                NewPhotoRepository(db),
		Scrap:                NewScrapRepository(db),
		Testimonial:          NewTestimonialRepository(db),
		Topic:                NewTopicRepository(db),
		ProfileVisit:         NewProfileVisitRepository(db),
		Poll:                 NewPollRepository(db),
		PollOption:           NewPollOptionRepository(db),
		PollVote:             NewPollVoteRepository(db),
	}
}
