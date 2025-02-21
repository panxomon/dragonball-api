package create

import (
	"context"
	"dragonball-test/internal/character/domain"
	"github.com/rs/zerolog/log"
)

type CreateCharacterService struct {
	repository domain.CharacterRepository
}

func NewCreateCharacterService(r domain.CharacterRepository) *CreateCharacterService {
	return &CreateCharacterService{
		repository: r,
	}
}

func (s *CreateCharacterService) CreateCharacter(ctx context.Context, name string) (*domain.Character, error) {

	character, err := s.repository.FindByName(ctx, name)

	if err != nil {
		log.Ctx(ctx).Info().Str("character", name).Msg("Error: found character")
		return nil, err
	}

	if character != nil {
		log.Ctx(ctx).Info().Str("character", name).Msg("Character: found in database")
		return character, nil
	}

	character = &domain.Character{
		Name: name, // Asignas el nombre que recibiste
		// Otros campos necesarios para crear el personaje
	}

	err = s.repository.Save(ctx, character)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("character", name).Msg("Failed to save character in database")
		return nil, err
	}

	return character, nil
}
