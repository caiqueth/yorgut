package service

import (
	"context"
	"errors"
	"testing"
	"yorgut/internal/model"
	"yorgut/internal/repository"
)

func TestCommunityInvitations(t *testing.T) {
	type env struct {
		s              *Services
		repos          *repository.Repositories
		ana, bia, caio *model.User
		communityId    uint
	}
	// setup creates a community owned by ana with bia as a member. Ana and
	// bia are friends, and so are bia and caio; ana and caio aren't.
	setup := func(t *testing.T, private bool) env {
		t.Helper()
		s, repos, ana, bia, caio := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		befriend(t, s, repos, bia, caio)
		id, err := s.Community.CreateCommunity(context.Background(), ana.ID, "Comunidade", "other", "", private, nil)
		if err != nil {
			t.Fatal(err)
		}
		e := env{s, repos, ana, bia, caio, id}
		if private {
			// The owner's invitation lets bia straight in.
			acceptInvitation(t, e.s, e.ana, e.bia, id)
		} else if _, err := s.CommunityMembership.JoinCommunity(bia.ID, id); err != nil {
			t.Fatal(err)
		}
		return e
	}
	invitationsOf := func(t *testing.T, e env, user *model.User) []uint {
		t.Helper()
		got, err := e.s.CommunityMembership.FindMyInvitations(user.ID)
		if err != nil {
			t.Fatal(err)
		}
		ids := []uint{}
		for _, inv := range got.Invitations {
			ids = append(ids, inv.ID)
		}
		return ids
	}
	isMember := func(t *testing.T, e env, user *model.User) (member bool, pending bool) {
		t.Helper()
		got, err := e.s.CommunityMembership.FindMembershipStatus(e.communityId, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		return got.AmIMember, got.Pending
	}

	t.Run("a member invites a friend, who accepts and joins", func(t *testing.T) {
		e := setup(t, false)
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Fatal(err)
		}

		friends, err := e.s.CommunityMembership.FindInvitableFriends(e.bia.ID, e.communityId, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range friends.Friends {
			wantMember, wantInvited := f.ID == e.ana.ID, f.ID == e.caio.ID
			if f.Member != wantMember || f.Invited != wantInvited {
				t.Errorf("friend %d: member %v, invited %v; want %v, %v", f.ID, f.Member, f.Invited, wantMember, wantInvited)
			}
		}

		invitations, err := e.s.CommunityMembership.FindMyInvitations(e.caio.ID)
		if err != nil || len(invitations.Invitations) != 1 {
			t.Fatalf("invitations: got %+v, err %v; want 1", invitations, err)
		}
		inv := invitations.Invitations[0]
		if inv.CommunityID != e.communityId || inv.InviterID != e.bia.ID {
			t.Errorf("invitation %+v, want community %d from bia", inv, e.communityId)
		}

		pending, err := e.s.CommunityMembership.AcceptInvitation(e.caio.ID, inv.ID)
		if err != nil || pending {
			t.Fatalf("AcceptInvitation: pending %v, err %v; want joined", pending, err)
		}
		if member, _ := isMember(t, e, e.caio); !member {
			t.Error("caio is not a member after accepting")
		}
		if got := memberCountOf(t, e.repos, e.communityId); got != 3 {
			t.Errorf("member count %d, want 3", got)
		}
		if ids := invitationsOf(t, e, e.caio); len(ids) != 0 {
			t.Errorf("invitations after accepting: %v, want none", ids)
		}
	})

	t.Run("declining drops the invitation, and the friend may be invited again", func(t *testing.T) {
		e := setup(t, false)
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Fatal(err)
		}
		ids := invitationsOf(t, e, e.caio)
		if err := e.s.CommunityMembership.DeclineInvitation(e.caio.ID, ids[0]); err != nil {
			t.Fatal(err)
		}
		if member, _ := isMember(t, e, e.caio); member {
			t.Error("caio became a member by declining")
		}
		if ids := invitationsOf(t, e, e.caio); len(ids) != 0 {
			t.Errorf("invitations after declining: %v, want none", ids)
		}
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Errorf("inviting again: %v", err)
		}
	})

	t.Run("only members invite, and only their friends", func(t *testing.T) {
		e := setup(t, false)
		if err := e.s.CommunityMembership.InviteToCommunity(e.caio.ID, e.communityId, e.bia.ID); !errors.Is(err, ErrInviterNotMember) {
			t.Errorf("non-member inviting: got %v, want ErrInviterNotMember", err)
		}
		if _, err := e.s.CommunityMembership.FindInvitableFriends(e.caio.ID, e.communityId, 1, 10); !errors.Is(err, ErrInviterNotMember) {
			t.Errorf("non-member listing friends: got %v, want ErrInviterNotMember", err)
		}
		if err := e.s.CommunityMembership.InviteToCommunity(e.ana.ID, e.communityId, e.caio.ID); !errors.Is(err, ErrInviteeNotFriend) {
			t.Errorf("inviting a non-friend: got %v, want ErrInviteeNotFriend", err)
		}
		if err := e.s.CommunityMembership.InviteToCommunity(e.ana.ID, e.communityId, e.ana.ID); !errors.Is(err, ErrInviteeNotFriend) {
			t.Errorf("inviting oneself: got %v, want ErrInviteeNotFriend", err)
		}
	})

	t.Run("members and the already invited can't be invited", func(t *testing.T) {
		e := setup(t, false)
		if err := e.s.CommunityMembership.InviteToCommunity(e.ana.ID, e.communityId, e.bia.ID); !errors.Is(err, ErrInviteeAlreadyMember) {
			t.Errorf("inviting a member: got %v, want ErrInviteeAlreadyMember", err)
		}
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Fatal(err)
		}
		befriend(t, e.s, e.repos, e.ana, e.caio)
		if err := e.s.CommunityMembership.InviteToCommunity(e.ana.ID, e.communityId, e.caio.ID); !errors.Is(err, ErrAlreadyInvited) {
			t.Errorf("inviting twice: got %v, want ErrAlreadyInvited", err)
		}
	})

	t.Run("only the invitee answers an invitation", func(t *testing.T) {
		e := setup(t, false)
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Fatal(err)
		}
		id := invitationsOf(t, e, e.caio)[0]
		if _, err := e.s.CommunityMembership.AcceptInvitation(e.ana.ID, id); !errors.Is(err, ErrInvitationNotFound) {
			t.Errorf("accepting someone else's: got %v, want ErrInvitationNotFound", err)
		}
		if err := e.s.CommunityMembership.DeclineInvitation(e.bia.ID, id); !errors.Is(err, ErrInvitationNotFound) {
			t.Errorf("declining someone else's: got %v, want ErrInvitationNotFound", err)
		}
	})

	t.Run("joining another way drops the invitation", func(t *testing.T) {
		e := setup(t, false)
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := e.s.CommunityMembership.JoinCommunity(e.caio.ID, e.communityId); err != nil {
			t.Fatal(err)
		}
		if ids := invitationsOf(t, e, e.caio); len(ids) != 0 {
			t.Errorf("invitations after joining: %v, want none", ids)
		}
	})

	t.Run("in a private community a member's invitation asks the owner", func(t *testing.T) {
		e := setup(t, true)
		if member, _ := isMember(t, e, e.bia); !member {
			t.Fatal("the owner's invitation didn't let bia in")
		}
		if err := e.s.CommunityMembership.InviteToCommunity(e.bia.ID, e.communityId, e.caio.ID); err != nil {
			t.Fatal(err)
		}
		pending, err := e.s.CommunityMembership.AcceptInvitation(e.caio.ID, invitationsOf(t, e, e.caio)[0])
		if err != nil || !pending {
			t.Fatalf("AcceptInvitation: pending %v, err %v; want a pending request", pending, err)
		}
		if member, isPending := isMember(t, e, e.caio); member || !isPending {
			t.Errorf("caio: member %v, pending %v; want pending", member, isPending)
		}
		if ids := invitationsOf(t, e, e.caio); len(ids) != 0 {
			t.Errorf("invitations after accepting: %v, want none", ids)
		}
		if err := e.s.CommunityMembership.ReplyJoinRequest(e.ana.ID, e.communityId, e.caio.ID, true); err != nil {
			t.Fatal(err)
		}
		if member, _ := isMember(t, e, e.caio); !member {
			t.Error("caio is not a member after the owner accepted")
		}
	})
}

// acceptInvitation has inviter invite invitee to the community, and invitee
// accept it.
func acceptInvitation(t *testing.T, s *Services, inviter, invitee *model.User, communityId uint) {
	t.Helper()
	if err := s.CommunityMembership.InviteToCommunity(inviter.ID, communityId, invitee.ID); err != nil {
		t.Fatalf("InviteToCommunity: %v", err)
	}
	invitations, err := s.CommunityMembership.FindMyInvitations(invitee.ID)
	if err != nil || len(invitations.Invitations) == 0 {
		t.Fatalf("FindMyInvitations: got %+v, err %v", invitations, err)
	}
	if _, err := s.CommunityMembership.AcceptInvitation(invitee.ID, invitations.Invitations[0].ID); err != nil {
		t.Fatalf("AcceptInvitation: %v", err)
	}
}
