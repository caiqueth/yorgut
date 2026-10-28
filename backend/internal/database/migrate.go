package database

import (
	"context"
	"database/sql"
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
//   - MySQL (production, and local dev with docker-compose): the versioned
//     SQL files in migrations/, applied by goose (MigrateMySQL). They are the
//     source of truth for production.
//   - SQLite (DEV_MODE=true without Docker, and the service tests): GORM
//     AutoMigrate from the models (AutoMigrateSQLite).
//
// A schema change therefore needs both a new migration file and the matching
// model change. `go test -tags mysqltest ./internal/database/` checks that the
// two produce the same MySQL schema.

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
		&model.CommunityJoinRequest{},
		&model.CommunityInvitation{},
		&model.ProfileVisit{},
		&model.Friendship{},
		&model.PhotoAlbum{},
		&model.Photo{},
		&model.Scrap{},
		&model.Testimonial{},
		&model.Topic{},
		&model.UserReview{},
		&model.UserFan{},
		&model.Poll{},
		&model.PollOption{},
		&model.PollVote{},
	}
}

// Migrate brings the schema up to date with the path that matches the
// database: MigrateMySQL or AutoMigrateSQLite. It panics on failure:
// running on a half-migrated schema breaks in less obvious ways.
func Migrate(db *gorm.DB) {
	var err error
	switch name := db.Dialector.Name(); name {
	case "mysql":
		err = MigrateMySQL(db)
	case "sqlite":
		err = AutoMigrateSQLite(db)
	default:
		err = fmt.Errorf("no migration path for database %q", name)
	}
	if err != nil {
		panic(fmt.Errorf("database migration failed: %w", err))
	}
}

// MigrateMySQL applies the pending migrations in migrations/ with goose. A
// named lock (mysqlLocker) keeps two instances starting at once from running
// them twice. Each migration runs in a transaction, but MySQL commits DDL
// (CREATE, ALTER, DROP) on its own: a migration that fails halfway keeps the
// statements before the failing one, so keep schema changes small.
func MigrateMySQL(db *gorm.DB) error {
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
	// Never call provider.Close(): it would close GORM's connection pool.
	return goose.NewProvider(goose.DialectMySQL, sqlDB, files, goose.WithSessionLocker(mysqlLocker{}))
}

// migrationLockName is the MySQL named lock (GET_LOCK) held while migrating.
const migrationLockName = "yorgut_goose_migrations"

// mysqlLocker holds a MySQL named lock for the whole migration session; goose
// only ships a Postgres one. The lock belongs to the connection, so MySQL
// releases it if the process dies.
type mysqlLocker struct{}

var _ lock.SessionLocker = mysqlLocker{}

func (mysqlLocker) SessionLock(ctx context.Context, conn *sql.Conn) error {
	var acquired sql.NullInt64
	// Waits up to 5 minutes for another instance to finish migrating.
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK(?, 300)", migrationLockName).Scan(&acquired); err != nil {
		return fmt.Errorf("taking the migration lock: %w", err)
	}
	if acquired.Int64 != 1 {
		return fmt.Errorf("timed out waiting for the migration lock %q", migrationLockName)
	}
	return nil
}

func (mysqlLocker) SessionUnlock(ctx context.Context, conn *sql.Conn) error {
	_, err := conn.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", migrationLockName)
	return err
}

// friendshipActivePairSQLite mirrors the idx_friendships_active_pair
// expression index of the baseline migration, which GORM tags can't express.
// SQLite's two-argument min/max are its LEAST/GREATEST.
const friendshipActivePairSQLite = `CREATE UNIQUE INDEX IF NOT EXISTS "idx_friendships_active_pair" ` +
	`ON "friendships" (min("user_requested_id", "user_requestee_id"), max("user_requested_id", "user_requestee_id")) ` +
	`WHERE "deleted_at" IS NULL`

// AutoMigrateSQLite creates or updates the SQLite schema from the models.
// It refuses other databases: on MySQL the migrations are the source of
// truth, and AutoMigrate would drift from them.
func AutoMigrateSQLite(db *gorm.DB) error {
	if name := db.Dialector.Name(); name != "sqlite" {
		return fmt.Errorf("AutoMigrateSQLite only runs on sqlite, not %q; use MigrateMySQL", name)
	}
	if err := db.AutoMigrate(Models()...); err != nil {
		return err
	}
	return db.Exec(friendshipActivePairSQLite).Error
}
