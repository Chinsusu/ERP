package postgres

import (
	"context"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/erp-cosmetics/reporting-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reportExecutionRepository struct {
	db *gorm.DB
}

// NewReportExecutionRepository creates new repository
func NewReportExecutionRepository(db *gorm.DB) repository.ReportExecutionRepository {
	return &reportExecutionRepository{db: db}
}

func (r *reportExecutionRepository) Create(ctx context.Context, execution *entity.ReportExecution) error {
	return r.db.WithContext(ctx).Create(execution).Error
}

func (r *reportExecutionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportExecution, error) {
	var execution entity.ReportExecution
	err := r.db.WithContext(ctx).Preload("Report").First(&execution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *reportExecutionRepository) Update(ctx context.Context, execution *entity.ReportExecution) error {
	return r.db.WithContext(ctx).Save(execution).Error
}

func (r *reportExecutionRepository) GetByReportID(ctx context.Context, reportID uuid.UUID, limit, offset int) ([]*entity.ReportExecution, int64, error) {
	var executions []*entity.ReportExecution
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ReportExecution{}).Where("report_id = ?", reportID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&executions).Error
	return executions, total, err
}

func (r *reportExecutionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.ReportExecution, int64, error) {
	var executions []*entity.ReportExecution
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ReportExecution{}).Where("created_by = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Report").Order("created_at DESC").Limit(limit).Offset(offset).Find(&executions).Error
	return executions, total, err
}

func (r *reportExecutionRepository) GetPending(ctx context.Context, limit int) ([]*entity.ReportExecution, error) {
	var executions []*entity.ReportExecution
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.ExecutionStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&executions).Error
	return executions, err
}
