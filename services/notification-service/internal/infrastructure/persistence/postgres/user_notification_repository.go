package postgres

import (
	"context"
	"time"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userNotificationRepository struct {
	db *gorm.DB
}

// NewUserNotificationRepository creates a new user notification repository
func NewUserNotificationRepository(db *gorm.DB) repository.UserNotificationRepository {
	return &userNotificationRepository{db: db}
}

func (r *userNotificationRepository) Create(ctx context.Context, notification *entity.UserNotification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *userNotificationRepository) CreateBatch(ctx context.Context, notifications []*entity.UserNotification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(notifications).Error
}

func (r *userNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UserNotification, error) {
	var notification entity.UserNotification
	err := r.db.WithContext(ctx).First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *userNotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.UserNotification, int64, error) {
	var notifications []*entity.UserNotification
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.UserNotification{}).
		Where("user_id = ? AND is_dismissed = ?", userID, false)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, total, err
}

func (r *userNotificationRepository) GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserNotification, error) {
	var notifications []*entity.UserNotification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ?", userID, false, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *userNotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserNotification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ?", userID, false, false).
		Count(&count).Error
	return count, err
}

func (r *userNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.UserNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *userNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *userNotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.UserNotification{}, "id = ?", id).Error
}

func (r *userNotificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.UserNotification{}, "user_id = ?", userID).Error
}
