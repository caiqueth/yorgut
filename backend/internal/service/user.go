package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

// Delete implements UserService.
func (u *UserServiceImpl) Delete(userId int) {
	panic("unimplemented")
}

// FindMany implements UserService.
func (u *UserServiceImpl) FindMany() response.Response {
	panic("unimplemented")
}

// Update implements UserService.
func (u *UserServiceImpl) Update(user request.UpdateUserRequest) response.Response {
	panic("unimplemented")
}

// CreateUser implements UserService.
func (u *UserServiceImpl) Create(user request.CreateUserRequest) response.Response {
	panic("unimplemented")
}

// FindById implements UserService.
func (u *UserServiceImpl) FindById(userId int) response.FindUserByIdResponse {
	userData, err := u.UserRepository.FindById(userId)
	if err != nil {
		panic(err)
	}

	userResponse := response.FindUserByIdResponse{
		Id:        userData.ID,
		Email:     userData.Email,
		Name:      userData.Name,
		BirthDate: userData.BirthDate,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}

	return userResponse

}

func NewUserServiceImpl(userRepository repository.UserRepository, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}
