package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bomRepository struct {
	db *gorm.DB
}

// NewBOMRepository creates a new BOM repository
func NewBOMRepository(db *gorm.DB) repository.BOMRepository {
	return &bomRepository{db: db}
}

func (r *bomRepository) Create(ctx context.Context, bom *entity.BOM) error {
	return r.db.WithContext(ctx).Create(bom).Error
}

func (r *bomRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.BOM, error) {
	var bom entity.BOM
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&bom, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

func (r *bomRepository) GetByNumber(ctx context.Context, bomNumber string) (*entity.BOM, error) {
	var bom entity.BOM
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&bom, "bom_number = ?", bomNumber).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

func (r *bomRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.BOM, error) {
	var boms []*entity.BOM
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("version DESC").
		Find(&boms).Error
	return boms, err
}

func (r *bomRepository) GetActiveBOMForProduct(ctx context.Context, productID uuid.UUID) (*entity.BOM, error) {
	var bom entity.BOM
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("product_id = ? AND status = ?", productID, entity.BOMStatusApproved).
		Where("effective_from <= ? OR effective_from IS NULL", time.Now()).
		Where("effective_to >= ? OR effective_to IS NULL", time.Now()).
		Order("version DESC").
		First(&bom).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

func (r *bomRepository) List(ctx context.Context, filter repository.BOMFilter) ([]*entity.BOM, int64, error) {
	var boms []*entity.BOM
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BOM{})

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("bom_number ILIKE ? OR name ILIKE ?", search, search)
	}

	query.Count(&total)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Order("created_at DESC").Find(&boms).Error
	return boms, total, err
}

func (r *bomRepository) Update(ctx context.Context, bom *entity.BOM) error {
	return r.db.WithContext(ctx).Save(bom).Error
}

func (r *bomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.BOM{}, "id = ?", id).Error
}

// Line items
func (r *bomRepository) CreateLineItem(ctx context.Context, item *entity.BOMLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *bomRepository) GetLineItems(ctx context.Context, bomID uuid.UUID) ([]*entity.BOMLineItem, error) {
	var items []*entity.BOMLineItem
	err := r.db.WithContext(ctx).
		Where("bom_id = ?", bomID).
		Order("line_number ASC").
		Find(&items).Error
	return items, err
}

func (r *bomRepository) UpdateLineItem(ctx context.Context, item *entity.BOMLineItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *bomRepository) DeleteLineItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.BOMLineItem{}, "id = ?", id).Error
}

// Versioning
func (r *bomRepository) CreateVersion(ctx context.Context, version *entity.BOMVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *bomRepository) GetVersions(ctx context.Context, bomID uuid.UUID) ([]*entity.BOMVersion, error) {
	var versions []*entity.BOMVersion
	err := r.db.WithContext(ctx).
		Where("bom_id = ?", bomID).
		Order("version DESC").
		Find(&versions).Error
	return versions, err
}

// WorkOrder Repository
type workOrderRepository struct {
	db *gorm.DB
}

// NewWorkOrderRepository creates a new work order repository
func NewWorkOrderRepository(db *gorm.DB) repository.WorkOrderRepository {
	return &workOrderRepository{db: db}
}

func (r *workOrderRepository) Create(ctx context.Context, wo *entity.WorkOrder) error {
	return r.db.WithContext(ctx).Create(wo).Error
}

func (r *workOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.WorkOrder, error) {
	var wo entity.WorkOrder
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("MaterialIssues").
		First(&wo, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wo, nil
}

func (r *workOrderRepository) GetByNumber(ctx context.Context, woNumber string) (*entity.WorkOrder, error) {
	var wo entity.WorkOrder
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&wo, "wo_number = ?", woNumber).Error
	if err != nil {
		return nil, err
	}
	return &wo, nil
}

func (r *workOrderRepository) List(ctx context.Context, filter repository.WOFilter) ([]*entity.WorkOrder, int64, error) {
	var wos []*entity.WorkOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.WorkOrder{})

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.BOMID != nil {
		query = query.Where("bom_id = ?", *filter.BOMID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("wo_number ILIKE ? OR batch_number ILIKE ?", search, search)
	}

	query.Count(&total)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Order("created_at DESC").Find(&wos).Error
	return wos, total, err
}

func (r *workOrderRepository) Update(ctx context.Context, wo *entity.WorkOrder) error {
	return r.db.WithContext(ctx).Save(wo).Error
}

func (r *workOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.WorkOrder{}, "id = ?", id).Error
}

func (r *workOrderRepository) CreateLineItems(ctx context.Context, items []*entity.WOLineItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *workOrderRepository) GetLineItems(ctx context.Context, woID uuid.UUID) ([]*entity.WOLineItem, error) {
	var items []*entity.WOLineItem
	err := r.db.WithContext(ctx).
		Where("work_order_id = ?", woID).
		Order("line_number ASC").
		Find(&items).Error
	return items, err
}

func (r *workOrderRepository) UpdateLineItem(ctx context.Context, item *entity.WOLineItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *workOrderRepository) CreateMaterialIssue(ctx context.Context, issue *entity.WOMaterialIssue) error {
	return r.db.WithContext(ctx).Create(issue).Error
}

func (r *workOrderRepository) GetMaterialIssues(ctx context.Context, woID uuid.UUID) ([]*entity.WOMaterialIssue, error) {
	var issues []*entity.WOMaterialIssue
	err := r.db.WithContext(ctx).
		Where("work_order_id = ?", woID).
		Order("created_at ASC").
		Find(&issues).Error
	return issues, err
}

func (r *workOrderRepository) GenerateWONumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	r.db.WithContext(ctx).Model(&entity.WorkOrder{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count)
	return fmt.Sprintf("WO-%d-%04d", year, count+1), nil
}

func (r *workOrderRepository) GenerateIssueNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	r.db.WithContext(ctx).Model(&entity.WOMaterialIssue{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count)
	return fmt.Sprintf("ISS-%d-%04d", year, count+1), nil
}
