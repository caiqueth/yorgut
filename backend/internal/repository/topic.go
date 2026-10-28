package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

type TopicRepository struct {
	repository[*model.Topic]
}

func NewTopicRepository(Db *gorm.DB) *TopicRepository {
	return &TopicRepository{
		repository: repository[*model.Topic]{db: Db},
	}
}

// FindOfCommunity returns one page of the community's topics (not their
// replies), newest first, plus how many topics it has.
func (t *TopicRepository) FindOfCommunity(communityId uint, page int, pageSize int) ([]model.Topic, int64, error) {
	filter := func(db *gorm.DB) *gorm.DB {
		return db.Where("community_id = ? and parent_topic_id is null", communityId)
	}

	var count int64
	if err := t.db.Model(&model.Topic{}).Scopes(filter).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	topics := []model.Topic{}
	result := t.db.
		Scopes(filter, newPaginate(page, pageSize).paginatedResult).
		Preload("Author").
		Preload("Author.Profile").
		Order("created_at desc, id desc").
		Find(&topics)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return topics, count, nil
}

// FindRepliesOfTopic returns one page of the topic's replies, oldest first,
// plus how many replies it has.
func (t *TopicRepository) FindRepliesOfTopic(topicId uint, page int, pageSize int) ([]model.Topic, int64, error) {
	count, err := t.CountRepliesOfTopic(topicId)
	if err != nil {
		return nil, 0, err
	}

	replies := []model.Topic{}
	result := t.db.
		Scopes(newPaginate(page, pageSize).paginatedResult).
		Where("parent_topic_id = ?", topicId).
		Preload("Author").
		Preload("Author.Profile").
		Order("created_at asc, id asc").
		Find(&replies)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return replies, count, nil
}

func (t *TopicRepository) CountRepliesOfTopic(topicId uint) (int64, error) {
	var count int64
	if err := t.db.Model(&model.Topic{}).Where("parent_topic_id = ?", topicId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindLastReplyOfTopic returns the topic's newest reply, or nil when it has
// none.
func (t *TopicRepository) FindLastReplyOfTopic(topicId uint) (*model.Topic, error) {
	replies := []model.Topic{}
	result := t.db.
		Where("parent_topic_id = ?", topicId).
		Order("created_at desc, id desc").
		Limit(1).
		Find(&replies)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(replies) == 0 {
		return nil, nil
	}
	return &replies[0], nil
}

func (t *TopicRepository) FindByIdWithAuthorProfile(topicId uint) (*model.Topic, error) {
	var topic model.Topic
	result := t.db.
		Preload("Author").
		Preload("Author.Profile").
		First(&topic, topicId)

	if result.Error != nil {
		return nil, translate(result.Error)
	}
	return &topic, nil
}
