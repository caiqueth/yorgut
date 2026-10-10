//go:build !pgtest

package service

import (
	"testing"
	"yorgut/internal/database"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// openTestDB returns a fresh, migrated in-memory SQLite database. Build with
// the pgtest tag to run the same tests on Postgres (db_pg_test.go).
func openTestDB(t *testing.T) *gorm.DB {
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
	return db
}
