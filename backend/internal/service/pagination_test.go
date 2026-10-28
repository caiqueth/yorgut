package service

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
)

// seenAgo sets the user's last_seen_at to d before now.
func seenAgo(t *testing.T, repos *repository.Repositories, user *model.User, d time.Duration) {
	t.Helper()
	now := time.Now().Add(-d)
	if err := repos.User.TouchLastSeen(user.ID, now, now.Add(time.Second)); err != nil {
		t.Fatalf("TouchLastSeen: %v", err)
	}
}

// createTestUsers creates n users named user0@test.com, user1@test.com, ...
func createTestUsers(t *testing.T, repos *repository.Repositories, n int) []*model.User {
	t.Helper()
	users := make([]*model.User, n)
	for i := range users {
		users[i] = createTestUser(t, repos, fmt.Sprintf("user%d@test.com", i))
	}
	return users
}

func TestFriendsPagination(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	ana := createTestUser(t, repos, "ana@test.com")
	others := createTestUsers(t, repos, 5)

	// Friendships in both directions count.
	for i, u := range others {
		if i%2 == 0 {
			befriend(t, s, repos, ana, u)
		} else {
			befriend(t, s, repos, u, ana)
		}
	}
	// A pending request and a removed friend don't.
	stranger := createTestUser(t, repos, "stranger@test.com")
	if err := s.Friendship.CreateFriendship(stranger.ID, ana.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.Friendship.RemoveFriend(ana.ID, others[4].ID); err != nil {
		t.Fatal(err)
	}

	// others[2] was seen most recently, then others[0]; 1 and 3 never.
	seenAgo(t, repos, others[0], time.Hour)
	seenAgo(t, repos, others[2], time.Minute)

	first, err := s.Friendship.FindFriendsOfUser(ana.ID, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if first.Count != 4 || first.Page != 1 || first.PageSize != 3 {
		t.Fatalf("got count %d, page %d, size %d; want 4, 1, 3", first.Count, first.Page, first.PageSize)
	}
	wantIds := []uint{others[2].ID, others[0].ID}
	if len(first.Friends) != 3 || first.Friends[0].ID != wantIds[0] || first.Friends[1].ID != wantIds[1] {
		t.Fatalf("first page %+v: want 3 friends starting with ids %v", first.Friends, wantIds)
	}

	second, err := s.Friendship.FindFriendsOfUser(ana.ID, 2, 3)
	if err != nil || len(second.Friends) != 1 {
		t.Fatalf("second page: got %+v, err %v, want 1 friend", second, err)
	}

	seen := map[uint]bool{}
	for _, f := range append(first.Friends, second.Friends...) {
		if seen[f.ID] {
			t.Errorf("friend %d on two pages", f.ID)
		}
		seen[f.ID] = true
	}
	if seen[others[4].ID] || seen[stranger.ID] {
		t.Error("removed friend or pending requester listed")
	}
}

func TestCommunityMembersPagination(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	owner := createTestUser(t, repos, "owner@test.com")
	members := createTestUsers(t, repos, 4)

	communityId, err := s.Community.CreateCommunity(context.Background(), owner.ID, "Eu odeio acordar cedo", "other", "", false, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range members {
		if _, err := s.CommunityMembership.JoinCommunity(m.ID, communityId); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.CommunityMembership.LeaveCommunity(members[3].ID, communityId); err != nil {
		t.Fatal(err)
	}
	seenAgo(t, repos, members[1], time.Minute)

	got, err := s.Community.FindCommunityMembers(communityId, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got.Count != 4 || len(got.Members) != 2 || got.Members[0].ID != members[1].ID {
		t.Fatalf("got count %d, members %+v; want 4 and the most recently seen first", got.Count, got.Members)
	}

	last, err := s.Community.FindCommunityMembers(communityId, 2, 2)
	if err != nil || len(last.Members) != 2 {
		t.Fatalf("second page: got %+v, err %v, want 2 members", last, err)
	}
	for _, m := range append(got.Members, last.Members...) {
		if m.ID == members[3].ID {
			t.Error("member who left is listed")
		}
	}

	if _, err := s.Community.FindCommunityMembers(9999, 1, 2); err == nil {
		t.Error("missing community: got no error")
	}
}

func TestUserCommunitiesPagination(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	ana := createTestUser(t, repos, "ana@test.com")
	bia := createTestUser(t, repos, "bia@test.com")

	names := []string{"Sorvete", "Arrasa", "Borracha", "Campainha"}
	for _, name := range names {
		if _, err := s.Community.CreateCommunity(context.Background(), ana.ID, name, "other", "", false, nil); err != nil {
			t.Fatal(err)
		}
	}
	// bia's community: ana isn't a member, and didn't create it.
	if _, err := s.Community.CreateCommunity(context.Background(), bia.ID, "Duplas", "other", "", false, nil); err != nil {
		t.Fatal(err)
	}

	byName := append([]string{}, names...)
	sort.Strings(byName)

	t.Run("memberships by name", func(t *testing.T) {
		got, err := s.CommunityMembership.FindUserMemberships(ana.ID, ana.ID, 2, 3, false)
		if err != nil {
			t.Fatal(err)
		}
		if got.Count != 4 || len(got.Communities) != 1 || got.Communities[0].Name != byName[3] {
			t.Fatalf("got count %d, communities %+v; want 4 and only %q", got.Count, got.Communities, byName[3])
		}
	})

	t.Run("memberships in random order", func(t *testing.T) {
		got, err := s.CommunityMembership.FindUserMemberships(ana.ID, ana.ID, 1, 10, true)
		if err != nil {
			t.Fatal(err)
		}
		// Reloading gives the same order.
		again, err := s.CommunityMembership.FindUserMemberships(ana.ID, ana.ID, 1, 10, true)
		if err != nil || fmt.Sprint(again.Communities) != fmt.Sprint(got.Communities) {
			t.Errorf("reload changed the order: %v, then %v (err %v)", got.Communities, again.Communities, err)
		}
		gotNames := []string{}
		for _, c := range got.Communities {
			gotNames = append(gotNames, c.Name)
		}
		sort.Strings(gotNames)
		if fmt.Sprint(gotNames) != fmt.Sprint(byName) {
			t.Fatalf("got %v, want %v in any order", gotNames, byName)
		}
	})

	t.Run("created communities", func(t *testing.T) {
		got, err := s.Community.FindCommunitiesCreatedByUser(ana.ID, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if got.Count != 4 || len(got.Communities) != 2 || got.Communities[0].Name != byName[0] {
			t.Fatalf("got count %d, communities %+v; want 4, starting with %q", got.Count, got.Communities, byName[0])
		}
	})
}

func TestTouchLastSeen(t *testing.T) {
	repos := newTestRepositories(t)
	user := createTestUser(t, repos, "ana@test.com")

	// MySQL keeps microseconds (DATETIME(6)).
	first := time.Now().Add(-time.Minute).Truncate(time.Microsecond)
	if err := repos.User.TouchLastSeen(user.ID, first, first.Add(-lastSeenEvery)); err != nil {
		t.Fatal(err)
	}
	// Still fresh: a second touch doesn't write.
	if err := repos.User.TouchLastSeen(user.ID, time.Now(), time.Now().Add(-lastSeenEvery)); err != nil {
		t.Fatal(err)
	}
	stored, _ := repos.User.FindById(user.ID)
	if stored.LastSeenAt == nil || !stored.LastSeenAt.Equal(first) {
		t.Fatalf("LastSeenAt = %v, want %v", stored.LastSeenAt, first)
	}
	if !stored.UpdatedAt.Equal(user.UpdatedAt) {
		t.Errorf("UpdatedAt changed from %v to %v", user.UpdatedAt, stored.UpdatedAt)
	}
}

func TestScrapbookPagination(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	ana := createTestUser(t, repos, "ana@test.com")
	bia := createTestUser(t, repos, "bia@test.com")

	for i := 1; i <= 12; i++ {
		if err := s.Scrap.CreateScrap(bia.ID, ana.ID, fmt.Sprintf("recado %d", i)); err != nil {
			t.Fatal(err)
		}
	}
	// A scrap to someone else isn't in ana's scrapbook.
	if err := s.Scrap.CreateScrap(ana.ID, bia.ID, "para a bia"); err != nil {
		t.Fatal(err)
	}

	first, err := s.Scrap.FindScrapOfUserRecipient(ana.ID, ana.ID, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if first.Count != 12 || len(first.Scraps) != 10 || first.Scraps[0].Content != "recado 12" {
		t.Fatalf("first page: count %d, %d scraps starting %q; want 12, 10, recado 12", first.Count, len(first.Scraps), first.Scraps[0].Content)
	}
	second, err := s.Scrap.FindScrapOfUserRecipient(ana.ID, ana.ID, 2, 10)
	if err != nil || len(second.Scraps) != 2 || second.Scraps[1].Content != "recado 1" {
		t.Fatalf("second page: got %+v, err %v; want 2 scraps ending with recado 1", second, err)
	}
}

func TestRelatedCommunities(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	owner := createTestUser(t, repos, "owner@test.com")

	ids := []uint{}
	for i := 0; i < 12; i++ {
		id, err := s.Community.CreateCommunity(context.Background(), owner.ID, fmt.Sprintf("Comunidade %d", i), "other", "", false, nil)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}

	got, err := s.Community.FindRelatedCommunities(ids[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Communities) != relatedCommunitiesShown {
		t.Fatalf("got %d communities, want %d", len(got.Communities), relatedCommunitiesShown)
	}
	again, err := s.Community.FindRelatedCommunities(ids[0])
	if err != nil || fmt.Sprint(again.Communities) != fmt.Sprint(got.Communities) {
		t.Errorf("reload changed the communities: %v, then %v (err %v)", got.Communities, again.Communities, err)
	}
	seen := map[uint]bool{}
	for _, c := range got.Communities {
		if c.ID == ids[0] {
			t.Error("the community itself is listed")
		}
		if seen[c.ID] {
			t.Errorf("community %d listed twice", c.ID)
		}
		seen[c.ID] = true
	}

	if _, err := s.Community.FindRelatedCommunities(9999); err == nil {
		t.Error("missing community: got no error")
	}
}

func TestFanListPagination(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	idol := createTestUser(t, repos, "idol@test.com")
	fans := createTestUsers(t, repos, 5)

	for _, fan := range fans {
		befriend(t, s, repos, fan, idol)
		if err := s.UserFan.BecomeFan(fan.ID, idol.ID); err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Millisecond)
	}
	if err := s.UserFan.StopBeingFan(fans[2].ID, idol.ID); err != nil {
		t.Fatal(err)
	}

	first, err := s.UserFan.FindFansOfUser(idol.ID, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if first.Count != 4 || len(first.Fans) != 3 || first.Fans[0].ID != fans[4].ID {
		t.Fatalf("first page: count %d, fans %+v; want 4, newest (%d) first", first.Count, first.Fans, fans[4].ID)
	}
	second, err := s.UserFan.FindFansOfUser(idol.ID, 2, 3)
	if err != nil || len(second.Fans) != 1 || second.Fans[0].ID != fans[0].ID {
		t.Fatalf("second page: got %+v, err %v; want only the oldest fan", second, err)
	}
	for _, f := range append(first.Fans, second.Fans...) {
		if f.ID == fans[2].ID {
			t.Error("a removed fan is listed")
		}
	}
}

func TestMyFriendsCarryRatings(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	ana := createTestUser(t, repos, "ana@test.com")
	friends := createTestUsers(t, repos, 3)
	for _, f := range friends {
		befriend(t, s, repos, ana, f)
	}
	// Ana rates and stars friends[0], only rates friends[1], ignores friends[2].
	// Someone else's rating of friends[0] mustn't show up as ana's.
	reviews := []model.UserReview{
		{ReviewerID: ana.ID, ReviewedID: friends[0].ID, Cool: 3, Sexy: 2, Trustworthy: 1},
		{ReviewerID: ana.ID, ReviewedID: friends[1].ID, Cool: 1},
		{ReviewerID: friends[2].ID, ReviewedID: friends[0].ID, Cool: 1, Sexy: 1, Trustworthy: 1},
	}
	for i := range reviews {
		if err := repos.UserReview.Upsert(&reviews[i]); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.UserFan.BecomeFan(ana.ID, friends[0].ID); err != nil {
		t.Fatal(err)
	}

	got, err := s.Friendship.FindMyFriends(ana.ID, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if got.Count != 3 || len(got.Friends) != 3 {
		t.Fatalf("got %+v, want 3 friends", got)
	}
	byId := map[uint]transport.MyFriendItem{}
	for _, f := range got.Friends {
		byId[f.ID] = f
	}
	want := map[uint]transport.MyFriendItem{
		friends[0].ID: {MyReview: transport.UserReviewItem{Cool: 3, Sexy: 2, Trustworthy: 1}, AmIFan: true},
		friends[1].ID: {MyReview: transport.UserReviewItem{Cool: 1}},
		friends[2].ID: {},
	}
	for id, w := range want {
		if g := byId[id]; g.MyReview != w.MyReview || g.AmIFan != w.AmIFan {
			t.Errorf("friend %d: review %+v fan %v, want %+v fan %v", id, g.MyReview, g.AmIFan, w.MyReview, w.AmIFan)
		}
	}
}
