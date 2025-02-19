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

func (u *UserService) Delete(userId int) {
	u.UserRepository.Delete(userId)
}

func (u *UserService) FindMany() response.Response {
	users := u.UserRepository.FindMany()
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
	}
}

func (u *UserService) Update(user request.UpdateUserRequest) error {
	err := u.UserRepository.Update(model.NewUser(user.Name, user.Email, user.BirthDate))
	if err != nil {
		panic(err)
	}
	return nil
}

func (u *UserService) Create(user request.CreateUserRequest) error {
	err := u.UserRepository.Create(model.NewUser(user.Name, user.Email, user.BirthDate))
	if err != nil {
		panic(err)
	}
	return nil
}

func (u *UserService) FindById(userId int) response.FindUserByIdResponse {
	userData, err := u.UserRepository.FindById(userId)
	if err != nil {
		panic(err)
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

	return userResponse

}

func NewUserService(userRepository repository.UserRepository, validate *validator.Validate) *UserService {
	return &UserService{
		UserRepository: userRepository,
		Validate:       validate,
	}
}
