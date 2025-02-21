package create

import (
	"context"

	"dragonball-test/internal/character/domain"
	"dragonball-test/internal/cqrs"
)

type CreateCharacterCommandHandler cqrs.CommandWithResultHandler[CreateCharacter, *domain.Character]

type createCharacterCommandHandler struct {
	creator domain.CharacterService
}

func NewCharacterCommandHandler(creator domain.CharacterService) CreateCharacterCommandHandler {
	return &createCharacterCommandHandler{creator: creator}
}

func (h *createCharacterCommandHandler) Handle(ctx context.Context, command CreateCharacter) (*domain.Character, error) {
	return h.creator.CreateCharacter(ctx, command.Name)
}
