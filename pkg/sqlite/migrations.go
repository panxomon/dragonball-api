package sqlite

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// RunMigrations applies migrations to the database
// please visit the following link for more information: https://github.com/golang-migrate/migrate/blob/master/database/sqlite3/sqlite3_test.go
func RunMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Err(err).Msg("failed to create migration driver")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "sqlite3", driver)
	if err != nil {
		log.Err(err).Msg("failed to create migration instance")
		return err
	}

	if err = m.Up(); err != nil {
		log.Err(err).Msg("failed to run migrations")
	}

	log.Info().Msg("migrations has been applied (if exists) ✈️")
	return nil
}
