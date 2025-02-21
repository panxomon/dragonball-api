package postgres

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDBClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		log.Err(err).Msg("failed to open postgres connection")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Err(err).Msg("failed to get sql db")
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		log.Err(err).Msg("failed to ping postgres")
		return nil, err
	}

	log.Info().Msg("postgres connection established 💾")
	return db, nil
}
