package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Service
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// NATS
	NATSUrl string `mapstructure:"NATS_URL"`

	// Logging
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// Business Rules
	AutoGenerateCodes bool `mapstructure:"AUTO_GENERATE_CODES"`

	// Cache TTL
	CacheTTL time.Duration `mapstructure:"CACHE_TTL"`
}

func Load() (*Config, error) {
	viper.SetDefault("SERVICE_NAME", "master-data-service")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", "8083")
	viper.SetDefault("GRPC_PORT", "9083")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "master_data_db")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "redis123")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("NATS_URL", "nats://localhost:4222")

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("AUTO_GENERATE_CODES", true)
	viper.SetDefault("CACHE_TTL", 15*time.Minute)

	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
