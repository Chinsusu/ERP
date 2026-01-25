package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/notification-service/internal/config"
	"github.com/erp-cosmetics/notification-service/internal/delivery/http"
	"github.com/erp-cosmetics/notification-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/notification-service/internal/infrastructure/email"
	"github.com/erp-cosmetics/notification-service/internal/infrastructure/event"
	"github.com/erp-cosmetics/notification-service/internal/infrastructure/persistence/postgres"
	alert_rule "github.com/erp-cosmetics/notification-service/internal/usecase/alert_rule"
	"github.com/erp-cosmetics/notification-service/internal/usecase/notification"
	"github.com/erp-cosmetics/notification-service/internal/usecase/template"
	user_notification "github.com/erp-cosmetics/notification-service/internal/usecase/user_notification"

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

	logger.Info("Starting Notification Service",
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

	// Initialize repositories
	notificationRepo := postgres.NewNotificationRepository(db)
	templateRepo := postgres.NewTemplateRepository(db)
	userNotificationRepo := postgres.NewUserNotificationRepository(db)
	alertRuleRepo := postgres.NewAlertRuleRepository(db)

	// Initialize email sender
	smtpConfig := cfg.GetSMTPConfig()
	var emailSender email.Sender
	if smtpConfig.Host != "" && smtpConfig.Host != "localhost" {
		emailSender = email.NewSMTPSender(&email.SMTPConfig{
			Host:      smtpConfig.Host,
			Port:      smtpConfig.Port,
			Username:  smtpConfig.Username,
			Password:  smtpConfig.Password,
			FromEmail: smtpConfig.FromEmail,
			FromName:  smtpConfig.FromName,
			UseTLS:    smtpConfig.UseTLS,
		}, logger.Get())
		logger.Info("SMTP email sender initialized", zap.String("host", smtpConfig.Host))
	} else {
		emailSender = email.NewMockSender(logger.Get())
		logger.Info("Mock email sender initialized (no SMTP configured)")
	}

	// Initialize event publisher
	eventPublisher := event.NewPublisher(natsClient, logger.Get())

	// Initialize use cases
	notificationUC := notification.NewUseCase(
		notificationRepo,
		templateRepo,
		userNotificationRepo,
		emailSender,
		eventPublisher,
		logger.Get(),
	)
	templateUC := template.NewUseCase(templateRepo, logger.Get())
	userNotificationUC := user_notification.NewUseCase(userNotificationRepo, logger.Get())
	alertRuleUC := alert_rule.NewUseCase(alertRuleRepo, logger.Get())

	// Initialize event subscriber
	eventSubscriber := event.NewSubscriber(
		natsClient,
		logger.Get(),
		templateRepo,
		userNotificationRepo,
		notificationRepo,
	)
	if err := eventSubscriber.SubscribeAll(); err != nil {
		logger.Error("Failed to subscribe to events", zap.Error(err))
	}

	// Initialize HTTP handlers
	notificationHandler := handler.NewNotificationHandler(notificationUC, templateUC)
	userNotificationHandler := handler.NewUserNotificationHandler(userNotificationUC)
	alertRuleHandler := handler.NewAlertRuleHandler(alertRuleUC)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := http.NewRouter(
		notificationHandler,
		userNotificationHandler,
		alertRuleHandler,
		healthHandler,
		logger.Get(),
	)
	httpServer := router.Setup()

	// Start background email processor
	go startEmailProcessor(notificationUC, logger.Get())

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

// startEmailProcessor starts background email processing
func startEmailProcessor(notificationUC notification.UseCase, log *zap.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()

		// Process pending notifications
		if err := notificationUC.ProcessPending(ctx, 50); err != nil {
			log.Error("Failed to process pending notifications", zap.Error(err))
		}

		// Process retryable notifications
		if err := notificationUC.ProcessRetryable(ctx, 20); err != nil {
			log.Error("Failed to process retryable notifications", zap.Error(err))
		}
	}
}
