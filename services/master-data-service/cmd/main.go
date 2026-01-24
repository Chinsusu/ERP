package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/master-data-service/internal/config"
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http"
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/master-data-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/master-data-service/internal/infrastructure/persistence/postgres"
	categoryUC "github.com/erp-cosmetics/master-data-service/internal/usecase/category"
	materialUC "github.com/erp-cosmetics/master-data-service/internal/usecase/material"
	productUC "github.com/erp-cosmetics/master-data-service/internal/usecase/product"
	unitUC "github.com/erp-cosmetics/master-data-service/internal/usecase/unit"

	"github.com/erp-cosmetics/shared/pkg/database"
	"github.com/erp-cosmetics/shared/pkg/logger"
	"github.com/erp-cosmetics/shared/pkg/nats"

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

	logger.Info("Starting Master Data Service",
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

	// Connect to NATS (optional - may not be available in all environments)
	var natsClient *nats.Client
	natsClient, err = nats.NewClient(&nats.Config{
		URL:           cfg.NATSUrl,
		MaxReconnects: 10,
		ReconnectWait: 2 * time.Second,
		Logger:        logger.Get(),
	})
	if err != nil {
		logger.Warn("Failed to connect to NATS, events will not be published", zap.Error(err))
	} else {
		defer natsClient.Close()
		logger.Info("Connected to NATS")
	}

	// Initialize event publisher
	eventPublisher := event.NewPublisher(natsClient, logger.Get())

	// Initialize repositories
	categoryRepo := postgres.NewCategoryRepository(db)
	unitRepo := postgres.NewUnitRepository(db)
	materialRepo := postgres.NewMaterialRepository(db)
	productRepo := postgres.NewProductRepository(db)

	// Initialize use cases
	categoryUseCase := categoryUC.NewUseCase(categoryRepo, eventPublisher)
	unitUseCase := unitUC.NewUseCase(unitRepo)
	materialUseCase := materialUC.NewUseCase(materialRepo, categoryRepo, unitRepo, eventPublisher, cfg.AutoGenerateCodes)
	productUseCase := productUC.NewUseCase(productRepo, categoryRepo, unitRepo, eventPublisher, cfg.AutoGenerateCodes)

	// Initialize HTTP handlers
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	unitHandler := handler.NewUnitHandler(unitUseCase)
	materialHandler := handler.NewMaterialHandler(materialUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(
		categoryHandler,
		unitHandler,
		materialHandler,
		productHandler,
		healthHandler,
		logger.Get(),
	)
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

	logger.Info("Server stopped")
}
