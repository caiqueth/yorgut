package service

import (
	"context"
	"testing"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/storage"
)

func TestProfileVisits(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	users := createTestUsers(t, repos, 4)
	ana, visitors := users[0], users[1:]

	// Own visits and repeated visits on one day count once at most.
	s.ProfileVisit.RecordVisit(ana.ID, ana.ID)
	for _, v := range visitors {
		s.ProfileVisit.RecordVisit(v.ID, ana.ID)
		s.ProfileVisit.RecordVisit(v.ID, ana.ID)
	}
	// visitors[0] also came 3 and 10 days ago.
	now := time.Now()
	for _, daysAgo := range []int{3, 10} {
		at := now.AddDate(0, 0, -daysAgo)
		if err := repos.ProfileVisit.Record(visitors[0].ID, ana.ID, dayOf(at), at); err != nil {
			t.Fatal(err)
		}
	}

	views, err := s.ProfileVisit.FindProfileViews(ana.ID, now)
	if err != nil {
		t.Fatal(err)
	}
	if views.AllTime != 5 || views.LastSevenDays != 4 {
		t.Errorf("views: all time %d, last 7 days %d; want 5 and 4", views.AllTime, views.LastSevenDays)
	}

	// visitors[0] comes back: they become the most recent.
	time.Sleep(10 * time.Millisecond)
	s.ProfileVisit.RecordVisit(visitors[0].ID, ana.ID)
	recent, err := s.ProfileVisit.FindRecentVisitors(ana.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 3 || recent[0].UserID != visitors[0].ID {
		t.Fatalf("recent visitors %+v: want 3, visitor %d first", recent, visitors[0].ID)
	}
	for _, r := range recent {
		if r.UserID == ana.ID {
			t.Error("ana listed as her own visitor")
		}
	}
}

func TestIndirectConnections(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	users := createTestUsers(t, repos, 6)
	ana, bia, caio, duda, edu, fabi := users[0], users[1], users[2], users[3], users[4], users[5]

	befriend(t, s, repos, ana, bia)
	befriend(t, s, repos, duda, ana)
	befriend(t, s, repos, bia, caio) // caio: friend of a friend
	befriend(t, s, repos, edu, bia)  // edu: too, in the other direction
	befriend(t, s, repos, bia, duda) // duda is ana's friend already
	befriend(t, s, repos, duda, caio)
	// A pending request doesn't connect.
	if err := s.Friendship.CreateFriendship(duda.ID, fabi.ID); err != nil {
		t.Fatal(err)
	}

	got, err := s.Friendship.CountIndirectConnections(ana.ID)
	if err != nil || got != 2 {
		t.Fatalf("got %d, err %v; want 2 (caio and edu)", got, err)
	}

	lonely := createTestUser(t, repos, "lonely@test.com")
	if got, err := s.Friendship.CountIndirectConnections(lonely.ID); err != nil || got != 0 {
		t.Errorf("user without friends: got %d, err %v; want 0", got, err)
	}
}

func TestNextBirthday(t *testing.T) {
	at := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 15, 30, 0, 0, time.Local)
	}
	born := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	}
	cases := []struct {
		name     string
		birth    time.Time
		now      time.Time
		wantNext time.Time
		wantDays int
	}{
		{"today", born(1990, time.March, 5), at(2026, time.March, 5), born(2026, time.March, 5), 0},
		{"later this year", born(1990, time.March, 15), at(2026, time.March, 5), born(2026, time.March, 15), 10},
		{"already passed", born(1990, time.March, 4), at(2026, time.March, 5), born(2027, time.March, 4), 364},
		{"across new year", born(1990, time.January, 2), at(2026, time.December, 30), born(2027, time.January, 2), 3},
		{"29 February, other year", born(2000, time.February, 29), at(2026, time.February, 20), born(2026, time.March, 1), 9},
	}
	for _, c := range cases {
		next, days := nextBirthday(c.birth, c.now)
		if !next.Equal(c.wantNext) || days != c.wantDays {
			t.Errorf("%s: got %v in %d days, want %v in %d", c.name, next, days, c.wantNext, c.wantDays)
		}
	}
}

func TestUpcomingBirthdays(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	ana := createTestUser(t, repos, "ana@test.com")

	now := time.Now()
	today := dayOf(now)
	friend := func(email string, firstName string, birthDate time.Time) *model.User {
		var birthDatePtr *time.Time
		if !birthDate.IsZero() {
			birthDatePtr = &birthDate
		}
		user := model.NewUser(firstName, "Test", email, birthDatePtr, nil)
		if err := repos.User.Create(user); err != nil {
			t.Fatal(err)
		}
		befriend(t, s, repos, ana, user)
		return user
	}
	inTen := friend("ten@test.com", "Dez", today.AddDate(-30, 0, 10))
	todays := friend("today@test.com", "Hoje", today.AddDate(-25, 0, 0))
	friend("far@test.com", "Longe", today.AddDate(-20, 0, 40))
	friend("unset@test.com", "Sem data", time.Time{})
	// Not a friend: not listed.
	strangerBirth := today.AddDate(-20, 0, 1)
	stranger := model.NewUser("Estranho", "Test", "stranger@test.com", &strangerBirth, nil)
	if err := repos.User.Create(stranger); err != nil {
		t.Fatal(err)
	}

	got, err := s.Friendship.FindUpcomingBirthdays(ana.ID, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Birthdays) != 2 {
		t.Fatalf("got %+v, want 2 birthdays", got.Birthdays)
	}
	if got.Birthdays[0].ID != todays.ID || got.Birthdays[0].DaysUntil != 0 {
		t.Errorf("first %+v: want today's birthday", got.Birthdays[0])
	}
	if got.Birthdays[1].ID != inTen.ID || got.Birthdays[1].DaysUntil != 10 {
		t.Errorf("second %+v: want the one in 10 days", got.Birthdays[1])
	}
}

func TestTodaysLuck(t *testing.T) {
	morning := time.Date(2026, time.October, 10, 8, 0, 0, 0, time.Local)
	evening := time.Date(2026, time.October, 10, 22, 0, 0, 0, time.Local)
	if TodaysLuck(1, morning) != TodaysLuck(1, evening) {
		t.Error("luck changed during the day")
	}

	// Over a month, a user should see more than one fortune.
	seen := map[string]bool{}
	for d := 0; d < 30; d++ {
		seen[TodaysLuck(1, morning.AddDate(0, 0, d))] = true
	}
	if len(seen) < 2 {
		t.Errorf("one fortune for 30 days: %v", seen)
	}
}

func TestDailySeed(t *testing.T) {
	morning := time.Date(2026, time.October, 10, 8, 0, 0, 0, time.Local)
	evening := time.Date(2026, time.October, 10, 22, 0, 0, 0, time.Local)
	if dailySeed(1, morning) != dailySeed(1, evening) {
		t.Error("seed changed during the day")
	}
	if dailySeed(1, morning) == dailySeed(1, morning.AddDate(0, 0, 1)) {
		t.Error("same seed on the next day")
	}
	if dailySeed(1, morning) == dailySeed(2, morning) {
		t.Error("same seed for two ids")
	}
}

func TestNextStepsData(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	users := createTestUsers(t, repos, 3)
	ana, bia, caio := users[0], users[1], users[2]

	rated := func() bool {
		t.Helper()
		got, err := s.Friendship.HasRatedFriends(ana.ID)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}

	befriend(t, s, repos, ana, bia)
	if rated() {
		t.Error("no ratings yet: got true")
	}
	if err := s.UserReview.UpsertReview(ana.ID, bia.ID, 2, 0, 0); err != nil {
		t.Fatal(err)
	}
	if !rated() {
		t.Error("rated a friend: got false")
	}
	// A rating of someone who is no longer a friend doesn't count.
	if err := s.Friendship.RemoveFriend(ana.ID, bia.ID); err != nil {
		t.Fatal(err)
	}
	befriend(t, s, repos, ana, caio)
	if rated() {
		t.Error("only rated a former friend: got true")
	}

	if got, err := s.CommunityMembership.CountUserMemberships(ana.ID); err != nil || got != 0 {
		t.Errorf("communities before: %d, err %v; want 0", got, err)
	}
	if _, err := s.Community.CreateCommunity(context.Background(), ana.ID, "Comunidade", "other", "", false, nil); err != nil {
		t.Fatal(err)
	}
	if got, err := s.CommunityMembership.CountUserMemberships(ana.ID); err != nil || got != 1 {
		t.Errorf("communities after creating one: %d, err %v; want 1", got, err)
	}
}

func TestRecordVisitsSetting(t *testing.T) {
	repos := newTestRepositories(t)
	s := NewServices(repos, storage.NewMemory())
	users := createTestUsers(t, repos, 3)
	ana, bia, caio := users[0], users[1], users[2]

	if on, err := repos.User.RecordsVisits(ana.ID); err != nil || !on {
		t.Fatalf("default: %v, err %v; want on", on, err)
	}

	// Bia visits ana while recording, then turns it off.
	s.ProfileVisit.RecordVisit(bia.ID, ana.ID)
	if err := setRecordVisits(s, bia.ID, false); err != nil {
		t.Fatal(err)
	}
	// Caio visits ana; bia visits caio, which isn't recorded.
	s.ProfileVisit.RecordVisit(caio.ID, ana.ID)
	s.ProfileVisit.RecordVisit(bia.ID, caio.ID)

	visitors, err := s.ProfileVisit.FindRecentVisitors(ana.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(visitors) != 1 || visitors[0].UserID != caio.ID {
		t.Errorf("ana's visitors %+v: want only caio (bia turned visits off)", visitors)
	}
	if views, err := s.ProfileVisit.FindProfileViews(ana.ID, time.Now()); err != nil || views.AllTime != 2 {
		t.Errorf("ana's views %+v, err %v; want 2 (counts keep earlier visits)", views, err)
	}
	if views, err := s.ProfileVisit.FindProfileViews(caio.ID, time.Now()); err != nil || views.AllTime != 0 {
		t.Errorf("caio's views %+v, err %v; want 0 (bia's visit wasn't recorded)", views, err)
	}

	// Bia doesn't see her own visitors while it's off.
	s.ProfileVisit.RecordVisit(caio.ID, bia.ID)
	if got, err := s.ProfileVisit.FindRecentVisitors(bia.ID); err != nil || len(got) != 0 {
		t.Errorf("bia's visitors while off: %+v, err %v; want none", got, err)
	}
	if err := setRecordVisits(s, bia.ID, true); err != nil {
		t.Fatal(err)
	}
	if got, err := s.ProfileVisit.FindRecentVisitors(bia.ID); err != nil || len(got) != 1 {
		t.Errorf("bia's visitors after turning it on: %+v, err %v; want caio", got, err)
	}
}

// setRecordVisits changes only "registrar minhas visitas", leaving the rest
// of the privacy settings at their defaults.
func setRecordVisits(s *Services, userId uint, on bool) error {
	return s.User.UpdatePrivacy(userId, PrivacyInput{
		RecordVisits:       on,
		PublicAlbum:        true,
		PublicCommunities:  true,
		PublicScraps:       true,
		PublicTestimonials: true,
	})
}
