package postgres

import (
	"context"
	"time"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.db.WithContext(ctx).Preload("Template").First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

func (r *notificationRepository) GetPending(ctx context.Context, limit int) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.NotificationStatusPending).
		Order("CASE WHEN priority = 'HIGH' THEN 1 WHEN priority = 'NORMAL' THEN 2 ELSE 3 END").
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) GetRetryable(ctx context.Context, limit int) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at <= ?", entity.NotificationStatusRetrying, time.Now()).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, int64, error) {
	var notifications []*entity.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("recipient_user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, total, err
}

func (r *notificationRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
