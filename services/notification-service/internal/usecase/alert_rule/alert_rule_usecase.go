package alert_rule

import (
	"context"
	"encoding/json"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateInput represents input for creating alert rule
type CreateInput struct {
	RuleCode               string                 `json:"rule_code" binding:"required"`
	Name                   string                 `json:"name" binding:"required"`
	Description            string                 `json:"description"`
	RuleType               string                 `json:"rule_type" binding:"required"`
	Conditions             map[string]interface{} `json:"conditions"`
	NotificationTemplateID *uuid.UUID             `json:"notification_template_id"`
	NotificationType       string                 `json:"notification_type"` // EMAIL, IN_APP, BOTH
	Recipients             []map[string]string    `json:"recipients"`
	CheckInterval          string                 `json:"check_interval"` // 1h, 30m, daily
	CreatedBy              *uuid.UUID
}

// UpdateInput represents input for updating alert rule
type UpdateInput struct {
	ID                     uuid.UUID
	Name                   string                 `json:"name"`
	Description            string                 `json:"description"`
	Conditions             map[string]interface{} `json:"conditions"`
	NotificationTemplateID *uuid.UUID             `json:"notification_template_id"`
	NotificationType       string                 `json:"notification_type"`
	Recipients             []map[string]string    `json:"recipients"`
	CheckInterval          string                 `json:"check_interval"`
	IsActive               *bool                  `json:"is_active"`
}

// ListOutput represents paginated alert rule list
type ListOutput struct {
	Rules    []*entity.AlertRule `json:"rules"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// UseCase defines alert rule use case interface
type UseCase interface {
	Create(ctx context.Context, input *CreateInput) (*entity.AlertRule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AlertRule, error)
	Update(ctx context.Context, input *UpdateInput) (*entity.AlertRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) (*ListOutput, error)
	ListActive(ctx context.Context) ([]*entity.AlertRule, error)
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type useCase struct {
	alertRuleRepo repository.AlertRuleRepository
	logger        *zap.Logger
}

// NewUseCase creates a new alert rule use case
func NewUseCase(alertRuleRepo repository.AlertRuleRepository, logger *zap.Logger) UseCase {
	return &useCase{
		alertRuleRepo: alertRuleRepo,
		logger:        logger,
	}
}

func (uc *useCase) Create(ctx context.Context, input *CreateInput) (*entity.AlertRule, error) {
	rule := &entity.AlertRule{
		RuleCode:               input.RuleCode,
		Name:                   input.Name,
		Description:            input.Description,
		RuleType:               input.RuleType,
		NotificationTemplateID: input.NotificationTemplateID,
		NotificationType:       input.NotificationType,
		CheckInterval:          input.CheckInterval,
		IsActive:               true,
		CreatedBy:              input.CreatedBy,
	}

	// Set defaults
	if rule.NotificationType == "" {
		rule.NotificationType = entity.NotificationTypeInApp
	}
	if rule.CheckInterval == "" {
		rule.CheckInterval = "1h"
	}

	// Convert conditions to JSON
	if input.Conditions != nil {
		conditionsJSON, err := json.Marshal(input.Conditions)
		if err != nil {
			return nil, err
		}
		rule.Conditions = conditionsJSON
	}

	// Convert recipients to JSON
	if input.Recipients != nil {
		recipientsJSON, err := json.Marshal(input.Recipients)
		if err != nil {
			return nil, err
		}
		rule.Recipients = recipientsJSON
	}

	if err := uc.alertRuleRepo.Create(ctx, rule); err != nil {
		uc.logger.Error("Failed to create alert rule", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Alert rule created",
		zap.String("id", rule.ID.String()),
		zap.String("code", rule.RuleCode),
	)

	return rule, nil
}

func (uc *useCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.AlertRule, error) {
	return uc.alertRuleRepo.GetByID(ctx, id)
}

func (uc *useCase) Update(ctx context.Context, input *UpdateInput) (*entity.AlertRule, error) {
	rule, err := uc.alertRuleRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		rule.Name = input.Name
	}
	if input.Description != "" {
		rule.Description = input.Description
	}
	if input.NotificationTemplateID != nil {
		rule.NotificationTemplateID = input.NotificationTemplateID
	}
	if input.NotificationType != "" {
		rule.NotificationType = input.NotificationType
	}
	if input.CheckInterval != "" {
		rule.CheckInterval = input.CheckInterval
	}
	if input.IsActive != nil {
		rule.IsActive = *input.IsActive
	}

	// Update conditions if provided
	if input.Conditions != nil {
		conditionsJSON, err := json.Marshal(input.Conditions)
		if err != nil {
			return nil, err
		}
		rule.Conditions = conditionsJSON
	}

	// Update recipients if provided
	if input.Recipients != nil {
		recipientsJSON, err := json.Marshal(input.Recipients)
		if err != nil {
			return nil, err
		}
		rule.Recipients = recipientsJSON
	}

	if err := uc.alertRuleRepo.Update(ctx, rule); err != nil {
		uc.logger.Error("Failed to update alert rule", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Alert rule updated", zap.String("id", rule.ID.String()))
	return rule, nil
}

func (uc *useCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.alertRuleRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete alert rule", zap.Error(err))
		return err
	}

	uc.logger.Info("Alert rule deleted", zap.String("id", id.String()))
	return nil
}

func (uc *useCase) List(ctx context.Context, page, pageSize int) (*ListOutput, error) {
	offset := (page - 1) * pageSize
	rules, total, err := uc.alertRuleRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Rules:    rules,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (uc *useCase) ListActive(ctx context.Context) ([]*entity.AlertRule, error) {
	return uc.alertRuleRepo.ListActive(ctx)
}

func (uc *useCase) Activate(ctx context.Context, id uuid.UUID) error {
	rule, err := uc.alertRuleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	rule.Activate()
	return uc.alertRuleRepo.Update(ctx, rule)
}

func (uc *useCase) Deactivate(ctx context.Context, id uuid.UUID) error {
	rule, err := uc.alertRuleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	rule.Deactivate()
	return uc.alertRuleRepo.Update(ctx, rule)
}
