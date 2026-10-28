package database

import (
	"fmt"
	"net"
	"strings"

	"gorm.io/gorm"
)

// ResetLocalMySQL deletes every table and row in the database DATABASE_DSN
// names (not the database itself, which a hosted user often can't recreate).
// It exists for the seed command only, and refuses unless DATABASE_DSN points
// at this machine (localhost, a loopback address or a Unix socket), so it
// can't wipe a hosted database such as production. Run Migrate afterwards to
// recreate the schema.
func ResetLocalMySQL(db *gorm.DB) error {
	if name := db.Dialector.Name(); name != "mysql" {
		return fmt.Errorf("ResetLocalMySQL only runs on mysql, not %q", name)
	}
	if err := requireLocalMySQL(dsn); err != nil {
		return err
	}
	// One connection: FOREIGN_KEY_CHECKS is per session.
	return db.Connection(func(conn *gorm.DB) error {
		var tables []string
		err := conn.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").
			Scan(&tables).Error
		if err != nil {
			return err
		}
		if err := conn.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
			return err
		}
		defer conn.Exec("SET FOREIGN_KEY_CHECKS = 1")
		for _, table := range tables {
			if err := conn.Exec("DROP TABLE IF EXISTS `" + table + "`").Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// requireLocalMySQL returns an error unless the DSN connects to this machine.
func requireLocalMySQL(dsn string) error {
	config, err := ParseMySQLDSN(dsn)
	if err != nil {
		return err
	}
	if config.Net == "unix" {
		return nil
	}
	host, _, err := net.SplitHostPort(config.Addr)
	if err != nil {
		host = config.Addr
	}
	if !isLocalHost(host) {
		return fmt.Errorf("refusing to reset the database at %q: only a database on localhost can be reset", host)
	}
	return nil
}

func isLocalHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
