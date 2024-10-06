package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"path/filepath"
	"runtime"
)

type (
	Config struct {
		App  `json:"app"`
		HTTP `json:"rest"`
		DB   `json:"db"`
		Log  `json:"logger"`
	}

	App struct {
		Name    string `env-required:"false" json:"name"    env:"APP_NAME"`
		Version string `env-required:"false" json:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"false" json:"port" env:"HTTP_PORT"`
	}

	DB struct {
		DBHost     string `env-required:"false" json:"host"     env:"DB_HOST"`
		DBPort     int    `env-required:"false" json:"port"     env:"DB_PORT"`
		DBUser     string `env-required:"false" json:"user"     env:"DB_USER"`
		DBPassword string `env-required:"true"  json:"password" env:"DB_PASSWORD"`
		DBName     string `env-required:"false" json:"name"     env:"DB_NAME"`
		PoolMax    int32  `env-required:"true"  json:"pool_max" env:"PG_POOL_MAX"`
	}

	Log struct {
		Level slog.Level `env-required:"false" json:"level"   env:"LOG_LEVEL"`
	}
)

// LoadConfig returns api config.
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	configPath := filepath.Join(basePath, "config.json")

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("env read error: %w", err)
	}

	return cfg, nil
}
