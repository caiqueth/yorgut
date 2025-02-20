package service

import (
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func (u *UserService) Delete(id uint) error {
	err := u.UserRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindMany() (response.Response, error) {
	users, err := u.UserRepository.FindMany()
	if err != nil {
		return response.Response{Data: nil}, err
	}

	userDtos := make([]response.UserDto, len(users))

	for i, user := range users {
		userDtos[i] = response.UserDto{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			BirthDate: user.BirthDate,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return response.Response{
		Data: map[string][]response.UserDto{"Users": userDtos},
	}, nil
}

func (u *UserService) Create(user request.CreateUserRequest) error {
	err := u.UserRepository.Create(model.NewUser(user.Name, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Update(user request.UpdateUserRequest) error {
	err := u.UserRepository.Update(model.NewUser(user.Name, user.Email, user.BirthDate))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindById(id uint) (response.FindUserByIdResponse, error) {
	userData, err := u.UserRepository.FindById(id)
	if err != nil {
		return response.FindUserByIdResponse{}, err
	}

	userResponse := response.FindUserByIdResponse{
		UserDto: response.UserDto{
			Id:        userData.ID,
			Email:     userData.Email,
			Name:      userData.Name,
			BirthDate: userData.BirthDate,
			CreatedAt: userData.CreatedAt,
			UpdatedAt: userData.UpdatedAt,
		}}

	return userResponse, err

}

func NewUserService(userRepository repository.UserRepository, validate *validator.Validate) *UserService {
	return &UserService{
		UserRepository: userRepository,
		Validate:       validate,
	}
}
