package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// EmailLogStatus constants
const (
	EmailStatusPending   = "PENDING"
	EmailStatusSent      = "SENT"
	EmailStatusDelivered = "DELIVERED"
	EmailStatusBounced   = "BOUNCED"
	EmailStatusFailed    = "FAILED"
)

// EmailLog represents an email delivery log
type EmailLog struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NotificationID *uuid.UUID     `gorm:"type:uuid" json:"notification_id,omitempty"`
	
	// Email details
	FromEmail      string         `gorm:"type:varchar(255);not null" json:"from_email"`
	FromName       string         `gorm:"type:varchar(200)" json:"from_name,omitempty"`
	ToEmail        string         `gorm:"type:varchar(255);not null" json:"to_email"`
	CCEmails       pq.StringArray `gorm:"type:text[]" json:"cc_emails,omitempty"`
	BCCEmails      pq.StringArray `gorm:"type:text[]" json:"bcc_emails,omitempty"`
	
	// Content
	Subject        string         `gorm:"type:text;not null" json:"subject"`
	BodyHTML       string         `gorm:"type:text" json:"body_html,omitempty"`
	BodyText       string         `gorm:"type:text" json:"body_text,omitempty"`
	
	// Delivery status
	Status         string         `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	SMTPResponse   string         `gorm:"type:text" json:"smtp_response,omitempty"`
	MessageID      string         `gorm:"type:varchar(255)" json:"message_id,omitempty"` // SMTP message ID
	
	// Tracking
	OpenedAt       *time.Time     `gorm:"type:timestamp" json:"opened_at,omitempty"`
	ClickedAt      *time.Time     `gorm:"type:timestamp" json:"clicked_at,omitempty"`
	BouncedAt      *time.Time     `gorm:"type:timestamp" json:"bounced_at,omitempty"`
	
	// Error handling
	ErrorCode      string         `gorm:"type:varchar(50)" json:"error_code,omitempty"`
	ErrorMessage   string         `gorm:"type:text" json:"error_message,omitempty"`
	
	SentAt         *time.Time     `gorm:"type:timestamp" json:"sent_at,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	
	// Relationships
	Notification   *Notification  `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
}

// TableName specifies the table name
func (EmailLog) TableName() string {
	return "email_logs"
}

// BeforeCreate sets defaults before creating
func (e *EmailLog) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.Status == "" {
		e.Status = EmailStatusPending
	}
	return nil
}

// MarkAsSent marks the email as sent
func (e *EmailLog) MarkAsSent(messageID, smtpResponse string) {
	now := time.Now()
	e.Status = EmailStatusSent
	e.SentAt = &now
	e.MessageID = messageID
	e.SMTPResponse = smtpResponse
}

// MarkAsDelivered marks the email as delivered
func (e *EmailLog) MarkAsDelivered() {
	e.Status = EmailStatusDelivered
}

// MarkAsBounced marks the email as bounced
func (e *EmailLog) MarkAsBounced(errorCode, errorMessage string) {
	now := time.Now()
	e.Status = EmailStatusBounced
	e.BouncedAt = &now
	e.ErrorCode = errorCode
	e.ErrorMessage = errorMessage
}

// MarkAsFailed marks the email as failed
func (e *EmailLog) MarkAsFailed(errorCode, errorMessage string) {
	e.Status = EmailStatusFailed
	e.ErrorCode = errorCode
	e.ErrorMessage = errorMessage
}

// RecordOpen records when the email was opened
func (e *EmailLog) RecordOpen() {
	if e.OpenedAt == nil {
		now := time.Now()
		e.OpenedAt = &now
	}
}

// RecordClick records when a link was clicked
func (e *EmailLog) RecordClick() {
	if e.ClickedAt == nil {
		now := time.Now()
		e.ClickedAt = &now
	}
}

// WasSent checks if email was sent
func (e *EmailLog) WasSent() bool {
	return e.Status == EmailStatusSent || e.Status == EmailStatusDelivered
}

// HasError checks if email has an error
func (e *EmailLog) HasError() bool {
	return e.Status == EmailStatusBounced || e.Status == EmailStatusFailed
}
