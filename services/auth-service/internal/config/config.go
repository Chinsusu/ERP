package config

import (
	"time"

	"github.com/erp-cosmetics/shared/pkg/config"
	"github.com/erp-cosmetics/shared/pkg/jwt"
)

// Config holds auth service configuration
type Config struct {
	*config.Config
	
	// JWT specific
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
}

// Load loads auth service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Parse token expiration durations
	accessExpire, err := jwt.ParseDuration(baseConfig.JWTAccessTokenExpire)
	if err != nil {
		accessExpire = 15 * time.Minute // Default
	}

	refreshExpire, err := jwt.ParseDuration(baseConfig.JWTRefreshTokenExpire)
	if err != nil {
		refreshExpire = 7 * 24 * time.Hour // Default
	}

	return &Config{
		Config:             baseConfig,
		AccessTokenExpire:  accessExpire,
		RefreshTokenExpire: refreshExpire,
	}, nil
}
