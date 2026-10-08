package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var ErrFirstNameRequired = errors.New("first name is required")

type UserService struct {
	userRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserService(
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		Validate:       validate,
	}
}

func (u *UserService) DeleteUser(id uint) error {
	err := u.userRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) FindManyUsers() (*transport.FindManyUsersResponse, error) {
	users, err := u.userRepository.FindMany()
	if err != nil {
		return nil, err
	}

	resp := &transport.FindManyUsersResponse{}
	resp.FromUsers(users)

	return resp, nil
}

func (u *UserService) CreateUser(firstName string, lastName string, email string, birthDate time.Time, profilePicture *string) error {
	err := u.userRepository.Create(model.NewUser(firstName, lastName, email, birthDate, profilePicture))
	if err != nil {
		return err
	}
	return nil
}

type UpdateUserInput struct {
	FirstName string
	LastName  string
	// BirthDate nil keeps the current one.
	BirthDate *time.Time
}

func (u *UserService) UpdateUser(userId uint, input UpdateUserInput) error {
	firstName := strings.TrimSpace(input.FirstName)
	if firstName == "" {
		return ErrFirstNameRequired
	}

	return u.userRepository.UpdateBasicData(userId, firstName, strings.TrimSpace(input.LastName), input.BirthDate)
}

func (u *UserService) FindUserById(id uint) (*transport.FindUserResponse, error) {
	userData, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	resp := transport.FindUserResponse{}
	resp.FromUser(userData)

	return &resp, nil

}

func (u *UserService) SearchUsers(name string, page int, pageSize int) *transport.SearchUsersResponse {
	users, count := u.userRepository.SearchUsers(name, page, pageSize)

	resp := &transport.SearchUsersResponse{}
	resp.FromUsers(users, page, pageSize, count)

	return resp
}

// GetOrCreateUser returns a model on purpose: the auth middleware stores it in
// the request context as the logged user.
func (u *UserService) GetOrCreateUser(email string, firstName string, lastName string, picture *string) (*model.User, error) {
	user, err := u.userRepository.FindByEmail(email)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		createErr := u.CreateUser(firstName, lastName, email, time.Time{}, picture)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create user: %w", createErr)
		}

		user, err = u.userRepository.FindByEmail(email)
		if err != nil {
			return nil, fmt.Errorf("user created but not found: %w", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("could not find user: %w", err)
}

func (u *UserService) FindUserWithProfileById(id uint) (*transport.FindUserWithProfileByIdResponse, error) {
	userData, err := u.userRepository.FindUserWithProfile(id)
	if err != nil {
		return nil, err
	}

	resp := transport.FindUserWithProfileByIdResponse{}
	resp.FromUser(userData)

	return &resp, nil

}
