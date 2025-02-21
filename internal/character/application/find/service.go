package find

import (
	"context"
	"dragonball-test/internal/character/domain"

	"github.com/rs/zerolog/log"
)

type CharacterFinderService struct {
	repository domain.CharacterRepository
}

func NewCharacterFinder(repository domain.CharacterRepository) *CharacterFinderService {
	return &CharacterFinderService{
		repository: repository,
	}
}

func (c *CharacterFinderService) Execute(ctx context.Context, name string) (*domain.Character, error) {
	// Buscar personaje
	character, err := c.repository.FindByName(ctx, name)
	if err == nil && character != nil {
		log.Ctx(ctx).Info().Str("character", name).Msg("Character found in database")
		return character, nil
	}

	// Si no está en la base, buscar en la API externa
	log.Ctx(ctx).Info().Str("character", name).Msg("Character not found in database, fetching from API")
	characterData, err := c.repository.Save(ctx, name)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("character", name).Msg("Failed to fetch character from API")
		return nil, err
	}

	// Guardar en la base de datos
	err = c.repository.Save(ctx, characterData)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("character", name).Msg("Failed to save character in database")
		return nil, err
	}

	return characterData, nil
}
