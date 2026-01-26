package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// NATS
	NATSURL string `mapstructure:"NATS_URL"`

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	// Logging
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// Business Rules
	AutoApproveMegaKOLs           bool    `mapstructure:"AUTO_APPROVE_MEGA_KOLS"`
	AutoApproveMacroKOLs          bool    `mapstructure:"AUTO_APPROVE_MACRO_KOLS"`
	SampleValueApprovalThreshold  float64 `mapstructure:"SAMPLE_VALUE_APPROVAL_THRESHOLD"`
	TrackKOLPosts                 bool    `mapstructure:"TRACK_KOL_POSTS"`
}

// LoadConfig loads configuration from environment
func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.HTTPPort == "" {
		config.HTTPPort = "8089"
	}
	if config.GRPCPort == "" {
		config.GRPCPort = "9089"
	}
	if config.DBSSLMode == "" {
		config.DBSSLMode = "disable"
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
