package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryType represents the type of category
type CategoryType string

const (
	CategoryTypeMaterial CategoryType = "MATERIAL"
	CategoryTypeProduct  CategoryType = "PRODUCT"
)

// Category represents a hierarchical category for materials or products
type Category struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code         string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	NameEN       string         `gorm:"type:varchar(255);column:name_en" json:"name_en,omitempty"`
	Description  string         `gorm:"type:text" json:"description,omitempty"`
	CategoryType CategoryType   `gorm:"type:varchar(50);not null;default:'MATERIAL'" json:"category_type"`
	ParentID     *uuid.UUID     `gorm:"type:uuid" json:"parent_id,omitempty"`
	Path         string         `gorm:"type:varchar(1000);not null;default:'/'" json:"path"`
	Level        int            `gorm:"not null;default:0" json:"level"`
	SortOrder    int            `gorm:"not null;default:0" json:"sort_order"`
	Status       string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy    *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy    *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName specifies the table name
func (Category) TableName() string {
	return "categories"
}

// Validate validates category data
func (c *Category) Validate() error {
	if c.Code == "" {
		return fmt.Errorf("code is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.CategoryType != CategoryTypeMaterial && c.CategoryType != CategoryTypeProduct {
		return fmt.Errorf("invalid category type: %s", c.CategoryType)
	}
	return nil
}

// IsActive checks if category is active
func (c *Category) IsActive() bool {
	return c.Status == "active"
}

// GeneratePath generates the materialized path based on parent
func (c *Category) GeneratePath(parentPath string) {
	if parentPath == "" || parentPath == "/" {
		c.Path = fmt.Sprintf("/%s/", c.Code)
	} else {
		c.Path = fmt.Sprintf("%s%s/", parentPath, c.Code)
	}
}

// BeforeCreate hook to set default values
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Path == "" {
		c.Path = fmt.Sprintf("/%s/", c.Code)
	}
	return nil
}
