package entity

import (
	"time"

	"github.com/google/uuid"
)

// ContactType represents the type of contact
type ContactType string

const (
	ContactTypePrimary    ContactType = "PRIMARY"
	ContactTypeSales      ContactType = "SALES"
	ContactTypeTechnical  ContactType = "TECHNICAL"
	ContactTypeQuality    ContactType = "QUALITY"
	ContactTypeAccounting ContactType = "ACCOUNTING"
	ContactTypeLogistics  ContactType = "LOGISTICS"
)

// Contact represents a supplier contact person
type Contact struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SupplierID  uuid.UUID   `json:"supplier_id" gorm:"type:uuid;not null"`
	ContactType ContactType `json:"contact_type" gorm:"type:varchar(20);not null"`
	FullName    string      `json:"full_name" gorm:"type:varchar(255);not null"`
	Position    string      `json:"position" gorm:"type:varchar(100)"`
	Department  string      `json:"department" gorm:"type:varchar(100)"`
	Email       string      `json:"email" gorm:"type:varchar(255)"`
	Phone       string      `json:"phone" gorm:"type:varchar(50)"`
	Mobile      string      `json:"mobile" gorm:"type:varchar(50)"`
	IsPrimary   bool        `json:"is_primary" gorm:"default:false"`
	Notes       string      `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time   `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for GORM
func (Contact) TableName() string {
	return "supplier_contacts"
}
