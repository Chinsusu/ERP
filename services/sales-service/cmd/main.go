package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/config"
	httpdelivery "github.com/erp-cosmetics/sales-service/internal/delivery/http"
	"github.com/erp-cosmetics/sales-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/infrastructure/event"
	postgresrepo "github.com/erp-cosmetics/sales-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/sales-service/internal/usecase/customer"
	"github.com/erp-cosmetics/sales-service/internal/usecase/quotation"
	salesorder "github.com/erp-cosmetics/sales-service/internal/usecase/sales_order"
	"github.com/erp-cosmetics/sales-service/internal/usecase/shipment"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	zapLogger := initLogger(cfg.LogLevel)
	defer zapLogger.Sync()

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Initialize NATS client
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		zapLogger.Warn("Failed to connect to NATS, events will be disabled", zap.Error(err))
	}
	if nc != nil {
		defer nc.Close()
	}

	// Initialize event publisher
	var eventPublisher *event.Publisher
	if nc != nil {
		eventPublisher, err = event.NewPublisher(nc, zapLogger)
		if err != nil {
			zapLogger.Warn("Failed to create event publisher", zap.Error(err))
		}
	}

	// Initialize repositories
	customerRepo := postgresrepo.NewCustomerRepository(db)
	customerGroupRepo := postgresrepo.NewCustomerGroupRepository(db)
	quotationRepo := postgresrepo.NewQuotationRepository(db)
	salesOrderRepo := postgresrepo.NewSalesOrderRepository(db)
	shipmentRepo := postgresrepo.NewShipmentRepository(db)
	// returnRepo := postgresrepo.NewReturnRepository(db) // For future use

	// Initialize use cases - Customer
	createCustomerUC := customer.NewCreateCustomerUseCase(customerRepo, eventPublisher)
	getCustomerUC := customer.NewGetCustomerUseCase(customerRepo)
	listCustomersUC := customer.NewListCustomersUseCase(customerRepo)
	updateCustomerUC := customer.NewUpdateCustomerUseCase(customerRepo)
	deleteCustomerUC := customer.NewDeleteCustomerUseCase(customerRepo)
	checkCreditUC := customer.NewCheckCreditUseCase(customerRepo)

	// Initialize use cases - Quotation
	createQuotationUC := quotation.NewCreateQuotationUseCase(quotationRepo, customerRepo)
	getQuotationUC := quotation.NewGetQuotationUseCase(quotationRepo)
	listQuotationsUC := quotation.NewListQuotationsUseCase(quotationRepo)
	sendQuotationUC := quotation.NewSendQuotationUseCase(quotationRepo, eventPublisher)
	convertToOrderUC := quotation.NewConvertToOrderUseCase(quotationRepo, salesOrderRepo, customerRepo, eventPublisher)

	// Initialize use cases - Sales Order
	createOrderUC := salesorder.NewCreateOrderUseCase(salesOrderRepo, customerRepo, eventPublisher)
	getOrderUC := salesorder.NewGetOrderUseCase(salesOrderRepo)
	listOrdersUC := salesorder.NewListOrdersUseCase(salesOrderRepo)
	confirmOrderUC := salesorder.NewConfirmOrderUseCase(salesOrderRepo, customerRepo, eventPublisher, cfg.EnableCreditCheck)
	cancelOrderUC := salesorder.NewCancelOrderUseCase(salesOrderRepo, customerRepo, eventPublisher)
	shipOrderUC := salesorder.NewShipOrderUseCase(salesOrderRepo, eventPublisher)
	deliverOrderUC := salesorder.NewDeliverOrderUseCase(salesOrderRepo, eventPublisher)

	// Initialize use cases - Shipment
	createShipmentUC := shipment.NewCreateShipmentUseCase(shipmentRepo, salesOrderRepo, eventPublisher)
	getShipmentUC := shipment.NewGetShipmentUseCase(shipmentRepo)
	listShipmentsUC := shipment.NewListShipmentsUseCase(shipmentRepo)
	shipShipmentUC := shipment.NewShipShipmentUseCase(shipmentRepo, salesOrderRepo, eventPublisher)
	deliverShipmentUC := shipment.NewDeliverShipmentUseCase(shipmentRepo, salesOrderRepo, eventPublisher)

	// Initialize HTTP handlers
	customerHandler := handler.NewCustomerHandler(
		createCustomerUC,
		getCustomerUC,
		listCustomersUC,
		updateCustomerUC,
		deleteCustomerUC,
		checkCreditUC,
		customerRepo,
		customerGroupRepo,
	)

	quotationHandler := handler.NewQuotationHandler(
		createQuotationUC,
		getQuotationUC,
		listQuotationsUC,
		sendQuotationUC,
		convertToOrderUC,
		quotationRepo,
	)

	salesOrderHandler := handler.NewSalesOrderHandler(
		createOrderUC,
		getOrderUC,
		listOrdersUC,
		confirmOrderUC,
		cancelOrderUC,
		shipOrderUC,
		deliverOrderUC,
		salesOrderRepo,
	)

	shipmentHandler := handler.NewShipmentHandler(
		createShipmentUC,
		getShipmentUC,
		listShipmentsUC,
		shipShipmentUC,
		deliverShipmentUC,
	)

	// Create HTTP router
	router := httpdelivery.NewRouter(
		customerHandler,
		quotationHandler,
		salesOrderHandler,
		shipmentHandler,
	)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	// Start HTTP server
	go func() {
		zapLogger.Info("Starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Start gRPC server
	grpcLis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		zapLogger.Fatal("Failed to listen for gRPC", zap.Error(err))
	}
	grpcServer := grpc.NewServer()

	// TODO: Register gRPC services here
	// pb.RegisterSalesServiceServer(grpcServer, grpcService)

	go func() {
		zapLogger.Info("Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(grpcLis); err != nil {
			zapLogger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	zapLogger.Info("Sales Service started successfully",
		zap.String("http_port", cfg.HTTPPort),
		zap.String("grpc_port", cfg.GRPCPort),
	)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("Shutting down servers...")

	// Shutdown gRPC server
	grpcServer.GracefulStop()
	zapLogger.Info("gRPC server stopped")

	// Shutdown HTTP server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		zapLogger.Error("HTTP server shutdown error", zap.Error(err))
	}
	zapLogger.Info("HTTP server stopped")

	zapLogger.Info("Sales Service stopped gracefully")
}

func initLogger(level string) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: zapLevel == zapcore.DebugLevel,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	gormLogger := logger.Default.LogMode(logger.Silent)
	if cfg.LogLevel == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate entities (optional - prefer migrations in production)
	if err := db.AutoMigrate(
		&entity.CustomerGroup{},
		&entity.Customer{},
		&entity.CustomerAddress{},
		&entity.CustomerContact{},
		&entity.Quotation{},
		&entity.QuotationLineItem{},
		&entity.SalesOrder{},
		&entity.SOLineItem{},
		&entity.Shipment{},
		&entity.Return{},
		&entity.ReturnLineItem{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
