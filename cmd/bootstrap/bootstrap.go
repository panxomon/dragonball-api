package bootstrap

import (
	"dragonball-test/internal/character/application"
	"dragonball-test/internal/character/application/create"
	"dragonball-test/internal/character/infrastructure/repository"
	"gorm.io/gorm"
)

type Bootstrap struct {
	Characters *application.App
}

func LoadComponents(db *gorm.DB, urlCharactersServer string) (*Bootstrap, error) {

	characterRepository := repository.NewCharacterRepository(db, urlCharactersServer)
	characterService := create.NewCreateCharacterService(characterRepository)
	handler := create.NewCharacterCommandHandler(characterService)

	charactersApp := application.NewApp(
		application.Commands{CreateCharacter: handler},
	)

	return &Bootstrap{
		Characters: charactersApp,
	}, nil

}
