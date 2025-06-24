package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
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
