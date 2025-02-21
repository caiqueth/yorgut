package controller

import "yorgut/internal/service"

type FriendshipController struct {
	FriendshipService service.FriendshipService
}

func NewFriendshipController(service service.FriendshipService) *FriendshipController {
	return &FriendshipController{
		FriendshipService: service,
	}
}

func (f *FriendshipController) Create() {
	panic("unimplemented")
}

func (f *FriendshipController) Delete() {
	panic("unimplemented")
}

func (f *FriendshipController) FindOfUser() {
	panic("unimplemented")
}
