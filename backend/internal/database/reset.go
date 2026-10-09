package database

import (
	"fmt"
	"net"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ResetLocalPostgres deletes every table and row by dropping and recreating
// the public schema. It exists for the seed command only, and refuses unless
// DATABASE_DSN points at this machine (localhost, a loopback address or a
// Unix socket), so it can't wipe a hosted database such as production.
// Run Migrate afterwards to recreate the schema.
func ResetLocalPostgres(db *gorm.DB) error {
	if name := db.Dialector.Name(); name != "postgres" {
		return fmt.Errorf("ResetLocalPostgres only runs on postgres, not %q", name)
	}
	if err := requireLocalPostgres(dsn); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DROP SCHEMA IF EXISTS public CASCADE").Error; err != nil {
			return err
		}
		return tx.Exec("CREATE SCHEMA public").Error
	})
}

// requireLocalPostgres returns an error unless every host the DSN can
// connect to (including PGHOST and fallbacks) is on this machine.
func requireLocalPostgres(dsn string) error {
	config, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parsing DATABASE_DSN: %w", err)
	}
	hosts := []string{config.Host}
	for _, fallback := range config.Fallbacks {
		hosts = append(hosts, fallback.Host)
	}
	for _, host := range hosts {
		if !isLocalHost(host) {
			return fmt.Errorf("refusing to reset the database at %q: only a database on localhost can be reset", host)
		}
	}
	return nil
}

func isLocalHost(host string) bool {
	if strings.HasPrefix(host, "/") {
		return true // Unix socket
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
