package repository

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/google/uuid"
)

// InventoryCountFilter represents filter for inventory counts
type InventoryCountFilter struct {
	WarehouseID *uuid.UUID
	ZoneID      *uuid.UUID
	Status      string
	CountType   string
	Search      string
	Page        int
	Limit       int
}

// InventoryCountRepository defines inventory count repository interface
type InventoryCountRepository interface {
	// Basic CRUD
	Create(ctx context.Context, count *entity.InventoryCount) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.InventoryCount, error)
	GetByNumber(ctx context.Context, countNumber string) (*entity.InventoryCount, error)
	List(ctx context.Context, filter *InventoryCountFilter) ([]*entity.InventoryCount, int64, error)
	Update(ctx context.Context, count *entity.InventoryCount) error

	// Line items
	CreateLineItem(ctx context.Context, item *entity.InventoryCountLineItem) error
	CreateLineItems(ctx context.Context, items []*entity.InventoryCountLineItem) error
	GetLineItemsByCountID(ctx context.Context, countID uuid.UUID) ([]*entity.InventoryCountLineItem, error)
	UpdateLineItem(ctx context.Context, item *entity.InventoryCountLineItem) error

	// Queries
	GetPendingItems(ctx context.Context, countID uuid.UUID) ([]*entity.InventoryCountLineItem, error)
	GetVarianceItems(ctx context.Context, countID uuid.UUID) ([]*entity.InventoryCountLineItem, error)
	GetNextCountNumber(ctx context.Context) (string, error)
}
