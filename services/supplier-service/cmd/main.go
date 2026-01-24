package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/config"
	"github.com/erp-cosmetics/supplier-service/internal/delivery/http"
	"github.com/erp-cosmetics/supplier-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/supplier-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/supplier-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/supplier-service/internal/usecase/certification"
	"github.com/erp-cosmetics/supplier-service/internal/usecase/evaluation"
	"github.com/erp-cosmetics/supplier-service/internal/usecase/supplier"

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

	logger.Info("Starting Supplier Service",
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

	// Connect to NATS
	natsClient, err := nats.NewClient(&nats.Config{
		URL:           cfg.NATSUrl,
		MaxReconnects: 10,
		ReconnectWait: cfg.GetNATSReconnectWait(),
		Logger:        logger.Get(),
	})
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsClient.Close()
	logger.Info("Connected to NATS")

	// Initialize event publisher
	eventPublisher := event.NewPublisher(natsClient, logger.Get())

	// Initialize repositories
	supplierRepo := postgres.NewSupplierRepository(db)
	addressRepo := postgres.NewAddressRepository(db)
	contactRepo := postgres.NewContactRepository(db)
	certRepo := postgres.NewCertificationRepository(db)
	evalRepo := postgres.NewEvaluationRepository(db)

	// Initialize use cases
	createSupplierUC := supplier.NewCreateSupplierUseCase(supplierRepo, eventPublisher)
	getSupplierUC := supplier.NewGetSupplierUseCase(supplierRepo)
	listSuppliersUC := supplier.NewListSuppliersUseCase(supplierRepo)
	updateSupplierUC := supplier.NewUpdateSupplierUseCase(supplierRepo)
	approveSupplierUC := supplier.NewApproveSupplierUseCase(supplierRepo, eventPublisher)
	blockSupplierUC := supplier.NewBlockSupplierUseCase(supplierRepo, eventPublisher)

	addCertUC := certification.NewAddCertificationUseCase(certRepo, supplierRepo, eventPublisher)
	getCertsUC := certification.NewGetCertificationsUseCase(certRepo)
	getExpiringCertsUC := certification.NewGetExpiringCertificationsUseCase(certRepo)

	createEvalUC := evaluation.NewCreateEvaluationUseCase(evalRepo, supplierRepo, eventPublisher)
	getEvalsUC := evaluation.NewGetEvaluationsUseCase(evalRepo)

	// Initialize handlers
	supplierHandler := handler.NewSupplierHandler(
		createSupplierUC, getSupplierUC, listSuppliersUC,
		updateSupplierUC, approveSupplierUC, blockSupplierUC,
		addressRepo, contactRepo,
	)
	certHandler := handler.NewCertificationHandler(addCertUC, getCertsUC, getExpiringCertsUC)
	evalHandler := handler.NewEvaluationHandler(createEvalUC, getEvalsUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(supplierHandler, certHandler, evalHandler, healthHandler, logger.Get())
	httpServer := router.Setup()

	// Start HTTP server
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
