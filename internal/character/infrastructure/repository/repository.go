package repository

import (
	"context"
	"dragonball-test/internal/character/domain"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

type CharacterRepository struct {
	db   *gorm.DB
	url  string
	data map[string]*domain.Character
}

func NewCharacterRepository(db *gorm.DB, urlDragonBall string) *CharacterRepository {
	return &CharacterRepository{
		db:  db,
		url: urlDragonBall,
	}
}

func (r *CharacterRepository) FindByName(ctx context.Context, name string) (*domain.Character, error) {

	if character, found := r.data[name]; found {
		return character, nil
	}

	var character domain.Character
	result := r.db.
		WithContext(ctx).
		Table("characters").
		Where("name = ?", name).
		First(&character)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Ctx(ctx).Info().Str("character", name).Msg("Character not found in database")
		} else {
			log.Ctx(ctx).Err(result.Error).Msg("Error al buscar en la base de datos")
			return nil, result.Error
		}
	}

	if result.RowsAffected > 0 {
		return &character, nil
	}

	log.Ctx(ctx).Info().Msg("Personaje no encontrado en base de datos, consultando API externa...")

	apiURL := r.url

	for {
		res, err := http.Get(apiURL)
		if err != nil {
			log.Ctx(ctx).Err(err).Str("project", "character").Msg("Error al obtener datos de la API")
			return nil, fmt.Errorf("error al obtener datos de la API: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			err := fmt.Errorf("Código de estado no OK recibido: %d", res.StatusCode)
			log.Ctx(ctx).Err(err).Str("project", "character").Int("status_code", res.StatusCode).Msg("Error al recibir respuesta de la API")
			return nil, err
		}

		var apiResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			log.Ctx(ctx).Err(err).Msg("Error al decodificar respuesta de la API")
			return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
		}

		meta, metaOk := apiResponse["meta"].(map[string]interface{})
		if !metaOk {
			err := fmt.Errorf("la respuesta de la API no contiene una clave 'meta' válida")
			log.Ctx(ctx).Err(err).Msg("Error al procesar la paginación")
			return nil, err
		}

		currentPage, ok := meta["currentPage"].(float64)
		if !ok {
			err := fmt.Errorf("la respuesta de la API no contiene un campo 'currentPage' válido")
			log.Ctx(ctx).Err(err).Msg("Error al procesar la página actual")
			return nil, err
		}

		totalPages, ok := meta["totalPages"].(float64)
		if !ok {
			err := fmt.Errorf("la respuesta de la API no contiene un campo 'totalPages' válido")
			log.Ctx(ctx).Err(err).Msg("Error al procesar el total de páginas")
			return nil, err
		}

		for _, characterData := range apiResponse["items"].([]interface{}) {
			characterMap := characterData.(map[string]interface{})
			if name == characterMap["name"].(string) {
				id := characterMap["id"].(float64)
				character := &domain.Character{
					ID:    uint(id),
					Name:  name,
					Race:  characterMap["race"].(string),
					Image: characterMap["image"].(string),
				}

				r.data[character.Name] = character
				return character, nil
			}
		}

		if currentPage >= totalPages {
			log.Ctx(ctx).Info().Str("character", name).Msg("Última página alcanzada sin encontrar el personaje")
			break
		}

		next, exists := meta["next"]
		if exists && next != nil {
			if nextStr, ok := next.(string); ok {
				apiURL = nextStr
			} else {
				err := fmt.Errorf("el valor de 'next' no es un string válido")
				log.Ctx(ctx).Err(err).Msg("Error en la paginación")
				return nil, err
			}
		} else {
			log.Ctx(ctx).Error().Msg("'next' no existe o es nil, finalizando búsqueda")
			break
		}
	}

	log.Ctx(ctx).Error().Msg(fmt.Sprintf("Personaje %s no encontrado en ninguna página", name))
	return nil, fmt.Errorf("personaje no encontrado: %s", name)

}

func (r *CharacterRepository) Save(ctx context.Context, character *domain.Character) error {
	result := r.db.WithContext(ctx).Table("characters").Save(character)
	if result.Error != nil {
		log.Ctx(ctx).Err(result.Error).Str("character", character.Name).Msg("Error al guardar personaje en la base de datos")
		return result.Error
	}
	if r.data == nil {
		r.data = make(map[string]*domain.Character)
	}
	r.data[character.Name] = character
	log.Ctx(ctx).Info().Str("character", character.Name).Msg("Character saved successfully in both database and memory cache")
	return nil
}
