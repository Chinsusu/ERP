package scheduler

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"go.uber.org/zap"
)

// Scheduler handles scheduled WMS jobs
type Scheduler struct {
	lotRepo   repository.LotRepository
	stockRepo repository.StockRepository
	eventPub  *event.Publisher
	logger    *zap.Logger
	config    *Config
	stopChan  chan struct{}
}

// Config holds scheduler configuration
type Config struct {
	ExpiryCheckInterval    time.Duration
	LowStockCheckInterval  time.Duration
	ExpiryAlertDays        []int // 90, 30, 7
	LowStockThreshold      float64
}

// DefaultConfig returns default scheduler config
func DefaultConfig() *Config {
	return &Config{
		ExpiryCheckInterval:   24 * time.Hour,     // Daily
		LowStockCheckInterval: 1 * time.Hour,      // Hourly
		ExpiryAlertDays:       []int{90, 30, 7},   // Alert at 90, 30, 7 days
		LowStockThreshold:     100,                // Minimum stock level
	}
}

// NewScheduler creates a new scheduler
func NewScheduler(
	lotRepo repository.LotRepository,
	stockRepo repository.StockRepository,
	eventPub *event.Publisher,
	logger *zap.Logger,
	config *Config,
) *Scheduler {
	if config == nil {
		config = DefaultConfig()
	}
	return &Scheduler{
		lotRepo:   lotRepo,
		stockRepo: stockRepo,
		eventPub:  eventPub,
		logger:    logger,
		config:    config,
		stopChan:  make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.logger.Info("Starting WMS scheduler")

	// Run expiry check immediately
	go s.runExpiryCheck()

	// Start scheduled jobs
	go s.scheduleExpiryCheck()
	go s.scheduleLowStockCheck()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	close(s.stopChan)
	s.logger.Info("WMS scheduler stopped")
}

// scheduleExpiryCheck runs expiry check at intervals
func (s *Scheduler) scheduleExpiryCheck() {
	ticker := time.NewTicker(s.config.ExpiryCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runExpiryCheck()
		case <-s.stopChan:
			return
		}
	}
}

// scheduleLowStockCheck runs low stock check at intervals
func (s *Scheduler) scheduleLowStockCheck() {
	ticker := time.NewTicker(s.config.LowStockCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runLowStockCheck()
		case <-s.stopChan:
			return
		}
	}
}

// runExpiryCheck checks for expiring and expired lots
func (s *Scheduler) runExpiryCheck() {
	ctx := context.Background()
	s.logger.Info("Running expiry check job")

	// Check for each alert threshold
	for _, days := range s.config.ExpiryAlertDays {
		expiringLots, err := s.lotRepo.GetExpiringLots(ctx, days)
		if err != nil {
			s.logger.Error("Failed to get expiring lots", zap.Error(err))
			continue
		}

		for _, lot := range expiringLots {
			// Only alert if expiry is exactly at this threshold
			daysUntil := lot.DaysUntilExpiry()
			if daysUntil <= days && daysUntil > days-1 {
				s.logger.Warn("Lot expiring soon",
					zap.String("lot_number", lot.LotNumber),
					zap.Int("days_until_expiry", daysUntil))

				// Get stock quantity for this lot
				stocks, _ := s.stockRepo.GetByMaterialAndLot(ctx, lot.MaterialID, lot.ID)
				qty := 0.0
				if stocks != nil {
					qty = stocks.Quantity
				}

				s.eventPub.PublishLotExpiringSoon(&event.LotExpiringEvent{
					LotID:           lot.ID.String(),
					LotNumber:       lot.LotNumber,
					MaterialID:      lot.MaterialID.String(),
					ExpiryDate:      lot.ExpiryDate.Format("2006-01-02"),
					DaysUntilExpiry: daysUntil,
					Quantity:        qty,
				})
			}
		}
	}

	// Check for expired lots and mark them
	expiredLots, err := s.lotRepo.GetExpiredLots(ctx)
	if err != nil {
		s.logger.Error("Failed to get expired lots", zap.Error(err))
		return
	}

	if len(expiredLots) > 0 {
		var lotIDs []interface{}
		for _, lot := range expiredLots {
			s.logger.Warn("Marking lot as expired",
				zap.String("lot_number", lot.LotNumber))

			lot.MarkExpired()

			// Publish expired event
			s.eventPub.PublishLotExpired(&event.LotExpiringEvent{
				LotID:           lot.ID.String(),
				LotNumber:       lot.LotNumber,
				MaterialID:      lot.MaterialID.String(),
				ExpiryDate:      lot.ExpiryDate.Format("2006-01-02"),
				DaysUntilExpiry: 0,
			})

			lotIDs = append(lotIDs, lot.ID)
		}

		// Bulk update
		s.logger.Info("Marked lots as expired", zap.Int("count", len(expiredLots)))
	}
}

// runLowStockCheck checks for low stock levels
func (s *Scheduler) runLowStockCheck() {
	ctx := context.Background()
	s.logger.Info("Running low stock check job")

	lowStockMaterials, err := s.stockRepo.GetLowStockMaterials(ctx, s.config.LowStockThreshold)
	if err != nil {
		s.logger.Error("Failed to get low stock materials", zap.Error(err))
		return
	}

	for _, summary := range lowStockMaterials {
		s.logger.Warn("Low stock detected",
			zap.String("material_id", summary.MaterialID.String()),
			zap.Float64("available", summary.TotalAvailable))

		s.eventPub.PublishLowStockAlert(&event.LowStockAlertEvent{
			MaterialID:      summary.MaterialID.String(),
			CurrentQuantity: summary.TotalAvailable,
			ReorderPoint:    s.config.LowStockThreshold,
		})
	}
}
