package database

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dsn string = os.Getenv("DATABASE_DSN")
var development bool = os.Getenv("DEV_MODE") == "true"

// DatabaseConnection opens SQLite when DEV_MODE=true (DATABASE_DSN is the
// file path), else Postgres (DATABASE_DSN is a Postgres DSN or URL). Schema
// setup is separate: see Migrate in migrate.go.
func DatabaseConnection() *gorm.DB {
	var db *gorm.DB
	var err error

	if development {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	}

	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		panic(err)
	}

	return db
}
