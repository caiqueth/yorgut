package service

import (
	"testing"
	"yorgut/internal/repository"
)

// newTestRepositories returns repositories backed by a fresh, migrated
// database: in-memory SQLite, or Postgres when built with the pgtest tag
// (see openTestDB).
func newTestRepositories(t *testing.T) *repository.Repositories {
	t.Helper()
	return repository.NewRepositories(openTestDB(t))
}
