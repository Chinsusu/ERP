package stats

import (
	"context"

	"github.com/erp-cosmetics/reporting-service/internal/infrastructure/aggregator"
	"go.uber.org/zap"
)

// UseCase defines stats use case interface
type UseCase interface {
	GetInventoryStats(ctx context.Context) (*aggregator.InventoryStats, error)
	GetSalesStats(ctx context.Context) (*aggregator.SalesStats, error)
	GetProductionStats(ctx context.Context) (*aggregator.ProductionStats, error)
	GetProcurementStats(ctx context.Context) (*aggregator.ProcurementStats, error)
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
}

// DashboardStats combines all stats for dashboard
type DashboardStats struct {
	Inventory   *aggregator.InventoryStats   `json:"inventory"`
	Sales       *aggregator.SalesStats       `json:"sales"`
	Production  *aggregator.ProductionStats  `json:"production"`
	Procurement *aggregator.ProcurementStats `json:"procurement"`
}

type useCase struct {
	aggregator *aggregator.StatsAggregator
	logger     *zap.Logger
}

// NewUseCase creates new stats use case
func NewUseCase(agg *aggregator.StatsAggregator, logger *zap.Logger) UseCase {
	return &useCase{
		aggregator: agg,
		logger:     logger,
	}
}

func (uc *useCase) GetInventoryStats(ctx context.Context) (*aggregator.InventoryStats, error) {
	return uc.aggregator.GetInventoryStats(ctx)
}

func (uc *useCase) GetSalesStats(ctx context.Context) (*aggregator.SalesStats, error) {
	return uc.aggregator.GetSalesStats(ctx)
}

func (uc *useCase) GetProductionStats(ctx context.Context) (*aggregator.ProductionStats, error) {
	return uc.aggregator.GetProductionStats(ctx)
}

func (uc *useCase) GetProcurementStats(ctx context.Context) (*aggregator.ProcurementStats, error) {
	return uc.aggregator.GetProcurementStats(ctx)
}

func (uc *useCase) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	inventory, _ := uc.aggregator.GetInventoryStats(ctx)
	sales, _ := uc.aggregator.GetSalesStats(ctx)
	production, _ := uc.aggregator.GetProductionStats(ctx)
	procurement, _ := uc.aggregator.GetProcurementStats(ctx)

	return &DashboardStats{
		Inventory:   inventory,
		Sales:       sales,
		Production:  production,
		Procurement: procurement,
	}, nil
}
