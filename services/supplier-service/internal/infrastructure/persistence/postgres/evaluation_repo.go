package postgres

import (
	"context"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type evaluationRepository struct {
	db *gorm.DB
}

// NewEvaluationRepository creates a new evaluation repository
func NewEvaluationRepository(db *gorm.DB) repository.EvaluationRepository {
	return &evaluationRepository{db: db}
}

func (r *evaluationRepository) Create(ctx context.Context, eval *entity.Evaluation) error {
	eval.CalculateOverallScore()
	return r.db.WithContext(ctx).Create(eval).Error
}

func (r *evaluationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Evaluation, error) {
	var eval entity.Evaluation
	err := r.db.WithContext(ctx).First(&eval, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &eval, nil
}

func (r *evaluationRepository) Update(ctx context.Context, eval *entity.Evaluation) error {
	eval.CalculateOverallScore()
	return r.db.WithContext(ctx).Save(eval).Error
}

func (r *evaluationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Evaluation{}, "id = ?", id).Error
}

func (r *evaluationRepository) GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Evaluation, error) {
	var evals []*entity.Evaluation
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("evaluation_date DESC").
		Find(&evals).Error
	return evals, err
}

func (r *evaluationRepository) List(ctx context.Context, filter *repository.EvaluationFilter) ([]*entity.Evaluation, int64, error) {
	var evals []*entity.Evaluation
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Evaluation{})

	if filter.SupplierID != nil {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}
	if filter.Period != "" {
		query = query.Where("evaluation_period = ?", filter.Period)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Page > 0 && filter.Limit > 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if err := query.Order("evaluation_date DESC").Find(&evals).Error; err != nil {
		return nil, 0, err
	}

	return evals, total, nil
}

func (r *evaluationRepository) GetLatestBySupplierID(ctx context.Context, supplierID uuid.UUID) (*entity.Evaluation, error) {
	var eval entity.Evaluation
	err := r.db.WithContext(ctx).
		Where("supplier_id = ? AND status = ?", supplierID, entity.EvaluationStatusApproved).
		Order("evaluation_date DESC").
		First(&eval).Error
	if err != nil {
		return nil, err
	}
	return &eval, nil
}

func (r *evaluationRepository) GetAverageScoresBySupplierID(ctx context.Context, supplierID uuid.UUID) (quality, delivery, service float64, err error) {
	var result struct {
		AvgQuality  float64
		AvgDelivery float64
		AvgService  float64
	}
	
	err = r.db.WithContext(ctx).
		Model(&entity.Evaluation{}).
		Select("AVG(quality_score) as avg_quality, AVG(delivery_score) as avg_delivery, AVG(service_score) as avg_service").
		Where("supplier_id = ? AND status = ?", supplierID, entity.EvaluationStatusApproved).
		Scan(&result).Error
	
	return result.AvgQuality, result.AvgDelivery, result.AvgService, err
}
