package service

import (
	"context"
	"errors"
	"testing"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
)

func newTestAuthService(t *testing.T) (*AuthService, *repository.Repositories) {
	t.Helper()
	repos := newTestRepositories(t)
	auth := NewAuthService(*repos.User, *repos.UserProfile, storage.NewMemory())
	auth.download = func(context.Context, string, int64) ([]byte, error) {
		t.Fatal("unexpected download")
		return nil, nil
	}
	return auth, repos
}

func TestGetOrCreateUser(t *testing.T) {
	t.Run("creates the user on first login", func(t *testing.T) {
		auth, _ := newTestAuthService(t)

		user, err := auth.getOrCreateUser("sub-1", "ana@test.com", "Ana", "Silva")
		if err != nil {
			t.Fatalf("getOrCreateUser: %v", err)
		}
		if user.ID == 0 || user.GoogleSub == nil || *user.GoogleSub != "sub-1" {
			t.Fatalf("unexpected user %+v", user)
		}

		again, err := auth.getOrCreateUser("sub-1", "ana@test.com", "Ana", "Silva")
		if err != nil || again.ID != user.ID {
			t.Fatalf("second login: got user %v, err %v; want id %d", again, err, user.ID)
		}
	})

	t.Run("links an existing user by email, then finds them by sub", func(t *testing.T) {
		auth, repos := newTestAuthService(t)
		seeded := model.NewUser("Bia", "Souza", "bia@test.com", nil, nil)
		if err := repos.User.Create(seeded); err != nil {
			t.Fatalf("seeding: %v", err)
		}
		existing, _ := repos.User.FindByEmail("bia@test.com")

		user, err := auth.getOrCreateUser("sub-2", "bia@test.com", "Bia", "Souza")
		if err != nil || user.ID != existing.ID {
			t.Fatalf("got user %v, err %v; want id %d", user, err, existing.ID)
		}

		// The Google account changed its email: the user is still found by sub.
		user, err = auth.getOrCreateUser("sub-2", "bia.new@test.com", "Bia", "Souza")
		if err != nil || user.ID != existing.ID {
			t.Fatalf("after email change: got user %v, err %v; want id %d", user, err, existing.ID)
		}
	})

	t.Run("rejects an email linked to another google account", func(t *testing.T) {
		auth, _ := newTestAuthService(t)
		if _, err := auth.getOrCreateUser("sub-3", "caio@test.com", "Caio", ""); err != nil {
			t.Fatalf("first login: %v", err)
		}

		_, err := auth.getOrCreateUser("sub-other", "caio@test.com", "Caio", "")
		if !errors.Is(err, ErrInvalidGoogleCredential) {
			t.Fatalf("got err %v, want ErrInvalidGoogleCredential", err)
		}
	})
}

func TestUserFromSession(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")
	auth, repos := newTestAuthService(t)
	user, err := auth.getOrCreateUser("sub-4", "duda@test.com", "Duda", "")
	if err != nil {
		t.Fatalf("getOrCreateUser: %v", err)
	}

	token, _ := helpers.NewSessionToken(user.ID)
	got, err := auth.UserFromSession(token)
	if err != nil || got.ID != user.ID {
		t.Fatalf("got user %v, err %v; want id %d", got, err, user.ID)
	}
	if got.LastSeenAt == nil {
		t.Fatal("LastSeenAt not set by the request")
	}
	stored, _ := repos.User.FindById(user.ID)
	if stored.LastSeenAt == nil {
		t.Fatal("LastSeenAt not saved")
	}

	deletedUserToken, _ := helpers.NewSessionToken(9999)
	if _, err := auth.UserFromSession(deletedUserToken); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("unknown user: got err %v, want ErrInvalidSession", err)
	}

	if _, err := auth.UserFromSession("garbage"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("garbage token: got err %v, want ErrInvalidSession", err)
	}
}

func TestCopyGooglePicture(t *testing.T) {
	ctx := context.Background()
	const googlePicture = "https://lh3.googleusercontent.com/a/photo=s96-c"

	// setup returns the service, its store and a user whose picture is
	// current; download serves a valid PNG and records the URL it was asked for.
	setup := func(t *testing.T, current *string) (*AuthService, *storage.Memory, *repository.Repositories, uint, *string) {
		t.Helper()
		repos := newTestRepositories(t)
		store := storage.NewMemory()
		auth := NewAuthService(*repos.User, *repos.UserProfile, store)
		requested := new(string)
		auth.download = func(_ context.Context, url string, _ int64) ([]byte, error) {
			*requested = url
			return testPNG(t), nil
		}
		user := model.NewUser("Ana", "", "ana@test.com", nil, current)
		if err := repos.User.Create(user); err != nil {
			t.Fatal(err)
		}
		return auth, store, repos, user.ID, requested
	}
	pictureOf := func(t *testing.T, repos *repository.Repositories, userId uint) *string {
		t.Helper()
		user, err := repos.User.FindUserWithProfile(userId)
		if err != nil {
			t.Fatal(err)
		}
		return user.Profile.ProfilePicture
	}

	t.Run("a user without a picture gets a stored copy", func(t *testing.T) {
		auth, store, repos, userId, requested := setup(t, nil)
		auth.copyGooglePicture(ctx, userId, googlePicture)

		got := pictureOf(t, repos, userId)
		if got == nil || !store.Has(*got) {
			t.Fatalf("picture = %v, want a stored key", got)
		}
		if *requested != "https://lh3.googleusercontent.com/a/photo=s400-c" {
			t.Errorf("downloaded %q, want the 400px version", *requested)
		}
	})

	t.Run("a picture still pointing at Google is replaced", func(t *testing.T) {
		old := googlePicture
		auth, store, repos, userId, _ := setup(t, &old)
		auth.copyGooglePicture(ctx, userId, googlePicture)

		got := pictureOf(t, repos, userId)
		if got == nil || isExternalPicture(*got) || !store.Has(*got) {
			t.Fatalf("picture = %v, want a stored key", got)
		}
	})

	t.Run("an uploaded picture is kept", func(t *testing.T) {
		uploaded := "uploaded.png"
		auth, _, repos, userId, requested := setup(t, &uploaded)
		auth.copyGooglePicture(ctx, userId, googlePicture)

		if got := pictureOf(t, repos, userId); got == nil || *got != uploaded {
			t.Fatalf("picture = %v, want %q kept", got, uploaded)
		}
		if *requested != "" {
			t.Error("downloaded the Google picture although the user has their own")
		}
	})

	t.Run("a failed download leaves the picture alone", func(t *testing.T) {
		auth, store, repos, userId, _ := setup(t, nil)
		auth.download = func(context.Context, string, int64) ([]byte, error) {
			return nil, errors.New("google is down")
		}
		auth.copyGooglePicture(ctx, userId, googlePicture)

		if got := pictureOf(t, repos, userId); got != nil {
			t.Fatalf("picture = %v, want none (retried at the next login)", *got)
		}
		if store.Len() != 0 {
			t.Errorf("store has %d objects, want 0", store.Len())
		}
	})

	t.Run("something that isn't a picture is refused", func(t *testing.T) {
		auth, store, repos, userId, _ := setup(t, nil)
		auth.download = func(context.Context, string, int64) ([]byte, error) {
			return []byte("<html>not a picture</html>"), nil
		}
		auth.copyGooglePicture(ctx, userId, googlePicture)

		if got := pictureOf(t, repos, userId); got != nil || store.Len() != 0 {
			t.Fatalf("picture = %v, store %d objects; want nothing stored", got, store.Len())
		}
	})
}
