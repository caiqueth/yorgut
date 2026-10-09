package service

import (
	"testing"
	"yorgut/internal/database"
	"yorgut/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestRepositories returns repositories backed by a fresh, migrated
// in-memory SQLite database.
func newTestRepositories(t *testing.T) *repository.Repositories {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{Logger: logger.Discard, TranslateError: true})
	if err != nil {
		t.Fatalf("opening sqlite: %v", err)
	}
	// Each connection to ":memory:" is a separate database: keep only one.
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := database.AutoMigrateSQLite(db); err != nil {
		t.Fatalf("migrating sqlite: %v", err)
	}
	return repository.NewRepositories(db)
}
