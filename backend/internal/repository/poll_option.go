package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PollOptionRepository struct {
	repository[*model.PollOption]
}

func NewPollOptionRepository(Db *gorm.DB) *PollOptionRepository {
	return &PollOptionRepository{
		repository: repository[*model.PollOption]{db: Db},
	}
}

// FindOfPoll returns the poll's options in the order they were given.
func (p *PollOptionRepository) FindOfPoll(pollId uint) ([]model.PollOption, error) {
	options := []model.PollOption{}
	result := p.db.Where("poll_id = ?", pollId).Order("id").Find(&options)
	if result.Error != nil {
		return nil, result.Error
	}
	return options, nil
}
