package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type PollVoteRepository struct {
	repository[*model.PollVote]
}

func NewPollVoteRepository(Db *gorm.DB) *PollVoteRepository {
	return &PollVoteRepository{
		repository: repository[*model.PollVote]{db: Db},
	}
}

// PollCount is a count of votes, keyed by poll or by option.
type PollCount struct {
	ID    uint
	Count int64
}

// CountByOption counts the poll's votes per option. Options without votes
// are missing.
func (p *PollVoteRepository) CountByOption(pollId uint) ([]PollCount, error) {
	counts := []PollCount{}
	result := p.db.
		Model(&model.PollVote{}).
		Select("poll_option_id as id, count(*) as count").
		Where("poll_id = ?", pollId).
		Group("poll_option_id").
		Scan(&counts)
	if result.Error != nil {
		return nil, result.Error
	}
	return counts, nil
}

// CountByPoll counts the votes of each of pollIds. Polls without votes are
// missing.
func (p *PollVoteRepository) CountByPoll(pollIds []uint) ([]PollCount, error) {
	counts := []PollCount{}
	if len(pollIds) == 0 {
		return counts, nil
	}
	result := p.db.
		Model(&model.PollVote{}).
		Select("poll_id as id, count(*) as count").
		Where("poll_id in ?", pollIds).
		Group("poll_id").
		Scan(&counts)
	if result.Error != nil {
		return nil, result.Error
	}
	return counts, nil
}

// FindOfUserAndPoll returns the user's vote in the poll.
func (p *PollVoteRepository) FindOfUserAndPoll(pollId uint, userId uint) (*model.PollVote, error) {
	var vote model.PollVote
	result := p.db.Where("poll_id = ? and user_id = ?", pollId, userId).First(&vote)
	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &vote, nil
}

// FindPollIdsVotedBy returns which of pollIds the user voted in.
func (p *PollVoteRepository) FindPollIdsVotedBy(userId uint, pollIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(pollIds) == 0 {
		return ids, nil
	}
	result := p.db.
		Model(&model.PollVote{}).
		Where("user_id = ? and poll_id in ?", userId, pollIds).
		Pluck("poll_id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}
