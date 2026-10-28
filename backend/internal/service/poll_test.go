package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"yorgut/internal/model"
)

func TestPolls(t *testing.T) {
	type env struct {
		s              *Services
		ana, bia, caio *model.User
		communityId    uint
	}
	// setup creates a community owned by ana, with bia as a member; caio is
	// an outsider.
	setup := func(t *testing.T, private bool) env {
		t.Helper()
		s, repos, ana, bia, caio := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		id, err := s.Community.CreateCommunity(context.Background(), ana.ID, "Comunidade", "other", "", private, nil)
		if err != nil {
			t.Fatal(err)
		}
		if private {
			acceptInvitation(t, s, ana, bia, id)
		} else if _, err := s.CommunityMembership.JoinCommunity(bia.ID, id); err != nil {
			t.Fatal(err)
		}
		return env{s, ana, bia, caio, id}
	}
	createPoll := func(t *testing.T, e env, author *model.User, options ...string) uint {
		t.Helper()
		id, err := e.s.Poll.CreatePoll(author.ID, e.communityId, CreatePollInput{Question: "Qual?", Options: options})
		if err != nil {
			t.Fatalf("CreatePoll: %v", err)
		}
		return id
	}
	findPoll := func(t *testing.T, e env, viewer *model.User, pollId uint) (optionIds []uint, votes []*int64, canVote, canClose, closed bool) {
		t.Helper()
		poll, err := e.s.Poll.FindPoll(viewer.ID, pollId)
		if err != nil {
			t.Fatalf("FindPoll: %v", err)
		}
		for _, o := range poll.Options {
			optionIds = append(optionIds, o.ID)
			votes = append(votes, o.Votes)
		}
		return optionIds, votes, poll.CanVote, poll.CanClose, poll.Closed
	}

	t.Run("members create polls with 2 to 10 different options", func(t *testing.T) {
		e := setup(t, false)
		createPoll(t, e, e.bia, "sim", "não")

		anHourAgo := time.Now().Add(-time.Hour)
		cases := map[string]struct {
			author *model.User
			input  CreatePollInput
			want   error
		}{
			"outsider":       {e.caio, CreatePollInput{Question: "Qual?", Options: []string{"a", "b"}}, ErrPollNotMember},
			"one option":     {e.bia, CreatePollInput{Question: "Qual?", Options: []string{"a"}}, ErrPollInvalidOptions},
			"repeated":       {e.bia, CreatePollInput{Question: "Qual?", Options: []string{"Sim", " sim "}}, ErrPollInvalidOptions},
			"blank option":   {e.bia, CreatePollInput{Question: "Qual?", Options: []string{"a", " "}}, ErrPollInvalidOptions},
			"blank question": {e.bia, CreatePollInput{Question: "  ", Options: []string{"a", "b"}}, ErrPollQuestionBlank},
			"closes in the past": {e.bia, CreatePollInput{
				Question: "Qual?", Options: []string{"a", "b"}, ClosesAt: &anHourAgo,
			}, ErrPollClosesInPast},
		}
		for name, c := range cases {
			if _, err := e.s.Poll.CreatePoll(c.author.ID, e.communityId, c.input); !errors.Is(err, c.want) {
				t.Errorf("%s: got %v, want %v", name, err, c.want)
			}
		}
	})

	t.Run("results show only after voting, and each member votes once", func(t *testing.T) {
		e := setup(t, false)
		pollId := createPoll(t, e, e.ana, "sim", "não")
		optionIds, votes, canVote, _, _ := findPoll(t, e, e.bia, pollId)
		if !canVote || votes[0] != nil || votes[1] != nil {
			t.Fatalf("before voting: canVote %v, votes %v; want canVote and no results", canVote, votes)
		}

		if err := e.s.Poll.Vote(e.bia.ID, pollId, optionIds[1]); err != nil {
			t.Fatal(err)
		}
		_, votes, canVote, _, _ = findPoll(t, e, e.bia, pollId)
		if canVote || votes[0] == nil || *votes[0] != 0 || *votes[1] != 1 {
			t.Errorf("after voting: canVote %v, votes %v; want 0 and 1", canVote, votes)
		}
		if _, votes, _, _, _ := findPoll(t, e, e.ana, pollId); votes[0] != nil {
			t.Error("ana sees results without voting")
		}

		if err := e.s.Poll.Vote(e.bia.ID, pollId, optionIds[0]); !errors.Is(err, ErrAlreadyVoted) {
			t.Errorf("voting twice: got %v, want ErrAlreadyVoted", err)
		}
		if err := e.s.Poll.Vote(e.caio.ID, pollId, optionIds[0]); !errors.Is(err, ErrPollNotMember) {
			t.Errorf("outsider voting: got %v, want ErrPollNotMember", err)
		}
		otherOptions, _, _, _, _ := findPoll(t, e, e.ana, createPoll(t, e, e.ana, "a", "b"))
		if err := e.s.Poll.Vote(e.ana.ID, pollId, otherOptions[0]); !errors.Is(err, ErrPollOptionNotFound) {
			t.Errorf("another poll's option: got %v, want ErrPollOptionNotFound", err)
		}

		polls, err := e.s.Poll.FindPolls(e.bia.ID, e.communityId, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range polls.Polls {
			wantVoted, wantCount := p.ID == pollId, map[bool]int64{true: 1, false: 0}[p.ID == pollId]
			if p.Voted != wantVoted || p.VoteCount != wantCount {
				t.Errorf("poll %d in the list: voted %v, %d votes; want %v, %d", p.ID, p.Voted, p.VoteCount, wantVoted, wantCount)
			}
		}
	})

	t.Run("the author or the owner closes a poll, which ends voting and shows results", func(t *testing.T) {
		e := setup(t, false)
		pollId := createPoll(t, e, e.bia, "sim", "não")
		if _, _, _, canClose, _ := findPoll(t, e, e.ana, pollId); !canClose {
			t.Error("the owner can't close bia's poll")
		}

		otherPoll := createPoll(t, e, e.ana, "a", "b")
		if err := e.s.Poll.ClosePoll(e.bia.ID, otherPoll); !errors.Is(err, ErrPollCloseForbidden) {
			t.Errorf("a member closing another's poll: got %v, want ErrPollCloseForbidden", err)
		}

		if err := e.s.Poll.ClosePoll(e.bia.ID, pollId); err != nil {
			t.Fatal(err)
		}
		optionIds, votes, canVote, canClose, closed := findPoll(t, e, e.ana, pollId)
		if !closed || canVote || canClose || votes[0] == nil {
			t.Errorf("closed poll: closed %v, canVote %v, canClose %v, votes %v; want closed with results", closed, canVote, canClose, votes)
		}
		if err := e.s.Poll.Vote(e.ana.ID, pollId, optionIds[0]); !errors.Is(err, ErrPollClosed) {
			t.Errorf("voting in a closed poll: got %v, want ErrPollClosed", err)
		}
		if err := e.s.Poll.ClosePoll(e.ana.ID, pollId); !errors.Is(err, ErrPollClosed) {
			t.Errorf("closing twice: got %v, want ErrPollClosed", err)
		}
	})

	t.Run("a private community's polls are for members only", func(t *testing.T) {
		e := setup(t, true)
		pollId := createPoll(t, e, e.bia, "sim", "não")
		if _, err := e.s.Poll.FindPolls(e.caio.ID, e.communityId, 1, 10); !errors.Is(err, ErrPrivatePolls) {
			t.Errorf("outsider listing: got %v, want ErrPrivatePolls", err)
		}
		if _, err := e.s.Poll.FindPoll(e.caio.ID, pollId); !errors.Is(err, ErrPrivatePolls) {
			t.Errorf("outsider reading: got %v, want ErrPrivatePolls", err)
		}
	})

	t.Run("anyone reads a public community's polls, but only members vote", func(t *testing.T) {
		e := setup(t, false)
		pollId := createPoll(t, e, e.bia, "sim", "não")
		if _, _, canVote, canClose, _ := findPoll(t, e, e.caio, pollId); canVote || canClose {
			t.Errorf("outsider: canVote %v, canClose %v; want neither", canVote, canClose)
		}
	})
}
