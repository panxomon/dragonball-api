package domain

import (
	"context"
)

type Character struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Race  string `json:"race"`
	Image string `json:"image"`
}

type CharacterResponse struct {
	Status string      `json:"status"`
	Data   []Character `json:"data"`
}

type CharacterService interface {
	CreateCharacter(ctx context.Context, name string) (*Character, error)
}

type CharacterRepository interface {
	FindByName(ctx context.Context, name string) (*Character, error)
	Save(ctx context.Context, character *Character) error
}
