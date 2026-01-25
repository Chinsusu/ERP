package postgres

import (
	"context"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/erp-cosmetics/reporting-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reportDefinitionRepository struct {
	db *gorm.DB
}

// NewReportDefinitionRepository creates new repository
func NewReportDefinitionRepository(db *gorm.DB) repository.ReportDefinitionRepository {
	return &reportDefinitionRepository{db: db}
}

func (r *reportDefinitionRepository) Create(ctx context.Context, report *entity.ReportDefinition) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportDefinitionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportDefinition, error) {
	var report entity.ReportDefinition
	err := r.db.WithContext(ctx).First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportDefinitionRepository) GetByCode(ctx context.Context, code string) (*entity.ReportDefinition, error) {
	var report entity.ReportDefinition
	err := r.db.WithContext(ctx).First(&report, "code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportDefinitionRepository) Update(ctx context.Context, report *entity.ReportDefinition) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportDefinitionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ReportDefinition{}, "id = ?", id).Error
}

func (r *reportDefinitionRepository) List(ctx context.Context, limit, offset int) ([]*entity.ReportDefinition, int64, error) {
	var reports []*entity.ReportDefinition
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ReportDefinition{}).Where("is_active = ?", true)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&reports).Error
	return reports, total, err
}

func (r *reportDefinitionRepository) ListByType(ctx context.Context, reportType string) ([]*entity.ReportDefinition, error) {
	var reports []*entity.ReportDefinition
	err := r.db.WithContext(ctx).
		Where("report_type = ? AND is_active = ?", reportType, true).
		Order("name ASC").
		Find(&reports).Error
	return reports, err
}

func (r *reportDefinitionRepository) ListActive(ctx context.Context) ([]*entity.ReportDefinition, error) {
	var reports []*entity.ReportDefinition
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&reports).Error
	return reports, err
}
