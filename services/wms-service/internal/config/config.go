package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	LogFormat   string `mapstructure:"LOG_FORMAT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	NATSUrl string `mapstructure:"NATS_URL"`

	JWTSecret string `mapstructure:"JWT_SECRET"`

	// WMS Specific
	EnableFEFO             bool   `mapstructure:"ENABLE_FEFO"`
	ExpiryAlertDays        string `mapstructure:"EXPIRY_ALERT_DAYS"`
	LowStockCheckInterval  string `mapstructure:"LOW_STOCK_CHECK_INTERVAL"`
	ColdStorageMinTemp     int    `mapstructure:"COLD_STORAGE_MIN_TEMP"`
	ColdStorageMaxTemp     int    `mapstructure:"COLD_STORAGE_MAX_TEMP"`
}

// Load loads configuration
func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("SERVICE_NAME", "wms-service")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", "8086")
	viper.SetDefault("GRPC_PORT", "9086")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "console")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "wms_db")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("NATS_URL", "nats://localhost:4222")

	viper.SetDefault("ENABLE_FEFO", true)
	viper.SetDefault("EXPIRY_ALERT_DAYS", "90,30,7")
	viper.SetDefault("LOW_STOCK_CHECK_INTERVAL", "1h")
	viper.SetDefault("COLD_STORAGE_MIN_TEMP", 2)
	viper.SetDefault("COLD_STORAGE_MAX_TEMP", 8)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
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

// GetNATSReconnectWait returns NATS reconnect wait
func (c *Config) GetNATSReconnectWait() time.Duration {
	return 2 * time.Second
}
