package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var (
	ErrTopicNotFound      = newError(ErrNotFound, "topic not found")
	ErrNotCommunityMember = newError(ErrForbidden, "only members of this community can post in its forum")
)

type TopicService struct {
	topicRepository               repository.TopicRepository
	communityMembershipRepository repository.CommunityMembershipRepository
	Validate                      *validator.Validate
}

func NewTopicService(
	topicRepository repository.TopicRepository,
	communityMembershipRepository repository.CommunityMembershipRepository,
	validate *validator.Validate,
) *TopicService {
	return &TopicService{
		topicRepository:               topicRepository,
		communityMembershipRepository: communityMembershipRepository,
		Validate:                      validate,
	}
}

func (ts *TopicService) isMember(userId uint, communityId uint) bool {
	membership, err := ts.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	return err == nil && membership != nil
}

func (ts *TopicService) FindTopics(communityId uint, limit int) (*transport.TopicsResponse, error) {
	response := transport.TopicsResponse{}
	response.Topics = []transport.TopicSummary{}

	topics, err := ts.topicRepository.FindMostRecentOfCommunity(communityId, limit)
	if err != nil {
		return nil, err
	}

	for _, t := range topics {
		count, replies, err := ts.topicRepository.FindRepliesOfTopic(t.ID)
		if err != nil {
			return nil, err
		}

		topicSummary := transport.TopicSummary{}
		topicSummary.FromTopic(&t, count, replies)

		response.Topics = append(response.Topics, topicSummary)
	}

	return &response, nil
}

func (ts *TopicService) FindTopicReplies(topicId uint) (*transport.TopicRepliesResponse, error) {
	response := transport.TopicRepliesResponse{}
	response.Replies = []transport.TopicReply{}

	topic, err := ts.topicRepository.FindByIdWithAuthorProfile(topicId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTopicNotFound
	}
	if err != nil {
		return nil, err
	}

	count, replies, err := ts.topicRepository.FindRepliesOfTopic(topicId)
	if err != nil {
		return nil, err
	}

	response.Topic = &transport.Topic{}
	response.Topic.FromTopic(topic, count)

	for _, r := range replies {
		reply := transport.TopicReply{}
		reply.FromTopic(&r)
		response.Replies = append(response.Replies, reply)
	}

	return &response, nil
}

func (ts *TopicService) CreateTopic(authorId uint, communityId uint, title string, text string) error {
	if !ts.isMember(authorId, communityId) {
		return ErrNotCommunityMember
	}

	if err := ts.topicRepository.Create(model.NewTopic(authorId, communityId, nil, title, text)); err != nil {
		return err
	}
	return nil
}

// ReplyTopic posts in the parent topic's community, never in one chosen by the
// client, so membership can't be bypassed by sending another community's id.
func (ts *TopicService) ReplyTopic(authorId uint, parentTopicId uint, title string, text string) error {
	parentTopic, err := ts.topicRepository.FindById(parentTopicId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrTopicNotFound
	}
	if err != nil {
		return err
	}

	if !ts.isMember(authorId, parentTopic.CommunityID) {
		return ErrNotCommunityMember
	}

	if err := ts.topicRepository.Create(model.NewTopic(authorId, parentTopic.CommunityID, &parentTopic.ID, title, text)); err != nil {
		return err
	}

	return nil
}
