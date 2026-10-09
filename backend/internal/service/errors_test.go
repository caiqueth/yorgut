package service

import (
	"errors"
	"testing"
	"time"

	"yorgut/internal/model"
	"yorgut/internal/repository"
)

func wantCategory(t *testing.T, err, category error) {
	t.Helper()
	if !errors.Is(err, category) {
		t.Fatalf("got err %v, want category %v", err, category)
	}
}

func TestRepositoryDuplicateIsTranslated(t *testing.T) {
	repos := newTestRepositories(t)
	createTestUser(t, repos, "dup@test.com")

	err := repos.User.Create(model.NewUser("Other", "User", "dup@test.com", time.Time{}, nil))
	if !errors.Is(err, repository.ErrDuplicate) {
		t.Fatalf("got err %v, want repository.ErrDuplicate", err)
	}
}

func TestErrorCategories(t *testing.T) {
	t.Run("user", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		user := createTestUser(t, repos, "u@test.com")

		_, err := s.User.FindUserById(9999)
		wantCategory(t, err, ErrNotFound)

		err = s.User.UpdateUser(user.ID, UpdateUserInput{FirstName: "   "})
		wantCategory(t, err, ErrInvalid)

		err = s.User.CreateUser("A", "B", "u@test.com", time.Time{}, nil)
		wantCategory(t, err, ErrConflict)
	})

	t.Run("community", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		owner := createTestUser(t, repos, "owner@test.com")
		other := createTestUser(t, repos, "other@test.com")
		id, err := s.Community.CreateCommunity(owner.ID, "Comunidade", "Outros", nil)
		if err != nil {
			t.Fatalf("CreateCommunity: %v", err)
		}

		wantCategory(t, s.Community.UpdateCommunity(other.ID, id, "n", "c", "d"), ErrForbidden)
		wantCategory(t, s.Community.UpdateCommunity(owner.ID, 9999, "n", "c", "d"), ErrNotFound)
		if err := s.Community.UpdateCommunity(owner.ID, id, "n", "c", "d"); err != nil {
			t.Fatalf("owner update: %v", err)
		}
	})

	t.Run("membership", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		owner := createTestUser(t, repos, "owner@test.com")
		member := createTestUser(t, repos, "member@test.com")
		id, err := s.Community.CreateCommunity(owner.ID, "Comunidade", "Outros", nil)
		if err != nil {
			t.Fatalf("CreateCommunity: %v", err)
		}

		if err := s.CommunityMembership.JoinCommunity(member.ID, id); err != nil {
			t.Fatalf("first join: %v", err)
		}
		err = s.CommunityMembership.JoinCommunity(member.ID, id)
		wantCategory(t, err, ErrConflict)
		if !errors.Is(err, ErrAlreadyMember) {
			t.Fatalf("got %v, want ErrAlreadyMember", err)
		}
		wantCategory(t, s.CommunityMembership.JoinCommunity(member.ID, 9999), ErrNotFound)

		wantCategory(t, s.CommunityMembership.LeaveCommunity(owner.ID, id), ErrForbidden)

		stranger := createTestUser(t, repos, "stranger@test.com")
		wantCategory(t, s.CommunityMembership.LeaveCommunity(stranger.ID, id), ErrNotFound)
	})

	t.Run("friendship", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		ana := createTestUser(t, repos, "ana@test.com")
		bia := createTestUser(t, repos, "bia@test.com")

		if err := s.Friendship.CreateFriendship(ana.ID, bia.ID); err != nil {
			t.Fatalf("first request: %v", err)
		}
		wantCategory(t, s.Friendship.CreateFriendship(ana.ID, bia.ID), ErrConflict)
		wantCategory(t, s.Friendship.CreateFriendship(bia.ID, ana.ID), ErrConflict)
		wantCategory(t, s.Friendship.CreateFriendship(ana.ID, 9999), ErrNotFound)
		wantCategory(t, s.Friendship.DeleteFriendship(9999), ErrNotFound)
	})

	t.Run("testimonial", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		ana := createTestUser(t, repos, "ana@test.com")
		bia := createTestUser(t, repos, "bia@test.com")
		caio := createTestUser(t, repos, "caio@test.com")

		wantCategory(t, s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "oi"), ErrForbidden)

		if err := s.Friendship.CreateFriendship(ana.ID, bia.ID); err != nil {
			t.Fatalf("CreateFriendship: %v", err)
		}
		pending, err := repos.Friendship.FindPendingFriendshipsOfUser(bia.ID)
		if err != nil || len(pending) != 1 {
			t.Fatalf("pending friendships: %v, err %v", pending, err)
		}
		if err := s.Friendship.ReplyFriendship(pending[0].ID, true); err != nil {
			t.Fatalf("ReplyFriendship: %v", err)
		}
		if err := s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "oi"); err != nil {
			t.Fatalf("CreateTestimonial between friends: %v", err)
		}
		testimonials, _, err := repos.Testimonial.FindOfUserRecipient(bia.ID, false, 0, 0)
		if err != nil || len(testimonials) != 1 {
			t.Fatalf("pending testimonials: %v, err %v", testimonials, err)
		}
		id := testimonials[0].ID

		// Only the recipient may approve: not the author, not a stranger.
		wantCategory(t, s.Testimonial.ApproveTestimonial(ana.ID, id), ErrForbidden)
		wantCategory(t, s.Testimonial.ApproveTestimonial(caio.ID, id), ErrForbidden)
		wantCategory(t, s.Testimonial.ApproveTestimonial(bia.ID, 9999), ErrNotFound)
		if err := s.Testimonial.ApproveTestimonial(bia.ID, id); err != nil {
			t.Fatalf("recipient approving: %v", err)
		}
	})

	t.Run("topic", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		owner := createTestUser(t, repos, "owner@test.com")
		stranger := createTestUser(t, repos, "stranger@test.com")
		id, err := s.Community.CreateCommunity(owner.ID, "Comunidade", "Outros", nil)
		if err != nil {
			t.Fatalf("CreateCommunity: %v", err)
		}

		wantCategory(t, s.Topic.CreateTopic(stranger.ID, id, "t", "text"), ErrForbidden)
		wantCategory(t, s.Topic.ReplyTopic(owner.ID, 9999, "t", "text"), ErrNotFound)
		if err := s.Topic.CreateTopic(owner.ID, id, "t", "text"); err != nil {
			t.Fatalf("member creating topic: %v", err)
		}
	})
}
