package entity

import (
	"bytes"
	"html/template"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TemplateCode     string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"template_code"`
	Name             string         `gorm:"type:varchar(200);not null" json:"name"`
	NotificationType string         `gorm:"type:varchar(50);not null" json:"notification_type"` // EMAIL, IN_APP, SMS, PUSH, BOTH
	SubjectTemplate  string         `gorm:"type:text" json:"subject_template,omitempty"`
	BodyTemplate     string         `gorm:"type:text;not null" json:"body_template"`
	Variables        datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"variables"` // Available template variables
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	CreatedBy        *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// BeforeCreate sets defaults before creating
func (t *NotificationTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// RenderSubject renders the subject template with provided data
func (t *NotificationTemplate) RenderSubject(data map[string]interface{}) (string, error) {
	if t.SubjectTemplate == "" {
		return "", nil
	}
	
	tmpl, err := template.New("subject").Parse(t.SubjectTemplate)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// RenderBody renders the body template with provided data
func (t *NotificationTemplate) RenderBody(data map[string]interface{}) (string, error) {
	tmpl, err := template.New("body").Parse(t.BodyTemplate)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// IsEmailTemplate checks if template is for email
func (t *NotificationTemplate) IsEmailTemplate() bool {
	return t.NotificationType == NotificationTypeEmail || t.NotificationType == NotificationTypeBoth
}

// IsInAppTemplate checks if template is for in-app notifications
func (t *NotificationTemplate) IsInAppTemplate() bool {
	return t.NotificationType == NotificationTypeInApp || t.NotificationType == NotificationTypeBoth
}

// Activate activates the template
func (t *NotificationTemplate) Activate() {
	t.IsActive = true
}

// Deactivate deactivates the template
func (t *NotificationTemplate) Deactivate() {
	t.IsActive = false
}
