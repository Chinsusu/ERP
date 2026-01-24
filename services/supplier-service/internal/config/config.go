package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Service
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	LogFormat   string `mapstructure:"LOG_FORMAT"`

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

	// JWT (for validating requests from gateway)
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// Certificate Settings
	CertExpiryAlertDays  int  `mapstructure:"CERT_EXPIRY_ALERT_DAYS"`
	AutoBlockOnGMPExpiry bool `mapstructure:"AUTO_BLOCK_ON_GMP_EXPIRY"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("SERVICE_NAME", "supplier-service")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", "8084")
	viper.SetDefault("GRPC_PORT", "9084")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "supplier_db")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("NATS_URL", "nats://localhost:4222")

	viper.SetDefault("CERT_EXPIRY_ALERT_DAYS", 90)
	viper.SetDefault("AUTO_BLOCK_ON_GMP_EXPIRY", true)

	// Bind environment variables explicitly
	viper.BindEnv("SERVICE_NAME")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("PORT")
	viper.BindEnv("GRPC_PORT")
	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("LOG_FORMAT")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSL_MODE")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("REDIS_DB")
	viper.BindEnv("NATS_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("CERT_EXPIRY_ALERT_DAYS")
	viper.BindEnv("AUTO_BLOCK_ON_GMP_EXPIRY")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there is no config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// GetNATSReconnectWait returns the NATS reconnect wait duration
func (c *Config) GetNATSReconnectWait() time.Duration {
	return 2 * time.Second
}
