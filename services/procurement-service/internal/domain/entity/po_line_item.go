package entity

import (
	"time"

	"github.com/google/uuid"
)

// POLineItemStatus represents the status of a PO line item
type POLineItemStatus string

const (
	POLineStatusPending   POLineItemStatus = "PENDING"
	POLineStatusPartial   POLineItemStatus = "PARTIAL"
	POLineStatusComplete  POLineItemStatus = "COMPLETE"
	POLineStatusCancelled POLineItemStatus = "CANCELLED"
)

// POLineItem represents a line item in a PO
type POLineItem struct {
	ID             uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	POID           uuid.UUID        `json:"po_id" gorm:"column:po_id;type:uuid;not null"`
	PRLineItemID   *uuid.UUID       `json:"pr_line_item_id" gorm:"type:uuid"`
	LineNumber     int              `json:"line_number" gorm:"type:int;not null"`
	MaterialID     uuid.UUID        `json:"material_id" gorm:"type:uuid;not null"`
	MaterialCode   string           `json:"material_code" gorm:"type:varchar(50)"`
	MaterialName   string           `json:"material_name" gorm:"type:varchar(255)"`
	Quantity       float64          `json:"quantity" gorm:"type:decimal(15,4);not null"`
	ReceivedQty    float64          `json:"received_qty" gorm:"type:decimal(15,4);not null;default:0"`
	PendingQty     float64          `json:"pending_qty" gorm:"type:decimal(15,4);not null"`
	UOMID          *uuid.UUID       `json:"uom_id" gorm:"type:uuid"`
	UOMCode        string           `json:"uom_code" gorm:"type:varchar(20)"`
	UnitPrice      float64          `json:"unit_price" gorm:"type:decimal(15,4);not null"`
	LineTotal      float64          `json:"line_total" gorm:"type:decimal(15,2);not null"`
	TaxRate        float64          `json:"tax_rate" gorm:"type:decimal(5,2);default:0"`
	TaxAmount      float64          `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	Currency       string           `json:"currency" gorm:"type:varchar(3);not null;default:'VND'"`
	ExpectedDate   *time.Time       `json:"expected_date" gorm:"type:date"`
	Specifications string           `json:"specifications" gorm:"type:text"`
	Notes          string           `json:"notes" gorm:"type:text"`
	Status         POLineItemStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt      time.Time        `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time        `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (POLineItem) TableName() string {
	return "po_line_items"
}

// CalculateLineTotal calculates line total and tax
func (li *POLineItem) CalculateLineTotal() {
	li.LineTotal = li.Quantity * li.UnitPrice
	li.TaxAmount = li.LineTotal * li.TaxRate / 100
	li.PendingQty = li.Quantity - li.ReceivedQty
}

// UpdateReceivedQty updates received qty and status
func (li *POLineItem) UpdateReceivedQty(qty float64) {
	li.ReceivedQty += qty
	li.PendingQty = li.Quantity - li.ReceivedQty
	
	if li.ReceivedQty >= li.Quantity {
		li.Status = POLineStatusComplete
	} else if li.ReceivedQty > 0 {
		li.Status = POLineStatusPartial
	}
}

// POAmendment represents an amendment to a PO
type POAmendment struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	POID            uuid.UUID `json:"po_id" gorm:"column:po_id;type:uuid;not null"`
	AmendmentNumber int       `json:"amendment_number" gorm:"type:int;not null"`
	FieldChanged    string    `json:"field_changed" gorm:"type:varchar(100);not null"`
	OldValue        string    `json:"old_value" gorm:"type:text"`
	NewValue        string    `json:"new_value" gorm:"type:text"`
	Reason          string    `json:"reason" gorm:"type:text"`
	AmendedBy       uuid.UUID `json:"amended_by" gorm:"type:uuid;not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (POAmendment) TableName() string {
	return "po_amendments"
}

// QCStatus represents QC inspection status
type QCStatus string

const (
	QCStatusPending QCStatus = "PENDING"
	QCStatusPassed  QCStatus = "PASSED"
	QCStatusFailed  QCStatus = "FAILED"
	QCStatusPartial QCStatus = "PARTIAL"
)

// POReceipt represents a receipt record (linked to WMS GRN)
type POReceipt struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	POID            uuid.UUID  `json:"po_id" gorm:"column:po_id;type:uuid;not null"`
	POLineItemID    uuid.UUID  `json:"po_line_item_id" gorm:"type:uuid;not null"`
	GRNID           *uuid.UUID `json:"grn_id" gorm:"type:uuid"`
	GRNNumber       string     `json:"grn_number" gorm:"type:varchar(50)"`
	ReceivedQty     float64    `json:"received_qty" gorm:"type:decimal(15,4);not null"`
	ReceivedDate    time.Time  `json:"received_date" gorm:"type:date;not null;default:CURRENT_DATE"`
	ReceivedBy      *uuid.UUID `json:"received_by" gorm:"type:uuid"`
	QCStatus        QCStatus   `json:"qc_status" gorm:"type:varchar(20);default:'PENDING'"`
	QCNotes         string     `json:"qc_notes" gorm:"type:text"`
	BatchNumber     string     `json:"batch_number" gorm:"type:varchar(100)"`
	LotNumber       string     `json:"lot_number" gorm:"type:varchar(100)"`
	ExpiryDate      *time.Time `json:"expiry_date" gorm:"type:date"`
	StorageLocation string     `json:"storage_location" gorm:"type:varchar(100)"`
	Notes           string     `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (POReceipt) TableName() string {
	return "po_receipts"
}
