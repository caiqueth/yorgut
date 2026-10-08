package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

type ScrapService struct {
	scrapRepository repository.ScrapRepository
	userRepository  repository.UserRepository
	Validate        *validator.Validate
}

func NewScrapService(scrapRepository repository.ScrapRepository, userRepository repository.UserRepository, validate *validator.Validate) *ScrapService {
	return &ScrapService{
		scrapRepository: scrapRepository,
		userRepository:  userRepository,
		Validate:        validate,
	}
}

func (s *ScrapService) CreateScrap(userAuthorID uint, userRecipientID uint, content string) error {
	author, err := s.userRepository.FindById(userAuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	recipient, err := s.userRepository.FindById(userRecipientID)
	if err != nil {
		return errors.New("recipient not found")
	}

	scrapModel := model.NewScrap(author, recipient, content)

	if err := s.scrapRepository.Create(scrapModel); err != nil {
		return err
	}

	return nil
}

func (s *ScrapService) DeleteScrap(id uint) error {
	err := s.scrapRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScrapService) FindScrapOfUserRecipient(userId uint) (*transport.FindScrapOfUserRecipientResponse, error) {
	scraps, err := s.scrapRepository.FindOfUserRecipient(userId)
	if err != nil {
		return nil, err
	}

	response := transport.FindScrapOfUserRecipientResponse{}
	response.FromScraps(scraps)

	return &response, nil
}

func (s *ScrapService) CountScrapsOfUser(userId uint) (uint, error) {
	count, err := s.scrapRepository.CountOfUserRecipient(userId)
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}
