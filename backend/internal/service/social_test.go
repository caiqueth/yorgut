package service

import (
	"testing"
	"time"

	"yorgut/internal/model"
	"yorgut/internal/repository"
)

// befriend makes two users friends through the service and returns the
// friendship id.
func befriend(t *testing.T, s *Services, repos *repository.Repositories, requester, requestee *model.User) uint {
	t.Helper()
	if err := s.Friendship.CreateFriendship(requester.ID, requestee.ID); err != nil {
		t.Fatalf("CreateFriendship: %v", err)
	}
	id := friendshipIdBetween(t, repos, requester, requestee)
	if err := s.Friendship.ReplyFriendship(id, true); err != nil {
		t.Fatalf("ReplyFriendship: %v", err)
	}
	return id
}

func friendshipIdBetween(t *testing.T, repos *repository.Repositories, a, b *model.User) uint {
	t.Helper()
	list, err := repos.Friendship.FindFriendshipsBetweenUsers(a.ID, b.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("friendships between users: %v, err %v", list, err)
	}
	return list[0].ID
}

// threeUsers builds the services and three users: ana, bia and caio.
func threeUsers(t *testing.T) (*Services, *repository.Repositories, *model.User, *model.User, *model.User) {
	t.Helper()
	repos := newTestRepositories(t)
	return NewServices(repos, nil), repos,
		createTestUser(t, repos, "ana@test.com"),
		createTestUser(t, repos, "bia@test.com"),
		createTestUser(t, repos, "caio@test.com")
}

func TestFriendshipPermissions(t *testing.T) {
	t.Run("only the requestee of a pending request replies", func(t *testing.T) {
		s, repos, ana, bia, caio := threeUsers(t)
		if err := s.Friendship.CreateFriendship(ana.ID, bia.ID); err != nil {
			t.Fatal(err)
		}
		id := friendshipIdBetween(t, repos, ana, bia)

		cases := []struct {
			name   string
			userId uint
			want   bool
		}{
			{"requestee", bia.ID, true},
			{"requester", ana.ID, false},
			{"third user", caio.ID, false},
		}
		for _, c := range cases {
			if got := s.Friendship.ShouldThisUserBeReplying(id, c.userId); got != c.want {
				t.Errorf("%s: got %v, want %v", c.name, got, c.want)
			}
		}
		if s.Friendship.ShouldThisUserBeReplying(9999, bia.ID) {
			t.Error("missing friendship: got true, want false")
		}
	})

	t.Run("nobody replies to an already accepted friendship", func(t *testing.T) {
		s, repos, ana, bia, _ := threeUsers(t)
		id := befriend(t, s, repos, ana, bia)
		if s.Friendship.ShouldThisUserBeReplying(id, bia.ID) {
			t.Error("requestee of an accepted friendship: got true, want false")
		}
	})

	t.Run("only the two parties delete", func(t *testing.T) {
		s, repos, ana, bia, caio := threeUsers(t)
		id := befriend(t, s, repos, ana, bia)
		cases := []struct {
			name   string
			userId uint
			want   bool
		}{
			{"requester", ana.ID, true},
			{"requestee", bia.ID, true},
			{"third user", caio.ID, false},
		}
		for _, c := range cases {
			if got := s.Friendship.ShouldThisUserBeDeleting(id, c.userId); got != c.want {
				t.Errorf("%s: got %v, want %v", c.name, got, c.want)
			}
		}
		if s.Friendship.ShouldThisUserBeDeleting(9999, ana.ID) {
			t.Error("missing friendship: got true, want false")
		}
	})

	t.Run("RemoveFriend on non-friends is not found", func(t *testing.T) {
		s, _, ana, bia, _ := threeUsers(t)
		wantCategory(t, s.Friendship.RemoveFriend(ana.ID, bia.ID), ErrNotFound)

		// A pending request is not a friendship either.
		if err := s.Friendship.CreateFriendship(ana.ID, bia.ID); err != nil {
			t.Fatal(err)
		}
		wantCategory(t, s.Friendship.RemoveFriend(ana.ID, bia.ID), ErrNotFound)
	})

	t.Run("RemoveFriend works in both directions", func(t *testing.T) {
		s, repos, ana, bia, _ := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		if err := s.Friendship.RemoveFriend(bia.ID, ana.ID); err != nil {
			t.Fatalf("RemoveFriend: %v", err)
		}
		if repos.Friendship.UsersAreFriends(ana.ID, bia.ID) {
			t.Error("still friends after RemoveFriend")
		}
	})
}

func TestUsersFriendshipStatus(t *testing.T) {
	s, repos, ana, bia, caio := threeUsers(t)
	dani := createTestUser(t, repos, "dani@test.com")

	if got := s.Friendship.UsersFriendshipStatus(ana.ID, ana.ID); got != nil {
		t.Errorf("self: got %+v, want nil", got)
	}

	strangers := s.Friendship.UsersFriendshipStatus(ana.ID, bia.ID)
	if strangers == nil || strangers.FriendshipId != nil || strangers.AmIFriend || strangers.IsPendingMyReply || strangers.IsPendingTheirReply {
		t.Errorf("strangers: got %+v, want all false and no id", strangers)
	}

	// Ana asks Bia.
	if err := s.Friendship.CreateFriendship(ana.ID, bia.ID); err != nil {
		t.Fatal(err)
	}
	id := friendshipIdBetween(t, repos, ana, bia)
	asker := s.Friendship.UsersFriendshipStatus(ana.ID, bia.ID)
	if asker == nil || asker.FriendshipId == nil || *asker.FriendshipId != id || asker.AmIFriend || !asker.IsPendingTheirReply || asker.IsPendingMyReply {
		t.Errorf("requester view: got %+v, want pending their reply only", asker)
	}
	asked := s.Friendship.UsersFriendshipStatus(bia.ID, ana.ID)
	if asked == nil || asked.FriendshipId == nil || *asked.FriendshipId != id || asked.AmIFriend || asked.IsPendingTheirReply || !asked.IsPendingMyReply {
		t.Errorf("requestee view: got %+v, want pending my reply only", asked)
	}

	if err := s.Friendship.ReplyFriendship(id, true); err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]uint{{ana.ID, bia.ID}, {bia.ID, ana.ID}} {
		got := s.Friendship.UsersFriendshipStatus(pair[0], pair[1])
		if got == nil || got.FriendshipId == nil || *got.FriendshipId != id || !got.AmIFriend || got.IsPendingMyReply || got.IsPendingTheirReply {
			t.Errorf("friends %v: got %+v, want AmIFriend only", pair, got)
		}
	}

	// A friendship with someone else doesn't leak into the status.
	befriend(t, s, repos, caio, dani)
	if got := s.Friendship.UsersFriendshipStatus(ana.ID, caio.ID); got == nil || got.AmIFriend || got.FriendshipId != nil {
		t.Errorf("unrelated friendships: got %+v, want strangers", got)
	}
}

func TestFans(t *testing.T) {
	// setup makes ana and bia friends; caio is a stranger.
	setup := func(t *testing.T) (*Services, *repository.Repositories, *model.User, *model.User, *model.User) {
		s, repos, ana, bia, caio := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		return s, repos, ana, bia, caio
	}

	t.Run("only friends can become fans", func(t *testing.T) {
		s, _, ana, bia, caio := setup(t)
		if !s.UserFan.CanBecomeFan(ana.ID, bia.ID) || !s.UserFan.CanBecomeFan(bia.ID, ana.ID) {
			t.Error("friends: got false, want true in both directions")
		}
		if s.UserFan.CanBecomeFan(ana.ID, caio.ID) {
			t.Error("stranger: got true, want false")
		}
		if s.UserFan.CanBecomeFan(ana.ID, ana.ID) {
			t.Error("self: got true, want false")
		}
	})

	t.Run("a pending request is not enough", func(t *testing.T) {
		s, _, ana, _, caio := setup(t)
		if err := s.Friendship.CreateFriendship(ana.ID, caio.ID); err != nil {
			t.Fatal(err)
		}
		if s.UserFan.CanBecomeFan(ana.ID, caio.ID) {
			t.Error("pending request: got true, want false")
		}
	})

	t.Run("becoming a fan is idempotent and can be undone", func(t *testing.T) {
		s, _, ana, bia, _ := setup(t)
		for i := 1; i <= 2; i++ {
			if err := s.UserFan.BecomeFan(ana.ID, bia.ID); err != nil {
				t.Fatalf("BecomeFan #%d: %v", i, err)
			}
		}
		if n, err := s.UserFan.CountFansOfUser(bia.ID); err != nil || n != 1 {
			t.Fatalf("fan count = %d, err %v, want 1", n, err)
		}

		if err := s.UserFan.StopBeingFan(ana.ID, bia.ID); err != nil {
			t.Fatalf("StopBeingFan: %v", err)
		}
		if n, _ := s.UserFan.CountFansOfUser(bia.ID); n != 0 {
			t.Fatalf("fan count after stopping = %d, want 0", n)
		}
		// Stopping again is harmless, and becoming a fan again works.
		if err := s.UserFan.StopBeingFan(ana.ID, bia.ID); err != nil {
			t.Fatalf("second StopBeingFan: %v", err)
		}
		if err := s.UserFan.BecomeFan(ana.ID, bia.ID); err != nil {
			t.Fatalf("BecomeFan again: %v", err)
		}
		if n, _ := s.UserFan.CountFansOfUser(bia.ID); n != 1 {
			t.Fatalf("fan count after returning = %d, want 1", n)
		}
	})

	t.Run("summary", func(t *testing.T) {
		s, repos, ana, bia, caio := setup(t)
		befriend(t, s, repos, caio, bia)
		if err := s.UserFan.BecomeFan(ana.ID, bia.ID); err != nil {
			t.Fatal(err)
		}
		dani := createTestUser(t, repos, "dani@test.com")

		cases := []struct {
			name     string
			viewer   uint
			amIFan   bool
			canFan   bool
			fanCount int64
		}{
			{"fan who is a friend", ana.ID, true, true, 1},
			{"friend who is not a fan", caio.ID, false, true, 1},
			{"the idol", bia.ID, false, false, 1},
			{"stranger", dani.ID, false, false, 1},
		}
		for _, c := range cases {
			got, err := s.UserFan.FindSummaryOfUser(c.viewer, bia.ID)
			if err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			if got.FanCount != c.fanCount || got.AmIFan != c.amIFan || got.CanFan != c.canFan {
				t.Errorf("%s: got %+v, want count %d amIFan %v canFan %v", c.name, got, c.fanCount, c.amIFan, c.canFan)
			}
		}
	})
}

func TestReviews(t *testing.T) {
	// setup makes ana and caio friends of bia.
	setup := func(t *testing.T) (*Services, *repository.Repositories, *model.User, *model.User, *model.User) {
		s, repos, ana, bia, caio := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		befriend(t, s, repos, caio, bia)
		return s, repos, ana, bia, caio
	}

	t.Run("only friends can review", func(t *testing.T) {
		s, repos, ana, bia, _ := setup(t)
		stranger := createTestUser(t, repos, "dani@test.com")
		if !s.UserReview.CanReview(ana.ID, bia.ID) {
			t.Error("friend: got false, want true")
		}
		if s.UserReview.CanReview(stranger.ID, bia.ID) {
			t.Error("stranger: got true, want false")
		}
		if s.UserReview.CanReview(ana.ID, ana.ID) {
			t.Error("self: got true, want false")
		}
	})

	t.Run("upsert updates instead of duplicating", func(t *testing.T) {
		s, _, ana, bia, _ := setup(t)
		if err := s.UserReview.UpsertReview(ana.ID, bia.ID, 1, 2, 3); err != nil {
			t.Fatal(err)
		}
		if err := s.UserReview.UpsertReview(ana.ID, bia.ID, 3, 0, 1); err != nil {
			t.Fatal(err)
		}
		got, err := s.UserReview.FindSummaryOfUser(ana.ID, bia.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ReviewCount != 1 {
			t.Fatalf("review count = %d, want 1", got.ReviewCount)
		}
		if got.MyReview == nil || got.MyReview.Cool != 3 || got.MyReview.Sexy != 0 || got.MyReview.Trustworthy != 1 {
			t.Fatalf("my review = %+v, want 3/0/1", got.MyReview)
		}
		if got.Cool != 3 || got.Sexy != 0 || got.Trustworthy != 1 {
			t.Errorf("averages = %v/%v/%v, want 3/0/1", got.Cool, got.Sexy, got.Trustworthy)
		}
	})

	t.Run("averages and myReview", func(t *testing.T) {
		s, repos, ana, bia, caio := setup(t)
		if err := s.UserReview.UpsertReview(ana.ID, bia.ID, 3, 2, 1); err != nil {
			t.Fatal(err)
		}
		if err := s.UserReview.UpsertReview(caio.ID, bia.ID, 1, 0, 2); err != nil {
			t.Fatal(err)
		}

		got, err := s.UserReview.FindSummaryOfUser(ana.ID, bia.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ReviewCount != 2 || got.Cool != 2 || got.Sexy != 1 || got.Trustworthy != 1.5 {
			t.Errorf("got %+v, want count 2 and averages 2/1/1.5", got)
		}
		if !got.CanReview || got.MyReview == nil || got.MyReview.Cool != 3 || got.MyReview.Sexy != 2 || got.MyReview.Trustworthy != 1 {
			t.Errorf("reviewer view: got %+v, want canReview and own review 3/2/1", got)
		}

		// A viewer who hasn't reviewed sees nil, but still the averages.
		stranger := createTestUser(t, repos, "dani@test.com")
		got, err = s.UserReview.FindSummaryOfUser(stranger.ID, bia.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.MyReview != nil || got.CanReview || got.ReviewCount != 2 {
			t.Errorf("stranger view: got %+v, want nil myReview, canReview false, count 2", got)
		}
	})

	t.Run("no reviews", func(t *testing.T) {
		s, _, ana, bia, _ := setup(t)
		got, err := s.UserReview.FindSummaryOfUser(ana.ID, bia.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ReviewCount != 0 || got.Cool != 0 || got.Sexy != 0 || got.Trustworthy != 0 || got.MyReview != nil {
			t.Errorf("got %+v, want zeros and nil myReview", got)
		}
	})
}

func TestTestimonials(t *testing.T) {
	// setup makes ana and bia friends; caio is a stranger.
	setup := func(t *testing.T) (*Services, *repository.Repositories, *model.User, *model.User, *model.User) {
		s, repos, ana, bia, caio := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		return s, repos, ana, bia, caio
	}

	t.Run("DeleteTestimonial permissions", func(t *testing.T) {
		cases := []struct {
			name      string
			viewer    func(ana, bia, caio *model.User) uint
			wantError error
		}{
			{"author", func(ana, _, _ *model.User) uint { return ana.ID }, nil},
			{"recipient", func(_, bia, _ *model.User) uint { return bia.ID }, nil},
			{"third user", func(_, _, caio *model.User) uint { return caio.ID }, ErrForbidden},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				s, _, ana, bia, caio := setup(t)
				if err := s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "legal"); err != nil {
					t.Fatal(err)
				}
				pending, err := s.Testimonial.FindTestimonialsOfUser(bia.ID, bia.ID, 0, 0)
				if err != nil || len(pending.Pending) != 1 {
					t.Fatalf("pending testimonials: %+v, err %v", pending, err)
				}

				err = s.Testimonial.DeleteTestimonial(c.viewer(ana, bia, caio), pending.Pending[0].ID)
				wantLeft := 0
				if c.wantError != nil {
					wantCategory(t, err, c.wantError)
					wantLeft = 1
				} else if err != nil {
					t.Fatalf("DeleteTestimonial: %v", err)
				}

				after, _ := s.Testimonial.FindTestimonialsOfUser(bia.ID, bia.ID, 0, 0)
				if len(after.Pending) != wantLeft {
					t.Errorf("testimonials left = %d, want %d", len(after.Pending), wantLeft)
				}
			})
		}

		s, _, ana, _, _ := setup(t)
		wantCategory(t, s.Testimonial.DeleteTestimonial(ana.ID, 9999), ErrNotFound)
	})

	t.Run("approval is recipient-only", func(t *testing.T) {
		s, _, ana, bia, _ := setup(t)
		if err := s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "legal"); err != nil {
			t.Fatal(err)
		}
		pending, _ := s.Testimonial.FindTestimonialsOfUser(bia.ID, bia.ID, 0, 0)
		id := pending.Pending[0].ID
		wantCategory(t, s.Testimonial.ApproveTestimonial(ana.ID, id), ErrForbidden)
		if err := s.Testimonial.ApproveTestimonial(bia.ID, id); err != nil {
			t.Fatalf("ApproveTestimonial: %v", err)
		}
	})

	t.Run("only friends write", func(t *testing.T) {
		s, _, ana, bia, caio := setup(t)
		wantCategory(t, s.Testimonial.CreateTestimonial(caio.ID, bia.ID, "oi"), ErrForbidden)
		wantCategory(t, s.Testimonial.CreateTestimonial(ana.ID, ana.ID, "eu"), ErrForbidden)
	})

	t.Run("pending is private, approved is public", func(t *testing.T) {
		s, repos, ana, bia, caio := setup(t)
		befriend(t, s, repos, caio, bia)
		if err := s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "pendente"); err != nil {
			t.Fatal(err)
		}
		if err := s.Testimonial.CreateTestimonial(caio.ID, bia.ID, "aprovado"); err != nil {
			t.Fatal(err)
		}
		all, _ := s.Testimonial.FindTestimonialsOfUser(bia.ID, bia.ID, 0, 0)
		for _, item := range all.Pending {
			if item.Content == "aprovado" {
				if err := s.Testimonial.ApproveTestimonial(bia.ID, item.ID); err != nil {
					t.Fatal(err)
				}
			}
		}

		owner, err := s.Testimonial.FindTestimonialsOfUser(bia.ID, bia.ID, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(owner.Testimonials) != 1 || owner.Testimonials[0].Content != "aprovado" || !owner.Testimonials[0].Approved {
			t.Errorf("owner approved = %+v, want only the approved one", owner.Testimonials)
		}
		if len(owner.Pending) != 1 || owner.Pending[0].Content != "pendente" {
			t.Errorf("owner pending = %+v, want only the pending one", owner.Pending)
		}
		if owner.CanWrite {
			t.Error("owner CanWrite: got true, want false")
		}

		visitor, err := s.Testimonial.FindTestimonialsOfUser(ana.ID, bia.ID, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(visitor.Testimonials) != 1 || len(visitor.Pending) != 0 || !visitor.CanWrite {
			t.Errorf("friend view: approved %d, pending %d, canWrite %v; want 1, 0, true",
				len(visitor.Testimonials), len(visitor.Pending), visitor.CanWrite)
		}
	})
}

func TestTopics(t *testing.T) {
	type env struct {
		s                       *Services
		repos                   *repository.Repositories
		owner, member, outsider *model.User
		communityId             uint
	}
	setup := func(t *testing.T) env {
		s, repos, owner, member, outsider := threeUsers(t)
		id, err := s.Community.CreateCommunity(owner.ID, "Comunidade", "Outros", nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := s.CommunityMembership.JoinCommunity(member.ID, id); err != nil {
			t.Fatal(err)
		}
		return env{s, repos, owner, member, outsider, id}
	}
	// onlyTopicId returns the id of the community's only top-level topic.
	onlyTopicId := func(t *testing.T, e env) uint {
		t.Helper()
		topics, err := e.s.Topic.FindTopics(e.communityId, 10)
		if err != nil || len(topics.Topics) != 1 {
			t.Fatalf("topics: %v, err %v", topics, err)
		}
		return topics.Topics[0].ID
	}

	t.Run("reply by a member belongs to the topic's community", func(t *testing.T) {
		e := setup(t)
		if err := e.s.Topic.CreateTopic(e.owner.ID, e.communityId, "titulo", "texto"); err != nil {
			t.Fatal(err)
		}
		topicId := onlyTopicId(t, e)

		if err := e.s.Topic.ReplyTopic(e.member.ID, topicId, "", "resposta"); err != nil {
			t.Fatalf("ReplyTopic: %v", err)
		}
		count, replies, err := e.repos.Topic.FindRepliesOfTopic(topicId)
		if err != nil || count != 1 || len(replies) != 1 {
			t.Fatalf("replies: count %d, %v, err %v", count, replies, err)
		}
		r := replies[0]
		if r.CommunityID != e.communityId || r.AuthorID != e.member.ID || r.ParentTopicID == nil || *r.ParentTopicID != topicId || r.Text != "resposta" {
			t.Errorf("reply = %+v, want community %d, author %d, parent %d", r, e.communityId, e.member.ID, topicId)
		}
		// Replies don't show up as top-level topics.
		onlyTopicId(t, e)
	})

	t.Run("non-members and missing topics", func(t *testing.T) {
		e := setup(t)
		wantCategory(t, e.s.Topic.CreateTopic(e.outsider.ID, e.communityId, "t", "x"), ErrForbidden)
		if err := e.s.Topic.CreateTopic(e.owner.ID, e.communityId, "titulo", "texto"); err != nil {
			t.Fatal(err)
		}
		topicId := onlyTopicId(t, e)
		wantCategory(t, e.s.Topic.ReplyTopic(e.outsider.ID, topicId, "", "x"), ErrForbidden)
		wantCategory(t, e.s.Topic.ReplyTopic(e.member.ID, 9999, "", "x"), ErrNotFound)
		if count, _, _ := e.repos.Topic.FindRepliesOfTopic(topicId); count != 0 {
			t.Errorf("replies = %d, want 0", count)
		}
	})

	t.Run("a member of another community can't reply", func(t *testing.T) {
		e := setup(t)
		if _, err := e.s.Community.CreateCommunity(e.outsider.ID, "Outra", "Outros", nil); err != nil {
			t.Fatal(err)
		}
		if err := e.s.Topic.CreateTopic(e.owner.ID, e.communityId, "titulo", "texto"); err != nil {
			t.Fatal(err)
		}
		wantCategory(t, e.s.Topic.ReplyTopic(e.outsider.ID, onlyTopicId(t, e), "", "x"), ErrForbidden)
	})

	t.Run("most recent topics first, limited", func(t *testing.T) {
		e := setup(t)
		for _, title := range []string{"primeiro", "segundo", "terceiro"} {
			if err := e.s.Topic.CreateTopic(e.owner.ID, e.communityId, title, "texto"); err != nil {
				t.Fatal(err)
			}
			time.Sleep(5 * time.Millisecond)
		}
		resp, err := e.s.Topic.FindTopics(e.communityId, 2)
		if err != nil {
			t.Fatal(err)
		}
		topics := resp.Topics
		if len(topics) != 2 || topics[0].Title != "terceiro" || topics[1].Title != "segundo" {
			t.Errorf("got %+v, want terceiro then segundo", topics)
		}
	})
}

func TestProfileUpdates(t *testing.T) {
	t.Run("blank fields are stored as NULL", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		user := createTestUser(t, repos, "ana@test.com")

		if err := s.UserProfile.UpdateProfile(user.ID, UpdateUserProfileInput{Description: "  oi  ", Music: "rock", Hobbies: "xadrez"}); err != nil {
			t.Fatal(err)
		}
		loaded, err := repos.User.FindUserWithProfile(user.ID)
		if err != nil || loaded.Profile == nil {
			t.Fatalf("profile: %+v, err %v", loaded, err)
		}
		p := loaded.Profile
		if p.Description == nil || *p.Description != "oi" || p.Music == nil || *p.Music != "rock" {
			t.Fatalf("stored profile = %+v, want trimmed description and music", p)
		}
		if p.Books != nil || p.Website != nil {
			t.Errorf("untouched fields should be NULL, got books %v, website %v", p.Books, p.Website)
		}

		// Clearing a field (or sending only whitespace) stores NULL, not "".
		if err := s.UserProfile.UpdateProfile(user.ID, UpdateUserProfileInput{Description: "", Music: "   ", Hobbies: "xadrez"}); err != nil {
			t.Fatal(err)
		}
		loaded, _ = repos.User.FindUserWithProfile(user.ID)
		p = loaded.Profile
		if p.Description != nil || p.Music != nil {
			t.Errorf("cleared fields: description %v, music %v, want nil", p.Description, p.Music)
		}
		if p.Hobbies == nil || *p.Hobbies != "xadrez" {
			t.Errorf("hobbies = %v, want xadrez", p.Hobbies)
		}
	})

	t.Run("UpdateUser keeps the birth date when nil", func(t *testing.T) {
		repos := newTestRepositories(t)
		s := NewServices(repos, nil)
		birth := time.Date(1990, 5, 17, 0, 0, 0, 0, time.UTC)
		if err := s.User.CreateUser("Ana", "Silva", "ana@test.com", birth, nil); err != nil {
			t.Fatal(err)
		}
		user, err := repos.User.FindByEmail("ana@test.com")
		if err != nil {
			t.Fatal(err)
		}

		if err := s.User.UpdateUser(user.ID, UpdateUserInput{FirstName: " Ana Maria ", LastName: " Souza "}); err != nil {
			t.Fatal(err)
		}
		got, _ := repos.User.FindById(user.ID)
		if got.FirstName != "Ana Maria" || got.LastName != "Souza" {
			t.Errorf("names = %q %q, want trimmed update", got.FirstName, got.LastName)
		}
		if !got.BirthDate.Equal(birth) {
			t.Errorf("birth date = %v, want %v kept", got.BirthDate, birth)
		}

		newBirth := time.Date(1985, 1, 2, 0, 0, 0, 0, time.UTC)
		if err := s.User.UpdateUser(user.ID, UpdateUserInput{FirstName: "Ana", BirthDate: &newBirth}); err != nil {
			t.Fatal(err)
		}
		got, _ = repos.User.FindById(user.ID)
		if !got.BirthDate.Equal(newBirth) {
			t.Errorf("birth date = %v, want %v", got.BirthDate, newBirth)
		}
		if got.LastName != "" {
			t.Errorf("last name = %q, want cleared", got.LastName)
		}
	})
}
