package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// Config the base config structure.
type Config struct {
	Log        ConfLog
	API        ConfAPI
	Monitoring ConfMonitoring
	Postgres   ConfPostgres
	Redis      ConfRedis
	Auth       ConfAuth
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

// ConfMonitoring monitoring related configurations.
type ConfMonitoring struct {
	Port string `env:"MONITORING_PORT" env-default:"8086"`
}

// ConfPostgres Postgres DB related configurations.
type ConfPostgres struct {
	Host                string        `env:"DB_HOST" env-default:"localhost"`
	Port                string        `env:"DB_PORT" env-default:"5432"`
	DbName              string        `env:"DB_NAME" env-default:"springfield-bank-dev"`
	User                string        `env:"DB_USER" env-default:"postgres"`
	Password            string        `env:"DB_PASSWORD" env-default:"postgres"`
	SslMode             string        `env:"DB_SSL_MODE" env-default:"prefer"`
	PoolMaxConn         int32         `env:"DB_POOL_MAX_CONN" env-default:"5"`
	PoolMaxConnLifetime time.Duration `env:"DB_POOL_MAX_CONN_LIFETIME" env-default:"5m"`
	Migrate             bool          `env:"DB_MIGRATE" env-default:"true"`
}

// ConfRedis Redis related configurations.
type ConfRedis struct {
	Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:"Redis2021!"`
}

// ConfAuth Authentication related configurations.
type ConfAuth struct {
	SecretKey      string        `env:"AUTH_SECRET_KEY" env-default:"YOU-SHOULD-CHANGE-ME"`
	AccessTokenDur time.Duration `env:"AUTH_ACCESS_TOKEN_DURATION" env-default:"15m"`
}

// GetDSN returns the database DSN, also known as Keyword/Value Connection String.
func (c ConfPostgres) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s pool_max_conns=%d pool_max_conn_lifetime=%s sslmode=%s",
		c.Host,
		c.Port,
		c.DbName,
		c.User,
		c.Password,
		c.PoolMaxConn,
		c.PoolMaxConnLifetime,
		c.SslMode,
	)
}

// ReadConfigFromFile reads configurations from OS env vars and file path and parses them into Config type.
// Supported extensions: .yaml, .yml, .json, .toml, .edn and .env.
//
// ATTENTION: OS env vars will take precedence over file vars.
func ReadConfigFromFile(filename string) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error reading configs file")
	}

	return &cfg
}

// ReadConfigFromEnv reads OS env variables and parses them into Config type.
func ReadConfigFromEnv() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error reading env")
	}

	return &cfg
}

// ReadConfig reads configurations from OS and file, if it exists, and parses them into Config type.
// Supported extensions: .yaml, .yml, .json, .toml, .edn and .env.
//
// ATTENTION: OS env vars will take precedence over file vars.
func ReadConfig(filename string) *Config {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Warn().Msgf("Config file not found %s", filename)
		return ReadConfigFromEnv()
	}

	return ReadConfigFromFile(filename)
}
