package database

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dsn string = os.Getenv("DATABASE_DSN")
var development bool = os.Getenv("DEV_MODE") == "true"

// DatabaseConnection opens SQLite when DEV_MODE=true (DATABASE_DSN is the
// file path), else MySQL (DATABASE_DSN is a MySQL DSN or URL, see
// ParseMySQLDSN). Schema setup is separate: see Migrate in migrate.go.
func DatabaseConnection() *gorm.DB {
	var db *gorm.DB
	var err error

	if development {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	} else {
		db, err = OpenMySQL(dsn, &gorm.Config{TranslateError: true})
	}

	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		panic(err)
	}

	return db
}

// datetimePrecision keeps microseconds in DATETIME columns (GORM's MySQL
// default is milliseconds): pagination and the tests compare stored times.
var datetimePrecision = 6

// OpenMySQL opens the MySQL database at dsn (see ParseMySQLDSN).
func OpenMySQL(dsn string, config *gorm.Config) (*gorm.DB, error) {
	dsnConfig, err := ParseMySQLDSN(dsn)
	if err != nil {
		return nil, err
	}
	return gorm.Open(mysql.New(mysql.Config{
		DSNConfig:                dsnConfig,
		DefaultDatetimePrecision: &datetimePrecision,
	}), config)
}

// ParseMySQLDSN reads a Go MySQL driver DSN
// ("user:pass@tcp(host:3306)/db?tls=true") or a URL
// ("mysql://user:pass@host:3306/db?tls=true"), the form most hosts give, and
// sets what the app relies on whatever the DSN says:
//   - times are stored and read as UTC (the services convert to time.Local
//     where the date matters);
//   - DATETIME columns are scanned into time.Time;
//   - RowsAffected counts matched rows, not only changed ones, as on SQLite.
func ParseMySQLDSN(dsn string) (*gomysql.Config, error) {
	if strings.HasPrefix(dsn, "mysql://") {
		converted, err := mysqlURLToDSN(dsn)
		if err != nil {
			return nil, err
		}
		dsn = converted
	}
	config, err := gomysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing DATABASE_DSN: %w", err)
	}
	config.ParseTime = true
	config.Loc = time.UTC
	config.ClientFoundRows = true
	return config, nil
}

func mysqlURLToDSN(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing DATABASE_DSN: %w", err)
	}
	config := gomysql.NewConfig()
	config.User = u.User.Username()
	config.Passwd, _ = u.User.Password()
	config.Net = "tcp"
	config.Addr = u.Host
	if u.Port() == "" {
		config.Addr = net.JoinHostPort(u.Hostname(), "3306")
	}
	config.DBName = strings.TrimPrefix(u.Path, "/")
	dsn := config.FormatDSN()
	if u.RawQuery != "" {
		dsn += "?" + u.RawQuery
	}
	return dsn, nil
}
