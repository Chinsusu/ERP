package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/reporting-service/internal/config"
	"github.com/erp-cosmetics/reporting-service/internal/delivery/http"
	"github.com/erp-cosmetics/reporting-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/reporting-service/internal/infrastructure/aggregator"
	"github.com/erp-cosmetics/reporting-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/reporting-service/internal/usecase/dashboard"
	"github.com/erp-cosmetics/reporting-service/internal/usecase/report"
	"github.com/erp-cosmetics/reporting-service/internal/usecase/stats"

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

	logger.Info("Starting Reporting Service",
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

	// Initialize repositories
	reportDefRepo := postgres.NewReportDefinitionRepository(db)
	reportExecRepo := postgres.NewReportExecutionRepository(db)
	dashboardRepo := postgres.NewDashboardRepository(db)
	widgetRepo := postgres.NewWidgetRepository(db)

	// Initialize aggregator
	serviceURLs := map[string]string{
		"wms":         "http://wms-service:8086",
		"sales":       "http://sales-service:8087",
		"procurement": "http://procurement-service:8085",
		"production":  "http://manufacturing-service:8088",
	}
	serviceClient := aggregator.NewServiceClient(serviceURLs, logger.Get())
	statsAggregator := aggregator.NewStatsAggregator(serviceClient, logger.Get())

	// Initialize use cases
	dashboardUC := dashboard.NewUseCase(dashboardRepo, widgetRepo, logger.Get())
	reportUC := report.NewUseCase(reportDefRepo, reportExecRepo, logger.Get())
	statsUC := stats.NewUseCase(statsAggregator, logger.Get())

	// Initialize handlers
	dashboardHandler := handler.NewDashboardHandler(dashboardUC)
	reportHandler := handler.NewReportHandler(reportUC)
	statsHandler := handler.NewStatsHandler(statsUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(
		dashboardHandler,
		reportHandler,
		statsHandler,
		healthHandler,
		logger.Get(),
	)
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
