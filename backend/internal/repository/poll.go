package repository

import (
	"time"
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PollRepository struct {
	repository[*model.Poll]
}

func NewPollRepository(Db *gorm.DB) *PollRepository {
	return &PollRepository{
		repository: repository[*model.Poll]{db: Db},
	}
}

// FindOfCommunity returns one page of the community's polls with their
// authors, newest first, plus how many polls it has.
func (p *PollRepository) FindOfCommunity(communityId uint, page int, pageSize int) ([]model.Poll, int64, error) {
	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("community_id = ?", communityId)
	}

	var count int64
	if err := p.db.Model(&model.Poll{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	polls := []model.Poll{}
	result := p.db.
		Scopes(filter, newPaginate(page, pageSize).paginatedResult).
		Preload("Author").
		Order("created_at desc, id desc").
		Find(&polls)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return polls, count, nil
}

// FindByIdWithAuthor returns the poll with its author's profile and its
// community.
func (p *PollRepository) FindByIdWithAuthor(pollId uint) (*model.Poll, error) {
	var poll model.Poll
	result := p.db.
		Preload("Author").
		Preload("Author.Profile").
		Preload("Community").
		First(&poll, pollId)
	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &poll, nil
}

// UpdateClosesAt sets when the poll closes.
func (p *PollRepository) UpdateClosesAt(pollId uint, closesAt time.Time) error {
	return p.db.Model(&model.Poll{}).Where("id = ?", pollId).Update("closes_at", closesAt).Error
}
