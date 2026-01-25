package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ReportType constants
const (
	ReportTypeInventory   = "INVENTORY"
	ReportTypeSales       = "SALES"
	ReportTypeProcurement = "PROCUREMENT"
	ReportTypeProduction  = "PRODUCTION"
	ReportTypeFinancial   = "FINANCIAL"
)

// ReportDefinition represents a report template
type ReportDefinition struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code               string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
	Name               string         `gorm:"type:varchar(200);not null" json:"name"`
	Description        string         `gorm:"type:text" json:"description,omitempty"`
	ReportType         string         `gorm:"type:varchar(50);not null" json:"report_type"`
	Category           string         `gorm:"type:varchar(50)" json:"category,omitempty"`
	QueryTemplate      string         `gorm:"type:text;not null" json:"query_template"`
	DataSource         string         `gorm:"type:varchar(100)" json:"data_source,omitempty"`
	Parameters         datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"parameters"`
	Columns            datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"columns"`
	DefaultSort        datatypes.JSON `gorm:"type:jsonb" json:"default_sort,omitempty"`
	GroupBy            datatypes.JSON `gorm:"type:jsonb" json:"group_by,omitempty"`
	RequiredPermission string         `gorm:"type:varchar(100)" json:"required_permission,omitempty"`
	IsSystem           bool           `gorm:"default:false" json:"is_system"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	CreatedBy          *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies table name
func (ReportDefinition) TableName() string {
	return "report_definitions"
}

// BeforeCreate sets defaults
func (r *ReportDefinition) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// ParameterDef represents a report parameter definition
type ParameterDef struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`     // string, integer, date, uuid, boolean
	Required bool        `json:"required"`
	Default  interface{} `json:"default,omitempty"`
	Label    string      `json:"label,omitempty"`
	Options  []string    `json:"options,omitempty"`
}

// ColumnDef represents a report column definition
type ColumnDef struct {
	Field  string `json:"field"`
	Header string `json:"header"`
	Type   string `json:"type,omitempty"`   // string, number, date, datetime, currency, percentage
	Format string `json:"format,omitempty"` // Custom format
	Width  int    `json:"width,omitempty"`
}
