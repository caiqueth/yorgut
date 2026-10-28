//go:build mysqltest

// These tests run against a real MySQL: the docker-compose one by default
// (`make db_up`), or the server TEST_MYSQL_DSN names, as a user that may
// create and drop databases. Never point it at production. They are behind
// a build tag so `go test ./...` stays fast and offline:
//
//	go test -tags mysqltest ./internal/database/
package database

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testServerDSN() string {
	if dsn := os.Getenv("TEST_MYSQL_DSN"); dsn != "" {
		return dsn
	}
	return "root:localdev@tcp(127.0.0.1:3306)/"
}

// testDSN is testServerDSN with dbName as the database.
func testDSN(t *testing.T, dbName string) string {
	t.Helper()
	config, err := ParseMySQLDSN(testServerDSN())
	if err != nil {
		t.Fatal(err)
	}
	config.DBName = dbName
	return config.FormatDSN()
}

// newMySQL creates an empty database named after the test and suffix.
func newMySQL(t *testing.T, suffix string) *gorm.DB {
	t.Helper()
	name := strings.ToLower(strings.NewReplacer("/", "_", "-", "_").Replace(t.Name() + "_" + suffix))
	admin, err := OpenMySQL(testServerDSN(), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	admin.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", name))
	if err := admin.Exec(fmt.Sprintf("CREATE DATABASE `%s`", name)).Error; err != nil {
		t.Fatal(err)
	}

	db, err := OpenMySQL(testDSN(t, name), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		admin.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", name))
		sqlAdmin, _ := admin.DB()
		sqlAdmin.Close()
	})
	return db
}

// dumpSchema describes tables, columns (type, nullability, default), index
// parts and foreign keys of the current database, one sorted line each,
// ignoring goose's own table. Column order is ignored: ALTER TABLE ADD COLUMN
// appends.
func dumpSchema(t *testing.T, db *gorm.DB) []string {
	t.Helper()
	queries := map[string]string{
		"column": `SELECT CONCAT(table_name, '.', column_name, ' ', column_type, ' nullable=', is_nullable,
				' default=', COALESCE(column_default, ''), ' ', extra)
			FROM information_schema.columns
			WHERE table_schema = DATABASE() AND table_name <> 'goose_db_version'`,
		"index": `SELECT CONCAT(table_name, ' ', index_name, IF(non_unique, '', ' unique'), ' #', seq_in_index, ' ',
				COALESCE(column_name, expression))
			FROM information_schema.statistics
			WHERE table_schema = DATABASE() AND table_name <> 'goose_db_version'`,
		"foreign key": `SELECT CONCAT(k.table_name, ' ', k.constraint_name, ' ', k.column_name, ' -> ',
				k.referenced_table_name, '.', k.referenced_column_name, ' on delete ', r.delete_rule, ' on update ', r.update_rule)
			FROM information_schema.key_column_usage k
			JOIN information_schema.referential_constraints r
				ON r.constraint_schema = k.constraint_schema AND r.constraint_name = k.constraint_name
			WHERE k.table_schema = DATABASE()`,
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

// migrationOnly are the indexes the migrations create in a way GORM tags
// can't describe: unique among live rows only (see the baseline migration).
// AutoMigrate makes the first three plain unique indexes on MySQL, and SQLite
// gets real partial ones.
var migrationOnly = []string{
	"idx_community_memberships_active_pair",
	"idx_community_join_requests_active_pair",
	"idx_community_invitations_active_pair",
	"idx_friendships_active_pair",
}

func withoutMigrationOnly(lines []string) []string {
	return slices.DeleteFunc(slices.Clone(lines), func(line string) bool {
		return slices.ContainsFunc(migrationOnly, func(name string) bool { return strings.Contains(line, " "+name+" ") })
	})
}

// The migrations and the models must describe the same schema.
func TestMigrationsMatchAutoMigrate(t *testing.T) {
	migrated := newMySQL(t, "migrated")
	if err := MigrateMySQL(migrated); err != nil {
		t.Fatal(err)
	}
	auto := newMySQL(t, "auto")
	if err := auto.AutoMigrate(Models()...); err != nil {
		t.Fatal(err)
	}

	migratedSchema := dumpSchema(t, migrated)
	for _, name := range migrationOnly {
		if !slices.ContainsFunc(migratedSchema, func(l string) bool { return strings.Contains(l, " "+name+" ") }) {
			t.Errorf("migrations did not create %s", name)
		}
	}
	onlyMigrated, onlyAuto := diffSchemas(withoutMigrationOnly(migratedSchema), withoutMigrationOnly(dumpSchema(t, auto)))
	if len(onlyMigrated) > 0 || len(onlyAuto) > 0 {
		t.Errorf("schemas differ.\nonly in migrations:\n  %s\nonly in AutoMigrate:\n  %s",
			strings.Join(onlyMigrated, "\n  "), strings.Join(onlyAuto, "\n  "))
	}

	// Running again is a no-op.
	if err := MigrateMySQL(migrated); err != nil {
		t.Fatal(err)
	}
}

// The active-pair indexes reject a second live row, in either direction for
// friendships, but not a new row next to soft-deleted ones.
func TestActivePairIndexes(t *testing.T) {
	db := newMySQL(t, "db")
	if err := MigrateMySQL(db); err != nil {
		t.Fatal(err)
	}
	exec := func(sql string) error { return db.Exec(sql).Error }
	mustExec := func(sql string) {
		t.Helper()
		if err := exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	mustExec(`INSERT INTO users (id, email) VALUES (1, 'a@x'), (2, 'b@x')`)
	mustExec(`INSERT INTO communities (id, name, owner_id) VALUES (1, 'c', 1)`)

	mustExec(`INSERT INTO community_memberships (community_id, user_id) VALUES (1, 2)`)
	if exec(`INSERT INTO community_memberships (community_id, user_id) VALUES (1, 2)`) == nil {
		t.Error("duplicate membership was accepted")
	}
	mustExec(`UPDATE community_memberships SET deleted_at = NOW()`)
	mustExec(`INSERT INTO community_memberships (community_id, user_id) VALUES (1, 2)`)
	mustExec(`UPDATE community_memberships SET deleted_at = NOW() WHERE deleted_at IS NULL`)
	mustExec(`INSERT INTO community_memberships (community_id, user_id) VALUES (1, 2)`)

	mustExec(`INSERT INTO community_join_requests (community_id, user_id) VALUES (1, 2)`)
	if exec(`INSERT INTO community_join_requests (community_id, user_id) VALUES (1, 2)`) == nil {
		t.Error("duplicate join request was accepted")
	}

	mustExec(`INSERT INTO friendships (user_requested_id, user_requestee_id, accepted) VALUES (1, 2, false)`)
	if exec(`INSERT INTO friendships (user_requested_id, user_requestee_id, accepted) VALUES (2, 1, false)`) == nil {
		t.Error("duplicate friendship in the other direction was accepted")
	}
	mustExec(`UPDATE friendships SET deleted_at = NOW()`)
	mustExec(`INSERT INTO friendships (user_requested_id, user_requestee_id, accepted) VALUES (2, 1, false)`)
}

// The migration lock is exclusive, and released for the next run.
func TestMigrationLock(t *testing.T) {
	db := newMySQL(t, "db")
	sqlDB, _ := db.DB()
	ctx := context.Background()
	first, err := sqlDB.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	if err := (mysqlLocker{}).SessionLock(ctx, first); err != nil {
		t.Fatal(err)
	}
	var free int
	db.Raw("SELECT IS_FREE_LOCK(?)", migrationLockName).Scan(&free)
	if free != 0 {
		t.Error("lock is free while held")
	}
	if err := (mysqlLocker{}).SessionUnlock(ctx, first); err != nil {
		t.Fatal(err)
	}
	if err := MigrateMySQL(db); err != nil {
		t.Fatalf("migrating after unlock: %v", err)
	}
}

func TestResetLocalMySQL(t *testing.T) {
	db := newMySQL(t, "db")
	if err := MigrateMySQL(db); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO users (email) VALUES ('a@x')`).Error; err != nil {
		t.Fatal(err)
	}
	// ResetLocalMySQL checks the package's DATABASE_DSN.
	old := dsn
	dsn = "u:p@tcp(127.0.0.1:3306)/unused"
	defer func() { dsn = old }()

	if err := ResetLocalMySQL(db); err != nil {
		t.Fatal(err)
	}
	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tables)
	if len(tables) > 0 {
		t.Errorf("tables left after reset: %v", tables)
	}
	if err := MigrateMySQL(db); err != nil {
		t.Fatalf("migrating after reset: %v", err)
	}
}
