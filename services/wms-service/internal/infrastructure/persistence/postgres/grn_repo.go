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

type grnRepository struct {
	db *gorm.DB
}

// NewGRNRepository creates a new GRN repository
func NewGRNRepository(db *gorm.DB) repository.GRNRepository {
	return &grnRepository{db: db}
}

func (r *grnRepository) Create(ctx context.Context, grn *entity.GRN) error {
	return r.db.WithContext(ctx).Create(grn).Error
}

func (r *grnRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.GRN, error) {
	var grn entity.GRN
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("LineItems").
		Preload("LineItems.Lot").
		Preload("LineItems.Location").
		First(&grn, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &grn, nil
}

func (r *grnRepository) GetByNumber(ctx context.Context, grnNumber string) (*entity.GRN, error) {
	var grn entity.GRN
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		First(&grn, "grn_number = ?", grnNumber).Error
	if err != nil {
		return nil, err
	}
	return &grn, nil
}

func (r *grnRepository) GetByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.GRN, error) {
	var grns []*entity.GRN
	err := r.db.WithContext(ctx).
		Where("po_id = ?", poID).
		Order("created_at DESC").
		Find(&grns).Error
	return grns, err
}

func (r *grnRepository) List(ctx context.Context, filter *repository.GRNFilter) ([]*entity.GRN, int64, error) {
	var grns []*entity.GRN
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.GRN{})

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.SupplierID != nil {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}
	if filter.POID != nil {
		query = query.Where("po_id = ?", *filter.POID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.QCStatus != "" {
		query = query.Where("qc_status = ?", filter.QCStatus)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("grn_number ILIKE ? OR po_number ILIKE ?", search, search)
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
		Find(&grns).Error; err != nil {
		return nil, 0, err
	}

	return grns, total, nil
}

func (r *grnRepository) Update(ctx context.Context, grn *entity.GRN) error {
	grn.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(grn).Error
}

func (r *grnRepository) CreateLineItem(ctx context.Context, item *entity.GRNLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *grnRepository) GetLineItemsByGRNID(ctx context.Context, grnID uuid.UUID) ([]*entity.GRNLineItem, error) {
	var items []*entity.GRNLineItem
	err := r.db.WithContext(ctx).
		Where("grn_id = ?", grnID).
		Preload("Lot").
		Preload("Location").
		Order("line_number").
		Find(&items).Error
	return items, err
}

func (r *grnRepository) UpdateLineItem(ctx context.Context, item *entity.GRNLineItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *grnRepository) GetNextGRNNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()

	r.db.WithContext(ctx).
		Model(&entity.GRN{}).
		Where("grn_number LIKE ?", fmt.Sprintf("GRN-%d-%%", year)).
		Count(&count)

	return fmt.Sprintf("GRN-%d-%04d", year, count+1), nil
}

// Goods Issue Repository
type goodsIssueRepository struct {
	db *gorm.DB
}

// NewGoodsIssueRepository creates a new goods issue repository
func NewGoodsIssueRepository(db *gorm.DB) repository.GoodsIssueRepository {
	return &goodsIssueRepository{db: db}
}

func (r *goodsIssueRepository) Create(ctx context.Context, issue *entity.GoodsIssue) error {
	return r.db.WithContext(ctx).Create(issue).Error
}

func (r *goodsIssueRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.GoodsIssue, error) {
	var issue entity.GoodsIssue
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("LineItems").
		Preload("LineItems.Lot").
		Preload("LineItems.Location").
		First(&issue, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *goodsIssueRepository) GetByNumber(ctx context.Context, issueNumber string) (*entity.GoodsIssue, error) {
	var issue entity.GoodsIssue
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		First(&issue, "issue_number = ?", issueNumber).Error
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *goodsIssueRepository) List(ctx context.Context, filter *repository.GoodsIssueFilter) ([]*entity.GoodsIssue, int64, error) {
	var issues []*entity.GoodsIssue
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.GoodsIssue{})

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.IssueType != "" {
		query = query.Where("issue_type = ?", filter.IssueType)
	}
	if filter.ReferenceType != "" {
		query = query.Where("reference_type = ?", filter.ReferenceType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("issue_number ILIKE ? OR reference_number ILIKE ?", search, search)
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
		Find(&issues).Error; err != nil {
		return nil, 0, err
	}

	return issues, total, nil
}

func (r *goodsIssueRepository) Update(ctx context.Context, issue *entity.GoodsIssue) error {
	issue.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(issue).Error
}

func (r *goodsIssueRepository) CreateLineItem(ctx context.Context, item *entity.GILineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *goodsIssueRepository) CreateLineItems(ctx context.Context, items []*entity.GILineItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *goodsIssueRepository) GetLineItemsByIssueID(ctx context.Context, issueID uuid.UUID) ([]*entity.GILineItem, error) {
	var items []*entity.GILineItem
	err := r.db.WithContext(ctx).
		Where("goods_issue_id = ?", issueID).
		Preload("Lot").
		Preload("Location").
		Order("line_number").
		Find(&items).Error
	return items, err
}

func (r *goodsIssueRepository) GetNextIssueNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()

	r.db.WithContext(ctx).
		Model(&entity.GoodsIssue{}).
		Where("issue_number LIKE ?", fmt.Sprintf("GI-%d-%%", year)).
		Count(&count)

	return fmt.Sprintf("GI-%d-%04d", year, count+1), nil
}
