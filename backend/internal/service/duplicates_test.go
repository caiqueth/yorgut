package service

import (
	"testing"
	"yorgut/internal/model"
)

// The partial unique indexes (migration 00003, mirrored for SQLite by
// database.AutoMigrateSQLite) reject a second live row even when the
// services' "already exists" check is bypassed, as two concurrent requests
// could do, while soft-deleted rows don't count.

func TestDuplicateMembershipRejected(t *testing.T) {
	repos := newTestRepositories(t)
	services := NewServices(repos, nil)
	owner := createTestUser(t, repos, "owner@test.com")
	member := createTestUser(t, repos, "member@test.com")
	communityId, err := services.Community.CreateCommunity(owner.ID, "Eu odeio acordar cedo", "Outros", nil)
	if err != nil {
		t.Fatalf("CreateCommunity: %v", err)
	}

	// The owner is already a member.
	if err := repos.CommunityMembership.Create(&model.CommunityMembership{CommunityID: communityId, UserID: owner.ID}); err == nil {
		t.Error("duplicate membership was inserted")
	}

	// Leaving keeps the old row soft-deleted; rejoining must still work.
	memberships := services.CommunityMembership
	if err := memberships.JoinCommunity(member.ID, communityId); err != nil {
		t.Fatalf("JoinCommunity: %v", err)
	}
	if err := memberships.LeaveCommunity(member.ID, communityId); err != nil {
		t.Fatalf("LeaveCommunity: %v", err)
	}
	if err := memberships.JoinCommunity(member.ID, communityId); err != nil {
		t.Fatalf("rejoining: %v", err)
	}
}

func TestDuplicateFriendshipRejected(t *testing.T) {
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

	// Same pair, in either direction.
	if err := repos.Friendship.Create(&model.Friendship{UserRequestedID: ana.ID, UserRequesteeID: bia.ID}); err == nil {
		t.Error("duplicate friendship was inserted")
	}
	if err := repos.Friendship.Create(&model.Friendship{UserRequestedID: bia.ID, UserRequesteeID: ana.ID}); err == nil {
		t.Error("reverse duplicate friendship was inserted")
	}

	// Once the first one is (soft) deleted, asking again works.
	if err := repos.Friendship.Delete(pending[0].ID); err != nil {
		t.Fatalf("deleting friendship: %v", err)
	}
	if err := friendships.CreateFriendship(bia.ID, ana.ID); err != nil {
		t.Errorf("asking again after deletion: %v", err)
	}
}
