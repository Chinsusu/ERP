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

type salesOrderRepository struct {
	db *gorm.DB
}

// NewSalesOrderRepository creates a new sales order repository
func NewSalesOrderRepository(db *gorm.DB) repository.SalesOrderRepository {
	return &salesOrderRepository{db: db}
}

func (r *salesOrderRepository) Create(ctx context.Context, order *entity.SalesOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *salesOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SalesOrder, error) {
	var order entity.SalesOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Quotation").
		Preload("LineItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_number ASC")
		}).
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *salesOrderRepository) GetByNumber(ctx context.Context, number string) (*entity.SalesOrder, error) {
	var order entity.SalesOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("LineItems").
		First(&order, "so_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *salesOrderRepository) Update(ctx context.Context, order *entity.SalesOrder) error {
	order.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *salesOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SalesOrder{}, "id = ?", id).Error
}

func (r *salesOrderRepository) List(ctx context.Context, filter *repository.SalesOrderFilter) ([]*entity.SalesOrder, int64, error) {
	var orders []*entity.SalesOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SalesOrder{})

	// Apply filters
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}
	if filter.DateFrom != "" {
		query = query.Where("so_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("so_date <= ?", filter.DateTo)
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
	err := query.Preload("Customer").Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *salesOrderRepository) GetNextSONumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).
		Model(&entity.SalesOrder{}).
		Where("EXTRACT(YEAR FROM so_date) = ?", year).
		Count(&count)
	return fmt.Sprintf("SO-%d-%04d", year, count+1), nil
}

// Line items
func (r *salesOrderRepository) CreateLineItem(ctx context.Context, item *entity.SOLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *salesOrderRepository) UpdateLineItem(ctx context.Context, item *entity.SOLineItem) error {
	item.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *salesOrderRepository) DeleteLineItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SOLineItem{}, "id = ?", itemID).Error
}

func (r *salesOrderRepository) GetLineItems(ctx context.Context, orderID uuid.UUID) ([]*entity.SOLineItem, error) {
	var items []*entity.SOLineItem
	err := r.db.WithContext(ctx).
		Where("sales_order_id = ?", orderID).
		Order("line_number ASC").
		Find(&items).Error
	return items, err
}

func (r *salesOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SOStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.SalesOrder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *salesOrderRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.SalesOrder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status": status,
			"updated_at":     time.Now(),
		}).Error
}

func (r *salesOrderRepository) UpdateShippedQuantity(ctx context.Context, lineItemID uuid.UUID, shippedQty float64) error {
	return r.db.WithContext(ctx).
		Model(&entity.SOLineItem{}).
		Where("id = ?", lineItemID).
		Updates(map[string]interface{}{
			"shipped_quantity": gorm.Expr("shipped_quantity + ?", shippedQty),
			"updated_at":       time.Now(),
		}).Error
}

func (r *salesOrderRepository) UpdateLineItemReservation(ctx context.Context, lineItemID uuid.UUID, reservationID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.SOLineItem{}).
		Where("id = ?", lineItemID).
		Updates(map[string]interface{}{
			"reservation_id": reservationID,
			"updated_at":     time.Now(),
		}).Error
}

func (r *salesOrderRepository) GetByCustomer(ctx context.Context, customerID uuid.UUID, limit int) ([]*entity.SalesOrder, error) {
	var orders []*entity.SalesOrder
	query := r.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&orders).Error
	return orders, err
}

func (r *salesOrderRepository) GetPendingOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]*entity.SalesOrder, error) {
	var orders []*entity.SalesOrder
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND status NOT IN ('DELIVERED', 'CANCELLED')", customerID).
		Find(&orders).Error
	return orders, err
}
