package repository

import (
	"context"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/google/uuid"
)

// ReportDefinitionRepository interface
type ReportDefinitionRepository interface {
	Create(ctx context.Context, report *entity.ReportDefinition) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportDefinition, error)
	GetByCode(ctx context.Context, code string) (*entity.ReportDefinition, error)
	Update(ctx context.Context, report *entity.ReportDefinition) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.ReportDefinition, int64, error)
	ListByType(ctx context.Context, reportType string) ([]*entity.ReportDefinition, error)
	ListActive(ctx context.Context) ([]*entity.ReportDefinition, error)
}

// ReportExecutionRepository interface
type ReportExecutionRepository interface {
	Create(ctx context.Context, execution *entity.ReportExecution) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportExecution, error)
	Update(ctx context.Context, execution *entity.ReportExecution) error
	GetByReportID(ctx context.Context, reportID uuid.UUID, limit, offset int) ([]*entity.ReportExecution, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.ReportExecution, int64, error)
	GetPending(ctx context.Context, limit int) ([]*entity.ReportExecution, error)
}

// DashboardRepository interface
type DashboardRepository interface {
	Create(ctx context.Context, dashboard *entity.Dashboard) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error)
	GetByCode(ctx context.Context, code string) (*entity.Dashboard, error)
	Update(ctx context.Context, dashboard *entity.Dashboard) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Dashboard, int64, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Dashboard, error)
	GetDefault(ctx context.Context) (*entity.Dashboard, error)
	GetWithWidgets(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error)
}

// WidgetRepository interface
type WidgetRepository interface {
	Create(ctx context.Context, widget *entity.Widget) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Widget, error)
	Update(ctx context.Context, widget *entity.Widget) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]*entity.Widget, error)
	DeleteByDashboardID(ctx context.Context, dashboardID uuid.UUID) error
}
