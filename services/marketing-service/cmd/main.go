package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/marketing-service/internal/config"
	httpdelivery "github.com/erp-cosmetics/marketing-service/internal/delivery/http"
	"github.com/erp-cosmetics/marketing-service/internal/delivery/http/handler"
	// entity import removed - using SQL migrations instead of auto-migrate
	"github.com/erp-cosmetics/marketing-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/marketing-service/internal/infrastructure/persistence/postgres"
	campaignuc "github.com/erp-cosmetics/marketing-service/internal/usecase/campaign"
	koluc "github.com/erp-cosmetics/marketing-service/internal/usecase/kol"
	sampleuc "github.com/erp-cosmetics/marketing-service/internal/usecase/sample"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer zapLogger.Sync()

	// Connect to database
	db, err := gorm.Open(gormpostgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Skip auto migrate since SQL migrations already applied
	// if err := db.AutoMigrate(
	// 	&entity.KOLTier{},
	// 	&entity.KOL{},
	// 	&entity.Campaign{},
	// 	&entity.KOLCollaboration{},
	// 	&entity.SampleRequest{},
	// 	&entity.SampleItem{},
	// 	&entity.SampleShipment{},
	// 	&entity.KOLPost{},
	// ); err != nil {
	// 	log.Fatalf("Failed to auto migrate: %v", err)
	// }

	// Connect to NATS
	var publisher *event.Publisher
	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		zapLogger.Warn("Failed to connect to NATS, events will not be published", zap.Error(err))
	} else {
		defer nc.Close()
		publisher, err = event.NewPublisher(nc, zapLogger)
		if err != nil {
			zapLogger.Warn("Failed to create event publisher", zap.Error(err))
		}
	}

	// Initialize repositories
	kolTierRepo := postgres.NewKOLTierRepository(db)
	kolRepo := postgres.NewKOLRepository(db)
	campaignRepo := postgres.NewCampaignRepository(db)
	collabRepo := postgres.NewKOLCollaborationRepository(db)
	sampleRequestRepo := postgres.NewSampleRequestRepository(db)
	sampleShipmentRepo := postgres.NewSampleShipmentRepository(db)
	kolPostRepo := postgres.NewKOLPostRepository(db)

	// Initialize use cases - KOL
	createKOLUC := koluc.NewCreateKOLUseCase(kolRepo, publisher)
	getKOLUC := koluc.NewGetKOLUseCase(kolRepo)
	listKOLsUC := koluc.NewListKOLsUseCase(kolRepo)
	updateKOLUC := koluc.NewUpdateKOLUseCase(kolRepo)
	deleteKOLUC := koluc.NewDeleteKOLUseCase(kolRepo)

	// Initialize use cases - Campaign
	createCampaignUC := campaignuc.NewCreateCampaignUseCase(campaignRepo, publisher)
	getCampaignUC := campaignuc.NewGetCampaignUseCase(campaignRepo)
	listCampaignsUC := campaignuc.NewListCampaignsUseCase(campaignRepo)
	launchCampaignUC := campaignuc.NewLaunchCampaignUseCase(campaignRepo, publisher)
	updateCampaignUC := campaignuc.NewUpdateCampaignUseCase(campaignRepo)

	// Initialize use cases - Sample
	createSampleRequestUC := sampleuc.NewCreateSampleRequestUseCase(sampleRequestRepo, kolRepo, publisher)
	getSampleRequestUC := sampleuc.NewGetSampleRequestUseCase(sampleRequestRepo)
	listSampleRequestsUC := sampleuc.NewListSampleRequestsUseCase(sampleRequestRepo)
	approveSampleRequestUC := sampleuc.NewApproveSampleRequestUseCase(sampleRequestRepo, kolRepo, publisher)
	rejectSampleRequestUC := sampleuc.NewRejectSampleRequestUseCase(sampleRequestRepo)
	shipSampleUC := sampleuc.NewShipSampleUseCase(sampleRequestRepo, sampleShipmentRepo, publisher)

	// Initialize handlers
	kolHandler := handler.NewKOLHandler(
		createKOLUC, getKOLUC, listKOLsUC, updateKOLUC, deleteKOLUC,
		kolTierRepo, kolPostRepo,
	)
	campaignHandler := handler.NewCampaignHandler(
		createCampaignUC, getCampaignUC, listCampaignsUC, launchCampaignUC, updateCampaignUC,
		collabRepo,
	)
	sampleHandler := handler.NewSampleHandler(
		createSampleRequestUC, getSampleRequestUC, listSampleRequestsUC,
		approveSampleRequestUC, rejectSampleRequestUC, shipSampleUC,
		sampleShipmentRepo,
	)

	// Create HTTP router
	router := httpdelivery.NewRouter(kolHandler, campaignHandler, sampleHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Start servers
	errChan := make(chan error, 2)

	// Start HTTP server
	go func() {
		zapLogger.Info("Starting HTTP server", zap.String("port", cfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			errChan <- fmt.Errorf("failed to listen for gRPC: %w", err)
			return
		}
		zapLogger.Info("Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	zapLogger.Info("Marketing Service started successfully",
		zap.String("http_port", cfg.Port),
		zap.String("grpc_port", cfg.GRPCPort),
	)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		zapLogger.Error("Server error", zap.Error(err))
	case sig := <-sigChan:
		zapLogger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	zapLogger.Info("Shutting down servers...")
	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(ctx); err != nil {
		zapLogger.Error("HTTP server shutdown error", zap.Error(err))
	}

	zapLogger.Info("Marketing Service stopped")
}
