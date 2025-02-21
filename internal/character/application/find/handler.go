package find

import (
	"context"
	"dragonball-test/internal/character/domain"
	"dragonball-test/internal/cqrs"
)

type CharacterFindQueryHandler cqrs.QueryHandler[CharacterFindQuery, *domain.Character]

type characterFindQueryHandler struct {
	finder *CharacterFinderService
}

func NewCharacterFindQueryHandler(finder *CharacterFinderService) CharacterFindQueryHandler {
	return &characterFindQueryHandler{finder: finder}
}

func (h *characterFindQueryHandler) Handle(ctx context.Context, query CharacterFindQuery) (*domain.Character, error) {
	return h.finder.Execute(ctx, query.Name)
}
