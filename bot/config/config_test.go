package config

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set up example environment variables
	os.Setenv("APP_NAME", "example-api")
	os.Setenv("APP_VERSION", "1.0.0")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("PG_POOL_MAX", "10")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("TRACING_URL", "http://localhost:14268/api/traces")

	// Load the configuration
	cfg, err := LoadConfig()

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the configuration was loaded correctly
	assert.NotNil(t, cfg)
	assert.Equal(t, "example-api", cfg.App.Name)
	assert.Equal(t, "1.0.0", cfg.App.Version)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "localhost", cfg.DB.DBHost)
	assert.Equal(t, 5432, cfg.DB.DBPort)
	assert.Equal(t, "testuser", cfg.DB.DBUser)
	assert.Equal(t, "testpass", cfg.DB.DBPassword)
	assert.Equal(t, "testdb", cfg.DB.DBName)
	assert.Equal(t, int32(10), cfg.DB.PoolMax)
	assert.Equal(t, slog.LevelInfo, cfg.Log.Level)
}

func TestLoadConfigMissingRequiredField(t *testing.T) {
	// Clear all environment variables
	os.Clearenv()

	// Set all required fields except DB_PASSWORD
	os.Setenv("PG_POOL_MAX", "10")

	// Attempt to load the configuration
	_, err := LoadConfig()

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DBPassword")
}
