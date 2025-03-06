package service

import (
	"errors"
	"yorgut/internal/dto"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type ScrapService struct {
	ScrapRepository repository.ScrapRepository
	UserRepository  repository.UserRepository
	Validate        *validator.Validate
}

func (s *ScrapService) Create(scrap dto.CreateScrapRequest) error {
	author, err := s.UserRepository.FindById(scrap.UserAuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	recipient, err := s.UserRepository.FindById(scrap.UserRecipientID)
	if err != nil {
		return errors.New("recipient not found")
	}

	scrapModel := model.NewScrap(author, recipient, scrap.Content)

	err = s.ScrapRepository.Create(scrapModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScrapService) Delete(id uint) error {
	err := s.ScrapRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScrapService) FindOfUserRecipient(userId int) (dto.FindScrapOfUserRecipientResponse, error) {
	panic("unimplemented")
}

func NewScrapService(scrapRepository repository.ScrapRepository, validate *validator.Validate) *ScrapService {
	return &ScrapService{
		ScrapRepository: scrapRepository,
		Validate:        validate,
	}
}
