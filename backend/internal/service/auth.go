package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
	"yorgut/internal/transport"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"
)

var (
	ErrInvalidGoogleCredential = errors.New("Não foi possível confirmar seu login com o Google. Tente de novo.")
	ErrInvalidSession          = errors.New("Sua sessão expirou. Entre de novo.")
)

// AuthService handles login. Google is the only identity provider: it proves
// who the user is once, at login, and from then on the user is identified by
// a Yorgut session token (see helpers.NewSessionToken).
type AuthService struct {
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
	store                 storage.ObjectStore
	// download fetches the Google profile picture (helpers.Download; tests
	// replace it).
	download func(ctx context.Context, url string, maxBytes int64) ([]byte, error)
}

func NewAuthService(
	userRepository repository.UserRepository,
	userProfileRepository repository.UserProfileRepository,
	store storage.ObjectStore,
) *AuthService {
	return &AuthService{
		userRepository:        userRepository,
		userProfileRepository: userProfileRepository,
		store:                 store,
		download:              helpers.Download,
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
	googlePicture, _ := payload.Claims["picture"].(string)

	user, err := a.getOrCreateUser(payload.Subject, email, firstName, lastName)
	if err != nil {
		return nil, err
	}
	a.copyGooglePicture(ctx, user.ID, googlePicture)

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
	if errors.Is(err, repository.ErrNotFound) {
		// The user was deleted after the token was issued.
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, fmt.Errorf("could not load session user: %w", err)
	}

	a.touchLastSeen(user)

	return user, nil
}

// lastSeenEvery is how stale User.LastSeenAt may get before a request
// refreshes it, so most requests don't write.
const lastSeenEvery = 5 * time.Minute

// touchLastSeen records that the user is using the app. A failure is only
// logged: it must not log the user out.
func (a *AuthService) touchLastSeen(user *model.User) {
	now := time.Now()
	staleBefore := now.Add(-lastSeenEvery)
	if user.LastSeenAt != nil && user.LastSeenAt.After(staleBefore) {
		return
	}
	if err := a.userRepository.TouchLastSeen(user.ID, now, staleBefore); err != nil {
		log.Warn().Err(err).Uint("userId", user.ID).Msg("could not update last seen")
		return
	}
	user.LastSeenAt = &now
}

func (a *AuthService) getOrCreateUser(googleSub string, email string, firstName string, lastName string) (*model.User, error) {
	user, err := a.userRepository.FindByGoogleSub(googleSub)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
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
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	newUser := model.NewUser(firstName, lastName, email, nil, nil)
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

// googlePictureSize matches the size suffix of a Google profile picture URL
// ("...=s96-c"), which asks for a 96px square.
var googlePictureSize = regexp.MustCompile(`=s\d+-c$`)

// copyGooglePicture stores the user's Google profile picture as their own
// picture, so it is served from our bucket like any upload. It only fills a
// missing picture or replaces one still pointing at Google (users created
// before pictures were copied); a picture the user uploaded is kept.
//
// Login must not fail because of it: errors are logged, and the next login
// tries again.
func (a *AuthService) copyGooglePicture(ctx context.Context, userId uint, googlePicture string) {
	if googlePicture == "" {
		return
	}
	user, err := a.userRepository.FindUserWithProfile(userId)
	if err != nil {
		log.Warn().Err(err).Uint("userId", userId).Msg("could not load the user to copy their Google picture")
		return
	}
	if current := user.Profile; current != nil && current.ProfilePicture != nil && !isExternalPicture(*current.ProfilePicture) {
		return
	}

	// 400px instead of Google's default 96px.
	url := googlePictureSize.ReplaceAllString(googlePicture, "=s400-c")
	picture, err := a.download(ctx, url, MaxUploadBytes)
	if err == nil {
		picture, err = preparePicture(picture)
	}
	if err != nil {
		log.Warn().Err(err).Uint("userId", userId).Msg("could not download the Google picture")
		return
	}

	key, err := storage.PutImage(ctx, a.store, picture)
	if err != nil {
		log.Warn().Err(err).Uint("userId", userId).Msg("could not store the Google picture")
		return
	}
	previous, err := a.userProfileRepository.UpdatePictureOfUser(userId, key)
	if err != nil {
		deletePicture(ctx, a.store, &key)
		log.Warn().Err(err).Uint("userId", userId).Msg("could not save the Google picture")
		return
	}
	deletePicture(ctx, a.store, previous)
}

// isExternalPicture tells a picture stored as a full URL (Google pictures of
// users created before they were copied) from a key in our store.
func isExternalPicture(picture string) bool {
	return strings.HasPrefix(picture, "http://") || strings.HasPrefix(picture, "https://")
}
