package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	// Server
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`
	ServiceName string `mapstructure:"SERVICE_NAME"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// NATS
	NatsURL string `mapstructure:"NATS_URL"`

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	// Logging
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// Business Rules
	EnableCreditCheck    bool `mapstructure:"ENABLE_CREDIT_CHECK"`
	AutoReserveOnConfirm bool `mapstructure:"AUTO_RESERVE_ON_CONFIRM"`
	AllowNegativeStock   bool `mapstructure:"ALLOW_NEGATIVE_STOCK"`
}

// LoadConfig loads configuration from environment
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("HTTP_PORT", "8088")
	viper.SetDefault("GRPC_PORT", "9088")
	viper.SetDefault("SERVICE_NAME", "sales-service")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "sales_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "console")
	viper.SetDefault("ENABLE_CREDIT_CHECK", true)
	viper.SetDefault("AUTO_RESERVE_ON_CONFIRM", true)
	viper.SetDefault("ALLOW_NEGATIVE_STOCK", false)

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, using defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}
