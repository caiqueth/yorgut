package service

import (
	"context"
	"errors"
	"testing"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
)

// TestDeleteAccount deletes a user who did a bit of everything and checks
// that nothing of theirs is left, that their communities were handed over or
// deleted, and that the others' counts follow.
func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	// SQLite only checks foreign keys when asked: ask, so a delete in the
	// wrong order fails here as it would on MySQL.
	if db.Dialector.Name() == "sqlite" {
		if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			t.Fatal(err)
		}
	}
	repos := repository.NewRepositories(db)
	store := storage.NewMemory()
	s := NewServices(repos, store)

	ana := createTestUser(t, repos, "ana@test.com")
	bia := createTestUser(t, repos, "bia@test.com")
	caio := createTestUser(t, repos, "caio@test.com")
	befriend(t, s, repos, ana, bia)
	befriend(t, s, repos, caio, ana)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Pictures: profile, album, and the community that gets deleted.
	must(s.UserProfile.UpdateProfilePicture(ctx, ana.ID, testPNG(t)))
	must(s.Photo.UploadPhoto(ctx, ana.ID, testPNG(t), "foto"))
	lonely, err := s.Community.CreateCommunity(ctx, ana.ID, "Só a Ana", "other", "", false, testPNG(t))
	must(err)

	// A community Bia joined after Ana created it: it goes to Bia.
	shared, err := s.Community.CreateCommunity(ctx, ana.ID, "Da Ana e da Bia", "other", "", true, nil)
	must(err)
	_, err = s.CommunityMembership.JoinCommunity(bia.ID, shared)
	must(err)
	must(s.CommunityMembership.ReplyJoinRequest(ana.ID, shared, bia.ID, true))

	// Caio's community, where Ana is a member and writes.
	caios, err := s.Community.CreateCommunity(ctx, caio.ID, "Do Caio", "other", "", false, nil)
	must(err)
	_, err = s.CommunityMembership.JoinCommunity(ana.ID, caios)
	must(err)
	must(s.Topic.CreateTopic(caio.ID, caios, "do Caio", "texto"))
	must(s.Topic.CreateTopic(ana.ID, caios, "da Ana", "texto"))
	topics, _, err := repos.Topic.FindOfCommunity(caios, 1, 10)
	must(err)
	var caioTopic, anaTopic uint
	for _, topic := range topics {
		if topic.AuthorID == caio.ID {
			caioTopic = topic.ID
		} else {
			anaTopic = topic.ID
		}
	}
	must(s.Topic.ReplyTopic(ana.ID, caioTopic, "", "resposta da Ana"))
	must(s.Topic.ReplyTopic(caio.ID, anaTopic, "", "resposta do Caio"))

	anaPoll, err := s.Poll.CreatePoll(ana.ID, caios, CreatePollInput{Question: "Da Ana?", Options: []string{"sim", "não"}})
	must(err)
	caioPoll, err := s.Poll.CreatePoll(caio.ID, caios, CreatePollInput{Question: "Do Caio?", Options: []string{"sim", "não"}})
	must(err)
	for _, vote := range []struct {
		user *model.User
		poll uint
	}{{caio, anaPoll}, {ana, caioPoll}, {caio, caioPoll}} {
		poll, err := s.Poll.FindPoll(vote.user.ID, vote.poll)
		must(err)
		must(s.Poll.Vote(vote.user.ID, vote.poll, poll.Options[0].ID))
	}

	// Leaving and rejoining leaves a soft-deleted membership behind.
	must(s.CommunityMembership.LeaveCommunity(ana.ID, caios))
	_, err = s.CommunityMembership.JoinCommunity(ana.ID, caios)
	must(err)

	must(s.CommunityMembership.InviteToCommunity(ana.ID, caios, bia.ID))
	must(s.Scrap.CreateScrap(ana.ID, bia.ID, "oi Bia"))
	must(s.Scrap.CreateScrap(bia.ID, ana.ID, "oi Ana"))
	must(s.Scrap.CreateScrap(bia.ID, caio.ID, "oi Caio"))
	must(s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "a Bia é legal"))
	must(s.Testimonial.CreateTestimonial(bia.ID, ana.ID, "a Ana é legal"))
	must(s.UserFan.BecomeFan(ana.ID, bia.ID))
	must(s.UserFan.BecomeFan(bia.ID, ana.ID))
	must(s.UserReview.UpsertReview(ana.ID, bia.ID, 3, 3, 3))
	must(s.UserReview.UpsertReview(bia.ID, ana.ID, 3, 3, 3))
	s.ProfileVisit.RecordVisit(ana.ID, bia.ID)
	s.ProfileVisit.RecordVisit(bia.ID, ana.ID)

	if store.Len() != 3 {
		t.Fatalf("stored pictures before: %d, want 3", store.Len())
	}

	must(s.User.DeleteAccount(ctx, ana.ID))

	if _, err := repos.User.FindById(ana.ID); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("the user is still there: %v", err)
	}
	if store.Len() != 0 {
		t.Errorf("stored pictures after: %d, want 0", store.Len())
	}

	// No row anywhere still points at Ana, soft-deleted ones included.
	for _, c := range []struct {
		value any
		query string
	}{
		{&model.Scrap{}, "user_author_id = ? or user_recipient_id = ?"},
		{&model.Testimonial{}, "user_author_id = ? or user_recipient_id = ?"},
		{&model.Friendship{}, "user_requested_id = ? or user_requestee_id = ?"},
		{&model.UserFan{}, "fan_id = ? or idol_id = ?"},
		{&model.UserReview{}, "reviewer_id = ? or reviewed_id = ?"},
		{&model.ProfileVisit{}, "visitor_id = ? or visited_id = ?"},
		{&model.CommunityInvitation{}, "inviter_id = ? or invitee_id = ?"},
		{&model.CommunityMembership{}, "user_id = ? or user_id = ?"},
		{&model.CommunityJoinRequest{}, "user_id = ? or user_id = ?"},
		{&model.Topic{}, "author_id = ? or author_id = ?"},
		{&model.Poll{}, "author_id = ? or author_id = ?"},
		{&model.PollVote{}, "user_id = ? or user_id = ?"},
		{&model.PhotoAlbum{}, "user_id = ? or user_id = ?"},
		{&model.UserProfile{}, "user_id = ? or user_id = ?"},
		{&model.Community{}, "owner_id = ? or created_by_id = ?"},
	} {
		var count int64
		must(db.Unscoped().Model(c.value).Where(c.query, ana.ID, ana.ID).Count(&count).Error)
		if count != 0 {
			t.Errorf("%T: %d rows left", c.value, count)
		}
	}

	// The community only Ana was in is gone; the shared one is Bia's now.
	if _, err := repos.Community.FindById(lonely); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("the community only Ana was in is still there: %v", err)
	}
	community, err := repos.Community.FindById(shared)
	must(err)
	if community.OwnerID != bia.ID || community.CreatedByID != bia.ID {
		t.Errorf("shared community owner %d, creator %d; want Bia (%d)", community.OwnerID, community.CreatedByID, bia.ID)
	}
	if got := memberCountOf(t, repos, shared); got != 1 {
		t.Errorf("shared community members: %d, want 1", got)
	}
	if got := memberCountOf(t, repos, caios); got != 1 {
		t.Errorf("Caio's community members: %d, want 1", got)
	}

	// Ana's topic went with Caio's reply to it; Caio's topic stays, without
	// Ana's reply.
	var replies int64
	must(db.Model(&model.Topic{}).Where("parent_topic_id = ?", caioTopic).Count(&replies).Error)
	if replies != 0 {
		t.Errorf("replies left on Caio's topic: %d, want 0", replies)
	}
	if _, err := repos.Topic.FindById(caioTopic); err != nil {
		t.Errorf("Caio's topic: %v", err)
	}
	if _, err := repos.Topic.FindById(anaTopic); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Ana's topic is still there: %v", err)
	}
	// Caio's poll keeps his own vote only.
	poll, err := s.Poll.FindPoll(caio.ID, caioPoll)
	must(err)
	if votes := poll.Options[0].Votes; votes == nil || *votes != 1 {
		t.Errorf("votes on Caio's poll: %v, want 1", votes)
	}

	if got := friendCountOf(t, repos, bia.ID); got != 0 {
		t.Errorf("Bia's friends: %d, want 0", got)
	}
	if got := friendCountOf(t, repos, caio.ID); got != 0 {
		t.Errorf("Caio's friends: %d, want 0", got)
	}
	// Others' things stay.
	scraps, err := repos.Scrap.CountOfUserRecipient(caio.ID)
	must(err)
	if scraps != 1 {
		t.Errorf("Caio's scraps: %d, want 1", scraps)
	}
}
