package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AlertRuleType constants
const (
	RuleTypeStockLow        = "STOCK_LOW"
	RuleTypeLotExpiry       = "LOT_EXPIRY"
	RuleTypeCertExpiry      = "CERT_EXPIRY"
	RuleTypeApprovalPending = "APPROVAL_PENDING"
	RuleTypeTempOutOfRange  = "TEMP_OUT_OF_RANGE"
)

// AlertRule represents a configurable alert rule
type AlertRule struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RuleCode               string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"rule_code"`
	Name                   string         `gorm:"type:varchar(200);not null" json:"name"`
	Description            string         `gorm:"type:text" json:"description,omitempty"`
	
	// Rule type and conditions
	RuleType               string         `gorm:"type:varchar(50);not null" json:"rule_type"`
	Conditions             datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"conditions"`
	
	// Template and recipients
	NotificationTemplateID *uuid.UUID     `gorm:"type:uuid" json:"notification_template_id,omitempty"`
	NotificationType       string         `gorm:"type:varchar(50);default:'IN_APP'" json:"notification_type"` // EMAIL, IN_APP, BOTH
	Recipients             datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"recipients"` // Array of {role: string} or {user_id: uuid}
	
	// Schedule
	CheckInterval          string         `gorm:"type:varchar(20);default:'1h'" json:"check_interval"` // 1h, 30m, daily, etc.
	LastCheckedAt          *time.Time     `gorm:"type:timestamp" json:"last_checked_at,omitempty"`
	NextCheckAt            *time.Time     `gorm:"type:timestamp" json:"next_check_at,omitempty"`
	
	// Status
	IsActive               bool           `gorm:"default:true" json:"is_active"`
	
	CreatedBy              *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relationships
	Template               *NotificationTemplate `gorm:"foreignKey:NotificationTemplateID" json:"template,omitempty"`
}

// TableName specifies the table name
func (AlertRule) TableName() string {
	return "alert_rules"
}

// BeforeCreate sets defaults before creating
func (r *AlertRule) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Activate activates the rule
func (r *AlertRule) Activate() {
	r.IsActive = true
}

// Deactivate deactivates the rule
func (r *AlertRule) Deactivate() {
	r.IsActive = false
}

// UpdateLastChecked updates the last checked timestamp
func (r *AlertRule) UpdateLastChecked() {
	now := time.Now()
	r.LastCheckedAt = &now
	
	// Calculate next check time based on interval
	var nextCheck time.Time
	switch r.CheckInterval {
	case "30m":
		nextCheck = now.Add(30 * time.Minute)
	case "1h":
		nextCheck = now.Add(1 * time.Hour)
	case "daily":
		nextCheck = now.Add(24 * time.Hour)
	default:
		nextCheck = now.Add(1 * time.Hour)
	}
	r.NextCheckAt = &nextCheck
}

// ShouldCheck checks if the rule should be checked now
func (r *AlertRule) ShouldCheck() bool {
	if !r.IsActive {
		return false
	}
	if r.NextCheckAt == nil {
		return true
	}
	return time.Now().After(*r.NextCheckAt)
}

// IsStockRule checks if this is a stock-related rule
func (r *AlertRule) IsStockRule() bool {
	return r.RuleType == RuleTypeStockLow
}

// IsExpiryRule checks if this is an expiry-related rule
func (r *AlertRule) IsExpiryRule() bool {
	return r.RuleType == RuleTypeLotExpiry || r.RuleType == RuleTypeCertExpiry
}

// RequiresEmail checks if rule requires email notification
func (r *AlertRule) RequiresEmail() bool {
	return r.NotificationType == NotificationTypeEmail || r.NotificationType == NotificationTypeBoth
}

// RequiresInApp checks if rule requires in-app notification
func (r *AlertRule) RequiresInApp() bool {
	return r.NotificationType == NotificationTypeInApp || r.NotificationType == NotificationTypeBoth
}
