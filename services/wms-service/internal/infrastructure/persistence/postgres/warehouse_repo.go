package postgres

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

// NewWarehouseRepository creates a new warehouse repository
func NewWarehouseRepository(db *gorm.DB) repository.WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(ctx context.Context, warehouse *entity.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

func (r *warehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.WithContext(ctx).
		Preload("Zones").
		First(&warehouse, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) GetByCode(ctx context.Context, code string) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.WithContext(ctx).
		Preload("Zones").
		First(&warehouse, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) List(ctx context.Context, filter *repository.WarehouseFilter) ([]*entity.Warehouse, int64, error) {
	var warehouses []*entity.Warehouse
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Warehouse{})

	if filter.Type != "" {
		query = query.Where("warehouse_type = ?", filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(20)
	}
	if filter.Page > 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if err := query.Preload("Zones").Order("code").Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (r *warehouseRepository) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

func (r *warehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Warehouse{}, "id = ?", id).Error
}

// Zone Repository
type zoneRepository struct {
	db *gorm.DB
}

// NewZoneRepository creates a new zone repository
func NewZoneRepository(db *gorm.DB) repository.ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(ctx context.Context, zone *entity.Zone) error {
	return r.db.WithContext(ctx).Create(zone).Error
}

func (r *zoneRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Zone, error) {
	var zone entity.Zone
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("Locations").
		First(&zone, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) GetByWarehouseID(ctx context.Context, warehouseID uuid.UUID) ([]*entity.Zone, error) {
	var zones []*entity.Zone
	err := r.db.WithContext(ctx).
		Where("warehouse_id = ?", warehouseID).
		Order("code").
		Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) GetQuarantineZone(ctx context.Context, warehouseID uuid.UUID) (*entity.Zone, error) {
	var zone entity.Zone
	err := r.db.WithContext(ctx).
		Where("warehouse_id = ? AND zone_type = ?", warehouseID, entity.ZoneTypeQuarantine).
		First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) GetStorageZone(ctx context.Context, warehouseID uuid.UUID) (*entity.Zone, error) {
	var zone entity.Zone
	err := r.db.WithContext(ctx).
		Where("warehouse_id = ? AND zone_type = ?", warehouseID, entity.ZoneTypeStorage).
		First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) Update(ctx context.Context, zone *entity.Zone) error {
	return r.db.WithContext(ctx).Save(zone).Error
}

func (r *zoneRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Zone{}, "id = ?", id).Error
}

// Location Repository
type locationRepository struct {
	db *gorm.DB
}

// NewLocationRepository creates a new location repository
func NewLocationRepository(db *gorm.DB) repository.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(ctx context.Context, location *entity.Location) error {
	return r.db.WithContext(ctx).Create(location).Error
}

func (r *locationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Location, error) {
	var location entity.Location
	err := r.db.WithContext(ctx).
		Preload("Zone").
		Preload("Zone.Warehouse").
		First(&location, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) GetByZoneID(ctx context.Context, zoneID uuid.UUID) ([]*entity.Location, error) {
	var locations []*entity.Location
	err := r.db.WithContext(ctx).
		Where("zone_id = ?", zoneID).
		Order("code").
		Find(&locations).Error
	return locations, err
}

func (r *locationRepository) GetByCode(ctx context.Context, zoneID uuid.UUID, code string) (*entity.Location, error) {
	var location entity.Location
	err := r.db.WithContext(ctx).
		Where("zone_id = ? AND code = ?", zoneID, code).
		First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) Update(ctx context.Context, location *entity.Location) error {
	return r.db.WithContext(ctx).Save(location).Error
}

func (r *locationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Location{}, "id = ?", id).Error
}
