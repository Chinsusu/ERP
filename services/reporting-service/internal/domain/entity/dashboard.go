package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// LayoutType constants
const (
	LayoutTypeGrid     = "GRID"
	LayoutTypeFreeform = "FREEFORM"
)

// Visibility constants
const (
	VisibilityPrivate = "PRIVATE"
	VisibilityShared  = "SHARED"
	VisibilityPublic  = "PUBLIC"
)

// Dashboard represents a user dashboard
type Dashboard struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	LayoutType  string         `gorm:"type:varchar(50);default:'GRID'" json:"layout_type"`
	Layout      datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"layout,omitempty"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Visibility  string         `gorm:"type:varchar(50);default:'PRIVATE'" json:"visibility"`
	SharedWith  datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"shared_with,omitempty"`
	CreatedBy   *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Widgets []Widget `gorm:"foreignKey:DashboardID" json:"widgets,omitempty"`
}

// TableName specifies table name
func (Dashboard) TableName() string {
	return "dashboards"
}

// BeforeCreate sets defaults
func (d *Dashboard) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// IsPublic checks if dashboard is public
func (d *Dashboard) IsPublic() bool {
	return d.Visibility == VisibilityPublic
}

// IsUserDashboard checks if dashboard belongs to a user
func (d *Dashboard) IsUserDashboard(userID uuid.UUID) bool {
	return d.CreatedBy != nil && *d.CreatedBy == userID
}

// CanView checks if user can view dashboard
func (d *Dashboard) CanView(userID uuid.UUID) bool {
	if d.IsPublic() || d.IsUserDashboard(userID) {
		return true
	}
	// Check shared_with (simplified - should parse JSON)
	return d.Visibility == VisibilityShared
}
