package repository

import (
	"context"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/google/uuid"
)

// NotificationRepository defines notification repository interface
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *entity.Notification) error
	
	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	
	// Update updates a notification
	Update(ctx context.Context, notification *entity.Notification) error
	
	// GetPending retrieves pending notifications for processing
	GetPending(ctx context.Context, limit int) ([]*entity.Notification, error)
	
	// GetRetryable retrieves notifications ready for retry
	GetRetryable(ctx context.Context, limit int) ([]*entity.Notification, error)
	
	// GetByUserID retrieves notifications for a user
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, int64, error)
	
	// CountByStatus counts notifications by status
	CountByStatus(ctx context.Context, status string) (int64, error)
}

// TemplateRepository defines template repository interface
type TemplateRepository interface {
	// Create creates a new template
	Create(ctx context.Context, template *entity.NotificationTemplate) error
	
	// GetByID retrieves a template by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error)
	
	// GetByCode retrieves a template by code
	GetByCode(ctx context.Context, code string) (*entity.NotificationTemplate, error)
	
	// Update updates a template
	Update(ctx context.Context, template *entity.NotificationTemplate) error
	
	// Delete deletes a template
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves all templates with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.NotificationTemplate, int64, error)
	
	// ListActive retrieves only active templates
	ListActive(ctx context.Context) ([]*entity.NotificationTemplate, error)
	
	// ListByType retrieves templates by notification type
	ListByType(ctx context.Context, notificationType string) ([]*entity.NotificationTemplate, error)
}

// UserNotificationRepository defines user notification repository interface
type UserNotificationRepository interface {
	// Create creates a new user notification
	Create(ctx context.Context, notification *entity.UserNotification) error
	
	// CreateBatch creates multiple notifications
	CreateBatch(ctx context.Context, notifications []*entity.UserNotification) error
	
	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UserNotification, error)
	
	// GetByUserID retrieves notifications for a user
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.UserNotification, int64, error)
	
	// GetUnreadByUserID retrieves unread notifications for a user
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserNotification, error)
	
	// CountUnread counts unread notifications for a user
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
	
	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	
	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	
	// Delete deletes a notification
	Delete(ctx context.Context, id uuid.UUID) error
	
	// DeleteByUserID deletes all notifications for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// AlertRuleRepository defines alert rule repository interface
type AlertRuleRepository interface {
	// Create creates a new alert rule
	Create(ctx context.Context, rule *entity.AlertRule) error
	
	// GetByID retrieves a rule by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AlertRule, error)
	
	// GetByCode retrieves a rule by code
	GetByCode(ctx context.Context, code string) (*entity.AlertRule, error)
	
	// Update updates a rule
	Update(ctx context.Context, rule *entity.AlertRule) error
	
	// Delete deletes a rule
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves all rules with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.AlertRule, int64, error)
	
	// ListActive retrieves only active rules
	ListActive(ctx context.Context) ([]*entity.AlertRule, error)
	
	// ListDueForCheck retrieves rules that need to be checked
	ListDueForCheck(ctx context.Context) ([]*entity.AlertRule, error)
	
	// ListByType retrieves rules by type
	ListByType(ctx context.Context, ruleType string) ([]*entity.AlertRule, error)
}

// EmailLogRepository defines email log repository interface
type EmailLogRepository interface {
	// Create creates a new email log
	Create(ctx context.Context, log *entity.EmailLog) error
	
	// GetByID retrieves an email log by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.EmailLog, error)
	
	// Update updates an email log
	Update(ctx context.Context, log *entity.EmailLog) error
	
	// GetByNotificationID retrieves logs for a notification
	GetByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*entity.EmailLog, error)
	
	// GetByEmail retrieves logs for an email address
	GetByEmail(ctx context.Context, email string, limit, offset int) ([]*entity.EmailLog, int64, error)
	
	// CountByStatus counts emails by status
	CountByStatus(ctx context.Context, status string) (int64, error)
}
