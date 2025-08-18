package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServerConfig struct {
	Port       string `yaml:"port" env:"PORT" env-default:"8080" validate:"required,numeric"`
	MaxRetries int    `yaml:"max_retries" env:"MAX_RETRIES" env-default:"5" validate:"gte=1"`
	RetryDelay int    `yaml:"retry_delay" env:"RETRY_DELAY" env-default:"5" validate:"gte=1"`
}

type PostgresConfig struct {
	Host       string `yaml:"host" env:"PG_HOST" validate:"required"`
	Port       string `yaml:"port" env:"PG_PORT" env-default:"5432" validate:"required,numeric"`
	User       string `yaml:"user" env:"PG_USER" env-default:"postgres" validate:"required"`
	Password   string `yaml:"password" env:"PG_PASSWORD" env-default:""`
	DBName     string `yaml:"dbname" env:"PG_DBNAME" validate:"required"`
	SSLMode    string `yaml:"sslmode" env:"PG_SSLMODE" env-default:"disable" validate:"oneof=disable require"`
	MaxConns   int32  `yaml:"max_conns" env:"PG_MAX_CONNS" env-default:"50" validate:"gte=1"`
	MinConns   int32  `yaml:"min_conns" env:"PG_MIN_CONNS" env-default:"10" validate:"gte=1"`
	Timeout    int    `yaml:"timeout" env:"PG_TIMEOUT" env-default:"5" validate:"gte=1"`
	MaxRetries int    `yaml:"max_retries" env:"PG_MAX_RETRIES" env-default:"5" validate:"gte=1"`
	RetryDelay int    `yaml:"retry_delay" env:"PG_RETRY_DELAY" env-default:"2" validate:"gte=1"`
}

type JWT struct {
	Secret        string        `env:"JWT_SECRET" validate:"required"`
	TokenExpiry   time.Duration `env:"JWT_TOKEN_EXPIRY" env-default:"1h"`
	RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" env-default:"24h"`
}

type Config struct {
	Env        string           `env:"ENV" env-default:"development" validate:"oneof=development production"`
	JWT        JWT              `env-prefix:"JWT_"`
	HTTPServer HTTPServerConfig `env-prefix:"HTTP_SERVER_"`
	Postgres   PostgresConfig   `env-prefix:"POSTGRES_"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config from env: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &cfg, nil
}
