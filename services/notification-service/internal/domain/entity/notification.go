package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationType constants
const (
	NotificationTypeEmail = "EMAIL"
	NotificationTypeInApp = "IN_APP"
	NotificationTypeSMS   = "SMS"
	NotificationTypePush  = "PUSH"
	NotificationTypeBoth  = "BOTH"
)

// NotificationStatus constants
const (
	NotificationStatusPending  = "PENDING"
	NotificationStatusSent     = "SENT"
	NotificationStatusFailed   = "FAILED"
	NotificationStatusRetrying = "RETRYING"
)

// NotificationPriority constants
const (
	PriorityHigh   = "HIGH"
	PriorityNormal = "NORMAL"
	PriorityLow    = "LOW"
)

// Notification represents an outbound notification (email, SMS, push)
type Notification struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NotificationType string         `gorm:"type:varchar(50);not null" json:"notification_type"`
	TemplateID       *uuid.UUID     `gorm:"type:uuid" json:"template_id,omitempty"`
	
	// Recipient information
	RecipientUserID  *uuid.UUID     `gorm:"type:uuid" json:"recipient_user_id,omitempty"`
	RecipientEmail   string         `gorm:"type:varchar(255)" json:"recipient_email,omitempty"`
	RecipientPhone   string         `gorm:"type:varchar(50)" json:"recipient_phone,omitempty"`
	
	// Content
	Subject          string         `gorm:"type:text" json:"subject,omitempty"`
	Body             string         `gorm:"type:text;not null" json:"body"`
	
	// Additional data
	Metadata         datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	
	// Status tracking
	Status           string         `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	Priority         string         `gorm:"type:varchar(20);default:'NORMAL'" json:"priority"`
	
	// Timestamps and retry info
	SentAt           *time.Time     `gorm:"type:timestamp" json:"sent_at,omitempty"`
	FailedAt         *time.Time     `gorm:"type:timestamp" json:"failed_at,omitempty"`
	RetryCount       int            `gorm:"default:0" json:"retry_count"`
	MaxRetries       int            `gorm:"default:3" json:"max_retries"`
	NextRetryAt      *time.Time     `gorm:"type:timestamp" json:"next_retry_at,omitempty"`
	ErrorMessage     string         `gorm:"type:text" json:"error_message,omitempty"`
	
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	Template         *NotificationTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}

// TableName specifies the table name
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate sets defaults before creating
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.Status == "" {
		n.Status = NotificationStatusPending
	}
	if n.Priority == "" {
		n.Priority = PriorityNormal
	}
	return nil
}

// MarkAsSent marks the notification as sent
func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.Status = NotificationStatusSent
	n.SentAt = &now
}

// MarkAsFailed marks the notification as failed
func (n *Notification) MarkAsFailed(err string) {
	now := time.Now()
	n.Status = NotificationStatusFailed
	n.FailedAt = &now
	n.ErrorMessage = err
}

// ScheduleRetry schedules a retry with exponential backoff
func (n *Notification) ScheduleRetry(err string) bool {
	n.RetryCount++
	if n.RetryCount >= n.MaxRetries {
		n.MarkAsFailed(err)
		return false
	}
	
	n.Status = NotificationStatusRetrying
	n.ErrorMessage = err
	
	// Exponential backoff: 5min, 15min, 45min
	backoff := time.Duration(n.RetryCount*n.RetryCount*5) * time.Minute
	nextRetry := time.Now().Add(backoff)
	n.NextRetryAt = &nextRetry
	
	return true
}

// CanRetry checks if the notification can be retried
func (n *Notification) CanRetry() bool {
	return n.RetryCount < n.MaxRetries
}

// IsPending checks if notification is pending
func (n *Notification) IsPending() bool {
	return n.Status == NotificationStatusPending
}

// IsHighPriority checks if notification is high priority
func (n *Notification) IsHighPriority() bool {
	return n.Priority == PriorityHigh
}
