package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"yorgut/internal/model"
)

// The schema has two sources, one per database engine:
//
//   - Postgres (production, and local dev with docker-compose): the versioned
//     SQL files in migrations/, applied by goose (MigratePostgres). They are
//     the source of truth for production.
//   - SQLite (DEV_MODE=true without Docker, and the service tests): GORM
//     AutoMigrate from the models (AutoMigrateSQLite).
//
// A schema change therefore needs both a new migration file and the matching
// model change. `go test -tags pgtest ./internal/database/` checks that the
// two produce the same Postgres schema.

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Models lists every GORM model, in an order AutoMigrate can create them.
func Models() []any {
	return []any{
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
		&model.Topic{},
		&model.UserReview{},
		&model.UserFan{},
	}
}

// Migrate brings the schema up to date with the path that matches the
// database: MigratePostgres or AutoMigrateSQLite. It panics on failure:
// running on a half-migrated schema breaks in less obvious ways.
func Migrate(db *gorm.DB) {
	var err error
	switch name := db.Dialector.Name(); name {
	case "postgres":
		err = MigratePostgres(db)
	case "sqlite":
		err = AutoMigrateSQLite(db)
	default:
		err = fmt.Errorf("no migration path for database %q", name)
	}
	if err != nil {
		panic(fmt.Errorf("database migration failed: %w", err))
	}
}

// MigratePostgres applies the pending migrations in migrations/ with goose.
// A Postgres advisory lock keeps two instances starting at once from running
// them twice. Each migration runs in its own transaction.
func MigratePostgres(db *gorm.DB) error {
	provider, err := newMigrationProvider(db)
	if err != nil {
		return err
	}
	results, err := provider.Up(context.Background())
	for _, r := range results {
		if r.Error != nil {
			continue
		}
		log.Info().Msgf("migration %s applied in %s", r.Source.Path, r.Duration)
	}
	if err != nil {
		return err
	}
	version, err := provider.GetDBVersion(context.Background())
	if err != nil {
		return err
	}
	log.Info().Msgf("database schema at migration version %d", version)
	return nil
}

func newMigrationProvider(db *gorm.DB) (*goose.Provider, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	files, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return nil, err
	}
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return nil, err
	}
	// Never call provider.Close(): it would close GORM's connection pool.
	return goose.NewProvider(goose.DialectPostgres, sqlDB, files, goose.WithSessionLocker(locker))
}

// friendshipActivePairSQLite mirrors the idx_friendships_active_pair
// expression index of migration 00003, which GORM tags can't express.
// SQLite's two-argument min/max are its LEAST/GREATEST.
const friendshipActivePairSQLite = `CREATE UNIQUE INDEX IF NOT EXISTS "idx_friendships_active_pair" ` +
	`ON "friendships" (min("user_requested_id", "user_requestee_id"), max("user_requested_id", "user_requestee_id")) ` +
	`WHERE "deleted_at" IS NULL`

// AutoMigrateSQLite creates or updates the SQLite schema from the models.
// It refuses other databases: on Postgres the migrations are the source of
// truth, and AutoMigrate would drift from them.
func AutoMigrateSQLite(db *gorm.DB) error {
	if name := db.Dialector.Name(); name != "sqlite" {
		return fmt.Errorf("AutoMigrateSQLite only runs on sqlite, not %q; use MigratePostgres", name)
	}
	if err := db.AutoMigrate(Models()...); err != nil {
		return err
	}
	return db.Exec(friendshipActivePairSQLite).Error
}
