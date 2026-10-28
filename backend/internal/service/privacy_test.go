package service

import (
	"errors"
	"slices"
	"testing"
	"time"
	"yorgut/internal/repository"
)

// closeEverything keeps all of userId's sections, the birthday and the given
// profile fields to friends.
func closeEverything(t *testing.T, s *Services, userId uint, hiddenFields ...string) {
	t.Helper()
	err := s.User.UpdatePrivacy(userId, PrivacyInput{RecordVisits: true, HiddenProfileFields: hiddenFields})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProfilePrivacy(t *testing.T) {
	// ana owns the profile, bia is her friend, caio a stranger; 0 is the
	// preview of ana's profile as non-friends see it.
	setup := func(t *testing.T) (*Services, *repository.Repositories, uint, uint, uint) {
		s, repos, ana, bia, caio := threeUsers(t)
		befriend(t, s, repos, ana, bia)
		birth := time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC)
		if err := repos.User.UpdateBasicData(ana.ID, "Ana", "", &birth); err != nil {
			t.Fatal(err)
		}
		description := "oi"
		if err := s.UserProfile.UpdateProfile(ana.ID, UpdateUserProfileInput{Description: description, Religion: "x"}); err != nil {
			t.Fatal(err)
		}
		return s, repos, ana.ID, bia.ID, caio.ID
	}

	sections := func(s *Services, viewerId, ownerId uint) map[string]error {
		_, album := s.Photo.FindPhotosOfUser(viewerId, ownerId)
		_, scraps := s.Scrap.FindScrapOfUserRecipient(viewerId, ownerId, 1, 10)
		_, testimonials := s.Testimonial.FindTestimonialsOfUser(viewerId, ownerId, 0, 0)
		_, communities := s.CommunityMembership.FindUserMemberships(viewerId, ownerId, 1, 10, false)
		return map[string]error{"album": album, "scraps": scraps, "testimonials": testimonials, "communities": communities}
	}

	t.Run("by default sections are open and the birthday is friends-only", func(t *testing.T) {
		s, _, ana, bia, caio := setup(t)
		for name, err := range sections(s, caio, ana) {
			if err != nil {
				t.Errorf("stranger, %s: %v", name, err)
			}
		}
		stranger, err := s.User.FindUserWithProfileById(caio, ana)
		if err != nil {
			t.Fatal(err)
		}
		if stranger.BirthDate != nil || stranger.Description == nil || !slices.Equal(stranger.HiddenFields, []string{"birthDate"}) {
			t.Errorf("stranger: birthday %v, description %v, hidden %v; want no birthday, the description, only the birthday hidden",
				stranger.BirthDate, stranger.Description, stranger.HiddenFields)
		}
		friend, _ := s.User.FindUserWithProfileById(bia, ana)
		if friend.BirthDate == nil {
			t.Error("friend: no birthday, want it")
		}
	})

	t.Run("closed sections are refused to strangers and the preview, not to friends", func(t *testing.T) {
		s, _, ana, bia, caio := setup(t)
		closeEverything(t, s, ana, "religion")

		for _, viewer := range []struct {
			name string
			id   uint
		}{{"stranger", caio}, {"preview", 0}} {
			for name, err := range sections(s, viewer.id, ana) {
				if !errors.Is(err, ErrPrivateSection) {
					t.Errorf("%s, %s: got %v, want ErrPrivateSection", viewer.name, name, err)
				}
			}
			profile, err := s.User.FindUserWithProfileById(viewer.id, ana)
			if err != nil {
				t.Fatal(err)
			}
			if profile.Religion != nil || !slices.Equal(profile.HiddenFields, []string{"birthDate", "religion"}) || profile.Description == nil {
				t.Errorf("%s: religion %v, hidden %v, description %v; want religion hidden, description shown",
					viewer.name, profile.Religion, profile.HiddenFields, profile.Description)
			}
			if v := profile.Visibility; v.Album || v.Scraps || v.Testimonials || v.Communities {
				t.Errorf("%s: visibility %+v, want everything closed", viewer.name, v)
			}
		}

		for _, viewer := range []uint{bia, ana} {
			for name, err := range sections(s, viewer, ana) {
				if err != nil {
					t.Errorf("viewer %d, %s: %v", viewer, name, err)
				}
			}
			profile, _ := s.User.FindUserWithProfileById(viewer, ana)
			if profile.Religion == nil || len(profile.HiddenFields) != 0 {
				t.Errorf("viewer %d: religion %v, hidden %v; want everything", viewer, profile.Religion, profile.HiddenFields)
			}
		}
	})

	t.Run("a closed scrapbook only takes scraps from friends", func(t *testing.T) {
		s, _, ana, bia, caio := setup(t)
		closeEverything(t, s, ana)
		wantCategory(t, s.Scrap.CreateScrap(caio, ana, "oi"), ErrForbidden)
		if err := s.Scrap.CreateScrap(bia, ana, "oi"); err != nil {
			t.Errorf("friend: %v", err)
		}
	})

	t.Run("an opened birthday shows to strangers", func(t *testing.T) {
		s, _, ana, _, caio := setup(t)
		err := s.User.UpdatePrivacy(ana, PrivacyInput{RecordVisits: true, PublicAlbum: true, PublicCommunities: true,
			PublicScraps: true, PublicTestimonials: true, PublicBirthday: true})
		if err != nil {
			t.Fatal(err)
		}
		profile, _ := s.User.FindUserWithProfileById(caio, ana)
		if profile.BirthDate == nil || profile.BirthDate.Month() != time.March {
			t.Errorf("birthday %v, want March 15", profile.BirthDate)
		}
	})

	t.Run("unknown profile fields are refused", func(t *testing.T) {
		s, _, ana, _, _ := setup(t)
		err := s.User.UpdatePrivacy(ana, PrivacyInput{HiddenProfileFields: []string{"email"}})
		wantCategory(t, err, ErrInvalid)
	})
}
