package service

import (
	"errors"

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
// naming what was looked up or written ("user", "community", ...). Other
// errors (and nil) pass through unchanged.
func fromRepository(err error, what string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return newError(ErrNotFound, what+" not found")
	case errors.Is(err, repository.ErrDuplicate):
		return newError(ErrConflict, what+" already exists")
	default:
		return err
	}
}
