package repository

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/google/uuid"
)

// WarehouseRepository defines warehouse repository interface
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *entity.Warehouse) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*entity.Warehouse, error)
	List(ctx context.Context, filter *WarehouseFilter) ([]*entity.Warehouse, int64, error)
	Update(ctx context.Context, warehouse *entity.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// WarehouseFilter defines filter options for warehouses
type WarehouseFilter struct {
	Type     string
	IsActive *bool
	Search   string
	Page     int
	Limit    int
}

// ZoneRepository defines zone repository interface
type ZoneRepository interface {
	Create(ctx context.Context, zone *entity.Zone) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Zone, error)
	GetByWarehouseID(ctx context.Context, warehouseID uuid.UUID) ([]*entity.Zone, error)
	GetQuarantineZone(ctx context.Context, warehouseID uuid.UUID) (*entity.Zone, error)
	GetStorageZone(ctx context.Context, warehouseID uuid.UUID) (*entity.Zone, error)
	Update(ctx context.Context, zone *entity.Zone) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// LocationRepository defines location repository interface
type LocationRepository interface {
	Create(ctx context.Context, location *entity.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Location, error)
	GetByZoneID(ctx context.Context, zoneID uuid.UUID) ([]*entity.Location, error)
	GetByCode(ctx context.Context, zoneID uuid.UUID, code string) (*entity.Location, error)
	Update(ctx context.Context, location *entity.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
}
