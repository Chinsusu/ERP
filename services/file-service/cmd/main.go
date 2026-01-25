package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/file-service/internal/config"
	"github.com/erp-cosmetics/file-service/internal/delivery/http"
	"github.com/erp-cosmetics/file-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/file-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/file-service/internal/infrastructure/storage"
	"github.com/erp-cosmetics/file-service/internal/usecase/file"

	"github.com/erp-cosmetics/shared/pkg/database"
	"github.com/erp-cosmetics/shared/pkg/logger"

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

	logger.Info("Starting File Service",
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

	// Connect to MinIO
	minioCfg := cfg.GetMinIOConfig()
	storageClient, err := storage.NewMinIOClient(&storage.Config{
		Endpoint:        minioCfg.Endpoint,
		AccessKeyID:     minioCfg.AccessKeyID,
		SecretAccessKey: minioCfg.SecretAccessKey,
		UseSSL:          minioCfg.UseSSL,
		Region:          minioCfg.Region,
	}, logger.Get())
	if err != nil {
		logger.Fatal("Failed to connect to MinIO", zap.Error(err))
	}
	logger.Info("Connected to MinIO", zap.String("endpoint", minioCfg.Endpoint))

	// Initialize default buckets
	ctx := context.Background()
	if err := storageClient.InitDefaultBuckets(ctx); err != nil {
		logger.Error("Failed to initialize buckets", zap.Error(err))
	}

	// Initialize repositories
	fileRepo := postgres.NewFileRepository(db)
	categoryRepo := postgres.NewFileCategoryRepository(db)

	// Initialize use cases
	fileUC := file.NewUseCase(fileRepo, categoryRepo, storageClient, logger.Get())

	// Initialize handlers
	fileHandler := handler.NewFileHandler(fileUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(fileHandler, healthHandler, logger.Get())
	httpServer := router.Setup()

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		logger.Info("HTTP server listening", zap.String("address", addr))
		if err := httpServer.Run(addr); err != nil {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	logger.Info("Server stopped")
}
