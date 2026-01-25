package notification

import (
	"context"
	"encoding/json"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/erp-cosmetics/notification-service/internal/infrastructure/email"
	"github.com/erp-cosmetics/notification-service/internal/infrastructure/event"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SendNotificationInput represents input for sending notification
type SendNotificationInput struct {
	NotificationType  string                 `json:"notification_type" binding:"required,oneof=EMAIL IN_APP BOTH"`
	TemplateCode      string                 `json:"template_code"`
	RecipientUserID   *uuid.UUID             `json:"recipient_user_id"`
	RecipientEmail    string                 `json:"recipient_email"`
	Subject           string                 `json:"subject"`
	Body              string                 `json:"body"`
	Variables         map[string]interface{} `json:"variables"`
	Priority          string                 `json:"priority"`
	LinkURL           string                 `json:"link_url"`
	EntityType        string                 `json:"entity_type"`
	EntityID          *uuid.UUID             `json:"entity_id"`
}

// SendNotificationOutput represents output of sending notification
type SendNotificationOutput struct {
	NotificationID     *uuid.UUID `json:"notification_id,omitempty"`
	UserNotificationID *uuid.UUID `json:"user_notification_id,omitempty"`
	Status             string     `json:"status"`
	Message            string     `json:"message"`
}

// UseCase defines notification use case interface
type UseCase interface {
	Send(ctx context.Context, input *SendNotificationInput) (*SendNotificationOutput, error)
	ProcessPending(ctx context.Context, limit int) error
	ProcessRetryable(ctx context.Context, limit int) error
}

type useCase struct {
	notificationRepo     repository.NotificationRepository
	templateRepo         repository.TemplateRepository
	userNotificationRepo repository.UserNotificationRepository
	emailSender          email.Sender
	eventPublisher       *event.Publisher
	logger               *zap.Logger
}

// NewUseCase creates a new notification use case
func NewUseCase(
	notificationRepo repository.NotificationRepository,
	templateRepo repository.TemplateRepository,
	userNotificationRepo repository.UserNotificationRepository,
	emailSender email.Sender,
	eventPublisher *event.Publisher,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		notificationRepo:     notificationRepo,
		templateRepo:         templateRepo,
		userNotificationRepo: userNotificationRepo,
		emailSender:          emailSender,
		eventPublisher:       eventPublisher,
		logger:               logger,
	}
}

func (uc *useCase) Send(ctx context.Context, input *SendNotificationInput) (*SendNotificationOutput, error) {
	output := &SendNotificationOutput{}

	// Get template if specified
	var template *entity.NotificationTemplate
	if input.TemplateCode != "" {
		var err error
		template, err = uc.templateRepo.GetByCode(ctx, input.TemplateCode)
		if err != nil {
			uc.logger.Error("Template not found", zap.String("code", input.TemplateCode), zap.Error(err))
			return nil, err
		}
	}

	// Render content from template if available
	subject := input.Subject
	body := input.Body

	if template != nil && input.Variables != nil {
		if template.SubjectTemplate != "" && subject == "" {
			rendered, err := template.RenderSubject(input.Variables)
			if err == nil {
				subject = rendered
			}
		}
		if template.BodyTemplate != "" && body == "" {
			rendered, err := template.RenderBody(input.Variables)
			if err == nil {
				body = rendered
			}
		}
	}

	// Handle different notification types
	switch input.NotificationType {
	case entity.NotificationTypeEmail:
		notifID, err := uc.sendEmail(ctx, input, template, subject, body)
		if err != nil {
			return nil, err
		}
		output.NotificationID = notifID
		output.Status = "QUEUED"
		output.Message = "Email notification queued"

	case entity.NotificationTypeInApp:
		userNotifID, err := uc.createInAppNotification(ctx, input, subject, body)
		if err != nil {
			return nil, err
		}
		output.UserNotificationID = userNotifID
		output.Status = "CREATED"
		output.Message = "In-app notification created"

	case entity.NotificationTypeBoth:
		// Send email
		notifID, err := uc.sendEmail(ctx, input, template, subject, body)
		if err != nil {
			uc.logger.Error("Failed to queue email", zap.Error(err))
		} else {
			output.NotificationID = notifID
		}

		// Create in-app notification
		userNotifID, err := uc.createInAppNotification(ctx, input, subject, body)
		if err != nil {
			uc.logger.Error("Failed to create in-app notification", zap.Error(err))
		} else {
			output.UserNotificationID = userNotifID
		}

		output.Status = "CREATED"
		output.Message = "Notifications created"
	}

	return output, nil
}

func (uc *useCase) sendEmail(
	ctx context.Context,
	input *SendNotificationInput,
	template *entity.NotificationTemplate,
	subject, body string,
) (*uuid.UUID, error) {
	if input.RecipientEmail == "" {
		return nil, nil
	}

	// Convert variables to JSON
	var metadata []byte
	if input.Variables != nil {
		metadata, _ = json.Marshal(input.Variables)
	}

	notification := &entity.Notification{
		NotificationType: entity.NotificationTypeEmail,
		RecipientUserID:  input.RecipientUserID,
		RecipientEmail:   input.RecipientEmail,
		Subject:          subject,
		Body:             body,
		Priority:         input.Priority,
	}

	if template != nil {
		notification.TemplateID = &template.ID
	}

	if metadata != nil {
		notification.Metadata = metadata
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	uc.logger.Info("Email notification queued",
		zap.String("id", notification.ID.String()),
		zap.String("to", input.RecipientEmail),
	)

	return &notification.ID, nil
}

func (uc *useCase) createInAppNotification(
	ctx context.Context,
	input *SendNotificationInput,
	title, message string,
) (*uuid.UUID, error) {
	if input.RecipientUserID == nil {
		return nil, nil
	}

	notification := &entity.UserNotification{
		UserID:           *input.RecipientUserID,
		Title:            title,
		Message:          message,
		NotificationType: entity.UserNotifTypeInfo,
		Category:         entity.CategorySystem,
		LinkURL:          input.LinkURL,
		EntityType:       input.EntityType,
		EntityID:         input.EntityID,
	}

	if err := uc.userNotificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	uc.logger.Info("In-app notification created",
		zap.String("id", notification.ID.String()),
		zap.String("user_id", input.RecipientUserID.String()),
	)

	return &notification.ID, nil
}

func (uc *useCase) ProcessPending(ctx context.Context, limit int) error {
	notifications, err := uc.notificationRepo.GetPending(ctx, limit)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := uc.processNotification(ctx, notification); err != nil {
			uc.logger.Error("Failed to process notification",
				zap.String("id", notification.ID.String()),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (uc *useCase) ProcessRetryable(ctx context.Context, limit int) error {
	notifications, err := uc.notificationRepo.GetRetryable(ctx, limit)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := uc.processNotification(ctx, notification); err != nil {
			uc.logger.Error("Failed to retry notification",
				zap.String("id", notification.ID.String()),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (uc *useCase) processNotification(ctx context.Context, notification *entity.Notification) error {
	switch notification.NotificationType {
	case entity.NotificationTypeEmail:
		return uc.processEmailNotification(ctx, notification)
	default:
		return nil
	}
}

func (uc *useCase) processEmailNotification(ctx context.Context, notification *entity.Notification) error {
	// Send email
	_, err := uc.emailSender.Send(
		notification.RecipientEmail,
		notification.Subject,
		notification.Body,
		"", // Plain text version
	)

	if err != nil {
		// Schedule retry
		if notification.ScheduleRetry(err.Error()) {
			uc.notificationRepo.Update(ctx, notification)
			uc.logger.Warn("Email send failed, scheduled retry",
				zap.String("id", notification.ID.String()),
				zap.Int("retry_count", notification.RetryCount),
			)
		} else {
			uc.notificationRepo.Update(ctx, notification)
			uc.logger.Error("Email send failed permanently",
				zap.String("id", notification.ID.String()),
				zap.Error(err),
			)
			// Publish failed event
			uc.eventPublisher.PublishNotificationFailed(
				notification.ID.String(),
				notification.NotificationType,
				err.Error(),
			)
		}
		return err
	}

	// Mark as sent
	notification.MarkAsSent()
	if err := uc.notificationRepo.Update(ctx, notification); err != nil {
		return err
	}

	// Publish sent event
	uc.eventPublisher.PublishNotificationSent(
		notification.ID.String(),
		notification.NotificationType,
		notification.RecipientEmail,
	)

	uc.logger.Info("Email sent successfully",
		zap.String("id", notification.ID.String()),
		zap.String("to", notification.RecipientEmail),
	)

	return nil
}
