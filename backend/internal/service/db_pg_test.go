//go:build pgtest

package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"yorgut/internal/database"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The service tests run on a real Postgres, started once per test binary by
// embedded-postgres (no Docker; the first run downloads ~30 MB). Every test
// gets its own database, migrated with goose exactly like production.
// The port differs from internal/database's tests so both packages can run
// at the same time.
const testPGPort = 54331

var databaseCounter atomic.Int64

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "yorgut-service-pgtest-")
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

func openPG(t *testing.T, dbName string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.Open(pgDSN(dbName)), &gorm.Config{Logger: logger.Discard, TranslateError: true})
	if err != nil {
		t.Fatalf("opening postgres: %v", err)
	}
	return db
}

// openTestDB returns a fresh database on the embedded Postgres, migrated with
// database.MigratePostgres.
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := fmt.Sprintf("service_test_%d", databaseCounter.Add(1))

	admin := openPG(t, "postgres")
	if err := admin.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, name)).Error; err != nil {
		t.Fatalf("creating database: %v", err)
	}
	db := openPG(t, name)
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		admin.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name))
		sqlAdmin, _ := admin.DB()
		sqlAdmin.Close()
	})

	if err := database.MigratePostgres(db); err != nil {
		t.Fatalf("migrating postgres: %v", err)
	}
	return db
}
