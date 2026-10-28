package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var (
	ErrScrapNotFound  = newError(ErrNotFound, "O recado não existe.")
	ErrScrapForbidden = newError(ErrForbidden, "Só o dono do perfil ou quem escreveu o recado pode apagá-lo.")
)

type ScrapService struct {
	scrapRepository      repository.ScrapRepository
	userRepository       repository.UserRepository
	friendshipRepository repository.FriendshipRepository
}

func NewScrapService(
	scrapRepository repository.ScrapRepository,
	userRepository repository.UserRepository,
	friendshipRepository repository.FriendshipRepository,
) *ScrapService {
	return &ScrapService{
		scrapRepository:      scrapRepository,
		userRepository:       userRepository,
		friendshipRepository: friendshipRepository,
	}
}

// CreateScrap writes a scrap in the recipient's scrapbook. A scrapbook kept
// to friends only takes scraps from friends.
func (s *ScrapService) CreateScrap(userAuthorID uint, userRecipientID uint, content string) error {
	if err := requireAccess(s.userRepository, s.friendshipRepository, userAuthorID, userRecipientID, sectionScraps); err != nil {
		return err
	}

	author, err := s.userRepository.FindById(userAuthorID)
	if err != nil {
		return fromRepository(err, "o autor")
	}

	recipient, err := s.userRepository.FindById(userRecipientID)
	if err != nil {
		return fromRepository(err, "o destinatário")
	}

	scrapModel := model.NewScrap(author, recipient, content)

	if err := s.scrapRepository.Create(scrapModel); err != nil {
		return err
	}

	return nil
}

// DeleteScrap deletes a scrap for good. Only the owner of the scrapbook it
// is in (its recipient) and the person who wrote it may.
func (s *ScrapService) DeleteScrap(viewerId uint, scrapId uint) error {
	scrap, err := s.scrapRepository.FindById(scrapId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrScrapNotFound
	}
	if err != nil {
		return err
	}

	if viewerId != scrap.UserRecipientID && viewerId != scrap.UserAuthorID {
		return ErrScrapForbidden
	}

	return s.scrapRepository.DeleteForever(scrapId)
}

// FindScrapOfUserRecipient returns one page of the user's scrapbook, newest
// first.
func (s *ScrapService) FindScrapOfUserRecipient(viewerId uint, userId uint, page int, pageSize int) (*transport.FindScrapOfUserRecipientResponse, error) {
	if err := requireAccess(s.userRepository, s.friendshipRepository, viewerId, userId, sectionScraps); err != nil {
		return nil, err
	}

	scraps, count, err := s.scrapRepository.FindOfUserRecipient(userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := transport.FindScrapOfUserRecipientResponse{}
	response.FromScraps(scraps, page, pageSize, count)

	return &response, nil
}

func (s *ScrapService) CountScrapsOfUser(userId uint) (uint, error) {
	count, err := s.scrapRepository.CountOfUserRecipient(userId)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}
