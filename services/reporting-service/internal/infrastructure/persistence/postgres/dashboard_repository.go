package postgres

import (
	"context"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/erp-cosmetics/reporting-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates new repository
func NewDashboardRepository(db *gorm.DB) repository.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) Create(ctx context.Context, dashboard *entity.Dashboard) error {
	return r.db.WithContext(ctx).Create(dashboard).Error
}

func (r *dashboardRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error) {
	var dashboard entity.Dashboard
	err := r.db.WithContext(ctx).First(&dashboard, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (r *dashboardRepository) GetByCode(ctx context.Context, code string) (*entity.Dashboard, error) {
	var dashboard entity.Dashboard
	err := r.db.WithContext(ctx).First(&dashboard, "code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (r *dashboardRepository) Update(ctx context.Context, dashboard *entity.Dashboard) error {
	return r.db.WithContext(ctx).Save(dashboard).Error
}

func (r *dashboardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Dashboard{}, "id = ?", id).Error
}

func (r *dashboardRepository) List(ctx context.Context, limit, offset int) ([]*entity.Dashboard, int64, error) {
	var dashboards []*entity.Dashboard
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Dashboard{}).Where("is_active = ?", true)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&dashboards).Error
	return dashboards, total, err
}

func (r *dashboardRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Dashboard, error) {
	var dashboards []*entity.Dashboard
	err := r.db.WithContext(ctx).
		Where("(created_by = ? OR visibility = 'PUBLIC') AND is_active = ?", userID, true).
		Order("name ASC").
		Find(&dashboards).Error
	return dashboards, err
}

func (r *dashboardRepository) GetDefault(ctx context.Context) (*entity.Dashboard, error) {
	var dashboard entity.Dashboard
	err := r.db.WithContext(ctx).First(&dashboard, "is_default = ? AND is_active = ?", true, true).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (r *dashboardRepository) GetWithWidgets(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error) {
	var dashboard entity.Dashboard
	err := r.db.WithContext(ctx).
		Preload("Widgets", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_visible = ?", true).Order("position_y, position_x")
		}).
		First(&dashboard, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}
