package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var (
	ErrTopicNotFound      = newError(ErrNotFound, "O tópico não existe.")
	ErrNotCommunityMember = newError(ErrForbidden, "Só membros da comunidade podem escrever no fórum.")
	ErrPrivateForum       = newError(ErrForbidden, "Só membros podem ler o fórum de uma comunidade privada.")
)

type TopicService struct {
	topicRepository               repository.TopicRepository
	communityRepository           repository.CommunityRepository
	communityMembershipRepository repository.CommunityMembershipRepository
}

func NewTopicService(
	topicRepository repository.TopicRepository,
	communityRepository repository.CommunityRepository,
	communityMembershipRepository repository.CommunityMembershipRepository,
) *TopicService {
	return &TopicService{
		topicRepository:               topicRepository,
		communityRepository:           communityRepository,
		communityMembershipRepository: communityMembershipRepository,
	}
}

func (ts *TopicService) isMember(userId uint, communityId uint) bool {
	membership, err := ts.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	return err == nil && membership != nil
}

// checkCanRead returns ErrPrivateForum when the community is private and the
// viewer isn't a member. Anyone reads the forum of a public community.
func (ts *TopicService) checkCanRead(viewerId uint, communityId uint) error {
	community, err := ts.communityRepository.FindById(communityId)
	if err != nil {
		return fromRepository(err, "a comunidade")
	}
	if community.Private && !ts.isMember(viewerId, communityId) {
		return ErrPrivateForum
	}
	return nil
}

// FindTopics returns one page of the community's topics, newest first, each
// with its reply count and last reply.
func (ts *TopicService) FindTopics(viewerId uint, communityId uint, page int, pageSize int) (*transport.TopicsResponse, error) {
	if err := ts.checkCanRead(viewerId, communityId); err != nil {
		return nil, err
	}

	topics, topicCount, err := ts.topicRepository.FindOfCommunity(communityId, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := transport.TopicsResponse{}
	response.PagedResponse = transport.PagedResponse{Page: page, PageSize: pageSize, Count: topicCount}
	response.Topics = []transport.TopicSummary{}

	for _, t := range topics {
		count, err := ts.topicRepository.CountRepliesOfTopic(t.ID)
		if err != nil {
			return nil, err
		}
		lastReply, err := ts.topicRepository.FindLastReplyOfTopic(t.ID)
		if err != nil {
			return nil, err
		}

		topicSummary := transport.TopicSummary{}
		topicSummary.FromTopic(&t, count, lastReply)

		response.Topics = append(response.Topics, topicSummary)
	}

	return &response, nil
}

// FindTopicReplies returns the topic and one page of its replies, oldest
// first.
func (ts *TopicService) FindTopicReplies(viewerId uint, topicId uint, page int, pageSize int) (*transport.TopicRepliesResponse, error) {
	response := transport.TopicRepliesResponse{}
	response.Replies = []transport.TopicReply{}

	topic, err := ts.topicRepository.FindByIdWithAuthorProfile(topicId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTopicNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := ts.checkCanRead(viewerId, topic.CommunityID); err != nil {
		return nil, err
	}

	replies, count, err := ts.topicRepository.FindRepliesOfTopic(topicId, page, pageSize)
	if err != nil {
		return nil, err
	}

	response.PagedResponse = transport.PagedResponse{Page: page, PageSize: pageSize, Count: count}
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
