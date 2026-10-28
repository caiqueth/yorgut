//go:build mysqltest

package service

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"
	"yorgut/internal/database"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The service tests run on a real MySQL: the docker-compose one by default
// (`make db_up`), or the server TEST_MYSQL_DSN names, as a user that may
// create and drop databases. Every test gets its own database, migrated with
// goose exactly like production.

var databaseCounter atomic.Int64

// Several test binaries may run at once: the names include the process id.
var runID = fmt.Sprintf("%d_%d", os.Getpid(), time.Now().Unix()%100000)

func testServerDSN() string {
	if dsn := os.Getenv("TEST_MYSQL_DSN"); dsn != "" {
		return dsn
	}
	return "root:localdev@tcp(127.0.0.1:3306)/"
}

func openMySQL(t *testing.T, dbName string) *gorm.DB {
	t.Helper()
	config, err := database.ParseMySQLDSN(testServerDSN())
	if err != nil {
		t.Fatal(err)
	}
	config.DBName = dbName
	db, err := database.OpenMySQL(config.FormatDSN(), &gorm.Config{Logger: logger.Discard, TranslateError: true})
	if err != nil {
		t.Fatalf("opening mysql: %v", err)
	}
	return db
}

// openTestDB returns a fresh database on the test MySQL, migrated with
// database.MigrateMySQL.
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := fmt.Sprintf("service_test_%s_%d", runID, databaseCounter.Add(1))

	admin := openMySQL(t, "")
	if err := admin.Exec(fmt.Sprintf("CREATE DATABASE `%s`", name)).Error; err != nil {
		t.Fatalf("creating database: %v", err)
	}
	db := openMySQL(t, name)
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		admin.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", name))
		sqlAdmin, _ := admin.DB()
		sqlAdmin.Close()
	})

	if err := database.MigrateMySQL(db); err != nil {
		t.Fatalf("migrating mysql: %v", err)
	}
	return db
}
