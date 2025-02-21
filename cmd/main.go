package main

import (
	router "dragonball-test/cmd/api"
	"dragonball-test/cmd/bootstrap"
	"dragonball-test/config"
	"dragonball-test/pkg/sqlite"

	"github.com/rs/zerolog/log"
)

func main() {
	env, err := config.LoadEnvVars()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env vars")
	}

	db, err := sqlite.CreateDBClient(env.DatabaseConnection)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create db client")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql db")
	}

	// Cambiar la carga de componentes para personajes en lugar de días festivos
	components, err := bootstrap.LoadComponents(db, env.UrlDragonBall)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load components")
	}

	if err := sqlite.RunMigrations(sqlDB); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Configuración del enrutador con la base de la ruta y los nuevos componentes
	r := router.SetupRouter(env.BasePath, components)

	log.Info().Msg("Starting server on :8080")

	// Iniciar el servidor en la dirección configurada
	if err := r.Run(env.ServerAddress); err != nil {
		log.Fatal().Err(err).Msg("could not start server")
	}
}
