package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserNotificationType constants (visual type)
const (
	UserNotifTypeInfo    = "INFO"
	UserNotifTypeSuccess = "SUCCESS"
	UserNotifTypeWarning = "WARNING"
	UserNotifTypeError   = "ERROR"
)

// UserNotificationCategory constants
const (
	CategoryApproval = "APPROVAL"
	CategoryAlert    = "ALERT"
	CategoryMessage  = "MESSAGE"
	CategorySystem   = "SYSTEM"
)

// UserNotification represents an in-app notification for a user
type UserNotification struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	
	// Content
	Title            string     `gorm:"type:varchar(300);not null" json:"title"`
	Message          string     `gorm:"type:text;not null" json:"message"`
	
	// Categorization
	NotificationType string     `gorm:"type:varchar(50);default:'INFO'" json:"notification_type"` // INFO, SUCCESS, WARNING, ERROR
	Category         string     `gorm:"type:varchar(50);default:'SYSTEM'" json:"category"`        // APPROVAL, ALERT, MESSAGE, SYSTEM
	
	// Deep linking
	LinkURL          string     `gorm:"type:text" json:"link_url,omitempty"`
	LinkText         string     `gorm:"type:varchar(100)" json:"link_text,omitempty"`
	
	// Related entity
	EntityType       string     `gorm:"type:varchar(50)" json:"entity_type,omitempty"` // PO, PR, WO, LOT, etc.
	EntityID         *uuid.UUID `gorm:"type:uuid" json:"entity_id,omitempty"`
	
	// Read status
	IsRead           bool       `gorm:"default:false" json:"is_read"`
	ReadAt           *time.Time `gorm:"type:timestamp" json:"read_at,omitempty"`
	
	// Dismiss status
	IsDismissed      bool       `gorm:"default:false" json:"is_dismissed"`
	DismissedAt      *time.Time `gorm:"type:timestamp" json:"dismissed_at,omitempty"`
	
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (UserNotification) TableName() string {
	return "user_notifications"
}

// BeforeCreate sets defaults before creating
func (n *UserNotification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.NotificationType == "" {
		n.NotificationType = UserNotifTypeInfo
	}
	if n.Category == "" {
		n.Category = CategorySystem
	}
	return nil
}

// MarkAsRead marks the notification as read
func (n *UserNotification) MarkAsRead() {
	if !n.IsRead {
		n.IsRead = true
		now := time.Now()
		n.ReadAt = &now
	}
}

// MarkAsUnread marks the notification as unread
func (n *UserNotification) MarkAsUnread() {
	n.IsRead = false
	n.ReadAt = nil
}

// Dismiss dismisses the notification
func (n *UserNotification) Dismiss() {
	n.IsDismissed = true
	now := time.Now()
	n.DismissedAt = &now
}

// IsAlert checks if notification is an alert
func (n *UserNotification) IsAlert() bool {
	return n.Category == CategoryAlert
}

// IsApproval checks if notification requires approval action
func (n *UserNotification) IsApproval() bool {
	return n.Category == CategoryApproval
}

// IsWarningOrError checks if notification is warning or error
func (n *UserNotification) IsWarningOrError() bool {
	return n.NotificationType == UserNotifTypeWarning || n.NotificationType == UserNotifTypeError
}

// HasLink checks if notification has a deep link
func (n *UserNotification) HasLink() bool {
	return n.LinkURL != ""
}
