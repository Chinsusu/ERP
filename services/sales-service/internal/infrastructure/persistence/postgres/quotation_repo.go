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

type quotationRepository struct {
	db *gorm.DB
}

// NewQuotationRepository creates a new quotation repository
func NewQuotationRepository(db *gorm.DB) repository.QuotationRepository {
	return &quotationRepository{db: db}
}

func (r *quotationRepository) Create(ctx context.Context, quotation *entity.Quotation) error {
	return r.db.WithContext(ctx).Create(quotation).Error
}

func (r *quotationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Quotation, error) {
	var quotation entity.Quotation
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("LineItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_number ASC")
		}).
		First(&quotation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &quotation, nil
}

func (r *quotationRepository) GetByNumber(ctx context.Context, number string) (*entity.Quotation, error) {
	var quotation entity.Quotation
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("LineItems").
		First(&quotation, "quotation_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &quotation, nil
}

func (r *quotationRepository) Update(ctx context.Context, quotation *entity.Quotation) error {
	quotation.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(quotation).Error
}

func (r *quotationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Quotation{}, "id = ?", id).Error
}

func (r *quotationRepository) List(ctx context.Context, filter *repository.QuotationFilter) ([]*entity.Quotation, int64, error) {
	var quotations []*entity.Quotation
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Quotation{})

	// Apply filters
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		query = query.Where("quotation_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("quotation_date <= ?", filter.DateTo)
	}
	if filter.ExpiringSoon {
		threshold := time.Now().AddDate(0, 0, 7)
		query = query.Where("valid_until <= ? AND status IN ('DRAFT', 'SENT')", threshold)
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
	err := query.Preload("Customer").Order("created_at DESC").Find(&quotations).Error
	return quotations, total, err
}

func (r *quotationRepository) GetNextQuotationNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).
		Model(&entity.Quotation{}).
		Where("EXTRACT(YEAR FROM quotation_date) = ?", year).
		Count(&count)
	return fmt.Sprintf("QT-%d-%04d", year, count+1), nil
}

// Line items
func (r *quotationRepository) CreateLineItem(ctx context.Context, item *entity.QuotationLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *quotationRepository) UpdateLineItem(ctx context.Context, item *entity.QuotationLineItem) error {
	item.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *quotationRepository) DeleteLineItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.QuotationLineItem{}, "id = ?", itemID).Error
}

func (r *quotationRepository) GetLineItems(ctx context.Context, quotationID uuid.UUID) ([]*entity.QuotationLineItem, error) {
	var items []*entity.QuotationLineItem
	err := r.db.WithContext(ctx).
		Where("quotation_id = ?", quotationID).
		Order("line_number ASC").
		Find(&items).Error
	return items, err
}

func (r *quotationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.QuotationStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Quotation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *quotationRepository) GetExpiringQuotations(ctx context.Context, days int) ([]*entity.Quotation, error) {
	var quotations []*entity.Quotation
	threshold := time.Now().AddDate(0, 0, days)
	err := r.db.WithContext(ctx).
		Where("valid_until <= ? AND status IN ('DRAFT', 'SENT')", threshold).
		Preload("Customer").
		Find(&quotations).Error
	return quotations, err
}

func (r *quotationRepository) MarkExpiredQuotations(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&entity.Quotation{}).
		Where("valid_until < ? AND status IN ('DRAFT', 'SENT')", time.Now()).
		Updates(map[string]interface{}{
			"status":     entity.QuotationStatusExpired,
			"updated_at": time.Now(),
		})
	return result.RowsAffected, result.Error
}
