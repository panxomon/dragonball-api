package application

import "dragonball-test/internal/character/application/create"

type App struct {
	Commands Commands
}

func NewApp(commands Commands) *App {
	return &App{commands}
}

type Commands struct {
	CreateCharacter create.CreateCharacterCommandHandler
}
