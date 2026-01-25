package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/config"
	"github.com/erp-cosmetics/auth-service/internal/delivery/http"
	"github.com/erp-cosmetics/auth-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/auth-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/auth-service/internal/infrastructure/persistence/redis"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth"
	"github.com/erp-cosmetics/auth-service/internal/usecase/permission"
	
	"github.com/erp-cosmetics/shared/pkg/database"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/erp-cosmetics/shared/pkg/logger"
	"github.com/erp-cosmetics/shared/pkg/nats"
	
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

	logger.Info("Starting Auth Service",
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment),
	)

	// Connect to PostgreSQL
	dbConfig := database.NewDefaultConfig(cfg.GetDSN())
	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL")

	// Connect to Redis
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Connected to Redis")

	// Connect to NATS
	natsClient, err := nats.NewClient(&nats.Config{
		URL:            cfg.NATSUrl,
		MaxReconnects:  10,
		ReconnectWait:  2 * time.Second,
		Logger:         logger.Get(),
	})
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsClient.Close()
	logger.Info("Connected to NATS")

	// Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWTSecret,
		cfg.AccessTokenExpire,
		cfg.RefreshTokenExpire,
	)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permRepo := postgres.NewPermissionRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
	cacheRepo := redis.NewCacheRepository(rdb)

	// Initialize event publisher
	eventPublisher := event.NewPublisher(natsClient, logger.Get())

	// Initialize use cases
	loginUC := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPublisher)
	logoutUC := auth.NewLogoutUseCase(userRepo, tokenRepo, cacheRepo, jwtManager, eventPublisher)
	refreshUC := auth.NewTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)
	_ = permission.NewGetPermissionsUseCase(permRepo, cacheRepo) // Will be used later

	// Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(loginUC, logoutUC, refreshUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(authHandler, healthHandler, logger.Get())
	httpServer := router.Setup()

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		logger.Info("HTTP server listening", zap.String("address", addr))
		if err := httpServer.Run(addr); err != nil {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	// Close Redis
	rdb.Close()

	logger.Info("Server stopped")
}
