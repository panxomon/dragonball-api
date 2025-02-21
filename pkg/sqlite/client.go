package sqlite

import (
	"dragonball-test/internal/character/domain"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateDBClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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

	err = db.AutoMigrate(&domain.Character{})
	if err != nil {
		log.Err(err).Msg("failed to migrate the database")
		return nil, err
	}

	log.Info().Msg("sqlite connection established 💾")
	return db, nil
}
