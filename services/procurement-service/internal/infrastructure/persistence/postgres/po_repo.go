package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type poRepository struct {
	db *gorm.DB
}

// NewPORepository creates a new PO repository
func NewPORepository(db *gorm.DB) repository.PORepository {
	return &poRepository{db: db}
}

func (r *poRepository) Create(ctx context.Context, po *entity.PurchaseOrder) error {
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *poRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("Amendments").
		Preload("Receipts").
		First(&po, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *poRepository) GetByNumber(ctx context.Context, poNumber string) (*entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		First(&po, "po_number = ? AND deleted_at IS NULL", poNumber).Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *poRepository) Update(ctx context.Context, po *entity.PurchaseOrder) error {
	po.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *poRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.PurchaseOrder{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *poRepository) List(ctx context.Context, filter *repository.POFilter) ([]*entity.PurchaseOrder, int64, error) {
	var pos []*entity.PurchaseOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PurchaseOrder{}).Where("deleted_at IS NULL")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.SupplierID != nil {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}
	if filter.PRID != nil {
		query = query.Where("pr_id = ?", *filter.PRID)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("po_number ILIKE ? OR supplier_name ILIKE ?", search, search)
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

	if err := query.Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return pos, total, nil
}

func (r *poRepository) GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.PurchaseOrder, error) {
	var pos []*entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Where("supplier_id = ? AND deleted_at IS NULL", supplierID).
		Order("created_at DESC").
		Find(&pos).Error
	return pos, err
}

func (r *poRepository) GetNextPONumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	
	r.db.WithContext(ctx).
		Model(&entity.PurchaseOrder{}).
		Where("po_number LIKE ?", fmt.Sprintf("PO-%d-%%", year)).
		Count(&count)

	return fmt.Sprintf("PO-%d-%04d", year, count+1), nil
}

func (r *poRepository) CreateLineItem(ctx context.Context, item *entity.POLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *poRepository) UpdateLineItem(ctx context.Context, item *entity.POLineItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *poRepository) DeleteLineItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.POLineItem{}, "id = ?", id).Error
}

func (r *poRepository) GetLineItemsByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.POLineItem, error) {
	var items []*entity.POLineItem
	err := r.db.WithContext(ctx).Where("po_id = ?", poID).Order("line_number").Find(&items).Error
	return items, err
}

func (r *poRepository) GetLineItemByID(ctx context.Context, id uuid.UUID) (*entity.POLineItem, error) {
	var item entity.POLineItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *poRepository) CreateAmendment(ctx context.Context, amendment *entity.POAmendment) error {
	return r.db.WithContext(ctx).Create(amendment).Error
}

func (r *poRepository) GetAmendmentsByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.POAmendment, error) {
	var amendments []*entity.POAmendment
	err := r.db.WithContext(ctx).Where("po_id = ?", poID).Order("amendment_number DESC").Find(&amendments).Error
	return amendments, err
}

func (r *poRepository) CreateReceipt(ctx context.Context, receipt *entity.POReceipt) error {
	return r.db.WithContext(ctx).Create(receipt).Error
}

func (r *poRepository) GetReceiptsByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.POReceipt, error) {
	var receipts []*entity.POReceipt
	err := r.db.WithContext(ctx).Where("po_id = ?", poID).Order("received_date DESC").Find(&receipts).Error
	return receipts, err
}

func (r *poRepository) GetReceiptsByLineItemID(ctx context.Context, lineItemID uuid.UUID) ([]*entity.POReceipt, error) {
	var receipts []*entity.POReceipt
	err := r.db.WithContext(ctx).Where("po_line_item_id = ?", lineItemID).Order("received_date DESC").Find(&receipts).Error
	return receipts, err
}
