package entity

import (
	"time"

	"github.com/google/uuid"
)

// POStatus represents the status of a PO
type POStatus string

const (
	POStatusDraft             POStatus = "DRAFT"
	POStatusSubmitted         POStatus = "SUBMITTED"
	POStatusConfirmed         POStatus = "CONFIRMED"
	POStatusPartiallyReceived POStatus = "PARTIALLY_RECEIVED"
	POStatusFullyReceived     POStatus = "FULLY_RECEIVED"
	POStatusClosed            POStatus = "CLOSED"
	POStatusCancelled         POStatus = "CANCELLED"
)

// PurchaseOrder represents a PO entity
type PurchaseOrder struct {
	ID                   uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PONumber             string     `json:"po_number" gorm:"column:po_number;type:varchar(20);unique;not null"`
	PODate               time.Time  `json:"po_date" gorm:"type:date;not null;default:CURRENT_DATE"`
	PRID                 *uuid.UUID `json:"pr_id" gorm:"column:pr_id;type:uuid"`
	SupplierID           uuid.UUID  `json:"supplier_id" gorm:"type:uuid;not null"`
	SupplierCode         string     `json:"supplier_code" gorm:"type:varchar(50)"`
	SupplierName         string     `json:"supplier_name" gorm:"type:varchar(255)"`
	Status               POStatus   `json:"status" gorm:"type:varchar(30);not null;default:'DRAFT'"`
	DeliveryAddress      string     `json:"delivery_address" gorm:"type:text"`
	DeliveryTerms        string     `json:"delivery_terms" gorm:"type:varchar(50);default:'EXW'"`
	PaymentTerms         string     `json:"payment_terms" gorm:"type:varchar(50);default:'Net 30'"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date" gorm:"type:date"`
	ActualDeliveryDate   *time.Time `json:"actual_delivery_date" gorm:"type:date"`
	TotalAmount          float64    `json:"total_amount" gorm:"type:decimal(15,2);not null;default:0"`
	TaxAmount            float64    `json:"tax_amount" gorm:"type:decimal(15,2);not null;default:0"`
	GrandTotal           float64    `json:"grand_total" gorm:"type:decimal(15,2);not null;default:0"`
	Currency             string     `json:"currency" gorm:"type:varchar(3);not null;default:'VND'"`
	Notes                string     `json:"notes" gorm:"type:text"`
	CreatedBy            uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
	SubmittedAt          *time.Time `json:"submitted_at"`
	ConfirmedAt          *time.Time `json:"confirmed_at"`
	ConfirmedBy          *uuid.UUID `json:"confirmed_by" gorm:"type:uuid"`
	CancelledAt          *time.Time `json:"cancelled_at"`
	CancelledBy          *uuid.UUID `json:"cancelled_by" gorm:"type:uuid"`
	CancellationReason   string     `json:"cancellation_reason" gorm:"type:text"`
	ClosedAt             *time.Time `json:"closed_at"`
	AmendmentCount       int        `json:"amendment_count" gorm:"default:0"`
	CreatedAt            time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt            *time.Time `json:"deleted_at" gorm:"index"`

	// Relations
	LineItems  []POLineItem  `json:"line_items,omitempty" gorm:"foreignKey:POID"`
	Amendments []POAmendment `json:"amendments,omitempty" gorm:"foreignKey:POID"`
	Receipts   []POReceipt   `json:"receipts,omitempty" gorm:"foreignKey:POID"`
}

// TableName returns the table name
func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// CalculateTotals calculates total, tax, and grand total
func (po *PurchaseOrder) CalculateTotals() {
	total := 0.0
	tax := 0.0
	for _, item := range po.LineItems {
		total += item.LineTotal
		tax += item.TaxAmount
	}
	po.TotalAmount = total
	po.TaxAmount = tax
	po.GrandTotal = total + tax
}

// Submit submits the PO
func (po *PurchaseOrder) Submit() {
	now := time.Now()
	po.Status = POStatusSubmitted
	po.SubmittedAt = &now
}

// Confirm confirms the PO
func (po *PurchaseOrder) Confirm(confirmedBy uuid.UUID) {
	now := time.Now()
	po.Status = POStatusConfirmed
	po.ConfirmedAt = &now
	po.ConfirmedBy = &confirmedBy
}

// Cancel cancels the PO
func (po *PurchaseOrder) Cancel(cancelledBy uuid.UUID, reason string) {
	now := time.Now()
	po.Status = POStatusCancelled
	po.CancelledAt = &now
	po.CancelledBy = &cancelledBy
	po.CancellationReason = reason
}

// Close closes the PO
func (po *PurchaseOrder) Close() {
	now := time.Now()
	po.Status = POStatusClosed
	po.ClosedAt = &now
}

// UpdateReceivedStatus updates status based on received quantities
func (po *PurchaseOrder) UpdateReceivedStatus() {
	fullyReceived := true
	partiallyReceived := false
	
	for _, item := range po.LineItems {
		if item.ReceivedQty < item.Quantity {
			fullyReceived = false
		}
		if item.ReceivedQty > 0 {
			partiallyReceived = true
		}
	}
	
	if fullyReceived {
		po.Status = POStatusFullyReceived
	} else if partiallyReceived {
		po.Status = POStatusPartiallyReceived
	}
}

// CanConfirm returns true if PO can be confirmed
func (po *PurchaseOrder) CanConfirm() bool {
	return po.Status == POStatusSubmitted || po.Status == POStatusDraft
}

// CanCancel returns true if PO can be cancelled
func (po *PurchaseOrder) CanCancel() bool {
	return po.Status == POStatusDraft || po.Status == POStatusSubmitted || po.Status == POStatusConfirmed
}

// CanClose returns true if PO can be closed
func (po *PurchaseOrder) CanClose() bool {
	return po.Status == POStatusFullyReceived
}
