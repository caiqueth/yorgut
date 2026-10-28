package database

import "testing"

func TestRequireLocalMySQL(t *testing.T) {
	cases := []struct {
		dsn   string
		local bool
	}{
		{"localdev:localdev@tcp(localhost:3306)/yorgut-dev", true},
		{"u:p@tcp(127.0.0.1:3306)/db", true},
		{"u:p@tcp([::1]:3306)/db", true},
		{"u:p@unix(/var/run/mysqld/mysqld.sock)/db", true},
		{"mysql://u:p@localhost/db", true},
		{"u:p@tcp(db.example.com:3306)/db?tls=true", false},
		{"mysql://u:p@mysql-abc.example.com:25060/db?tls=true", false},
		{"u:p@tcp(10.0.0.5:3306)/db", false},
		{"u:p@tcp(db)/db", false},
	}
	for _, c := range cases {
		err := requireLocalMySQL(c.dsn)
		if c.local && err != nil {
			t.Errorf("%q: expected local, got %v", c.dsn, err)
		}
		if !c.local && err == nil {
			t.Errorf("%q: expected refusal", c.dsn)
		}
	}
}

func TestParseMySQLDSN(t *testing.T) {
	cases := []struct {
		dsn, addr, user, password, dbName string
		tls                               bool
	}{
		{"u:p@tcp(db.example.com:3306)/yorgut", "db.example.com:3306", "u", "p", "yorgut", false},
		{"mysql://u:p%40ss@db.example.com:25060/yorgut?tls=true", "db.example.com:25060", "u", "p@ss", "yorgut", true},
		{"mysql://u:p@db.example.com/yorgut", "db.example.com:3306", "u", "p", "yorgut", false},
	}
	for _, c := range cases {
		config, err := ParseMySQLDSN(c.dsn)
		if err != nil {
			t.Errorf("%q: %v", c.dsn, err)
			continue
		}
		if config.Addr != c.addr || config.User != c.user || config.Passwd != c.password || config.DBName != c.dbName {
			t.Errorf("%q: got addr %q user %q password %q db %q", c.dsn, config.Addr, config.User, config.Passwd, config.DBName)
		}
		if (config.TLS != nil) != c.tls {
			t.Errorf("%q: TLS = %v, want %v", c.dsn, config.TLS != nil, c.tls)
		}
		if !config.ParseTime || !config.ClientFoundRows || config.Loc.String() != "UTC" {
			t.Errorf("%q: forced options not set: %+v", c.dsn, config)
		}
	}
}
