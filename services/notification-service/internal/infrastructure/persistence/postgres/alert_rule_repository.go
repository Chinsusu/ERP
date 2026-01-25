package postgres

import (
	"context"
	"time"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type alertRuleRepository struct {
	db *gorm.DB
}

// NewAlertRuleRepository creates a new alert rule repository
func NewAlertRuleRepository(db *gorm.DB) repository.AlertRuleRepository {
	return &alertRuleRepository{db: db}
}

func (r *alertRuleRepository) Create(ctx context.Context, rule *entity.AlertRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *alertRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AlertRule, error) {
	var rule entity.AlertRule
	err := r.db.WithContext(ctx).Preload("Template").First(&rule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *alertRuleRepository) GetByCode(ctx context.Context, code string) (*entity.AlertRule, error) {
	var rule entity.AlertRule
	err := r.db.WithContext(ctx).Preload("Template").First(&rule, "rule_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *alertRuleRepository) Update(ctx context.Context, rule *entity.AlertRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

func (r *alertRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AlertRule{}, "id = ?", id).Error
}

func (r *alertRuleRepository) List(ctx context.Context, limit, offset int) ([]*entity.AlertRule, int64, error) {
	var rules []*entity.AlertRule
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AlertRule{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Template").Order("name ASC").Limit(limit).Offset(offset).Find(&rules).Error
	return rules, total, err
}

func (r *alertRuleRepository) ListActive(ctx context.Context) ([]*entity.AlertRule, error) {
	var rules []*entity.AlertRule
	err := r.db.WithContext(ctx).
		Preload("Template").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&rules).Error
	return rules, err
}

func (r *alertRuleRepository) ListDueForCheck(ctx context.Context) ([]*entity.AlertRule, error) {
	var rules []*entity.AlertRule
	err := r.db.WithContext(ctx).
		Preload("Template").
		Where("is_active = ? AND (next_check_at IS NULL OR next_check_at <= ?)", true, time.Now()).
		Find(&rules).Error
	return rules, err
}

func (r *alertRuleRepository) ListByType(ctx context.Context, ruleType string) ([]*entity.AlertRule, error) {
	var rules []*entity.AlertRule
	err := r.db.WithContext(ctx).
		Preload("Template").
		Where("rule_type = ? AND is_active = ?", ruleType, true).
		Find(&rules).Error
	return rules, err
}
