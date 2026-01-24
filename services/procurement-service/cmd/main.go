package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/config"
	"github.com/erp-cosmetics/procurement-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/procurement-service/internal/delivery/http/router"
	"github.com/erp-cosmetics/procurement-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/procurement-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/procurement-service/internal/usecase/po"
	"github.com/erp-cosmetics/procurement-service/internal/usecase/pr"
	"github.com/erp-cosmetics/shared/pkg/database"
	"github.com/erp-cosmetics/shared/pkg/logger"
	"github.com/erp-cosmetics/shared/pkg/nats"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.ServiceName, &logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	defer logger.Sync()

	log.Info("Starting Procurement Service",
		zap.String("service", cfg.ServiceName),
		zap.String("port", cfg.Port),
	)

	// Connect to PostgreSQL
	db, err := database.Connect(database.NewDefaultConfig(cfg.GetDSN()))
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	log.Info("Connected to PostgreSQL")

	// Connect to NATS
	natsClient, err := nats.NewClient(&nats.Config{
		URL: cfg.NATSUrl,
	})
	if err != nil {
		log.Warn("Failed to connect to NATS, events will not be published", zap.Error(err))
	} else {
		log.Info("Connected to NATS")
	}

	// Initialize repositories
	prRepo := postgres.NewPRRepository(db)
	poRepo := postgres.NewPORepository(db)

	// Initialize event publisher
	eventPub := event.NewPublisher(natsClient, log)

	// Initialize PR use cases
	createPRUC := pr.NewCreatePRUseCase(prRepo, eventPub)
	getPRUC := pr.NewGetPRUseCase(prRepo)
	listPRsUC := pr.NewListPRsUseCase(prRepo)
	submitPRUC := pr.NewSubmitPRUseCase(prRepo, eventPub)
	approvePRUC := pr.NewApprovePRUseCase(prRepo, eventPub)
	rejectPRUC := pr.NewRejectPRUseCase(prRepo, eventPub)

	// Initialize PO use cases
	createPOFromPRUC := po.NewCreatePOFromPRUseCase(prRepo, poRepo, eventPub)
	getPOUC := po.NewGetPOUseCase(poRepo)
	listPOsUC := po.NewListPOsUseCase(poRepo)
	confirmPOUC := po.NewConfirmPOUseCase(poRepo, eventPub)
	cancelPOUC := po.NewCancelPOUseCase(poRepo, eventPub)
	closePOUC := po.NewClosePOUseCase(poRepo, eventPub)
	getPOReceiptsUC := po.NewGetPOReceiptsUseCase(poRepo)

	// Initialize handlers
	prHandler := handler.NewPRHandler(createPRUC, getPRUC, listPRsUC, submitPRUC, approvePRUC, rejectPRUC)
	poHandler := handler.NewPOHandler(createPOFromPRUC, getPOUC, listPOsUC, confirmPOUC, cancelPOUC, closePOUC, getPOReceiptsUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	r := router.SetupRouter(prHandler, poHandler, healthHandler)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Info("HTTP server started", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	if natsClient != nil {
		natsClient.Close()
	}

	log.Info("Server exited properly")
}
