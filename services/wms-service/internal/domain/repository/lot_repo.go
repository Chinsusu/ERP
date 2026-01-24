package repository

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/google/uuid"
)

// LotRepository defines lot repository interface
type LotRepository interface {
	Create(ctx context.Context, lot *entity.Lot) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Lot, error)
	GetByLotNumber(ctx context.Context, lotNumber string) (*entity.Lot, error)
	List(ctx context.Context, filter *LotFilter) ([]*entity.Lot, int64, error)
	Update(ctx context.Context, lot *entity.Lot) error
	
	// FEFO specific
	GetAvailableLots(ctx context.Context, materialID uuid.UUID) ([]*entity.Lot, error)
	GetExpiringLots(ctx context.Context, days int) ([]*entity.Lot, error)
	GetExpiredLots(ctx context.Context) ([]*entity.Lot, error)
	
	// Lot number generation
	GetNextLotNumber(ctx context.Context) (string, error)
	
	// Bulk operations
	MarkExpired(ctx context.Context, lotIDs []uuid.UUID) error
}

// LotFilter defines filter options for lots
type LotFilter struct {
	MaterialID    *uuid.UUID
	SupplierID    *uuid.UUID
	Status        string
	QCStatus      string
	ExpiryBefore  *time.Time
	ExpiryAfter   *time.Time
	Search        string
	Page          int
	Limit         int
}
