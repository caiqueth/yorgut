package transport

import (
	"encoding/json"
	"strings"
	"testing"
	"yorgut/internal/model"
)

func TestPictureURLs(t *testing.T) {
	previous := mediaURL
	t.Cleanup(func() { mediaURL = previous })
	SetMediaURL(func(key string) string { return "http://media/" + key })

	key := "abc.jpg"
	google := "https://lh3.googleusercontent.com/a/photo"

	tests := []struct {
		name    string
		picture *string
		want    *string
	}{
		{"stored key gets the media URL", &key, ptr("http://media/abc.jpg")},
		{"Google picture passes through", &google, &google},
		{"no picture stays nil", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &model.User{Profile: &model.UserProfile{ProfilePicture: tt.picture}}
			got := profilePictureOf(user)
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Fatalf("profilePictureOf = %v, want %v", deref(got), deref(tt.want))
			}
		})
	}

	var photo PhotoItem
	photo.FromPhoto(&model.Photo{URL: key})
	if photo.URL != "http://media/abc.jpg" {
		t.Fatalf("photo URL = %q, want the media URL", photo.URL)
	}
}

func ptr(s string) *string { return &s }

func deref(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

// TestEmailIsPrivate checks that only the logged user's own response carries
// an email; responses about other users must never have one.
func TestEmailIsPrivate(t *testing.T) {
	user := &model.User{Email: "ana@test.com", FirstName: "Ana", Profile: &model.UserProfile{}}

	hasEmail := func(response any) bool {
		body, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}
		return strings.Contains(string(body), "ana@test.com")
	}

	other := FindUserResponse{}
	other.FromUser(user)
	profile := FindUserWithProfileByIdResponse{}
	profile.FromUser(user)
	if hasEmail(other) || hasEmail(profile) {
		t.Error("a response about a user carries their email")
	}

	me := FindMeResponse{}
	me.FromUser(user)
	if !hasEmail(me) {
		t.Error("the logged user's own response lost their email")
	}
}

// TestProfileFieldKeys keeps the list users pick from and the fields hidden
// for non-friends the same: a key without a field would be "hidden" while
// still sent.
func TestProfileFieldKeys(t *testing.T) {
	fields := (&FindUserWithProfileByIdResponse{}).profileFields()
	if len(fields) != len(ProfileFieldKeys) {
		t.Fatalf("%d fields, %d keys", len(fields), len(ProfileFieldKeys))
	}
	for _, key := range ProfileFieldKeys {
		if _, ok := fields[key]; !ok {
			t.Errorf("key %q hides no field", key)
		}
	}
}
