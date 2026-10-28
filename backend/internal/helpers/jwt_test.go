package helpers

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSessionTokenRoundTrip(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")

	token, err := NewSessionToken(42)
	if err != nil {
		t.Fatalf("NewSessionToken: %v", err)
	}

	userId, err := ParseSessionToken(token)
	if err != nil {
		t.Fatalf("ParseSessionToken: %v", err)
	}
	if userId != 42 {
		t.Fatalf("got user id %d, want 42", userId)
	}
}

func TestParseSessionTokenRejects(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")

	sign := func(claims jwt.Claims, method jwt.SigningMethod, key any) string {
		s, err := jwt.NewWithClaims(method, claims).SignedString(key)
		if err != nil {
			t.Fatalf("signing: %v", err)
		}
		return s
	}
	valid := func() jwt.RegisteredClaims {
		return jwt.RegisteredClaims{
			Issuer:    sessionIssuer,
			Subject:   "42",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		}
	}

	expired := valid()
	expired.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute))
	noExpiry := valid()
	noExpiry.ExpiresAt = nil
	otherIssuer := valid()
	otherIssuer.Issuer = "accounts.google.com"
	badSubject := valid()
	badSubject.Subject = "not-a-number"

	cases := map[string]string{
		"garbage":      "not.a.token",
		"wrong secret": sign(valid(), jwt.SigningMethodHS256, []byte("other-secret")),
		"alg none":     sign(valid(), jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType),
		"expired":      sign(expired, jwt.SigningMethodHS256, []byte("test-secret")),
		"no expiry":    sign(noExpiry, jwt.SigningMethodHS256, []byte("test-secret")),
		"other issuer": sign(otherIssuer, jwt.SigningMethodHS256, []byte("test-secret")),
		"bad subject":  sign(badSubject, jwt.SigningMethodHS256, []byte("test-secret")),
	}

	for name, token := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseSessionToken(token); !errors.Is(err, ErrInvalidSessionToken) {
				t.Fatalf("got err %v, want ErrInvalidSessionToken", err)
			}
		})
	}
}
