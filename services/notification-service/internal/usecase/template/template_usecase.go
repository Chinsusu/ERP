package template

import (
	"context"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateTemplateInput represents input for creating template
type CreateTemplateInput struct {
	TemplateCode     string   `json:"template_code" binding:"required"`
	Name             string   `json:"name" binding:"required"`
	NotificationType string   `json:"notification_type" binding:"required,oneof=EMAIL IN_APP SMS PUSH BOTH"`
	SubjectTemplate  string   `json:"subject_template"`
	BodyTemplate     string   `json:"body_template" binding:"required"`
	Variables        []string `json:"variables"`
	CreatedBy        *uuid.UUID
}

// UpdateTemplateInput represents input for updating template
type UpdateTemplateInput struct {
	ID               uuid.UUID
	Name             string   `json:"name"`
	NotificationType string   `json:"notification_type"`
	SubjectTemplate  string   `json:"subject_template"`
	BodyTemplate     string   `json:"body_template"`
	Variables        []string `json:"variables"`
	IsActive         *bool    `json:"is_active"`
}

// TemplateListOutput represents paginated template list
type TemplateListOutput struct {
	Templates []*entity.NotificationTemplate `json:"templates"`
	Total     int64                          `json:"total"`
	Page      int                            `json:"page"`
	PageSize  int                            `json:"page_size"`
}

// UseCase defines template use case interface
type UseCase interface {
	Create(ctx context.Context, input *CreateTemplateInput) (*entity.NotificationTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error)
	GetByCode(ctx context.Context, code string) (*entity.NotificationTemplate, error)
	Update(ctx context.Context, input *UpdateTemplateInput) (*entity.NotificationTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) (*TemplateListOutput, error)
	ListActive(ctx context.Context) ([]*entity.NotificationTemplate, error)
}

type useCase struct {
	templateRepo repository.TemplateRepository
	logger       *zap.Logger
}

// NewUseCase creates a new template use case
func NewUseCase(templateRepo repository.TemplateRepository, logger *zap.Logger) UseCase {
	return &useCase{
		templateRepo: templateRepo,
		logger:       logger,
	}
}

func (uc *useCase) Create(ctx context.Context, input *CreateTemplateInput) (*entity.NotificationTemplate, error) {
	template := &entity.NotificationTemplate{
		TemplateCode:     input.TemplateCode,
		Name:             input.Name,
		NotificationType: input.NotificationType,
		SubjectTemplate:  input.SubjectTemplate,
		BodyTemplate:     input.BodyTemplate,
		IsActive:         true,
		CreatedBy:        input.CreatedBy,
	}

	if input.Variables != nil {
		// Store variables as JSON
		template.Variables = []byte(`["` + joinVariables(input.Variables) + `"]`)
	}

	if err := uc.templateRepo.Create(ctx, template); err != nil {
		uc.logger.Error("Failed to create template", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Template created",
		zap.String("id", template.ID.String()),
		zap.String("code", template.TemplateCode),
	)

	return template, nil
}

func (uc *useCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error) {
	return uc.templateRepo.GetByID(ctx, id)
}

func (uc *useCase) GetByCode(ctx context.Context, code string) (*entity.NotificationTemplate, error) {
	return uc.templateRepo.GetByCode(ctx, code)
}

func (uc *useCase) Update(ctx context.Context, input *UpdateTemplateInput) (*entity.NotificationTemplate, error) {
	template, err := uc.templateRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		template.Name = input.Name
	}
	if input.NotificationType != "" {
		template.NotificationType = input.NotificationType
	}
	if input.SubjectTemplate != "" {
		template.SubjectTemplate = input.SubjectTemplate
	}
	if input.BodyTemplate != "" {
		template.BodyTemplate = input.BodyTemplate
	}
	if input.IsActive != nil {
		template.IsActive = *input.IsActive
	}

	if err := uc.templateRepo.Update(ctx, template); err != nil {
		uc.logger.Error("Failed to update template", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Template updated",
		zap.String("id", template.ID.String()),
	)

	return template, nil
}

func (uc *useCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.templateRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete template", zap.Error(err))
		return err
	}

	uc.logger.Info("Template deleted", zap.String("id", id.String()))
	return nil
}

func (uc *useCase) List(ctx context.Context, page, pageSize int) (*TemplateListOutput, error) {
	offset := (page - 1) * pageSize
	templates, total, err := uc.templateRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &TemplateListOutput{
		Templates: templates,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

func (uc *useCase) ListActive(ctx context.Context) ([]*entity.NotificationTemplate, error) {
	return uc.templateRepo.ListActive(ctx)
}

// Helper function to join variables
func joinVariables(vars []string) string {
	result := ""
	for i, v := range vars {
		if i > 0 {
			result += `", "`
		}
		result += v
	}
	return result
}
