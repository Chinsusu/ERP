package user_notification

import (
	"context"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateInput represents input for creating user notification
type CreateInput struct {
	UserID           uuid.UUID `json:"user_id" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	Message          string    `json:"message" binding:"required"`
	NotificationType string    `json:"notification_type"` // INFO, SUCCESS, WARNING, ERROR
	Category         string    `json:"category"`          // APPROVAL, ALERT, MESSAGE, SYSTEM
	LinkURL          string    `json:"link_url"`
	EntityType       string    `json:"entity_type"`
	EntityID         *uuid.UUID `json:"entity_id"`
}

// ListOutput represents paginated notification list
type ListOutput struct {
	Notifications []*entity.UserNotification `json:"notifications"`
	Total         int64                      `json:"total"`
	UnreadCount   int64                      `json:"unread_count"`
	Page          int                        `json:"page"`
	PageSize      int                        `json:"page_size"`
}

// UseCase defines user notification use case interface
type UseCase interface {
	Create(ctx context.Context, input *CreateInput) (*entity.UserNotification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) (*ListOutput, error)
	GetUnread(ctx context.Context, userID uuid.UUID) ([]*entity.UserNotification, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type useCase struct {
	userNotificationRepo repository.UserNotificationRepository
	logger               *zap.Logger
}

// NewUseCase creates a new user notification use case
func NewUseCase(userNotificationRepo repository.UserNotificationRepository, logger *zap.Logger) UseCase {
	return &useCase{
		userNotificationRepo: userNotificationRepo,
		logger:               logger,
	}
}

func (uc *useCase) Create(ctx context.Context, input *CreateInput) (*entity.UserNotification, error) {
	notification := &entity.UserNotification{
		UserID:           input.UserID,
		Title:            input.Title,
		Message:          input.Message,
		NotificationType: input.NotificationType,
		Category:         input.Category,
		LinkURL:          input.LinkURL,
		EntityType:       input.EntityType,
		EntityID:         input.EntityID,
	}

	// Set defaults
	if notification.NotificationType == "" {
		notification.NotificationType = entity.UserNotifTypeInfo
	}
	if notification.Category == "" {
		notification.Category = entity.CategorySystem
	}

	if err := uc.userNotificationRepo.Create(ctx, notification); err != nil {
		uc.logger.Error("Failed to create user notification", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("User notification created",
		zap.String("id", notification.ID.String()),
		zap.String("user_id", input.UserID.String()),
	)

	return notification, nil
}

func (uc *useCase) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) (*ListOutput, error) {
	offset := (page - 1) * pageSize
	notifications, total, err := uc.userNotificationRepo.GetByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	unreadCount, err := uc.userNotificationRepo.CountUnread(ctx, userID)
	if err != nil {
		unreadCount = 0
	}

	return &ListOutput{
		Notifications: notifications,
		Total:         total,
		UnreadCount:   unreadCount,
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func (uc *useCase) GetUnread(ctx context.Context, userID uuid.UUID) ([]*entity.UserNotification, error) {
	return uc.userNotificationRepo.GetUnreadByUserID(ctx, userID)
}

func (uc *useCase) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	return uc.userNotificationRepo.CountUnread(ctx, userID)
}

func (uc *useCase) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	if err := uc.userNotificationRepo.MarkAsRead(ctx, id); err != nil {
		uc.logger.Error("Failed to mark notification as read", zap.Error(err))
		return err
	}

	uc.logger.Debug("Notification marked as read", zap.String("id", id.String()))
	return nil
}

func (uc *useCase) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := uc.userNotificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		uc.logger.Error("Failed to mark all notifications as read", zap.Error(err))
		return err
	}

	uc.logger.Info("All notifications marked as read", zap.String("user_id", userID.String()))
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.userNotificationRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete notification", zap.Error(err))
		return err
	}

	uc.logger.Debug("Notification deleted", zap.String("id", id.String()))
	return nil
}
