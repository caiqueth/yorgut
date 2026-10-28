package service

import (
	"errors"
	"testing"
)

func TestUpdateTheme(t *testing.T) {
	s, repos, ana, _, _ := threeUsers(t)
	theme := func() string {
		t.Helper()
		user, err := repos.User.FindById(ana.ID)
		if err != nil {
			t.Fatal(err)
		}
		return user.Theme
	}

	if got := theme(); got != "2005" {
		t.Fatalf("new user theme = %q, want 2005", got)
	}
	if err := s.User.UpdateTheme(ana.ID, "2008"); err != nil {
		t.Fatal(err)
	}
	if got := theme(); got != "2008" {
		t.Fatalf("theme = %q, want 2008", got)
	}
	if err := s.User.UpdateTheme(ana.ID, "geocities"); !errors.Is(err, ErrUnknownTheme) {
		t.Fatalf("unknown theme: err = %v, want ErrUnknownTheme", err)
	}
	if got := theme(); got != "2008" {
		t.Fatalf("theme after a refused change = %q, want 2008", got)
	}
}
