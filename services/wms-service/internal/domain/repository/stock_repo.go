package repository

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/google/uuid"
)

// StockRepository defines stock repository interface
type StockRepository interface {
	// Basic CRUD
	Create(ctx context.Context, stock *entity.Stock) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Stock, error)
	Update(ctx context.Context, stock *entity.Stock) error
	
	// Stock queries
	GetByLocation(ctx context.Context, locationID uuid.UUID) ([]*entity.Stock, error)
	GetByMaterial(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error)
	GetByMaterialAndLot(ctx context.Context, materialID, lotID uuid.UUID) (*entity.Stock, error)
	GetByLocationMaterialLot(ctx context.Context, locationID, materialID uuid.UUID, lotID *uuid.UUID) (*entity.Stock, error)
	List(ctx context.Context, filter *StockFilter) ([]*entity.Stock, int64, error)
	
	// FEFO - Critical for cosmetics
	GetAvailableStockFEFO(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error)
	IssueStockFEFO(ctx context.Context, materialID uuid.UUID, quantity float64, createdBy uuid.UUID) ([]entity.LotIssued, error)
	
	// Stock operations
	ReceiveStock(ctx context.Context, stock *entity.Stock, movement *entity.StockMovement) error
	IssueStock(ctx context.Context, stock *entity.Stock, movement *entity.StockMovement) error
	TransferStock(ctx context.Context, fromStock, toStock *entity.Stock, movement *entity.StockMovement) error
	AdjustStock(ctx context.Context, stock *entity.Stock, adjustmentQty float64, movement *entity.StockMovement) error
	
	// Reservation
	ReserveStock(ctx context.Context, materialID uuid.UUID, quantity float64, reservation *entity.StockReservation) error
	ReleaseReservation(ctx context.Context, reservationID uuid.UUID) error
	
	// Aggregations
	GetMaterialSummary(ctx context.Context, materialID uuid.UUID) (*entity.StockSummary, error)
	GetLowStockMaterials(ctx context.Context, threshold float64) ([]*entity.StockSummary, error)
	GetExpiringStock(ctx context.Context, days int) ([]*entity.Stock, error)
	
	// Movement
	CreateMovement(ctx context.Context, movement *entity.StockMovement) error
	GetMovementsByLot(ctx context.Context, lotID uuid.UUID) ([]*entity.StockMovement, error)
	GetMovementsByMaterial(ctx context.Context, materialID uuid.UUID, limit int) ([]*entity.StockMovement, error)
	GetNextMovementNumber(ctx context.Context, movementType entity.MovementType) (string, error)
}

// StockFilter defines filter options for stock
type StockFilter struct {
	WarehouseID  *uuid.UUID
	ZoneID       *uuid.UUID
	LocationID   *uuid.UUID
	MaterialID   *uuid.UUID
	LotID        *uuid.UUID
	HasStock     *bool
	ExpiringDays *int
	Page         int
	Limit        int
}

// ReservationRepository defines reservation repository interface
type ReservationRepository interface {
	Create(ctx context.Context, reservation *entity.StockReservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.StockReservation, error)
	GetByReference(ctx context.Context, referenceID uuid.UUID) ([]*entity.StockReservation, error)
	GetActiveByMaterial(ctx context.Context, materialID uuid.UUID) ([]*entity.StockReservation, error)
	Update(ctx context.Context, reservation *entity.StockReservation) error
	Release(ctx context.Context, id uuid.UUID) error
	GetExpiredReservations(ctx context.Context) ([]*entity.StockReservation, error)
}
