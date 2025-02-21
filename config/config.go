package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	LogLevel           string        `required:"false" split_words:"true" default:":info"`
	ServerAddress      string        `required:"false" split_words:"true" default:":8080"`
	MsName             string        `required:"false" split_words:"true" default:"dragonball-test"`
	MsVersion          string        `required:"false" split_words:"true" default:"v1.0.0"`
	BasePath           string        `required:"false" split_words:"true" default:"/v1/characters"`
	ContextDeadline    time.Duration `required:"false" split_words:"true" default:"30s"`
	Environment        string        `required:"true" split_words:"true"`
	DatabaseConnection string        `required:"true" split_words:"true"`
	UrlDragonBall      string        `required:"true" split_words:"true"`
}

func LoadEnvVars() (Environment, error) {
	var env Environment

	godotenv.Load()

	if err := envconfig.Process("", &env); err != nil {
		return Environment{}, fmt.Errorf("error loading env vars: %w", err)
	}

	return env, nil
}
