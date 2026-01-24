package lot

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
)

// GetLotUseCase handles getting lot
type GetLotUseCase struct {
	lotRepo repository.LotRepository
}

// NewGetLotUseCase creates a new use case
func NewGetLotUseCase(lotRepo repository.LotRepository) *GetLotUseCase {
	return &GetLotUseCase{lotRepo: lotRepo}
}

// Execute gets a lot by ID
func (uc *GetLotUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Lot, error) {
	return uc.lotRepo.GetByID(ctx, id)
}

// ListLotsUseCase handles listing lots
type ListLotsUseCase struct {
	lotRepo repository.LotRepository
}

// NewListLotsUseCase creates a new use case
func NewListLotsUseCase(lotRepo repository.LotRepository) *ListLotsUseCase {
	return &ListLotsUseCase{lotRepo: lotRepo}
}

// Execute lists lots
func (uc *ListLotsUseCase) Execute(ctx context.Context, filter *repository.LotFilter) ([]*entity.Lot, int64, error) {
	return uc.lotRepo.List(ctx, filter)
}

// GetExpiringLotsUseCase handles getting expiring lots
type GetExpiringLotsUseCase struct {
	lotRepo repository.LotRepository
}

// NewGetExpiringLotsUseCase creates a new use case
func NewGetExpiringLotsUseCase(lotRepo repository.LotRepository) *GetExpiringLotsUseCase {
	return &GetExpiringLotsUseCase{lotRepo: lotRepo}
}

// Execute gets lots expiring within days
func (uc *GetExpiringLotsUseCase) Execute(ctx context.Context, days int) ([]*entity.Lot, error) {
	return uc.lotRepo.GetExpiringLots(ctx, days)
}

// GetLotMovementsUseCase handles getting lot movements
type GetLotMovementsUseCase struct {
	stockRepo repository.StockRepository
}

// NewGetLotMovementsUseCase creates a new use case
func NewGetLotMovementsUseCase(stockRepo repository.StockRepository) *GetLotMovementsUseCase {
	return &GetLotMovementsUseCase{stockRepo: stockRepo}
}

// Execute gets movements for a lot
func (uc *GetLotMovementsUseCase) Execute(ctx context.Context, lotID uuid.UUID) ([]*entity.StockMovement, error) {
	return uc.stockRepo.GetMovementsByLot(ctx, lotID)
}
