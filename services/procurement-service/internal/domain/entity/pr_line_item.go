package entity

import (
	"time"

	"github.com/google/uuid"
)

// PRLineItem represents a line item in a PR
type PRLineItem struct {
	ID                    uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PRID                  uuid.UUID  `json:"pr_id" gorm:"column:pr_id;type:uuid;not null"`
	LineNumber            int        `json:"line_number" gorm:"type:int;not null"`
	MaterialID            uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	MaterialCode          string     `json:"material_code" gorm:"type:varchar(50)"`
	MaterialName          string     `json:"material_name" gorm:"type:varchar(255)"`
	Quantity              float64    `json:"quantity" gorm:"type:decimal(15,4);not null"`
	UOMID                 *uuid.UUID `json:"uom_id" gorm:"type:uuid"`
	UOMCode               string     `json:"uom_code" gorm:"type:varchar(20)"`
	UnitPrice             float64    `json:"unit_price" gorm:"type:decimal(15,4);not null;default:0"`
	LineTotal             float64    `json:"line_total" gorm:"type:decimal(15,2);not null;default:0"`
	Currency              string     `json:"currency" gorm:"type:varchar(3);not null;default:'VND'"`
	RequiredDate          *time.Time `json:"required_date" gorm:"type:date"`
	Specifications        string     `json:"specifications" gorm:"type:text"`
	SuggestedSupplierID   *uuid.UUID `json:"suggested_supplier_id" gorm:"type:uuid"`
	SuggestedSupplierName string     `json:"suggested_supplier_name" gorm:"type:varchar(255)"`
	Notes                 string     `json:"notes" gorm:"type:text"`
	CreatedAt             time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (PRLineItem) TableName() string {
	return "pr_line_items"
}

// CalculateLineTotal calculates the line total
func (li *PRLineItem) CalculateLineTotal() {
	li.LineTotal = li.Quantity * li.UnitPrice
}

// PRApproval represents an approval record
type PRApproval struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PRID          uuid.UUID     `json:"pr_id" gorm:"column:pr_id;type:uuid;not null"`
	ApproverID    uuid.UUID     `json:"approver_id" gorm:"type:uuid;not null"`
	ApproverName  string        `json:"approver_name" gorm:"type:varchar(255)"`
	ApprovalLevel ApprovalLevel `json:"approval_level" gorm:"type:varchar(30);not null"`
	Action        string        `json:"action" gorm:"type:varchar(20);not null"`
	Notes         string        `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time     `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (PRApproval) TableName() string {
	return "pr_approvals"
}
