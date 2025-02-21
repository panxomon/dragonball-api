package server

import (
	"dragonball-test/cmd/config"
	"dragonball-test/internal/commons/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	ServiceName = os.Getenv("Dragonball-test")
)

func NewHttpServer(cfg *config.Config, logger logging.Logger, artifactResources model.ArtifactResources) *gin.Engine {
	app := gin.New()
	app.Use(
		otelgin.Middleware(ServiceName),
		commonsMiddlewares.Logger(logger, &artifactResources, cfg.EnablePayloadLogging),
		middlewares.ErrorHandler(),
		middlewares.HandlerPanic(),
		middlewares.CustomRecoveryGinPanic(),
		commonsMiddlewares.HeaderValidationV2(),
		commonsMiddlewares.SanitizeRequest(),
	)
	return app
}
