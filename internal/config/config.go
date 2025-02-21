package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const configPath = "internal/config/config.yaml"

type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Server   Server   `yaml:"server"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASSWORD"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("NewConfig (1): %w", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("NewConfig (2): %w", err)
	}
	return &cfg, nil
}
