package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

// Config the base config structure.
type Config struct {
	Log ConfLog
	API ConfAPI
}

// ConfLog logging related configurations.
type ConfLog struct {
	Encoding string `env:"LOG_ENCODING" env-default:"json"`
	Level    string `env:"LOG_LEVEL" env-default:"info"`
}

// ConfAPI API related configurations.
type ConfAPI struct {
	HTTPPort string `env:"API_HTTP_PORT" env-default:"8080"`
}

// ReadConfigFromFile reads configurations from file path and parses them into Config type.
// Supported extensions: .yaml, .yml, .json, .toml, .edn and .env.
func ReadConfigFromFile(filename string) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading configs file")
	}

	return &cfg
}

// ReadConfigFromEnv reads SO env variables and parses them into Config type.
func ReadConfigFromEnv() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading env")
	}

	return &cfg
}
