package config

import (
	"github.com/erp-cosmetics/shared/pkg/config"
)

// Config holds reporting service configuration
type Config struct {
	*config.Config

	// Export settings
	ExportPath       string
	MaxExportRows    int
	ExportRetention  int // Days to keep exports

	// Stats cache
	StatsCacheTTL    int // Seconds
}

// Load loads reporting service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Config:          baseConfig,
		ExportPath:      getEnvOrDefault("EXPORT_PATH", "/tmp/exports"),
		MaxExportRows:   getEnvAsIntOrDefault("MAX_EXPORT_ROWS", 100000),
		ExportRetention: getEnvAsIntOrDefault("EXPORT_RETENTION_DAYS", 7),
		StatsCacheTTL:   getEnvAsIntOrDefault("STATS_CACHE_TTL", 300),
	}, nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	// In production, use os.Getenv
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	// In production, parse from os.Getenv
	return defaultValue
}
