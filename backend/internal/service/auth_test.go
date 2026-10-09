package service

import (
	"errors"
	"testing"
	"time"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
)

func newTestAuthService(t *testing.T) (*AuthService, *repository.Repositories) {
	t.Helper()
	repos := newTestRepositories(t)
	return NewAuthService(*repos.User, nil), repos
}

func TestGetOrCreateUser(t *testing.T) {
	t.Run("creates the user on first login", func(t *testing.T) {
		auth, _ := newTestAuthService(t)

		user, err := auth.getOrCreateUser("sub-1", "ana@test.com", "Ana", "Silva", nil)
		if err != nil {
			t.Fatalf("getOrCreateUser: %v", err)
		}
		if user.ID == 0 || user.GoogleSub == nil || *user.GoogleSub != "sub-1" {
			t.Fatalf("unexpected user %+v", user)
		}

		again, err := auth.getOrCreateUser("sub-1", "ana@test.com", "Ana", "Silva", nil)
		if err != nil || again.ID != user.ID {
			t.Fatalf("second login: got user %v, err %v; want id %d", again, err, user.ID)
		}
	})

	t.Run("links an existing user by email, then finds them by sub", func(t *testing.T) {
		auth, repos := newTestAuthService(t)
		seeded := model.NewUser("Bia", "Souza", "bia@test.com", time.Time{}, nil)
		if err := repos.User.Create(seeded); err != nil {
			t.Fatalf("seeding: %v", err)
		}
		existing, _ := repos.User.FindByEmail("bia@test.com")

		user, err := auth.getOrCreateUser("sub-2", "bia@test.com", "Bia", "Souza", nil)
		if err != nil || user.ID != existing.ID {
			t.Fatalf("got user %v, err %v; want id %d", user, err, existing.ID)
		}

		// The Google account changed its email: the user is still found by sub.
		user, err = auth.getOrCreateUser("sub-2", "bia.new@test.com", "Bia", "Souza", nil)
		if err != nil || user.ID != existing.ID {
			t.Fatalf("after email change: got user %v, err %v; want id %d", user, err, existing.ID)
		}
	})

	t.Run("rejects an email linked to another google account", func(t *testing.T) {
		auth, _ := newTestAuthService(t)
		if _, err := auth.getOrCreateUser("sub-3", "caio@test.com", "Caio", "", nil); err != nil {
			t.Fatalf("first login: %v", err)
		}

		_, err := auth.getOrCreateUser("sub-other", "caio@test.com", "Caio", "", nil)
		if !errors.Is(err, ErrInvalidGoogleCredential) {
			t.Fatalf("got err %v, want ErrInvalidGoogleCredential", err)
		}
	})
}

func TestUserFromSession(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")
	auth, _ := newTestAuthService(t)
	user, err := auth.getOrCreateUser("sub-4", "duda@test.com", "Duda", "", nil)
	if err != nil {
		t.Fatalf("getOrCreateUser: %v", err)
	}

	token, _ := helpers.NewSessionToken(user.ID)
	got, err := auth.UserFromSession(token)
	if err != nil || got.ID != user.ID {
		t.Fatalf("got user %v, err %v; want id %d", got, err, user.ID)
	}

	deletedUserToken, _ := helpers.NewSessionToken(9999)
	if _, err := auth.UserFromSession(deletedUserToken); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("unknown user: got err %v, want ErrInvalidSession", err)
	}

	if _, err := auth.UserFromSession("garbage"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("garbage token: got err %v, want ErrInvalidSession", err)
	}
}
