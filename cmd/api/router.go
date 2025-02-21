package api

import (
	"context"
	bootstrap "dragonball-test/cmd/bootstrap"
	"dragonball-test/internal/endpoint"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Key string

const BootstrapKey Key = "bootstrap"

func SetupRouter(basePath string, components *bootstrap.Bootstrap) *gin.Engine {
	r := gin.Default()

	r.GET("/health", HealthCheck)
	
	r.POST("/characters", func(c *gin.Context) {		
		endpoint.NewCharacterEndpoint(components.Characters).Invoke(c)
	})
	return r
}

func HealthCheck(c *gin.Context) {
	log.Info().Str("project", "healthcheck").Msg("Health check invoked")
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func BootstrapMiddleware(components *bootstrap.Bootstrap) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pasa los componentes del bootstrap al contexto
		ctx := context.WithValue(c.Request.Context(), BootstrapKey, components)
		log.Ctx(ctx).Info().Str("project", "endpoint").Msg("Invoking endpoint")
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
