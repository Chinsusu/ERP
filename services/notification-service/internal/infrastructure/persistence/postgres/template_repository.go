package postgres

import (
	"context"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type templateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *gorm.DB) repository.TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, template *entity.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error) {
	var template entity.NotificationTemplate
	err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) GetByCode(ctx context.Context, code string) (*entity.NotificationTemplate, error) {
	var template entity.NotificationTemplate
	err := r.db.WithContext(ctx).First(&template, "template_code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) Update(ctx context.Context, template *entity.NotificationTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *templateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.NotificationTemplate{}, "id = ?", id).Error
}

func (r *templateRepository) List(ctx context.Context, limit, offset int) ([]*entity.NotificationTemplate, int64, error) {
	var templates []*entity.NotificationTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.NotificationTemplate{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&templates).Error
	return templates, total, err
}

func (r *templateRepository) ListActive(ctx context.Context) ([]*entity.NotificationTemplate, error) {
	var templates []*entity.NotificationTemplate
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&templates).Error
	return templates, err
}

func (r *templateRepository) ListByType(ctx context.Context, notificationType string) ([]*entity.NotificationTemplate, error) {
	var templates []*entity.NotificationTemplate
	err := r.db.WithContext(ctx).
		Where("notification_type = ? AND is_active = ?", notificationType, true).
		Order("name ASC").
		Find(&templates).Error
	return templates, err
}
