package entity

import (
	"time"

	"github.com/google/uuid"
)

// IssueType represents the type of goods issue
type IssueType string

const (
	IssueTypeProduction  IssueType = "PRODUCTION"
	IssueTypeSales       IssueType = "SALES"
	IssueTypeSample      IssueType = "SAMPLE"
	IssueTypeScrap       IssueType = "SCRAP"
	IssueTypeAdjustment  IssueType = "ADJUSTMENT"
)

// GoodsIssueStatus represents goods issue status
type GoodsIssueStatus string

const (
	GoodsIssueStatusDraft     GoodsIssueStatus = "DRAFT"
	GoodsIssueStatusConfirmed GoodsIssueStatus = "CONFIRMED"
	GoodsIssueStatusCompleted GoodsIssueStatus = "COMPLETED"
	GoodsIssueStatusCancelled GoodsIssueStatus = "CANCELLED"
)

// GoodsIssue represents a goods issue document
type GoodsIssue struct {
	ID              uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IssueNumber     string           `json:"issue_number" gorm:"type:varchar(30);unique;not null"` // GI-YYYY-XXXX
	IssueDate       time.Time        `json:"issue_date" gorm:"type:date;not null"`
	IssueType       IssueType        `json:"issue_type" gorm:"type:varchar(30);not null"`
	ReferenceType   ReferenceType    `json:"reference_type" gorm:"type:varchar(30)"`
	ReferenceID     *uuid.UUID       `json:"reference_id" gorm:"type:uuid"`
	ReferenceNumber string           `json:"reference_number" gorm:"type:varchar(30)"`
	WarehouseID     uuid.UUID        `json:"warehouse_id" gorm:"type:uuid;not null"`
	Status          GoodsIssueStatus `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	Notes           string           `json:"notes" gorm:"type:text"`
	IssuedBy        *uuid.UUID       `json:"issued_by" gorm:"type:uuid"`
	CreatedAt       time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Warehouse *Warehouse       `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	LineItems []GILineItem     `json:"line_items,omitempty" gorm:"foreignKey:GoodsIssueID"`
}

// TableName returns the table name
func (GoodsIssue) TableName() string {
	return "goods_issues"
}

// Confirm confirms the goods issue
func (g *GoodsIssue) Confirm() {
	g.Status = GoodsIssueStatusConfirmed
	g.UpdatedAt = time.Now()
}

// Complete completes the goods issue
func (g *GoodsIssue) Complete() {
	g.Status = GoodsIssueStatusCompleted
	g.UpdatedAt = time.Now()
}

// Cancel cancels the goods issue
func (g *GoodsIssue) Cancel() {
	g.Status = GoodsIssueStatusCancelled
	g.UpdatedAt = time.Now()
}

// GILineItem represents a line item in goods issue
type GILineItem struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GoodsIssueID uuid.UUID  `json:"goods_issue_id" gorm:"type:uuid;not null"`
	LineNumber   int        `json:"line_number" gorm:"not null"`
	MaterialID   uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	RequestedQty float64    `json:"requested_qty" gorm:"type:decimal(15,4);not null"`
	IssuedQty    float64    `json:"issued_qty" gorm:"type:decimal(15,4);not null"`
	UnitID       uuid.UUID  `json:"unit_id" gorm:"type:uuid;not null"`
	LotID        *uuid.UUID `json:"lot_id" gorm:"type:uuid"`
	LocationID   *uuid.UUID `json:"location_id" gorm:"type:uuid"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Lot      *Lot      `json:"lot,omitempty" gorm:"foreignKey:LotID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName returns the table name
func (GILineItem) TableName() string {
	return "gi_line_items"
}
