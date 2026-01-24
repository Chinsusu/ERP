package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
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

	// JWT
	JWTSecret             string `mapstructure:"JWT_SECRET"`
	JWTAccessTokenExpire  string `mapstructure:"JWT_ACCESS_TOKEN_EXPIRE"`
	JWTRefreshTokenExpire string `mapstructure:"JWT_REFRESH_TOKEN_EXPIRE"`

	// Logging
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// CORS
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

// Load reads configuration from environment variables and config file
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("ENVIRONMENT", "development")
	v.SetDefault("PORT", "8080")
	v.SetDefault("GRPC_PORT", "9090")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	// Read from environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Optional: Read from config file if exists
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// Ignore error if config file doesn't exist
	_ = v.ReadInConfig()

	// Bind environment variables explicitly
	v.BindEnv("SERVICE_NAME")
	v.BindEnv("DB_HOST")
	v.BindEnv("DB_PORT")
	v.BindEnv("DB_USER")
	v.BindEnv("DB_PASSWORD")
	v.BindEnv("DB_NAME")
	v.BindEnv("DB_SSL_MODE")
	v.BindEnv("REDIS_HOST")
	v.BindEnv("REDIS_PORT")
	v.BindEnv("REDIS_PASSWORD")
	v.BindEnv("NATS_URL")
	v.BindEnv("JWT_SECRET")
	v.BindEnv("JWT_ACCESS_TOKEN_EXPIRE")
	v.BindEnv("JWT_REFRESH_TOKEN_EXPIRE")
	v.BindEnv("LOG_LEVEL")
	v.BindEnv("LOG_FORMAT")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// GetDSN returns PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
