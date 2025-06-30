package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type AppConfig struct {
	Env string `env:"ENV" env-default:"development" yaml:"env"`
}

type GRPCConfig struct {
	Address string `env:"GRPC_ADDRESS" env-default:"40042" yaml:"address"`
}

type PostgresConfig struct {
	DSN string `env:"POSTGRES_DSN" yaml:"dsn"`
}

type JWTConfig struct {
	Secret     string        `env:"JWT_SECRET" yaml:"secret"`
	Expiration time.Duration `env:"JWT_EXPIRATION" yaml:"expiration"`
}

func NewConfig(configPath string, dotEnvPath string) (*Config, error) {
	if dotEnvPath != "" {
		if err := godotenv.Load(dotEnvPath); err != nil {
			return nil, err
		}
	}
	var cfg Config
	if _, err := os.Stat(configPath); err == nil {
		if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return &cfg, nil
}
