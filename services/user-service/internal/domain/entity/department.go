package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Department represents an organizational department
type Department struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	ParentID  *uuid.UUID     `gorm:"type:uuid" json:"parent_id,omitempty"`
	ManagerID *uuid.UUID     `gorm:"type:uuid" json:"manager_id,omitempty"`
	Level     int            `gorm:"not null;default:0" json:"level"`
	Path      string         `gorm:"type:varchar(500);not null;default:'/'" json:"path"`
	Status    string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Parent   *Department   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Manager  *User         `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Children []Department  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Users    []User        `gorm:"foreignKey:DepartmentID" json:"users,omitempty"`
}

// TableName specifies the table name
func (Department) TableName() string {
	return "departments"
}

// Validate validates department data
func (d *Department) Validate() error {
	if d.Code == "" {
		return fmt.Errorf("code is required")
	}
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if d.Status != "active" && d.Status != "inactive" {
		return fmt.Errorf("invalid status: %s", d.Status)
	}
	return nil
}

// UpdatePath updates the materialized path based on parent
func (d *Department) UpdatePath(parentPath string) {
	if d.ParentID == nil {
		d.Path = fmt.Sprintf("/%s/", d.Code)
		d.Level = 0
	} else {
		d.Path = fmt.Sprintf("%s%s/", parentPath, d.Code)
		d.Level = strings.Count(parentPath, "/")
	}
}

// IsActive checks if department is active
func (d *Department) IsActive() bool {
	return d.Status == "active"
}

// IsRoot checks if department is root (no parent)
func (d *Department) IsRoot() bool {
	return d.ParentID == nil
}
