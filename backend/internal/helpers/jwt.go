package helpers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Yorgut session tokens: a JWT signed with SESSION_SECRET whose subject is
// the user id. They are issued once at login (after Google validates the user)
// and sent back by the frontend as "Authorization: Bearer <token>".

const (
	sessionIssuer   = "yorgut"
	sessionDuration = 7 * 24 * time.Hour
)

var ErrInvalidSessionToken = errors.New("invalid session token")

// The secret is read on each call instead of at package init, so it is
// always read after the .env file has been loaded.
func sessionSecret() []byte {
	return []byte(os.Getenv("SESSION_SECRET"))
}

// SessionSecretConfigured reports whether SESSION_SECRET is set. The server
// refuses to start without it.
func SessionSecretConfigured() bool {
	return len(sessionSecret()) > 0
}

func NewSessionToken(userId uint) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    sessionIssuer,
		Subject:   strconv.FormatUint(uint64(userId), 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(sessionDuration)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(sessionSecret())
}

// ParseSessionToken checks the token's signature, issuer and expiry, and
// returns the user id it was issued for.
func ParseSessionToken(tokenString string) (uint, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(*jwt.Token) (any, error) { return sessionSecret(), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(sessionIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidSessionToken, err)
	}

	userId, err := strconv.ParseUint(claims.Subject, 10, 0)
	if err != nil || userId == 0 {
		return 0, ErrInvalidSessionToken
	}

	return uint(userId), nil
}
