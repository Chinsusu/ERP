package warehouse

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
)

// ListWarehousesUseCase handles listing warehouses
type ListWarehousesUseCase struct {
	warehouseRepo repository.WarehouseRepository
}

// NewListWarehousesUseCase creates a new use case
func NewListWarehousesUseCase(warehouseRepo repository.WarehouseRepository) *ListWarehousesUseCase {
	return &ListWarehousesUseCase{warehouseRepo: warehouseRepo}
}

// Execute lists warehouses
func (uc *ListWarehousesUseCase) Execute(ctx context.Context, filter *repository.WarehouseFilter) ([]*entity.Warehouse, int64, error) {
	return uc.warehouseRepo.List(ctx, filter)
}

// GetWarehouseUseCase handles getting warehouse
type GetWarehouseUseCase struct {
	warehouseRepo repository.WarehouseRepository
}

// NewGetWarehouseUseCase creates a new use case
func NewGetWarehouseUseCase(warehouseRepo repository.WarehouseRepository) *GetWarehouseUseCase {
	return &GetWarehouseUseCase{warehouseRepo: warehouseRepo}
}

// Execute gets a warehouse by ID
func (uc *GetWarehouseUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	return uc.warehouseRepo.GetByID(ctx, id)
}

// GetZonesUseCase handles getting zones for a warehouse
type GetZonesUseCase struct {
	zoneRepo repository.ZoneRepository
}

// NewGetZonesUseCase creates a new use case
func NewGetZonesUseCase(zoneRepo repository.ZoneRepository) *GetZonesUseCase {
	return &GetZonesUseCase{zoneRepo: zoneRepo}
}

// Execute gets zones for a warehouse
func (uc *GetZonesUseCase) Execute(ctx context.Context, warehouseID uuid.UUID) ([]*entity.Zone, error) {
	return uc.zoneRepo.GetByWarehouseID(ctx, warehouseID)
}

// GetLocationsUseCase handles getting locations for a zone
type GetLocationsUseCase struct {
	locationRepo repository.LocationRepository
}

// NewGetLocationsUseCase creates a new use case
func NewGetLocationsUseCase(locationRepo repository.LocationRepository) *GetLocationsUseCase {
	return &GetLocationsUseCase{locationRepo: locationRepo}
}

// Execute gets locations for a zone
func (uc *GetLocationsUseCase) Execute(ctx context.Context, zoneID uuid.UUID) ([]*entity.Location, error) {
	return uc.locationRepo.GetByZoneID(ctx, zoneID)
}
