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

// SetupRouter configura y devuelve el enrutador con las rutas necesarias.
func SetupRouter(basePath string, components *bootstrap.Bootstrap) *gin.Engine {
	r := gin.Default()

	// Ruta de HealthCheck
	r.GET("/health", HealthCheck)

	// Ruta para crear personajes
	r.POST("/characters", func(c *gin.Context) {
		// Aquí llamas al endpoint que maneja la creación de personajes
		endpoint.NewCharacterEndpoint(components.Characters).Invoke(c)
	})

	// Si es necesario configurar algo en el middlewares, puedes agregarlo aquí
	// Ejemplo: r.Use(BootstrapMiddleware(components))

	return r
}

// HealthCheck es un simple endpoint para verificar si la API está funcionando correctamente
func HealthCheck(c *gin.Context) {
	log.Info().Str("project", "healthcheck").Msg("Health check invoked")
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// Middleware para añadir el contexto de bootstrap si es necesario en el futuro
func BootstrapMiddleware(components *bootstrap.Bootstrap) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pasa los componentes del bootstrap al contexto
		ctx := context.WithValue(c.Request.Context(), BootstrapKey, components)
		log.Ctx(ctx).Info().Str("project", "endpoint").Msg("Invoking endpoint")
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
