package service

import (
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type ScrapServiceImpl struct {
	ScrapRepository repository.ScrapRepository
	Validate        *validator.Validate
}

// Create implements ScrapService.
func (s *ScrapServiceImpl) Create(scrap request.CreateScrapRequest) {
	panic("unimplemented")
}

// Delete implements ScrapService.
func (s *ScrapServiceImpl) Delete(id int) {
	panic("unimplemented")
}

// FindOfUserRecipient implements ScrapService.
func (s *ScrapServiceImpl) FindOfUserRecipient(userId int) response.FindScrapOfUserRecipientResponse {
	panic("unimplemented")
}

func NewScrapServiceImpl(scrapRepository repository.ScrapRepository, validate *validator.Validate) ScrapService {
	return &ScrapServiceImpl{
		ScrapRepository: scrapRepository,
		Validate:        validate,
	}
}
