package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/manufacturing-service/internal/config"
	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/router"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/bom"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/ncr"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/qc"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/traceability"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/workorder"
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

	log.Info("Starting Manufacturing Service",
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
	bomRepo := postgres.NewBOMRepository(db)
	woRepo := postgres.NewWorkOrderRepository(db)
	qcRepo := postgres.NewQCRepository(db)
	ncrRepo := postgres.NewNCRRepository(db)
	traceRepo := postgres.NewTraceabilityRepository(db)

	// Initialize event publisher
	eventPub := event.NewPublisher(natsClient, log)

	// Initialize BOM use cases
	createBOMUC := bom.NewCreateBOMUseCase(bomRepo, eventPub, cfg.BOMEncryptionKey)
	getBOMUC := bom.NewGetBOMUseCase(bomRepo, cfg.BOMEncryptionKey)
	listBOMsUC := bom.NewListBOMsUseCase(bomRepo)
	approveBOMUC := bom.NewApproveBOMUseCase(bomRepo, eventPub)
	getActiveBOMUC := bom.NewGetActiveBOMUseCase(bomRepo)

	// Initialize Work Order use cases
	createWOUC := workorder.NewCreateWOUseCase(woRepo, bomRepo, eventPub)
	getWOUC := workorder.NewGetWOUseCase(woRepo)
	listWOsUC := workorder.NewListWOsUseCase(woRepo)
	releaseWOUC := workorder.NewReleaseWOUseCase(woRepo, eventPub)
	startWOUC := workorder.NewStartWOUseCase(woRepo, eventPub)
	completeWOUC := workorder.NewCompleteWOUseCase(woRepo, traceRepo, eventPub)

	// Initialize QC use cases
	getCheckpointsUC := qc.NewGetCheckpointsUseCase(qcRepo)
	createInspectionUC := qc.NewCreateInspectionUseCase(qcRepo, eventPub)
	getInspectionUC := qc.NewGetInspectionUseCase(qcRepo)
	listInspectionsUC := qc.NewListInspectionsUseCase(qcRepo)
	approveInspectionUC := qc.NewApproveInspectionUseCase(qcRepo, eventPub)

	// Initialize NCR use cases
	createNCRUC := ncr.NewCreateNCRUseCase(ncrRepo, eventPub)
	getNCRUC := ncr.NewGetNCRUseCase(ncrRepo)
	listNCRsUC := ncr.NewListNCRsUseCase(ncrRepo)
	closeNCRUC := ncr.NewCloseNCRUseCase(ncrRepo)

	// Initialize Traceability use cases
	traceBackwardUC := traceability.NewTraceBackwardUseCase(traceRepo, woRepo)
	traceForwardUC := traceability.NewTraceForwardUseCase(traceRepo)

	// Initialize handlers
	bomHandler := handler.NewBOMHandler(createBOMUC, getBOMUC, listBOMsUC, approveBOMUC, getActiveBOMUC)
	woHandler := handler.NewWOHandler(createWOUC, getWOUC, listWOsUC, releaseWOUC, startWOUC, completeWOUC)
	qcHandler := handler.NewQCHandler(getCheckpointsUC, createInspectionUC, getInspectionUC, listInspectionsUC, approveInspectionUC)
	ncrHandler := handler.NewNCRHandler(createNCRUC, getNCRUC, listNCRsUC, closeNCRUC)
	traceHandler := handler.NewTraceHandler(traceBackwardUC, traceForwardUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	r := router.SetupRouter(bomHandler, woHandler, qcHandler, ncrHandler, traceHandler, healthHandler)

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
