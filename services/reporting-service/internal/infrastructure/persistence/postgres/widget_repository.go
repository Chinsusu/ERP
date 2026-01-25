package postgres

import (
	"context"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/erp-cosmetics/reporting-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type widgetRepository struct {
	db *gorm.DB
}

// NewWidgetRepository creates new repository
func NewWidgetRepository(db *gorm.DB) repository.WidgetRepository {
	return &widgetRepository{db: db}
}

func (r *widgetRepository) Create(ctx context.Context, widget *entity.Widget) error {
	return r.db.WithContext(ctx).Create(widget).Error
}

func (r *widgetRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Widget, error) {
	var widget entity.Widget
	err := r.db.WithContext(ctx).First(&widget, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &widget, nil
}

func (r *widgetRepository) Update(ctx context.Context, widget *entity.Widget) error {
	return r.db.WithContext(ctx).Save(widget).Error
}

func (r *widgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Widget{}, "id = ?", id).Error
}

func (r *widgetRepository) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]*entity.Widget, error) {
	var widgets []*entity.Widget
	err := r.db.WithContext(ctx).
		Where("dashboard_id = ? AND is_visible = ?", dashboardID, true).
		Order("position_y, position_x").
		Find(&widgets).Error
	return widgets, err
}

func (r *widgetRepository) DeleteByDashboardID(ctx context.Context, dashboardID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Widget{}, "dashboard_id = ?", dashboardID).Error
}
