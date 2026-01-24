package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type shipmentRepository struct {
	db *gorm.DB
}

// NewShipmentRepository creates a new shipment repository
func NewShipmentRepository(db *gorm.DB) repository.ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(ctx context.Context, shipment *entity.Shipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

func (r *shipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Shipment, error) {
	var shipment entity.Shipment
	err := r.db.WithContext(ctx).
		Preload("SalesOrder").
		First(&shipment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) GetByNumber(ctx context.Context, number string) (*entity.Shipment, error) {
	var shipment entity.Shipment
	err := r.db.WithContext(ctx).
		Preload("SalesOrder").
		First(&shipment, "shipment_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) Update(ctx context.Context, shipment *entity.Shipment) error {
	shipment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(shipment).Error
}

func (r *shipmentRepository) List(ctx context.Context, filter *repository.ShipmentFilter) ([]*entity.Shipment, int64, error) {
	var shipments []*entity.Shipment
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Shipment{})

	// Apply filters
	if filter.SalesOrderID != nil {
		query = query.Where("sales_order_id = ?", filter.SalesOrderID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Carrier != "" {
		query = query.Where("carrier = ?", filter.Carrier)
	}
	if filter.DateFrom != "" {
		query = query.Where("shipped_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("shipped_date <= ?", filter.DateTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Get results
	err := query.Preload("SalesOrder").Order("created_at DESC").Find(&shipments).Error
	return shipments, total, err
}

func (r *shipmentRepository) GetNextShipmentNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).
		Model(&entity.Shipment{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count)
	return fmt.Sprintf("SHP-%d-%04d", year, count+1), nil
}

func (r *shipmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ShipmentStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Shipment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *shipmentRepository) GetBySalesOrder(ctx context.Context, salesOrderID uuid.UUID) ([]*entity.Shipment, error) {
	var shipments []*entity.Shipment
	err := r.db.WithContext(ctx).
		Where("sales_order_id = ?", salesOrderID).
		Order("created_at DESC").
		Find(&shipments).Error
	return shipments, err
}

func (r *shipmentRepository) GetByTrackingNumber(ctx context.Context, trackingNumber string) (*entity.Shipment, error) {
	var shipment entity.Shipment
	err := r.db.WithContext(ctx).
		First(&shipment, "tracking_number = ?", trackingNumber).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

// Return Repository
type returnRepository struct {
	db *gorm.DB
}

// NewReturnRepository creates a new return repository
func NewReturnRepository(db *gorm.DB) repository.ReturnRepository {
	return &returnRepository{db: db}
}

func (r *returnRepository) Create(ctx context.Context, ret *entity.Return) error {
	return r.db.WithContext(ctx).Create(ret).Error
}

func (r *returnRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Return, error) {
	var ret entity.Return
	err := r.db.WithContext(ctx).
		Preload("SalesOrder").
		Preload("LineItems").
		First(&ret, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (r *returnRepository) GetByNumber(ctx context.Context, number string) (*entity.Return, error) {
	var ret entity.Return
	err := r.db.WithContext(ctx).
		Preload("SalesOrder").
		Preload("LineItems").
		First(&ret, "return_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (r *returnRepository) Update(ctx context.Context, ret *entity.Return) error {
	ret.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(ret).Error
}

func (r *returnRepository) List(ctx context.Context, filter *repository.ReturnFilter) ([]*entity.Return, int64, error) {
	var returns []*entity.Return
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Return{})

	// Apply filters
	if filter.SalesOrderID != nil {
		query = query.Where("sales_order_id = ?", filter.SalesOrderID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ReturnType != "" {
		query = query.Where("return_type = ?", filter.ReturnType)
	}
	if filter.DateFrom != "" {
		query = query.Where("return_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("return_date <= ?", filter.DateTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Get results
	err := query.Preload("SalesOrder").Order("created_at DESC").Find(&returns).Error
	return returns, total, err
}

func (r *returnRepository) GetNextReturnNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).
		Model(&entity.Return{}).
		Where("EXTRACT(YEAR FROM return_date) = ?", year).
		Count(&count)
	return fmt.Sprintf("RET-%d-%04d", year, count+1), nil
}

func (r *returnRepository) CreateLineItem(ctx context.Context, item *entity.ReturnLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *returnRepository) GetLineItems(ctx context.Context, returnID uuid.UUID) ([]*entity.ReturnLineItem, error) {
	var items []*entity.ReturnLineItem
	err := r.db.WithContext(ctx).
		Where("return_id = ?", returnID).
		Find(&items).Error
	return items, err
}

func (r *returnRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ReturnStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Return{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *returnRepository) GetBySalesOrder(ctx context.Context, salesOrderID uuid.UUID) ([]*entity.Return, error) {
	var returns []*entity.Return
	err := r.db.WithContext(ctx).
		Where("sales_order_id = ?", salesOrderID).
		Order("created_at DESC").
		Find(&returns).Error
	return returns, err
}
