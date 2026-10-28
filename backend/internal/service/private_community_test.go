package service

import (
	"context"
	"errors"
	"testing"
	"yorgut/internal/model"
	"yorgut/internal/repository"
)

func TestPrivateCommunities(t *testing.T) {
	type env struct {
		s                      *Services
		repos                  *repository.Repositories
		owner, asker, outsider *model.User
		communityId            uint
	}
	// setup creates a private community owned by owner, a topic in it, and a
	// pending request from asker.
	setup := func(t *testing.T) env {
		t.Helper()
		s, repos, owner, asker, outsider := threeUsers(t)
		id, err := s.Community.CreateCommunity(context.Background(), owner.ID, "Secreta", "other", "", true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := s.Topic.CreateTopic(owner.ID, id, "bem-vindos", "texto"); err != nil {
			t.Fatal(err)
		}
		pending, err := s.CommunityMembership.JoinCommunity(asker.ID, id)
		if err != nil || !pending {
			t.Fatalf("JoinCommunity: pending %v, err %v; want a pending request", pending, err)
		}
		return env{s, repos, owner, asker, outsider, id}
	}
	isMember := func(t *testing.T, e env, userId uint) (member bool, pending bool) {
		t.Helper()
		got, err := e.s.CommunityMembership.FindMembershipStatus(e.communityId, userId)
		if err != nil {
			t.Fatal(err)
		}
		return got.AmIMember, got.Pending
	}

	t.Run("creating makes the owner a member, without asking", func(t *testing.T) {
		e := setup(t)
		community, _ := e.repos.Community.FindById(e.communityId)
		if !community.Private {
			t.Error("community not private")
		}
		if member, pending := isMember(t, e, e.owner.ID); !member || pending {
			t.Errorf("owner: member %v, pending %v; want member", member, pending)
		}
	})

	t.Run("joining asks the owner instead", func(t *testing.T) {
		e := setup(t)
		if member, pending := isMember(t, e, e.asker.ID); member || !pending {
			t.Errorf("asker: member %v, pending %v; want pending", member, pending)
		}
		if got := memberCountOf(t, e.repos, e.communityId); got != 1 {
			t.Errorf("member count %d, want 1", got)
		}
		_, err := e.s.CommunityMembership.JoinCommunity(e.asker.ID, e.communityId)
		if !errors.Is(err, ErrAlreadyRequested) {
			t.Errorf("asking twice: got %v, want ErrAlreadyRequested", err)
		}
	})

	t.Run("only members read the forum", func(t *testing.T) {
		e := setup(t)
		topics, err := e.s.Topic.FindTopics(e.owner.ID, e.communityId, 1, 10)
		if err != nil || len(topics.Topics) != 1 {
			t.Fatalf("member: got %+v, err %v; want 1 topic", topics, err)
		}
		topicId := topics.Topics[0].ID

		for name, userId := range map[string]uint{"outsider": e.outsider.ID, "pending asker": e.asker.ID} {
			if _, err := e.s.Topic.FindTopics(userId, e.communityId, 1, 10); !errors.Is(err, ErrPrivateForum) {
				t.Errorf("%s topics: got %v, want ErrPrivateForum", name, err)
			}
			if _, err := e.s.Topic.FindTopicReplies(userId, topicId, 1, 10); !errors.Is(err, ErrPrivateForum) {
				t.Errorf("%s replies: got %v, want ErrPrivateForum", name, err)
			}
		}
		if _, err := e.s.Topic.FindTopics(e.outsider.ID, 9999, 1, 10); !errors.Is(err, ErrNotFound) {
			t.Errorf("missing community: got %v, want ErrNotFound", err)
		}
	})

	t.Run("only the owner sees and answers requests", func(t *testing.T) {
		e := setup(t)
		got, err := e.s.CommunityMembership.FindJoinRequests(e.owner.ID, e.communityId)
		if err != nil || len(got.Requests) != 1 || got.Requests[0].UserID != e.asker.ID {
			t.Fatalf("got %+v, err %v; want the asker's request", got, err)
		}
		if _, err := e.s.CommunityMembership.FindJoinRequests(e.asker.ID, e.communityId); !errors.Is(err, ErrForbidden) {
			t.Errorf("non-owner listing: got %v, want ErrForbidden", err)
		}
		err = e.s.CommunityMembership.ReplyJoinRequest(e.outsider.ID, e.communityId, e.asker.ID, true)
		if !errors.Is(err, ErrForbidden) {
			t.Errorf("non-owner replying: got %v, want ErrForbidden", err)
		}
		err = e.s.CommunityMembership.ReplyJoinRequest(e.owner.ID, e.communityId, e.outsider.ID, true)
		if !errors.Is(err, ErrJoinRequestNotFound) {
			t.Errorf("replying to nobody: got %v, want ErrJoinRequestNotFound", err)
		}
	})

	t.Run("accepting makes a member", func(t *testing.T) {
		e := setup(t)
		if err := e.s.CommunityMembership.ReplyJoinRequest(e.owner.ID, e.communityId, e.asker.ID, true); err != nil {
			t.Fatal(err)
		}
		if member, pending := isMember(t, e, e.asker.ID); !member || pending {
			t.Errorf("asker: member %v, pending %v; want member", member, pending)
		}
		if got := memberCountOf(t, e.repos, e.communityId); got != 2 {
			t.Errorf("member count %d, want 2", got)
		}
		if _, err := e.s.Topic.FindTopics(e.asker.ID, e.communityId, 1, 10); err != nil {
			t.Errorf("new member reading the forum: %v", err)
		}
		// Leaving and asking again works: the old request was deleted.
		if err := e.s.CommunityMembership.LeaveCommunity(e.asker.ID, e.communityId); err != nil {
			t.Fatal(err)
		}
		if pending, err := e.s.CommunityMembership.JoinCommunity(e.asker.ID, e.communityId); err != nil || !pending {
			t.Errorf("asking again: pending %v, err %v", pending, err)
		}
	})

	t.Run("declining drops the request", func(t *testing.T) {
		e := setup(t)
		if err := e.s.CommunityMembership.ReplyJoinRequest(e.owner.ID, e.communityId, e.asker.ID, false); err != nil {
			t.Fatal(err)
		}
		if member, pending := isMember(t, e, e.asker.ID); member || pending {
			t.Errorf("asker: member %v, pending %v; want neither", member, pending)
		}
	})

	t.Run("the asker cancels", func(t *testing.T) {
		e := setup(t)
		if err := e.s.CommunityMembership.CancelJoinRequest(e.asker.ID, e.communityId); err != nil {
			t.Fatal(err)
		}
		if _, pending := isMember(t, e, e.asker.ID); pending {
			t.Error("request still pending")
		}
		err := e.s.CommunityMembership.CancelJoinRequest(e.asker.ID, e.communityId)
		if !errors.Is(err, ErrJoinRequestNotFound) {
			t.Errorf("cancelling twice: got %v, want ErrJoinRequestNotFound", err)
		}
	})

	t.Run("making it public accepts everyone who asked", func(t *testing.T) {
		e := setup(t)
		if err := e.s.Community.UpdateCommunity(e.owner.ID, e.communityId, "Secreta", "other", "", false); err != nil {
			t.Fatal(err)
		}
		if member, pending := isMember(t, e, e.asker.ID); !member || pending {
			t.Errorf("asker: member %v, pending %v; want member", member, pending)
		}
		if got := memberCountOf(t, e.repos, e.communityId); got != 2 {
			t.Errorf("member count %d, want 2", got)
		}
		if _, err := e.s.Topic.FindTopics(e.outsider.ID, e.communityId, 1, 10); err != nil {
			t.Errorf("outsider reading a public forum: %v", err)
		}
		// Public now: joining is immediate.
		if pending, err := e.s.CommunityMembership.JoinCommunity(e.outsider.ID, e.communityId); err != nil || pending {
			t.Errorf("joining a public community: pending %v, err %v", pending, err)
		}
	})

	t.Run("making it private keeps the members", func(t *testing.T) {
		s, repos, owner, member, _ := threeUsers(t)
		id, err := s.Community.CreateCommunity(context.Background(), owner.ID, "Aberta", "other", "", false, nil)
		if err != nil {
			t.Fatal(err)
		}
		if pending, err := s.CommunityMembership.JoinCommunity(member.ID, id); err != nil || pending {
			t.Fatalf("joining: pending %v, err %v", pending, err)
		}
		if err := s.Community.UpdateCommunity(owner.ID, id, "Aberta", "other", "", true); err != nil {
			t.Fatal(err)
		}
		if got := memberCountOf(t, repos, id); got != 2 {
			t.Errorf("member count %d, want 2", got)
		}
		if _, err := s.Topic.FindTopics(member.ID, id, 1, 10); err != nil {
			t.Errorf("member reading: %v", err)
		}
	})

	t.Run("discovery shows private and pending", func(t *testing.T) {
		e := setup(t)
		got, err := e.s.Community.SearchCommunitiesToDiscover(e.asker.ID, "", 1, 10)
		if err != nil || len(got.Communities) != 1 {
			t.Fatalf("got %+v, err %v; want 1 community", got, err)
		}
		if c := got.Communities[0]; !c.Private || !c.Pending {
			t.Errorf("asker sees private %v, pending %v; want both", c.Private, c.Pending)
		}
		got, err = e.s.Community.SearchCommunitiesToDiscover(e.outsider.ID, "", 1, 10)
		if err != nil || len(got.Communities) != 1 || got.Communities[0].Pending {
			t.Errorf("outsider: got %+v, err %v; want 1 community, not pending", got, err)
		}
	})
}
