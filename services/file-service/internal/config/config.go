package config

import (
	"fmt"
	"os"

	"github.com/erp-cosmetics/shared/pkg/config"
)

// Config holds file service configuration
type Config struct {
	*config.Config

	// MinIO configuration
	MinIOEndpoint        string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinIORegion          string

	// Upload settings
	MaxUploadSize int64
}

// Load loads file service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Config:               baseConfig,
		MinIOEndpoint:        getEnvOrDefault("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKeyID:     getEnvOrDefault("MINIO_ACCESS_KEY_ID", "minioadmin"),
		MinIOSecretAccessKey: getEnvOrDefault("MINIO_SECRET_ACCESS_KEY", "minioadmin"),
		MinIOUseSSL:          getEnvAsBoolOrDefault("MINIO_USE_SSL", false),
		MinIORegion:          getEnvOrDefault("MINIO_REGION", "us-east-1"),
		MaxUploadSize:        getEnvAsInt64OrDefault("MAX_UPLOAD_SIZE", 52428800), // 50MB
	}, nil
}

// GetMinIOConfig returns MinIO configuration
func (c *Config) GetMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:        c.MinIOEndpoint,
		AccessKeyID:     c.MinIOAccessKeyID,
		SecretAccessKey: c.MinIOSecretAccessKey,
		UseSSL:          c.MinIOUseSSL,
		Region:          c.MinIORegion,
	}
}

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Region          string
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

func getEnvAsInt64OrDefault(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var res int64
	fmt.Sscanf(value, "%d", &res)
	return res
}
