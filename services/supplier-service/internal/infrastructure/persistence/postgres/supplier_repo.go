package postgres

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type supplierRepository struct {
	db *gorm.DB
}

// NewSupplierRepository creates a new supplier repository
func NewSupplierRepository(db *gorm.DB) repository.SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *entity.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

func (r *supplierRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		Preload("Contacts").
		Preload("Certifications").
		First(&supplier, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepository) GetByCode(ctx context.Context, code string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		Preload("Contacts").
		Preload("Certifications").
		First(&supplier, "code = ? AND deleted_at IS NULL", code).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *entity.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

func (r *supplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Supplier{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *supplierRepository) List(ctx context.Context, filter *repository.SupplierFilter) ([]*entity.Supplier, int64, error) {
	var suppliers []*entity.Supplier
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Supplier{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter.SupplierType != "" {
		query = query.Where("supplier_type = ?", filter.SupplierType)
	}
	if filter.BusinessType != "" {
		query = query.Where("business_type = ?", filter.BusinessType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR tax_code ILIKE ?", search, search, search)
	}
	if filter.MinRating != nil {
		query = query.Where("overall_rating >= ?", *filter.MinRating)
	}
	if filter.HasGMP != nil && *filter.HasGMP {
		query = query.Where(`EXISTS (
			SELECT 1 FROM supplier_certifications 
			WHERE supplier_certifications.supplier_id = suppliers.id 
			AND certification_type = 'GMP' 
			AND status = 'VALID'
		)`)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(20)
	}
	if filter.Page > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset)
	}

	// Execute query
	if err := query.Order("created_at DESC").Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (r *supplierRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SupplierStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Supplier{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *supplierRepository) GetNextCode(ctx context.Context) (string, error) {
	var maxCode string
	err := r.db.WithContext(ctx).
		Model(&entity.Supplier{}).
		Select("code").
		Order("code DESC").
		Limit(1).
		Pluck("code", &maxCode).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if maxCode == "" {
		return "SUP-0001", nil
	}

	// Parse code and increment
	var seq int
	fmt.Sscanf(maxCode, "SUP-%04d", &seq)
	return fmt.Sprintf("SUP-%04d", seq+1), nil
}

func (r *supplierRepository) UpdateRating(ctx context.Context, id uuid.UUID, quality, delivery, service, overall float64) error {
	return r.db.WithContext(ctx).
		Model(&entity.Supplier{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"quality_rating":  quality,
			"delivery_rating": delivery,
			"service_rating":  service,
			"overall_rating":  overall,
		}).Error
}
