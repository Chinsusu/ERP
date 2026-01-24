package entity

import (
	"time"

	"github.com/google/uuid"
)

// ASLStatus represents the status of an ASL entry
type ASLStatus string

const (
	ASLStatusActive   ASLStatus = "ACTIVE"
	ASLStatusInactive ASLStatus = "INACTIVE"
	ASLStatusExpired  ASLStatus = "EXPIRED"
)

// ApprovedSupplier represents an entry in the Approved Supplier List (ASL)
type ApprovedSupplier struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SupplierID   uuid.UUID  `json:"supplier_id" gorm:"type:uuid;not null"`
	MaterialID   uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	Priority     int        `json:"priority" gorm:"type:int;not null;default:1"`
	LeadTimeDays int        `json:"lead_time_days" gorm:"type:int;not null;default:14"`
	MinOrderQty  float64    `json:"min_order_qty" gorm:"type:decimal(15,4);not null;default:1"`
	UnitPrice    float64    `json:"unit_price" gorm:"type:decimal(15,4);not null"`
	Currency     string     `json:"currency" gorm:"type:varchar(3);not null;default:'VND'"`
	ValidFrom    time.Time  `json:"valid_from" gorm:"type:date;not null"`
	ValidTo      *time.Time `json:"valid_to" gorm:"type:date"`
	ApprovedBy   uuid.UUID  `json:"approved_by" gorm:"type:uuid;not null"`
	ApprovedAt   time.Time  `json:"approved_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status       ASLStatus  `json:"status" gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	Notes        string     `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Supplier *Supplier `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
}

// TableName returns the table name for GORM
func (ApprovedSupplier) TableName() string {
	return "approved_supplier_list"
}

// IsActive returns true if the ASL entry is active and valid
func (a *ApprovedSupplier) IsActive() bool {
	now := time.Now()
	if a.Status != ASLStatusActive {
		return false
	}
	if a.ValidFrom.After(now) {
		return false
	}
	if a.ValidTo != nil && a.ValidTo.Before(now) {
		return false
	}
	return true
}

// GetPriorityLabel returns human-readable priority label
func (a *ApprovedSupplier) GetPriorityLabel() string {
	switch a.Priority {
	case 1:
		return "Primary"
	case 2:
		return "Secondary"
	case 3:
		return "Backup"
	default:
		return "Unknown"
	}
}
