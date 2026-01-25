package event

import (
	"encoding/json"
	"time"

	"github.com/erp-cosmetics/shared/pkg/nats"
	"go.uber.org/zap"
)

// Event subjects
const (
	SubjectNotificationSent   = "notification.sent"
	SubjectNotificationFailed = "notification.failed"
	SubjectEmailSent          = "notification.email.sent"
	SubjectEmailFailed        = "notification.email.failed"
)

// Publisher handles event publishing
type Publisher struct {
	client *nats.Client
	logger *zap.Logger
}

// NewPublisher creates a new event publisher
func NewPublisher(client *nats.Client, logger *zap.Logger) *Publisher {
	return &Publisher{
		client: client,
		logger: logger,
	}
}

// Event represents a generic event
type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Publish publishes an event
func (p *Publisher) Publish(subject string, data map[string]interface{}) error {
	event := Event{
		Type:      subject,
		Timestamp: time.Now(),
		Data:      data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.Error(err))
		return err
	}

	if err := p.client.Publish(subject, payload); err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return err
	}

	p.logger.Debug("Event published",
		zap.String("subject", subject),
		zap.Any("data", data),
	)

	return nil
}

// PublishNotificationSent publishes notification sent event
func (p *Publisher) PublishNotificationSent(notificationID, notificationType, recipientEmail string) error {
	return p.Publish(SubjectNotificationSent, map[string]interface{}{
		"notification_id":   notificationID,
		"notification_type": notificationType,
		"recipient_email":   recipientEmail,
	})
}

// PublishNotificationFailed publishes notification failed event
func (p *Publisher) PublishNotificationFailed(notificationID, notificationType, errorMessage string) error {
	return p.Publish(SubjectNotificationFailed, map[string]interface{}{
		"notification_id":   notificationID,
		"notification_type": notificationType,
		"error_message":     errorMessage,
	})
}

// PublishEmailSent publishes email sent event
func (p *Publisher) PublishEmailSent(emailLogID, toEmail, subject string) error {
	return p.Publish(SubjectEmailSent, map[string]interface{}{
		"email_log_id": emailLogID,
		"to_email":     toEmail,
		"subject":      subject,
	})
}

// PublishEmailFailed publishes email failed event
func (p *Publisher) PublishEmailFailed(emailLogID, toEmail, errorMessage string) error {
	return p.Publish(SubjectEmailFailed, map[string]interface{}{
		"email_log_id":  emailLogID,
		"to_email":      toEmail,
		"error_message": errorMessage,
	})
}
