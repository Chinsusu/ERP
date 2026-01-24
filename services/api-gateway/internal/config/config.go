package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// Server
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// JWT
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// Rate Limiting
	RateLimitEnabled bool `mapstructure:"RATE_LIMIT_ENABLED"`
	RateLimitPerMin  int  `mapstructure:"RATE_LIMIT_PER_MIN"`
	RateLimitBurst   int  `mapstructure:"RATE_LIMIT_BURST"`
	IPRateLimitPerMin int `mapstructure:"IP_RATE_LIMIT_PER_MIN"`

	// Circuit Breaker
	CircuitBreakerEnabled   bool          `mapstructure:"CIRCUIT_BREAKER_ENABLED"`
	CircuitBreakerThreshold int           `mapstructure:"CIRCUIT_BREAKER_THRESHOLD"`
	CircuitBreakerTimeout   time.Duration `mapstructure:"CIRCUIT_BREAKER_TIMEOUT"`

	// Timeouts
	DefaultTimeout time.Duration `mapstructure:"DEFAULT_TIMEOUT"`

	// CORS
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	// Logging
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// Routes
	Routes []RouteConfig `mapstructure:"-"`
}

type RouteConfig struct {
	Prefix       string        `yaml:"prefix"`
	Service      string        `yaml:"service"`
	AuthRequired bool          `yaml:"auth_required"`
	RateLimit    string        `yaml:"rate_limit,omitempty"`
	Timeout      time.Duration `yaml:"timeout,omitempty"`
}

type RoutesFile struct {
	Routes []RouteConfig `yaml:"routes"`
}

func Load() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "redis123")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")

	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_PER_MIN", 100)
	viper.SetDefault("RATE_LIMIT_BURST", 20)
	viper.SetDefault("IP_RATE_LIMIT_PER_MIN", 30)

	viper.SetDefault("CIRCUIT_BREAKER_ENABLED", true)
	viper.SetDefault("CIRCUIT_BREAKER_THRESHOLD", 5)
	viper.SetDefault("CIRCUIT_BREAKER_TIMEOUT", 30*time.Second)

	viper.SetDefault("DEFAULT_TIMEOUT", 30*time.Second)

	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://erp.company.com")

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load routes from YAML file
	routes, err := LoadRoutes("config/routes.yaml")
	if err != nil {
		// Use default routes if file not found
		cfg.Routes = DefaultRoutes()
	} else {
		cfg.Routes = routes
	}

	return cfg, nil
}

func LoadRoutes(path string) ([]RouteConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var routesFile RoutesFile
	if err := yaml.Unmarshal(data, &routesFile); err != nil {
		return nil, err
	}

	return routesFile.Routes, nil
}

func DefaultRoutes() []RouteConfig {
	return []RouteConfig{
		// Auth Service (no auth required)
		{Prefix: "/api/v1/auth", Service: "auth-service:8081", AuthRequired: false, RateLimit: "30/min"},

		// User Service
		{Prefix: "/api/v1/users", Service: "user-service:8082", AuthRequired: true},
		{Prefix: "/api/v1/departments", Service: "user-service:8082", AuthRequired: true},

		// Master Data Service
		{Prefix: "/api/v1/categories", Service: "master-data-service:8083", AuthRequired: true},
		{Prefix: "/api/v1/units", Service: "master-data-service:8083", AuthRequired: true},
		{Prefix: "/api/v1/materials", Service: "master-data-service:8083", AuthRequired: true},
		{Prefix: "/api/v1/products", Service: "master-data-service:8083", AuthRequired: true},

		// Supplier Service (Phase 2.1)
		{Prefix: "/api/v1/suppliers", Service: "supplier-service:8084", AuthRequired: true},
		{Prefix: "/api/v1/certifications", Service: "supplier-service:8084", AuthRequired: true},

		// Procurement Service (Phase 2.2)
		{Prefix: "/api/v1/purchase-requisitions", Service: "procurement-service:8085", AuthRequired: true},
		{Prefix: "/api/v1/purchase-orders", Service: "procurement-service:8085", AuthRequired: true},

		// Future Services (placeholders)
		{Prefix: "/api/v1/warehouse", Service: "wms-service:8086", AuthRequired: true},
		{Prefix: "/api/v1/manufacturing", Service: "manufacturing-service:8087", AuthRequired: true},
		{Prefix: "/api/v1/sales", Service: "sales-service:8088", AuthRequired: true},
		{Prefix: "/api/v1/marketing", Service: "marketing-service:8089", AuthRequired: true},
		{Prefix: "/api/v1/notifications", Service: "notification-service:8090", AuthRequired: true},
		{Prefix: "/api/v1/files", Service: "file-service:8091", AuthRequired: true},
		{Prefix: "/api/v1/reports", Service: "reporting-service:8092", AuthRequired: true},
	}
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
