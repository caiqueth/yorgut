package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

var (
	ErrInvalidGoogleCredential = errors.New("invalid google credential")
	ErrInvalidSession          = errors.New("invalid or expired session")
)

// AuthService handles login. Google is the only identity provider: it proves
// who the user is once, at login, and from then on the user is identified by
// a Yorgut session token (see helpers.NewSessionToken).
type AuthService struct {
	userRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewAuthService(
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		Validate:       validate,
	}
}

// LoginWithGoogle validates a Google ID token, finds or creates the user it
// belongs to, and returns a Yorgut session token for them.
func (a *AuthService) LoginWithGoogle(ctx context.Context, credential string) (*transport.LoginResponse, error) {
	payload, err := idtoken.Validate(ctx, credential, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidGoogleCredential, err)
	}

	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	if payload.Subject == "" || email == "" || !emailVerified {
		return nil, ErrInvalidGoogleCredential
	}

	firstName, _ := payload.Claims["given_name"].(string)
	lastName, _ := payload.Claims["family_name"].(string)
	var picture *string
	if p, ok := payload.Claims["picture"].(string); ok && p != "" {
		picture = &p
	}

	user, err := a.getOrCreateUser(payload.Subject, email, firstName, lastName, picture)
	if err != nil {
		return nil, err
	}

	token, err := helpers.NewSessionToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not issue session token: %w", err)
	}

	return &transport.LoginResponse{Token: token}, nil
}

// UserFromSession returns the user a session token was issued for. It
// returns a model on purpose: the auth middleware stores it in the request
// context as the logged user.
func (a *AuthService) UserFromSession(token string) (*model.User, error) {
	userId, err := helpers.ParseSessionToken(token)
	if err != nil {
		return nil, ErrInvalidSession
	}

	user, err := a.userRepository.FindById(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// The user was deleted after the token was issued.
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, fmt.Errorf("could not load session user: %w", err)
	}

	return user, nil
}

func (a *AuthService) getOrCreateUser(googleSub string, email string, firstName string, lastName string, picture *string) (*model.User, error) {
	user, err := a.userRepository.FindByGoogleSub(googleSub)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	// Users created before GoogleSub existed (and seeded users) are matched
	// by email once, and linked to their Google account from then on.
	user, err = a.userRepository.FindByEmail(email)
	if err == nil {
		if user.GoogleSub != nil && *user.GoogleSub != googleSub {
			// The email is already linked to a different Google account.
			return nil, ErrInvalidGoogleCredential
		}
		if err := a.userRepository.SetGoogleSub(user.ID, googleSub); err != nil {
			return nil, fmt.Errorf("could not link google account: %w", err)
		}
		user.GoogleSub = &googleSub
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	newUser := model.NewUser(firstName, lastName, email, time.Time{}, picture)
	newUser.GoogleSub = &googleSub
	if err := a.userRepository.Create(newUser); err != nil {
		// Two first logins at the same time: the unique index rejected the
		// second insert, so the user exists now.
		if existing, findErr := a.userRepository.FindByGoogleSub(googleSub); findErr == nil {
			return existing, nil
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return a.userRepository.FindByGoogleSub(googleSub)
}
