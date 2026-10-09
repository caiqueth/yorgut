package database

import "testing"

func TestRequireLocalPostgres(t *testing.T) {
	t.Setenv("PGHOST", "")
	cases := []struct {
		dsn   string
		local bool
	}{
		{"postgres://localdev:localdev@localhost:5432/yorgut-dev?sslmode=disable", true},
		{"postgres://u:p@127.0.0.1:5432/db", true},
		{"postgres://u:p@[::1]:5432/db", true},
		{"host=localhost user=u password=p dbname=db port=5432 sslmode=disable", true},
		{"host=/var/run/postgresql dbname=db", true},
		{"postgres://u:p@dpg-abc123-a.oregon-postgres.render.com/db", false},
		{"postgres://u:p@dpg-abc123-a/db", false},
		{"host=localhost,db.example.com dbname=db", false},
		{"host=10.0.0.5 dbname=db", false},
	}
	for _, c := range cases {
		err := requireLocalPostgres(c.dsn)
		if c.local && err != nil {
			t.Errorf("%q: expected local, got %v", c.dsn, err)
		}
		if !c.local && err == nil {
			t.Errorf("%q: expected refusal", c.dsn)
		}
	}
}
