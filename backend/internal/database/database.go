package database

import (
	"os"
	"yorgut/internal/model"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dburl string = os.Getenv("BLUEPRINT_DB_URL")

func DatabaseConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dburl), &gorm.Config{})
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		panic(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.UserProfile{},
		&model.Community{},
		&model.CommunityProfile{},
		&model.CommunityMembership{},
		&model.Friendship{},
		&model.PhotoAlbum{},
		&model.Photo{},
		&model.Scrap{},
		&model.Testimonial{},
		&model.Forum{},
		&model.Topic{},
	)
}
