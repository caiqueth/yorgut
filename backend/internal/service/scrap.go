package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type ScrapService struct {
	ScrapRepository repository.ScrapRepository
	Validate        *validator.Validate
}

func (s *ScrapService) Create(scrap request.CreateScrapRequest) {
	panic("unimplemented")
}

func (s *ScrapService) Delete(id int) {
	panic("unimplemented")
}

func (s *ScrapService) FindOfUserRecipient(userId int) response.FindScrapOfUserRecipientResponse {
	panic("unimplemented")
}

func NewScrapService(scrapRepository repository.ScrapRepository, validate *validator.Validate) *ScrapService {
	return &ScrapService{
		ScrapRepository: scrapRepository,
		Validate:        validate,
	}
}
