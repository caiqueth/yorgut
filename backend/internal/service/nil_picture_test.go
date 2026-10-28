package service

import (
	"context"
	"testing"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
)

// TestNilProfilePictureDoesNotPanic guards the nil-picture bug: response
// mappers used to dereference User.Profile.ProfilePicture (and the community
// equivalent), panicking for users without a picture or without a profile row
// and turning list endpoints into 500s. Every response that carries a picture
// must now map it to nil instead.
func TestNilProfilePictureDoesNotPanic(t *testing.T) {
	scenarios := []struct {
		name         string
		dropProfiles bool
	}{
		{"profile row without a picture", false},
		{"no profile row at all", true},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			db := openTestDB(t)
			repos := repository.NewRepositories(db)
			s := NewServices(repos, storage.NewMemory())
			ana := createTestUser(t, repos, "ana@test.com")
			bia := createTestUser(t, repos, "bia@test.com")
			caio := createTestUser(t, repos, "caio@test.com")

			befriend(t, s, repos, ana, bia)
			if err := s.Friendship.CreateFriendship(caio.ID, bia.ID); err != nil {
				t.Fatal(err)
			}
			if err := s.Scrap.CreateScrap(ana.ID, bia.ID, "oi"); err != nil {
				t.Fatal(err)
			}
			if err := s.Testimonial.CreateTestimonial(ana.ID, bia.ID, "legal"); err != nil {
				t.Fatal(err)
			}
			communityId, err := s.Community.CreateCommunity(context.Background(), ana.ID, "Comunidade", "other", "", false, nil)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := s.CommunityMembership.JoinCommunity(bia.ID, communityId); err != nil {
				t.Fatal(err)
			}
			if err := s.Topic.CreateTopic(ana.ID, communityId, "titulo", "texto"); err != nil {
				t.Fatal(err)
			}
			topics, err := s.Topic.FindTopics(bia.ID, communityId, 1, 10)
			if err != nil || len(topics.Topics) != 1 {
				t.Fatalf("topics: %+v, err %v", topics, err)
			}
			topicId := topics.Topics[0].ID
			if err := s.Topic.ReplyTopic(bia.ID, topicId, "", "resposta"); err != nil {
				t.Fatal(err)
			}

			if sc.dropProfiles {
				for _, table := range []any{&model.UserProfile{}, &model.CommunityProfile{}} {
					if err := db.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
						t.Fatalf("dropping profile rows: %v", err)
					}
				}
			}

			wantNil := func(t *testing.T, what string, got *string) {
				t.Helper()
				if got != nil {
					t.Errorf("%s = %q, want nil", what, *got)
				}
			}
			// check runs f and fails the test, instead of crashing the whole
			// run, if it panics.
			check := func(name string, f func(t *testing.T)) {
				t.Run(name, func(t *testing.T) {
					defer func() {
						if r := recover(); r != nil {
							t.Fatalf("panic: %v", r)
						}
					}()
					f(t)
				})
			}

			check("friends list", func(t *testing.T) {
				got, err := s.Friendship.FindFriendsOfUser(bia.ID, 1, 10)
				if err != nil || len(got.Friends) != 1 {
					t.Fatalf("got %+v, err %v, want 1 friend", got, err)
				}
				wantNil(t, "friend picture", got.Friends[0].ProfilePicture)
			})
			check("pending friendships", func(t *testing.T) {
				got, err := s.Friendship.FindPendingFriendshipsOfUser(bia.ID)
				if err != nil || len(got.PendingFriendships) != 1 {
					t.Fatalf("got %+v, err %v, want 1 pending", got, err)
				}
				wantNil(t, "requester picture", got.PendingFriendships[0].ProfilePicture)
			})
			check("scraps", func(t *testing.T) {
				got, err := s.Scrap.FindScrapOfUserRecipient(bia.ID, bia.ID, 1, 10)
				if err != nil || len(got.Scraps) != 1 {
					t.Fatalf("got %+v, err %v, want 1 scrap", got, err)
				}
				wantNil(t, "scrap author picture", got.Scraps[0].AuthorProfilePicture)
			})
			check("testimonials", func(t *testing.T) {
				got, err := s.Testimonial.FindTestimonialsOfUser(bia.ID, bia.ID, 0, 0)
				if err != nil || len(got.Pending) != 1 {
					t.Fatalf("got %+v, err %v, want 1 pending testimonial", got, err)
				}
				wantNil(t, "testimonial author picture", got.Pending[0].AuthorProfilePicture)
			})
			check("community members", func(t *testing.T) {
				got, err := s.Community.FindCommunityMembers(communityId, 1, 10)
				if err != nil || len(got.Members) != 2 {
					t.Fatalf("got %+v, err %v, want 2 members", got, err)
				}
				for _, m := range got.Members {
					wantNil(t, "member picture", m.ProfilePicture)
				}
			})
			check("topics", func(t *testing.T) {
				got, err := s.Topic.FindTopics(bia.ID, communityId, 1, 10)
				if err != nil || len(got.Topics) != 1 {
					t.Fatalf("got %+v, err %v, want 1 topic", got, err)
				}
				wantNil(t, "topic author picture", got.Topics[0].AuthorProfilePicture)
			})
			check("topic and replies", func(t *testing.T) {
				got, err := s.Topic.FindTopicReplies(bia.ID, topicId, 1, 10)
				if err != nil || got.Topic == nil || len(got.Replies) != 1 {
					t.Fatalf("got %+v, err %v, want topic and 1 reply", got, err)
				}
				wantNil(t, "topic author picture", got.Topic.AuthorProfilePicture)
				wantNil(t, "reply author picture", got.Replies[0].AuthorProfilePicture)
			})
			check("user profile", func(t *testing.T) {
				got, err := s.User.FindUserWithProfileById(ana.ID, ana.ID)
				if err != nil {
					t.Fatal(err)
				}
				wantNil(t, "profile picture", got.ProfilePicture)
			})
			check("user search", func(t *testing.T) {
				got, err := s.User.SearchUsers("Test", 1, 10)
				if err != nil {
					t.Fatal(err)
				}
				if len(got.Users) != 3 || got.Count != 3 {
					t.Fatalf("got %d users, count %d; want 3 and 3", len(got.Users), got.Count)
				}
				for _, u := range got.Users {
					wantNil(t, "search result picture", u.ProfilePicture)
				}
			})
			check("user memberships", func(t *testing.T) {
				got, err := s.CommunityMembership.FindUserMemberships(bia.ID, bia.ID, 1, 10, false)
				if err != nil || len(got.Communities) != 1 {
					t.Fatalf("got %+v, err %v, want 1 membership", got, err)
				}
				wantNil(t, "community picture", got.Communities[0].ProfilePicture)
			})
			check("owned communities", func(t *testing.T) {
				got, err := s.Community.FindCommunitiesCreatedByUser(ana.ID, 1, 10)
				if err != nil || len(got.Communities) != 1 {
					t.Fatalf("got %+v, err %v, want 1 community", got, err)
				}
				wantNil(t, "community picture", got.Communities[0].ProfilePicture)
			})
			check("community discovery", func(t *testing.T) {
				got, err := s.Community.SearchCommunitiesToDiscover(caio.ID, "", 1, 10)
				if err != nil || len(got.Communities) != 1 {
					t.Fatalf("got %+v, err %v, want 1 community", got, err)
				}
				wantNil(t, "community picture", got.Communities[0].ProfilePicture)
			})
		})
	}
}
