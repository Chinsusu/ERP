package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/erp-cosmetics/api-gateway/internal/config"
	"github.com/erp-cosmetics/api-gateway/internal/router"
	"github.com/erp-cosmetics/shared/pkg/logger"

	redisClient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	if err := logger.Init(&logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	}); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting API Gateway",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment),
	)

	// Connect to Redis
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Failed to connect to Redis, rate limiting disabled", zap.Error(err))
		rdb = nil
	} else {
		logger.Info("Connected to Redis")
	}

	// Setup router
	r := router.NewRouter(cfg, logger.Get(), rdb)
	
	// Start health checks
	r.StartHealthChecks()

	// Setup HTTP server
	engine := r.Setup()

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		logger.Info("API Gateway listening", zap.String("address", addr))
		if err := engine.Run(addr); err != nil {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down API Gateway...")

	// Close Redis
	if rdb != nil {
		rdb.Close()
	}

	logger.Info("API Gateway stopped")
}
