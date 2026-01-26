package dashboard

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/erp-cosmetics/reporting-service/internal/domain/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines dashboard use case interface
type UseCase interface {
	Create(ctx context.Context, input *CreateInput) (*entity.Dashboard, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error)
	GetByCode(ctx context.Context, code string) (*entity.Dashboard, error)
	GetWithWidgets(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error)
	GetDefault(ctx context.Context) (*entity.Dashboard, error)
	Update(ctx context.Context, input *UpdateInput) (*entity.Dashboard, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) (*ListOutput, error)
	AddWidget(ctx context.Context, dashboardID uuid.UUID, input *AddWidgetInput) (*entity.Widget, error)
	UpdateWidget(ctx context.Context, widgetID uuid.UUID, input *UpdateWidgetInput) (*entity.Widget, error)
	RemoveWidget(ctx context.Context, widgetID uuid.UUID) error
}

// CreateInput for creating dashboard
type CreateInput struct {
	Code        string     `json:"code" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	LayoutType  string     `json:"layout_type"`
	Visibility  string     `json:"visibility"`
	CreatedBy   *uuid.UUID
}

// UpdateInput for updating dashboard
type UpdateInput struct {
	ID          uuid.UUID
	Name        string `json:"name"`
	Description string `json:"description"`
	LayoutType  string `json:"layout_type"`
	Visibility  string `json:"visibility"`
	IsDefault   *bool  `json:"is_default"`
}

// AddWidgetInput for adding widget
type AddWidgetInput struct {
	WidgetType      string                 `json:"widget_type" binding:"required"`
	Title           string                 `json:"title" binding:"required"`
	Subtitle        string                 `json:"subtitle"`
	DataSource      string                 `json:"data_source"`
	Config          map[string]interface{} `json:"config"`
	PositionX       int                    `json:"position_x"`
	PositionY       int                    `json:"position_y"`
	Width           int                    `json:"width"`
	Height          int                    `json:"height"`
	RefreshInterval int                    `json:"refresh_interval"`
}

// UpdateWidgetInput for updating widget
type UpdateWidgetInput struct {
	Title           string                 `json:"title"`
	Subtitle        string                 `json:"subtitle"`
	DataSource      string                 `json:"data_source"`
	Config          map[string]interface{} `json:"config"`
	PositionX       *int                   `json:"position_x"`
	PositionY       *int                   `json:"position_y"`
	Width           *int                   `json:"width"`
	Height          *int                   `json:"height"`
	RefreshInterval *int                   `json:"refresh_interval"`
	IsVisible       *bool                  `json:"is_visible"`
}

// ListOutput for listing dashboards
type ListOutput struct {
	Dashboards []*entity.Dashboard `json:"dashboards"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
}

type useCase struct {
	dashboardRepo repository.DashboardRepository
	widgetRepo    repository.WidgetRepository
	logger        *zap.Logger
}

// NewUseCase creates new dashboard use case
func NewUseCase(dashboardRepo repository.DashboardRepository, widgetRepo repository.WidgetRepository, logger *zap.Logger) UseCase {
	return &useCase{
		dashboardRepo: dashboardRepo,
		widgetRepo:    widgetRepo,
		logger:        logger,
	}
}

func (uc *useCase) Create(ctx context.Context, input *CreateInput) (*entity.Dashboard, error) {
	dashboard := &entity.Dashboard{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		LayoutType:  input.LayoutType,
		Visibility:  input.Visibility,
		CreatedBy:   input.CreatedBy,
	}

	if dashboard.LayoutType == "" {
		dashboard.LayoutType = entity.LayoutTypeGrid
	}
	if dashboard.Visibility == "" {
		dashboard.Visibility = entity.VisibilityPrivate
	}

	if err := uc.dashboardRepo.Create(ctx, dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (uc *useCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error) {
	return uc.dashboardRepo.GetByID(ctx, id)
}

func (uc *useCase) GetByCode(ctx context.Context, code string) (*entity.Dashboard, error) {
	return uc.dashboardRepo.GetByCode(ctx, code)
}

func (uc *useCase) GetWithWidgets(ctx context.Context, id uuid.UUID) (*entity.Dashboard, error) {
	return uc.dashboardRepo.GetWithWidgets(ctx, id)
}

func (uc *useCase) GetDefault(ctx context.Context) (*entity.Dashboard, error) {
	return uc.dashboardRepo.GetDefault(ctx)
}

func (uc *useCase) Update(ctx context.Context, input *UpdateInput) (*entity.Dashboard, error) {
	dashboard, err := uc.dashboardRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		dashboard.Name = input.Name
	}
	if input.Description != "" {
		dashboard.Description = input.Description
	}
	if input.LayoutType != "" {
		dashboard.LayoutType = input.LayoutType
	}
	if input.Visibility != "" {
		dashboard.Visibility = input.Visibility
	}
	if input.IsDefault != nil {
		dashboard.IsDefault = *input.IsDefault
	}

	if err := uc.dashboardRepo.Update(ctx, dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (uc *useCase) Delete(ctx context.Context, id uuid.UUID) error {
	dashboard, err := uc.dashboardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if dashboard.IsSystem {
		return fmt.Errorf("cannot delete system dashboard")
	}

	// Delete widgets first
	if err := uc.widgetRepo.DeleteByDashboardID(ctx, id); err != nil {
		return err
	}

	return uc.dashboardRepo.Delete(ctx, id)
}

func (uc *useCase) List(ctx context.Context, page, pageSize int) (*ListOutput, error) {
	offset := (page - 1) * pageSize
	dashboards, total, err := uc.dashboardRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Dashboards: dashboards,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (uc *useCase) AddWidget(ctx context.Context, dashboardID uuid.UUID, input *AddWidgetInput) (*entity.Widget, error) {
	widget := &entity.Widget{
		DashboardID:     dashboardID,
		WidgetType:      input.WidgetType,
		Title:           input.Title,
		Subtitle:        input.Subtitle,
		DataSource:      input.DataSource,
		PositionX:       input.PositionX,
		PositionY:       input.PositionY,
		Width:           input.Width,
		Height:          input.Height,
		RefreshInterval: input.RefreshInterval,
		IsVisible:       true,
	}

	if widget.Width == 0 {
		widget.Width = 4
	}
	if widget.Height == 0 {
		widget.Height = 2
	}
	if widget.RefreshInterval == 0 {
		widget.RefreshInterval = 300
	}

	if err := uc.widgetRepo.Create(ctx, widget); err != nil {
		return nil, err
	}

	return widget, nil
}

func (uc *useCase) UpdateWidget(ctx context.Context, widgetID uuid.UUID, input *UpdateWidgetInput) (*entity.Widget, error) {
	widget, err := uc.widgetRepo.GetByID(ctx, widgetID)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		widget.Title = input.Title
	}
	if input.Subtitle != "" {
		widget.Subtitle = input.Subtitle
	}
	if input.DataSource != "" {
		widget.DataSource = input.DataSource
	}
	if input.PositionX != nil {
		widget.PositionX = *input.PositionX
	}
	if input.PositionY != nil {
		widget.PositionY = *input.PositionY
	}
	if input.Width != nil {
		widget.Width = *input.Width
	}
	if input.Height != nil {
		widget.Height = *input.Height
	}
	if input.RefreshInterval != nil {
		widget.RefreshInterval = *input.RefreshInterval
	}
	if input.IsVisible != nil {
		widget.IsVisible = *input.IsVisible
	}

	if err := uc.widgetRepo.Update(ctx, widget); err != nil {
		return nil, err
	}

	return widget, nil
}

func (uc *useCase) RemoveWidget(ctx context.Context, widgetID uuid.UUID) error {
	return uc.widgetRepo.Delete(ctx, widgetID)
}
