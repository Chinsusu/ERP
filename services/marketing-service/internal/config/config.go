package config

import (
	"github.com/erp-cosmetics/shared/pkg/config"
	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	*config.Config

	// Business Rules
	AutoApproveMegaKOLs           bool    `mapstructure:"AUTO_APPROVE_MEGA_KOLS"`
	AutoApproveMacroKOLs          bool    `mapstructure:"AUTO_APPROVE_MACRO_KOLS"`
	SampleValueApprovalThreshold  float64 `mapstructure:"SAMPLE_VALUE_APPROVAL_THRESHOLD"`
	TrackKOLPosts                 bool    `mapstructure:"TRACK_KOL_POSTS"`
}

// LoadConfig loads configuration from environment
func LoadConfig() (*Config, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	viper.BindEnv("AUTO_APPROVE_MEGA_KOLS")
	viper.BindEnv("AUTO_APPROVE_MACRO_KOLS")
	viper.BindEnv("SAMPLE_VALUE_APPROVAL_THRESHOLD")
	viper.BindEnv("TRACK_KOL_POSTS")

	cfg := &Config{
		Config: baseConfig,
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
