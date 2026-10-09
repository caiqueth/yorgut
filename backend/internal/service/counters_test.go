package service

import (
	"errors"
	"testing"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
)

func createTestUser(t *testing.T, repos *repository.Repositories, email string) *model.User {
	t.Helper()
	if err := repos.User.Create(model.NewUser("Test", "User", email, time.Time{}, nil)); err != nil {
		t.Fatalf("creating user: %v", err)
	}
	user, err := repos.User.FindByEmail(email)
	if err != nil {
		t.Fatalf("loading user: %v", err)
	}
	return user
}

func friendCountOf(t *testing.T, repos *repository.Repositories, userId uint) uint {
	t.Helper()
	user, err := repos.User.FindById(userId)
	if err != nil {
		t.Fatalf("loading user: %v", err)
	}
	return user.FriendCount
}

func memberCountOf(t *testing.T, repos *repository.Repositories, communityId uint) uint {
	t.Helper()
	community, err := repos.Community.FindById(communityId)
	if err != nil {
		t.Fatalf("loading community: %v", err)
	}
	return community.MemberCount
}

func TestFriendCounts(t *testing.T) {
	setup := func(t *testing.T) (*FriendshipService, *repository.Repositories, *model.User, *model.User, uint) {
		repos := newTestRepositories(t)
		friendships := NewServices(repos, nil).Friendship
		ana := createTestUser(t, repos, "ana@test.com")
		bia := createTestUser(t, repos, "bia@test.com")
		if err := friendships.CreateFriendship(ana.ID, bia.ID); err != nil {
			t.Fatalf("CreateFriendship: %v", err)
		}
		pending, err := repos.Friendship.FindPendingFriendshipsOfUser(bia.ID)
		if err != nil || len(pending) != 1 {
			t.Fatalf("pending friendships: %v, err %v", pending, err)
		}
		return friendships, repos, ana, bia, pending[0].ID
	}

	t.Run("accepting and removing a friend", func(t *testing.T) {
		friendships, repos, ana, bia, friendshipId := setup(t)

		if err := friendships.ReplyFriendship(friendshipId, true); err != nil {
			t.Fatalf("ReplyFriendship: %v", err)
		}
		if a, b := friendCountOf(t, repos, ana.ID), friendCountOf(t, repos, bia.ID); a != 1 || b != 1 {
			t.Fatalf("after accepting: counts %d and %d, want 1 and 1", a, b)
		}

		if err := friendships.DeleteFriendship(friendshipId); err != nil {
			t.Fatalf("DeleteFriendship: %v", err)
		}
		if a, b := friendCountOf(t, repos, ana.ID), friendCountOf(t, repos, bia.ID); a != 0 || b != 0 {
			t.Fatalf("after removing: counts %d and %d, want 0 and 0", a, b)
		}
	})

	t.Run("cancelling a pending request keeps the counts", func(t *testing.T) {
		friendships, repos, ana, bia, friendshipId := setup(t)
		// Ana and Bia already have one friend each, through Caio.
		caio := createTestUser(t, repos, "caio@test.com")
		for _, user := range []*model.User{ana, bia} {
			if err := friendships.CreateFriendship(caio.ID, user.ID); err != nil {
				t.Fatalf("CreateFriendship: %v", err)
			}
			pending, _ := repos.Friendship.FindPendingFriendshipsOfUser(user.ID)
			for _, p := range pending {
				if p.UserRequestedID == caio.ID {
					if err := friendships.ReplyFriendship(p.ID, true); err != nil {
						t.Fatalf("ReplyFriendship: %v", err)
					}
				}
			}
		}

		// Used to decrement both counts even though the request was never accepted.
		if err := friendships.DeleteFriendship(friendshipId); err != nil {
			t.Fatalf("DeleteFriendship: %v", err)
		}
		if a, b := friendCountOf(t, repos, ana.ID), friendCountOf(t, repos, bia.ID); a != 1 || b != 1 {
			t.Fatalf("counts %d and %d, want 1 and 1", a, b)
		}
	})

	t.Run("declining a request keeps the counts", func(t *testing.T) {
		friendships, repos, ana, bia, friendshipId := setup(t)

		if err := friendships.ReplyFriendship(friendshipId, false); err != nil {
			t.Fatalf("ReplyFriendship: %v", err)
		}
		if a, b := friendCountOf(t, repos, ana.ID), friendCountOf(t, repos, bia.ID); a != 0 || b != 0 {
			t.Fatalf("counts %d and %d, want 0 and 0", a, b)
		}
	})
}

func TestMemberCounts(t *testing.T) {
	repos := newTestRepositories(t)
	services := NewServices(repos, nil)
	owner := createTestUser(t, repos, "owner@test.com")
	member := createTestUser(t, repos, "member@test.com")

	communityId, err := services.Community.CreateCommunity(owner.ID, "Eu odeio acordar cedo", "Outros", nil)
	if err != nil {
		t.Fatalf("CreateCommunity: %v", err)
	}
	if got := memberCountOf(t, repos, communityId); got != 1 {
		t.Fatalf("after creating: %d members, want 1", got)
	}

	if err := services.CommunityMembership.JoinCommunity(member.ID, communityId); err != nil {
		t.Fatalf("JoinCommunity: %v", err)
	}
	if got := memberCountOf(t, repos, communityId); got != 2 {
		t.Fatalf("after joining: %d members, want 2", got)
	}

	if err := services.CommunityMembership.LeaveCommunity(member.ID, communityId); err != nil {
		t.Fatalf("LeaveCommunity: %v", err)
	}
	if got := memberCountOf(t, repos, communityId); got != 1 {
		t.Fatalf("after leaving: %d members, want 1 (soft-deleted memberships must not count)", got)
	}

	if err := services.CommunityMembership.LeaveCommunity(owner.ID, communityId); !errors.Is(err, ErrOwnerCannotLeave) {
		t.Fatalf("owner leaving: got err %v, want ErrOwnerCannotLeave", err)
	}
}

func TestTransactionRollsBack(t *testing.T) {
	repos := newTestRepositories(t)
	boom := errors.New("boom")

	err := repos.Transactor.Transaction(func(tx *repository.Repositories) error {
		if err := tx.User.Create(model.NewUser("Rolled", "Back", "rollback@test.com", time.Time{}, nil)); err != nil {
			t.Fatalf("creating user: %v", err)
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("got err %v, want the callback's error", err)
	}

	if _, err := repos.User.FindByEmail("rollback@test.com"); err == nil {
		t.Fatalf("the user created inside the failed transaction was kept")
	}
}
