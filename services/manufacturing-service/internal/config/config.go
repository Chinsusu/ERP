package config

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration
type Config struct {
	ServiceName string
	Port        string
	GRPCPort    string
	LogLevel    string
	LogFormat   string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// NATS
	NATSUrl string

	// BOM Encryption
	BOMEncryptionKey []byte

	// WMS gRPC
	WMSGRPCAddress string
}

// Load loads configuration from environment
func Load() (*Config, error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("SERVICE_NAME", "manufacturing-service")
	viper.SetDefault("PORT", "8087")
	viper.SetDefault("GRPC_PORT", "9087")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "manufacturing_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("WMS_GRPC_ADDRESS", "localhost:9086")

	cfg := &Config{
		ServiceName:    viper.GetString("SERVICE_NAME"),
		Port:           viper.GetString("PORT"),
		GRPCPort:       viper.GetString("GRPC_PORT"),
		LogLevel:       viper.GetString("LOG_LEVEL"),
		LogFormat:      viper.GetString("LOG_FORMAT"),
		DBHost:         viper.GetString("DB_HOST"),
		DBPort:         viper.GetString("DB_PORT"),
		DBUser:         viper.GetString("DB_USER"),
		DBPassword:     viper.GetString("DB_PASSWORD"),
		DBName:         viper.GetString("DB_NAME"),
		DBSSLMode:      viper.GetString("DB_SSLMODE"),
		NATSUrl:        viper.GetString("NATS_URL"),
		WMSGRPCAddress: viper.GetString("WMS_GRPC_ADDRESS"),
	}

	// Load encryption key (32 bytes for AES-256)
	keyHex := os.Getenv("BOM_ENCRYPTION_KEY")
	if keyHex == "" {
		// Default key for development only - CHANGE IN PRODUCTION
		keyHex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	}
	
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid BOM_ENCRYPTION_KEY: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("BOM_ENCRYPTION_KEY must be 32 bytes (64 hex characters)")
	}
	cfg.BOMEncryptionKey = key

	return cfg, nil
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
