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

type prRepository struct {
	db *gorm.DB
}

// NewPRRepository creates a new PR repository
func NewPRRepository(db *gorm.DB) repository.PRRepository {
	return &prRepository{db: db}
}

func (r *prRepository) Create(ctx context.Context, pr *entity.PurchaseRequisition) error {
	return r.db.WithContext(ctx).Create(pr).Error
}

func (r *prRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PurchaseRequisition, error) {
	var pr entity.PurchaseRequisition
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("Approvals").
		First(&pr, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *prRepository) GetByNumber(ctx context.Context, prNumber string) (*entity.PurchaseRequisition, error) {
	var pr entity.PurchaseRequisition
	err := r.db.WithContext(ctx).
		Preload("LineItems").
		First(&pr, "pr_number = ? AND deleted_at IS NULL", prNumber).Error
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *prRepository) Update(ctx context.Context, pr *entity.PurchaseRequisition) error {
	pr.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(pr).Error
}

func (r *prRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.PurchaseRequisition{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *prRepository) List(ctx context.Context, filter *repository.PRFilter) ([]*entity.PurchaseRequisition, int64, error) {
	var prs []*entity.PurchaseRequisition
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PurchaseRequisition{}).Where("deleted_at IS NULL")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.RequesterID != nil {
		query = query.Where("requester_id = ?", *filter.RequesterID)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("pr_number ILIKE ?", search)
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

	if err := query.Order("created_at DESC").Find(&prs).Error; err != nil {
		return nil, 0, err
	}

	return prs, total, nil
}

func (r *prRepository) GetNextPRNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	
	r.db.WithContext(ctx).
		Model(&entity.PurchaseRequisition{}).
		Where("pr_number LIKE ?", fmt.Sprintf("PR-%d-%%", year)).
		Count(&count)

	return fmt.Sprintf("PR-%d-%04d", year, count+1), nil
}

func (r *prRepository) CreateLineItem(ctx context.Context, item *entity.PRLineItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *prRepository) UpdateLineItem(ctx context.Context, item *entity.PRLineItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *prRepository) DeleteLineItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.PRLineItem{}, "id = ?", id).Error
}

func (r *prRepository) GetLineItemsByPRID(ctx context.Context, prID uuid.UUID) ([]*entity.PRLineItem, error) {
	var items []*entity.PRLineItem
	err := r.db.WithContext(ctx).Where("pr_id = ?", prID).Order("line_number").Find(&items).Error
	return items, err
}

func (r *prRepository) CreateApproval(ctx context.Context, approval *entity.PRApproval) error {
	return r.db.WithContext(ctx).Create(approval).Error
}

func (r *prRepository) GetApprovalsByPRID(ctx context.Context, prID uuid.UUID) ([]*entity.PRApproval, error) {
	var approvals []*entity.PRApproval
	err := r.db.WithContext(ctx).Where("pr_id = ?", prID).Order("created_at DESC").Find(&approvals).Error
	return approvals, err
}
