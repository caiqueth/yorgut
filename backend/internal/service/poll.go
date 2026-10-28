package service

import (
	"errors"
	"strings"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var (
	ErrPollNotFound       = newError(ErrNotFound, "A enquete não existe.")
	ErrPrivatePolls       = newError(ErrForbidden, "Só membros podem ver as enquetes de uma comunidade privada.")
	ErrPollNotMember      = newError(ErrForbidden, "Só membros da comunidade podem criar enquetes e votar.")
	ErrPollCloseForbidden = newError(ErrForbidden, "Só quem criou a enquete ou o dono da comunidade pode encerrá-la.")
	ErrPollClosed         = newError(ErrConflict, "Esta enquete já foi encerrada.")
	ErrAlreadyVoted       = newError(ErrConflict, "Você já votou nesta enquete.")
	ErrPollOptionNotFound = newError(ErrInvalid, "Essa opção não é desta enquete.")
	ErrPollInvalidOptions = newError(ErrInvalid, "A enquete precisa de 2 a 10 opções diferentes, todas com texto.")
	ErrPollQuestionBlank  = newError(ErrInvalid, "Escreva a pergunta da enquete.")
	ErrPollClosesInPast   = newError(ErrInvalid, "A data de encerramento precisa ser no futuro.")
)

// PollService handles community polls: creating, listing, voting and
// closing. Members create polls and vote, once each; anyone who may read the
// community's forum may see its polls.
type PollService struct {
	transactor                    repository.Transactor
	pollRepository                repository.PollRepository
	pollOptionRepository          repository.PollOptionRepository
	pollVoteRepository            repository.PollVoteRepository
	communityRepository           repository.CommunityRepository
	communityMembershipRepository repository.CommunityMembershipRepository
}

func NewPollService(
	transactor repository.Transactor,
	pollRepository repository.PollRepository,
	pollOptionRepository repository.PollOptionRepository,
	pollVoteRepository repository.PollVoteRepository,
	communityRepository repository.CommunityRepository,
	communityMembershipRepository repository.CommunityMembershipRepository,
) *PollService {
	return &PollService{
		transactor:                    transactor,
		pollRepository:                pollRepository,
		pollOptionRepository:          pollOptionRepository,
		pollVoteRepository:            pollVoteRepository,
		communityRepository:           communityRepository,
		communityMembershipRepository: communityMembershipRepository,
	}
}

// CreatePollInput is a new poll, as the author wrote it.
type CreatePollInput struct {
	Question    string
	Description string
	Options     []string
	// ClosesAt nil leaves the poll open until someone closes it.
	ClosesAt *time.Time
}

// isClosed says whether voting has ended.
func isClosed(poll *model.Poll, now time.Time) bool {
	return poll.ClosesAt != nil && !poll.ClosesAt.After(now)
}

func (p *PollService) isMember(userId uint, communityId uint) (bool, error) {
	_, err := p.communityMembershipRepository.FindOfUserAndCommunity(communityId, userId)
	if errors.Is(err, repository.ErrNotFound) {
		return false, nil
	}
	return err == nil, err
}

// checkCanRead returns ErrPrivatePolls when the community is private and the
// viewer isn't a member, the same rule as the forum.
func (p *PollService) checkCanRead(viewerId uint, community *model.Community) error {
	if !community.Private {
		return nil
	}
	member, err := p.isMember(viewerId, community.ID)
	if err != nil {
		return err
	}
	if !member {
		return ErrPrivatePolls
	}
	return nil
}

// FindPolls returns one page of the community's polls, newest first.
func (p *PollService) FindPolls(viewerId uint, communityId uint, page int, pageSize int) (*transport.PollsResponse, error) {
	community, err := p.communityRepository.FindById(communityId)
	if err != nil {
		return nil, fromRepository(err, "a comunidade")
	}
	if err := p.checkCanRead(viewerId, community); err != nil {
		return nil, err
	}

	polls, count, err := p.pollRepository.FindOfCommunity(communityId, page, pageSize)
	if err != nil {
		return nil, err
	}
	pollIds := make([]uint, len(polls))
	for i, poll := range polls {
		pollIds[i] = poll.ID
	}
	counts, err := p.pollVoteRepository.CountByPoll(pollIds)
	if err != nil {
		return nil, err
	}
	votedIds, err := p.pollVoteRepository.FindPollIdsVotedBy(viewerId, pollIds)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	resp := &transport.PollsResponse{}
	resp.FromPolls(polls, countsById(counts), votedIds, func(poll *model.Poll) bool { return isClosed(poll, now) }, page, pageSize, count)
	return resp, nil
}

// FindPoll returns the poll with its options. Each option's votes are shown
// only once the viewer voted or the poll closed.
func (p *PollService) FindPoll(viewerId uint, pollId uint) (*transport.PollResponse, error) {
	poll, err := p.pollRepository.FindByIdWithAuthor(pollId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrPollNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := p.checkCanRead(viewerId, poll.Community); err != nil {
		return nil, err
	}

	options, err := p.pollOptionRepository.FindOfPoll(pollId)
	if err != nil {
		return nil, err
	}
	counts, err := p.pollVoteRepository.CountByOption(pollId)
	if err != nil {
		return nil, err
	}
	member, err := p.isMember(viewerId, poll.CommunityID)
	if err != nil {
		return nil, err
	}

	var myOptionId *uint
	vote, err := p.pollVoteRepository.FindOfUserAndPoll(pollId, viewerId)
	if err == nil {
		myOptionId = &vote.PollOptionID
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	closed := isClosed(poll, time.Now())
	resp := &transport.PollResponse{}
	resp.FromPoll(poll, options, countsById(counts), myOptionId != nil || closed)
	resp.Closed = closed
	resp.MyOptionID = myOptionId
	resp.CanVote = member && !closed && myOptionId == nil
	resp.CanClose = !closed && (viewerId == poll.AuthorID || viewerId == poll.Community.OwnerID)
	return resp, nil
}

// CreatePoll creates a poll in the community and returns its id. Only
// members may.
func (p *PollService) CreatePoll(authorId uint, communityId uint, input CreatePollInput) (uint, error) {
	if _, err := p.communityRepository.FindById(communityId); err != nil {
		return 0, fromRepository(err, "a comunidade")
	}
	member, err := p.isMember(authorId, communityId)
	if err != nil {
		return 0, err
	}
	if !member {
		return 0, ErrPollNotMember
	}

	question := strings.TrimSpace(input.Question)
	if question == "" {
		return 0, ErrPollQuestionBlank
	}
	options, err := cleanPollOptions(input.Options)
	if err != nil {
		return 0, err
	}
	if input.ClosesAt != nil && !input.ClosesAt.After(time.Now()) {
		return 0, ErrPollClosesInPast
	}

	poll := model.NewPoll(communityId, authorId, question, nilIfBlank(input.Description), input.ClosesAt)
	err = p.transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.Poll.Create(poll); err != nil {
			return err
		}
		for _, text := range options {
			if err := tx.PollOption.Create(model.NewPollOption(poll.ID, text)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return poll.ID, nil
}

// cleanPollOptions trims the options and checks there are 2 to 10, all with
// text and all different (ignoring case).
func cleanPollOptions(options []string) ([]string, error) {
	cleaned := make([]string, 0, len(options))
	seen := map[string]bool{}
	for _, o := range options {
		text := strings.TrimSpace(o)
		key := strings.ToLower(text)
		if text == "" || seen[key] {
			return nil, ErrPollInvalidOptions
		}
		seen[key] = true
		cleaned = append(cleaned, text)
	}
	if len(cleaned) < 2 || len(cleaned) > 10 {
		return nil, ErrPollInvalidOptions
	}
	return cleaned, nil
}

// Vote records the user's vote for one of the poll's options. Only members
// vote, once, while the poll is open.
func (p *PollService) Vote(userId uint, pollId uint, optionId uint) error {
	poll, err := p.pollRepository.FindById(pollId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrPollNotFound
	}
	if err != nil {
		return err
	}
	member, err := p.isMember(userId, poll.CommunityID)
	if err != nil {
		return err
	}
	if !member {
		return ErrPollNotMember
	}
	if isClosed(poll, time.Now()) {
		return ErrPollClosed
	}

	option, err := p.pollOptionRepository.FindById(optionId)
	if errors.Is(err, repository.ErrNotFound) || (err == nil && option.PollID != pollId) {
		return ErrPollOptionNotFound
	}
	if err != nil {
		return err
	}

	err = p.pollVoteRepository.Create(model.NewPollVote(pollId, optionId, userId))
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrAlreadyVoted
	}
	return err
}

// ClosePoll ends voting now. Only the poll's author or the community's owner
// may.
func (p *PollService) ClosePoll(userId uint, pollId uint) error {
	poll, err := p.pollRepository.FindByIdWithAuthor(pollId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrPollNotFound
	}
	if err != nil {
		return err
	}
	if userId != poll.AuthorID && userId != poll.Community.OwnerID {
		return ErrPollCloseForbidden
	}
	now := time.Now()
	if isClosed(poll, now) {
		return ErrPollClosed
	}
	return p.pollRepository.UpdateClosesAt(pollId, now)
}

// countsById turns per-poll or per-option counts into a map.
func countsById(counts []repository.PollCount) map[uint]int64 {
	byId := make(map[uint]int64, len(counts))
	for _, c := range counts {
		byId[c.ID] = c.Count
	}
	return byId
}
