//go:build pgtest

// These tests run against a real Postgres, downloaded and started by
// embedded-postgres (no Docker needed; the first run downloads ~30 MB).
// They are behind a build tag so `go test ./...` stays fast and offline:
//
//	go test -tags pgtest ./internal/database/
package database

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testPGPort = 54330

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "yorgut-pgtest-")
	if err != nil {
		panic(err)
	}
	var pgLog bytes.Buffer
	pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Version(embeddedpostgres.V17).
		Port(testPGPort).
		RuntimePath(filepath.Join(dir, "runtime")).
		Logger(&pgLog))
	if err := pg.Start(); err != nil {
		fmt.Println(pgLog.String())
		panic(err)
	}
	code := m.Run()
	if err := pg.Stop(); err != nil {
		fmt.Println(err)
	}
	os.RemoveAll(dir)
	os.Exit(code)
}

func pgDSN(dbName string) string {
	return fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=%s sslmode=disable", testPGPort, dbName)
}

// newPostgres creates an empty database named after the test and suffix.
func newPostgres(t *testing.T, suffix string) *gorm.DB {
	t.Helper()
	name := strings.ToLower(strings.NewReplacer("/", "_", "-", "_").Replace(t.Name() + "_" + suffix))
	admin, err := gorm.Open(postgres.Open(pgDSN("postgres")), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	admin.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, name))
	if err := admin.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, name)).Error; err != nil {
		t.Fatal(err)
	}
	sqlAdmin, _ := admin.DB()
	sqlAdmin.Close()

	db, err := gorm.Open(postgres.Open(pgDSN(name)), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

// dumpSchema describes tables, columns (type, nullability, default),
// indexes, constraints and sequences, one sorted line each, ignoring goose's
// own table. Column order is ignored: ALTER TABLE ADD COLUMN appends.
func dumpSchema(t *testing.T, db *gorm.DB) []string {
	t.Helper()
	queries := map[string]string{
		"column": `SELECT table_name || '.' || column_name || ' ' || data_type || ' nullable=' || is_nullable || ' default=' || coalesce(column_default, '')
			FROM information_schema.columns WHERE table_schema = 'public' AND table_name <> 'goose_db_version'`,
		"index": `SELECT indexdef FROM pg_indexes WHERE schemaname = 'public' AND tablename <> 'goose_db_version'`,
		"constraint": `SELECT conrelid::regclass::text || ' ' || conname || ' ' || pg_get_constraintdef(oid)
			FROM pg_constraint WHERE connamespace = 'public'::regnamespace AND conrelid::regclass::text <> 'goose_db_version'`,
		"sequence": `SELECT sequence_name || ' ' || data_type FROM information_schema.sequences
			WHERE sequence_schema = 'public' AND sequence_name NOT LIKE 'goose%'`,
	}
	var lines []string
	for kind, query := range queries {
		var rows []string
		if err := db.Raw(query).Scan(&rows).Error; err != nil {
			t.Fatal(err)
		}
		for _, row := range rows {
			lines = append(lines, kind+": "+row)
		}
	}
	slices.Sort(lines)
	return lines
}

func diffSchemas(a, b []string) (onlyA, onlyB []string) {
	for _, line := range a {
		if !slices.Contains(b, line) {
			onlyA = append(onlyA, line)
		}
	}
	for _, line := range b {
		if !slices.Contains(a, line) {
			onlyB = append(onlyB, line)
		}
	}
	return onlyA, onlyB
}

// migrationOnly are objects the migrations create that GORM tags can't
// describe (AutoMigrateSQLite adds an SQLite equivalent by hand).
const migrationOnly = "idx_friendships_active_pair"

func withoutMigrationOnly(lines []string) []string {
	return slices.DeleteFunc(slices.Clone(lines), func(line string) bool {
		return strings.Contains(line, migrationOnly)
	})
}

// The migrations and the models must describe the same schema.
func TestMigrationsMatchAutoMigrate(t *testing.T) {
	migrated := newPostgres(t, "migrated")
	if err := MigratePostgres(migrated); err != nil {
		t.Fatal(err)
	}
	auto := newPostgres(t, "auto")
	if err := auto.AutoMigrate(Models()...); err != nil {
		t.Fatal(err)
	}

	migratedSchema := dumpSchema(t, migrated)
	if !slices.ContainsFunc(migratedSchema, func(l string) bool { return strings.Contains(l, migrationOnly) }) {
		t.Errorf("migrations did not create %s", migrationOnly)
	}
	onlyMigrated, onlyAuto := diffSchemas(withoutMigrationOnly(migratedSchema), dumpSchema(t, auto))
	if len(onlyMigrated) > 0 || len(onlyAuto) > 0 {
		t.Errorf("schemas differ.\nonly in migrations:\n  %s\nonly in AutoMigrate:\n  %s",
			strings.Join(onlyMigrated, "\n  "), strings.Join(onlyAuto, "\n  "))
	}

	// Running again is a no-op.
	if err := MigratePostgres(migrated); err != nil {
		t.Fatal(err)
	}
}

// A database AutoMigrate already created (like production before goose) can
// be taken over by the migrations without errors or changes besides the
// migration-only objects.
func TestMigrationsOnAutoMigratedDatabase(t *testing.T) {
	db := newPostgres(t, "db")
	if err := db.AutoMigrate(Models()...); err != nil {
		t.Fatal(err)
	}
	before := dumpSchema(t, db)
	if err := MigratePostgres(db); err != nil {
		t.Fatal(err)
	}
	onlyBefore, onlyAfter := diffSchemas(before, withoutMigrationOnly(dumpSchema(t, db)))
	if len(onlyBefore) > 0 || len(onlyAfter) > 0 {
		t.Errorf("migrations changed the schema.\nremoved:\n  %s\nadded:\n  %s",
			strings.Join(onlyBefore, "\n  "), strings.Join(onlyAfter, "\n  "))
	}
}

// Migration 3 soft-deletes duplicates before adding its unique indexes, and
// recounts friends and members.
func TestMigrationRemovesDuplicates(t *testing.T) {
	db := newPostgres(t, "db")
	provider, err := newMigrationProvider(db)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := provider.UpTo(context.Background(), 2); err != nil {
		t.Fatal(err)
	}

	exec := func(sql string) {
		t.Helper()
		if err := db.Exec(sql).Error; err != nil {
			t.Fatal(err)
		}
	}
	exec(`INSERT INTO users (id, email, friend_count) VALUES (1, 'a@x', 9), (2, 'b@x', 9), (3, 'c@x', 9)`)
	exec(`INSERT INTO communities (id, name, owner_id, member_count) VALUES (1, 'c', 1, 9)`)
	exec(`INSERT INTO community_memberships (id, community_id, user_id, deleted_at) VALUES
		(1, 1, 1, NULL), (2, 1, 1, NULL), (3, 1, 2, now()), (4, 1, 2, NULL)`)
	// 1-2: pending (10) and accepted (11) in opposite directions: the accepted one stays.
	// 1-3: two accepted (12, 13): the oldest stays.
	exec(`INSERT INTO friendships (id, user_requested_id, user_requestee_id, accepted, deleted_at) VALUES
		(10, 1, 2, false, NULL), (11, 2, 1, true, NULL), (12, 1, 3, true, NULL), (13, 3, 1, true, NULL)`)

	if _, err := provider.Up(context.Background()); err != nil {
		t.Fatal(err)
	}

	var liveMemberships, liveFriendships []int
	db.Raw(`SELECT id FROM community_memberships WHERE deleted_at IS NULL ORDER BY id`).Scan(&liveMemberships)
	db.Raw(`SELECT id FROM friendships WHERE deleted_at IS NULL ORDER BY id`).Scan(&liveFriendships)
	if !slices.Equal(liveMemberships, []int{1, 4}) {
		t.Errorf("live memberships = %v, want [1 4]", liveMemberships)
	}
	if !slices.Equal(liveFriendships, []int{11, 12}) {
		t.Errorf("live friendships = %v, want [11 12]", liveFriendships)
	}

	var friendCounts []int
	db.Raw(`SELECT friend_count FROM users ORDER BY id`).Scan(&friendCounts)
	if !slices.Equal(friendCounts, []int{2, 1, 1}) {
		t.Errorf("friend counts = %v, want [2 1 1]", friendCounts)
	}
	var memberCount int
	db.Raw(`SELECT member_count FROM communities WHERE id = 1`).Scan(&memberCount)
	if memberCount != 2 {
		t.Errorf("member count = %d, want 2", memberCount)
	}

	// The indexes now reject new duplicates, in either direction.
	if err := db.Exec(`INSERT INTO community_memberships (community_id, user_id) VALUES (1, 1)`).Error; err == nil {
		t.Error("duplicate membership was accepted")
	}
	if err := db.Exec(`INSERT INTO friendships (user_requested_id, user_requestee_id, accepted) VALUES (1, 2, false)`).Error; err == nil {
		t.Error("duplicate friendship was accepted")
	}
}

func TestResetLocalPostgres(t *testing.T) {
	db := newPostgres(t, "db")
	if err := MigratePostgres(db); err != nil {
		t.Fatal(err)
	}
	// ResetLocalPostgres checks the package's DATABASE_DSN.
	old := dsn
	dsn = pgDSN("unused")
	defer func() { dsn = old }()

	if err := ResetLocalPostgres(db); err != nil {
		t.Fatal(err)
	}
	if schema := dumpSchema(t, db); len(schema) > 0 {
		t.Errorf("schema not empty after reset: %v", schema)
	}
	if err := MigratePostgres(db); err != nil {
		t.Fatalf("migrating after reset: %v", err)
	}
}
