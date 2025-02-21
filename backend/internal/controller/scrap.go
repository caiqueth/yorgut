package controller

import "yorgut/internal/service"

type ScrapController struct {
	ScrapService service.ScrapService
}

func NewScrapController(service service.ScrapService) *ScrapController {
	return &ScrapController{
		ScrapService: service,
	}
}

func (s *ScrapController) Create() {
	panic("unimplemented")
}

func (s *ScrapController) Delete() {
	panic("unimplemented")
}

func (s *ScrapController) FindOfUserRecipient() {
	panic("unimplemented")
}
