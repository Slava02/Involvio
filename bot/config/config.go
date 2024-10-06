package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
	"runtime"
)

type (
	Config struct {
		BotConfig   `json:"bot_config"`
		CacheConfig `json:"cache_config"`
	}

	BotConfig struct {
		Token          string `json:"token"`
		RequestTimeout int    `json:"request_timeout"`
		UpdateTimeout  int64  `json:"update_timeout"`
	}

	CacheConfig struct {
		Address  string `json:"address"`
		DB       int    `json:"db"`
		Password string `json:"password"`
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
