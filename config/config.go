package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"   env:"PG_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
