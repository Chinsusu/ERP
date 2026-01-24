package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/config"
	"github.com/erp-cosmetics/wms-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/wms-service/internal/delivery/http/router"
	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/scheduler"
	adjustment_uc "github.com/erp-cosmetics/wms-service/internal/usecase/adjustment"
	grn_uc "github.com/erp-cosmetics/wms-service/internal/usecase/grn"
	issue_uc "github.com/erp-cosmetics/wms-service/internal/usecase/issue"
	lot_uc "github.com/erp-cosmetics/wms-service/internal/usecase/lot"
	reservation_uc "github.com/erp-cosmetics/wms-service/internal/usecase/reservation"
	stock_uc "github.com/erp-cosmetics/wms-service/internal/usecase/stock"
	warehouse_uc "github.com/erp-cosmetics/wms-service/internal/usecase/warehouse"
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

	log.Info("Starting WMS Service",
		zap.String("service", cfg.ServiceName),
		zap.String("port", cfg.Port),
	)

	// Connect to PostgreSQL
	db, err := database.Connect(database.NewDefaultConfig(cfg.GetDSN()))
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	log.Info("Connected to PostgreSQL")

	// Auto-migrate entities
	err = db.AutoMigrate(
		&entity.Warehouse{},
		&entity.Zone{},
		&entity.Location{},
		&entity.Lot{},
		&entity.Stock{},
		&entity.StockMovement{},
		&entity.StockReservation{},
		&entity.GRN{},
		&entity.GRNLineItem{},
		&entity.GoodsIssue{},
		&entity.GILineItem{},
	)
	if err != nil {
		log.Warn("Auto-migration warning", zap.Error(err))
	}

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
	warehouseRepo := postgres.NewWarehouseRepository(db)
	zoneRepo := postgres.NewZoneRepository(db)
	locationRepo := postgres.NewLocationRepository(db)
	lotRepo := postgres.NewLotRepository(db)
	stockRepo := postgres.NewStockRepository(db)
	grnRepo := postgres.NewGRNRepository(db)
	issueRepo := postgres.NewGoodsIssueRepository(db)

	// Initialize event publisher
	eventPub := event.NewPublisher(natsClient, log)

	// Initialize warehouse use cases
	listWarehousesUC := warehouse_uc.NewListWarehousesUseCase(warehouseRepo)
	getWarehouseUC := warehouse_uc.NewGetWarehouseUseCase(warehouseRepo)
	getZonesUC := warehouse_uc.NewGetZonesUseCase(zoneRepo)
	getLocationsUC := warehouse_uc.NewGetLocationsUseCase(locationRepo)

	// Initialize stock use cases
	getStockUC := stock_uc.NewGetStockUseCase(stockRepo)
	issueStockFEFOUC := stock_uc.NewIssueStockFEFOUseCase(stockRepo, eventPub)
	reserveStockUC := stock_uc.NewReserveStockUseCase(stockRepo, eventPub)
	releaseReservationUC := stock_uc.NewReleaseReservationUseCase(stockRepo)

	// Initialize lot use cases
	getLotUC := lot_uc.NewGetLotUseCase(lotRepo)
	listLotsUC := lot_uc.NewListLotsUseCase(lotRepo)
	getExpiringLotsUC := lot_uc.NewGetExpiringLotsUseCase(lotRepo)
	getLotMovementsUC := lot_uc.NewGetLotMovementsUseCase(stockRepo)

	// Initialize GRN use cases
	createGRNUC := grn_uc.NewCreateGRNUseCase(grnRepo, lotRepo, stockRepo, zoneRepo, locationRepo, eventPub)
	completeGRNUC := grn_uc.NewCompleteGRNUseCase(grnRepo, lotRepo, stockRepo, zoneRepo, eventPub)
	getGRNUC := grn_uc.NewGetGRNUseCase(grnRepo)
	listGRNsUC := grn_uc.NewListGRNsUseCase(grnRepo)

	// Initialize Goods Issue use cases
	createIssueUC := issue_uc.NewCreateGoodsIssueUseCase(issueRepo, stockRepo, eventPub)
	getIssueUC := issue_uc.NewGetGoodsIssueUseCase(issueRepo)
	listIssuesUC := issue_uc.NewListGoodsIssuesUseCase(issueRepo)

	// Initialize Reservation use cases
	createReservationUC := reservation_uc.NewCreateReservationUseCase(stockRepo, eventPub)
	releaseReservationUC2 := reservation_uc.NewReleaseReservationUseCase(stockRepo)
	checkAvailabilityUC := reservation_uc.NewCheckAvailabilityUseCase(stockRepo)

	// Initialize Adjustment use cases
	createAdjustmentUC := adjustment_uc.NewCreateAdjustmentUseCase(stockRepo)
	transferStockUC := adjustment_uc.NewTransferStockUseCase(stockRepo)

	// Initialize handlers
	warehouseHandler := handler.NewWarehouseHandler(listWarehousesUC, getWarehouseUC, getZonesUC, getLocationsUC)
	stockHandler := handler.NewStockHandler(getStockUC, issueStockFEFOUC, reserveStockUC, releaseReservationUC)
	lotHandler := handler.NewLotHandler(getLotUC, listLotsUC, getExpiringLotsUC, getLotMovementsUC)
	grnHandler := handler.NewGRNHandler(createGRNUC, completeGRNUC, getGRNUC, listGRNsUC)
	issueHandler := handler.NewGoodsIssueHandler(createIssueUC, getIssueUC, listIssuesUC)
	reservationHandler := handler.NewReservationHandler(createReservationUC, releaseReservationUC2, checkAvailabilityUC)
	adjustmentHandler := handler.NewAdjustmentHandler(createAdjustmentUC, transferStockUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	r := router.SetupRouter(
		warehouseHandler,
		stockHandler,
		lotHandler,
		grnHandler,
		issueHandler,
		reservationHandler,
		adjustmentHandler,
		healthHandler,
	)

	// Start scheduler for expiry alerts
	lowStockInterval, _ := time.ParseDuration(cfg.LowStockCheckInterval)
	if lowStockInterval == 0 {
		lowStockInterval = 1 * time.Hour
	}
	schedulerConfig := &scheduler.Config{
		ExpiryCheckInterval:   24 * time.Hour,
		LowStockCheckInterval: lowStockInterval,
		ExpiryAlertDays:       []int{90, 30, 7},
		LowStockThreshold:     100,
	}
	wmsScheduler := scheduler.NewScheduler(lotRepo, stockRepo, eventPub, log, schedulerConfig)
	wmsScheduler.Start()

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

	// Stop scheduler
	wmsScheduler.Stop()

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
