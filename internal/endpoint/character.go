package endpoint

import (
	dragonball "dragonball-test/internal/character/application"
	"dragonball-test/internal/character/application/create"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CharacterEndpoint struct {
	characterApp *dragonball.App
}

func NewCharacterEndpoint(characterApp *dragonball.App) *CharacterEndpoint {
	return &CharacterEndpoint{
		characterApp: characterApp,
	}
}

func (ep *CharacterEndpoint) Invoke(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	character, err := ep.characterApp.Commands.CreateCharacter.Handle(ctx, create.CreateCharacter{Name: request.Name})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch character")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch character"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
}
