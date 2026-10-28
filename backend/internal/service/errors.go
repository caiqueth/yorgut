package service

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"yorgut/internal/repository"
)

// The error categories. Controllers map them to HTTP statuses in one place
// (controller.RespondError), so a service reports what went wrong and never
// how to answer. Specific errors (ErrCommunityForbidden, ...) belong to one
// of these categories: check with errors.Is(err, ErrForbidden).
var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrConflict  = errors.New("conflict")
	ErrInvalid   = errors.New("invalid")
)

// categorizedError keeps a readable message while belonging to a category.
type categorizedError struct {
	message  string
	category error
}

func (e *categorizedError) Error() string { return e.message }
func (e *categorizedError) Unwrap() error { return e.category }

// newError creates an error with the given message that matches category
// through errors.Is.
func newError(category error, message string) error {
	return &categorizedError{message: message, category: category}
}

// fromRepository translates the repository's errors into service categories,
// naming what was looked up or written with its article, as users read it:
// "o usuário", "a comunidade". Other errors (and nil) pass through unchanged.
func fromRepository(err error, what string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return newError(ErrNotFound, capitalize(what)+" não existe.")
	case errors.Is(err, repository.ErrDuplicate):
		return newError(ErrConflict, capitalize(what)+" já existe.")
	default:
		return err
	}
}

func capitalize(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// UserMessage is what an error says to the user: the message of its
// categorized error, without what wrapped it on the way up (those details are
// for the logs).
func UserMessage(err error) string {
	var categorized *categorizedError
	if errors.As(err, &categorized) {
		return categorized.message
	}
	return err.Error()
}
