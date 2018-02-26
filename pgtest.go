package pgtest

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var Root = os.ExpandEnv("${HOME}/.pgtest")
var Version = "10.2-1"

var isLinuxSystem = runtime.GOOS == "linux"

type SetupFunc func(db *sql.DB) error

type TestFunc func(db *sql.DB)

func WithDatabase(t *testing.T, setup SetupFunc, test TestFunc) {
	withCurrentT(t, func() {
		if err := PreparePostgresInstallation(Root, Version, isLinuxSystem); err != nil {
			t.Fatalf("Could not prepare postgres installation: %s", err)
			return
		}

		config := postgresConfig{
			Binary:   filepath.Join(Root, Version, "unpacked/pgsql/bin/postgres"),
			Snapshot: filepath.Join(Root, Version, "initdb/pgdata"),
		}

		pg, err := startPostgresInstance(config)
		if err != nil {
			t.Fatalf("Could not start postgres instance: %s", err)
			return
		}

		defer pg.Close()

		db, err := connect(pg.URL)
		if err != nil {
			t.Fatalf("Could not open a database connection to postgres at %s: %s", pg.URL, err)
			return
		}

		defer db.Close()

		if err := setup(db); err != nil {
			t.Fatalf("Database setup failed: %s", err)
			return
		}

		test(db)
	})
}

func NoSetup(*sql.DB) error {
	return nil
}
