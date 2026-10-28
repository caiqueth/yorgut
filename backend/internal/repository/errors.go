package repository

import (
	"errors"

	"gorm.io/gorm"
)

var (
	// ErrNotFound means the requested row does not exist (or is soft-deleted).
	ErrNotFound = errors.New("not found")
	// ErrDuplicate means a write was rejected by a unique index.
	ErrDuplicate = errors.New("already exists")
)

// translate turns the GORM errors callers care about into this package's
// errors, so services never import gorm. Other errors pass through unchanged.
// It relies on gorm.Config.TranslateError, which maps driver-specific unique
// violations to gorm.ErrDuplicatedKey.
func translate(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDuplicate
	default:
		return err
	}
}
