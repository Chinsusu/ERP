package config

import (
	"github.com/erp-cosmetics/shared/pkg/config"
)

// Config holds user service configuration
type Config struct {
	*config.Config
	
	// Auth Service
	AuthServiceAddr string
}

// Load loads user service configuration
func Load() (*Config, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Config:          baseConfig,
		AuthServiceAddr: "localhost:9081", // Default auth service gRPC address
	}, nil
}
