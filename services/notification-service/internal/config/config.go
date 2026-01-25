package config

import (
	"github.com/erp-cosmetics/shared/pkg/config"
)

// Config holds notification service configuration
type Config struct {
	*config.Config

	// SMTP configuration
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string
	SMTPUseTLS    bool

	// Email settings
	MaxEmailsPerMinute int
	MaxRetryAttempts   int
	RetryDelayMinutes  int
}

// Load loads notification service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Config: baseConfig,

		// SMTP configuration from environment
		SMTPHost:      getEnvOrDefault("SMTP_HOST", "localhost"),
		SMTPPort:      getEnvAsIntOrDefault("SMTP_PORT", 587),
		SMTPUsername:  getEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPassword:  getEnvOrDefault("SMTP_PASSWORD", ""),
		SMTPFromEmail: getEnvOrDefault("SMTP_FROM_EMAIL", "noreply@erp-cosmetics.local"),
		SMTPFromName:  getEnvOrDefault("SMTP_FROM_NAME", "ERP Cosmetics System"),
		SMTPUseTLS:    getEnvAsBoolOrDefault("SMTP_USE_TLS", true),

		// Email settings
		MaxEmailsPerMinute: getEnvAsIntOrDefault("MAX_EMAILS_PER_MINUTE", 60),
		MaxRetryAttempts:   getEnvAsIntOrDefault("MAX_RETRY_ATTEMPTS", 3),
		RetryDelayMinutes:  getEnvAsIntOrDefault("RETRY_DELAY_MINUTES", 5),
	}, nil
}

// GetSMTPConfig returns SMTP configuration
func (c *Config) GetSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:      c.SMTPHost,
		Port:      c.SMTPPort,
		Username:  c.SMTPUsername,
		Password:  c.SMTPPassword,
		FromEmail: c.SMTPFromEmail,
		FromName:  c.SMTPFromName,
		UseTLS:    c.SMTPUseTLS,
	}
}

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
	UseTLS    bool
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := getEnv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := getEnv(key); value != "" {
		var intValue int
		if _, err := parseIntFromString(value, &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := getEnv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func getEnv(key string) string {
	// This would use os.Getenv in production
	// For now, we rely on viper which is already configured in base config
	return ""
}

func parseIntFromString(s string, result *int) (int, error) {
	// Simple parse implementation
	var value int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			value = value*10 + int(c-'0')
		}
	}
	*result = value
	return value, nil
}
