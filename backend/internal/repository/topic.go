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

func (t *TopicRepository) FindMostRecentOfCommunity(communityId uint, limit int) ([]model.Topic, error) {
	var topics []model.Topic
	result := t.db.
		Where("community_id = ? and parent_topic_id is null", communityId).
		Preload("Author").
		Preload("Author.Profile").
		Order("created_at desc").
		Limit(limit).
		Find(&topics)

	if result.Error != nil {
		return nil, result.Error
	}

	return topics, nil
}

func (t *TopicRepository) FindRepliesOfTopic(topicId uint) (int64, []model.Topic, error) {
	var replies []model.Topic
	var count int64
	resultCount := t.db.Model(&model.Topic{}).Where("parent_topic_id = ?", topicId).Count(&count)

	if resultCount.Error != nil {
		return 0, nil, resultCount.Error
	}

	resultReplies := t.db.
		Where("parent_topic_id = ?", topicId).
		Preload("Author").
		Preload("Author.Profile").
		Order("created_at asc").
		Find(&replies)

	if resultReplies.Error != nil {
		return 0, nil, resultReplies.Error
	}

	return count, replies, nil
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
