package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/user-service/internal/config"
	"github.com/erp-cosmetics/user-service/internal/delivery/http"
	"github.com/erp-cosmetics/user-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/user-service/internal/infrastructure/client"
	"github.com/erp-cosmetics/user-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/user-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/user-service/internal/usecase/department"
	"github.com/erp-cosmetics/user-service/internal/usecase/user"

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

	logger.Info("Starting User Service",
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
		ReconnectWait: 2 * time.Second,
		Logger:        logger.Get(),
	})
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsClient.Close()
	logger.Info("Connected to NATS")

	// Initialize Auth Service client
	authClient, err := client.NewAuthClient(cfg.AuthServiceAddr)
	if err != nil {
		logger.Warn("Failed to connect to Auth Service", zap.Error(err))
		// Continue without auth client - can be added later
	}
	if authClient != nil {
		defer authClient.Close()
		logger.Info("Connected to Auth Service")
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	deptRepo := postgres.NewDepartmentRepository(db)
	profileRepo := postgres.NewUserProfileRepository(db)

	// Initialize event publisher
	eventPub := event.NewPublisher(natsClient, logger.Get())

	// Initialize use cases
	createUserUC := user.NewCreateUserUseCase(userRepo, profileRepo, authClient, eventPub)
	getUserUC := user.NewGetUserUseCase(userRepo)
	listUsersUC := user.NewListUsersUseCase(userRepo)

	createDeptUC := department.NewCreateDepartmentUseCase(deptRepo, eventPub)
	getDeptTreeUC := department.NewGetDepartmentTreeUseCase(deptRepo)

	// Initialize HTTP handlers
	userHandler := handler.NewUserHandler(createUserUC, getUserUC, listUsersUC)
	deptHandler := handler.NewDepartmentHandler(createDeptUC, getDeptTreeUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(userHandler, deptHandler, healthHandler, logger.Get())
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
