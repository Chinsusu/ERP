package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type inventoryCountRepository struct {
	db *gorm.DB
}

// NewInventoryCountRepository creates a new inventory count repository
func NewInventoryCountRepository(db *gorm.DB) repository.InventoryCountRepository {
	return &inventoryCountRepository{db: db}
}

func (r *inventoryCountRepository) Create(ctx context.Context, count *entity.InventoryCount) error {
	return r.db.WithContext(ctx).Create(count).Error
}

func (r *inventoryCountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.InventoryCount, error) {
	var count entity.InventoryCount
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("Zone").
		Preload("LineItems").
		Preload("LineItems.Location").
		Preload("LineItems.Lot").
		First(&count, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *inventoryCountRepository) GetByNumber(ctx context.Context, countNumber string) (*entity.InventoryCount, error) {
	var count entity.InventoryCount
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		First(&count, "count_number = ?", countNumber).Error
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *inventoryCountRepository) List(ctx context.Context, filter *repository.InventoryCountFilter) ([]*entity.InventoryCount, int64, error) {
	var counts []*entity.InventoryCount
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.InventoryCount{})

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.ZoneID != nil {
		query = query.Where("zone_id = ?", *filter.ZoneID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CountType != "" {
		query = query.Where("count_type = ?", filter.CountType)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("count_number ILIKE ?", search)
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

	if err := query.
		Preload("Warehouse").
		Order("created_at DESC").
		Find(&counts).Error; err != nil {
		return nil, 0, err
	}

	return counts, total, nil
}

func (r *inventoryCountRepository) Update(ctx context.Context, count *entity.InventoryCount) error {
	count.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(count).Error
}

func (r *inventoryCountRepository) CreateLineItem(ctx context.Context, item *entity.InventoryCountLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *inventoryCountRepository) CreateLineItems(ctx context.Context, items []*entity.InventoryCountLineItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *inventoryCountRepository) GetLineItemsByCountID(ctx context.Context, countID uuid.UUID) ([]*entity.InventoryCountLineItem, error) {
	var items []*entity.InventoryCountLineItem
	err := r.db.WithContext(ctx).
		Where("inventory_count_id = ?", countID).
		Preload("Location").
		Preload("Lot").
		Find(&items).Error
	return items, err
}

func (r *inventoryCountRepository) UpdateLineItem(ctx context.Context, item *entity.InventoryCountLineItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *inventoryCountRepository) GetPendingItems(ctx context.Context, countID uuid.UUID) ([]*entity.InventoryCountLineItem, error) {
	var items []*entity.InventoryCountLineItem
	err := r.db.WithContext(ctx).
		Where("inventory_count_id = ? AND is_counted = ?", countID, false).
		Preload("Location").
		Preload("Lot").
		Find(&items).Error
	return items, err
}

func (r *inventoryCountRepository) GetVarianceItems(ctx context.Context, countID uuid.UUID) ([]*entity.InventoryCountLineItem, error) {
	var items []*entity.InventoryCountLineItem
	err := r.db.WithContext(ctx).
		Where("inventory_count_id = ? AND is_counted = ? AND variance != 0", countID, true).
		Preload("Location").
		Preload("Lot").
		Find(&items).Error
	return items, err
}

func (r *inventoryCountRepository) GetNextCountNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()

	r.db.WithContext(ctx).
		Model(&entity.InventoryCount{}).
		Where("count_number LIKE ?", fmt.Sprintf("IC-%d-%%", year)).
		Count(&count)

	return fmt.Sprintf("IC-%d-%04d", year, count+1), nil
}
