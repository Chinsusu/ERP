package entity

import (
	"time"

	"github.com/google/uuid"
)

// QuotationStatus represents quotation status
type QuotationStatus string

const (
	QuotationStatusDraft     QuotationStatus = "DRAFT"
	QuotationStatusSent      QuotationStatus = "SENT"
	QuotationStatusAccepted  QuotationStatus = "ACCEPTED"
	QuotationStatusRejected  QuotationStatus = "REJECTED"
	QuotationStatusExpired   QuotationStatus = "EXPIRED"
	QuotationStatusConverted QuotationStatus = "CONVERTED"
)

// Quotation represents a sales quotation
type Quotation struct {
	ID                 uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	QuotationNumber    string               `json:"quotation_number" gorm:"type:varchar(20);unique;not null"`
	CustomerID         uuid.UUID            `json:"customer_id" gorm:"type:uuid;not null"`
	Customer           *Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	QuotationDate      time.Time            `json:"quotation_date" gorm:"type:date;not null"`
	ValidUntil         time.Time            `json:"valid_until" gorm:"type:date;not null"`
	Subtotal           float64              `json:"subtotal" gorm:"type:decimal(18,2);default:0"`
	DiscountPercent    float64              `json:"discount_percent" gorm:"type:decimal(5,2);default:0"`
	DiscountAmount     float64              `json:"discount_amount" gorm:"type:decimal(18,2);default:0"`
	TaxPercent         float64              `json:"tax_percent" gorm:"type:decimal(5,2);default:10"`
	TaxAmount          float64              `json:"tax_amount" gorm:"type:decimal(18,2);default:0"`
	TotalAmount        float64              `json:"total_amount" gorm:"type:decimal(18,2);default:0"`
	Status             QuotationStatus      `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	Notes              string               `json:"notes" gorm:"type:text"`
	TermsAndConditions string               `json:"terms_and_conditions" gorm:"type:text"`
	CreatedBy          *uuid.UUID           `json:"created_by" gorm:"type:uuid"`
	UpdatedBy          *uuid.UUID           `json:"updated_by" gorm:"type:uuid"`
	CreatedAt          time.Time            `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time            `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	LineItems []QuotationLineItem `json:"line_items,omitempty" gorm:"foreignKey:QuotationID"`
}

func (Quotation) TableName() string {
	return "quotations"
}

// IsExpired checks if the quotation is expired
func (q *Quotation) IsExpired() bool {
	return time.Now().After(q.ValidUntil)
}

// CanBeSent checks if quotation can be sent
func (q *Quotation) CanBeSent() bool {
	return q.Status == QuotationStatusDraft && !q.IsExpired() && len(q.LineItems) > 0
}

// CanBeConverted checks if quotation can be converted to sales order
func (q *Quotation) CanBeConverted() bool {
	return (q.Status == QuotationStatusSent || q.Status == QuotationStatusAccepted) && !q.IsExpired()
}

// Send marks quotation as sent
func (q *Quotation) Send() {
	q.Status = QuotationStatusSent
	q.UpdatedAt = time.Now()
}

// Accept marks quotation as accepted
func (q *Quotation) Accept() {
	q.Status = QuotationStatusAccepted
	q.UpdatedAt = time.Now()
}

// Reject marks quotation as rejected
func (q *Quotation) Reject() {
	q.Status = QuotationStatusRejected
	q.UpdatedAt = time.Now()
}

// MarkConverted marks quotation as converted
func (q *Quotation) MarkConverted() {
	q.Status = QuotationStatusConverted
	q.UpdatedAt = time.Now()
}

// MarkExpired marks quotation as expired
func (q *Quotation) MarkExpired() {
	q.Status = QuotationStatusExpired
	q.UpdatedAt = time.Now()
}

// CalculateTotals calculates subtotal, tax, and total
func (q *Quotation) CalculateTotals() {
	q.Subtotal = 0
	for _, item := range q.LineItems {
		q.Subtotal += item.LineTotal
	}

	// Apply discount
	if q.DiscountPercent > 0 {
		q.DiscountAmount = q.Subtotal * q.DiscountPercent / 100
	}

	// Calculate tax
	taxableAmount := q.Subtotal - q.DiscountAmount
	q.TaxAmount = taxableAmount * q.TaxPercent / 100

	// Calculate total
	q.TotalAmount = taxableAmount + q.TaxAmount
}

// QuotationLineItem represents a line item in quotation
type QuotationLineItem struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	QuotationID     uuid.UUID  `json:"quotation_id" gorm:"type:uuid;not null"`
	LineNumber      int        `json:"line_number" gorm:"not null"`
	ProductID       uuid.UUID  `json:"product_id" gorm:"type:uuid;not null"`
	ProductCode     string     `json:"product_code" gorm:"type:varchar(50)"`
	ProductName     string     `json:"product_name" gorm:"type:varchar(200)"`
	Quantity        float64    `json:"quantity" gorm:"type:decimal(18,3);not null"`
	UomID           *uuid.UUID `json:"uom_id" gorm:"type:uuid"`
	UnitPrice       float64    `json:"unit_price" gorm:"type:decimal(18,2);not null"`
	DiscountPercent float64    `json:"discount_percent" gorm:"type:decimal(5,2);default:0"`
	DiscountAmount  float64    `json:"discount_amount" gorm:"type:decimal(18,2);default:0"`
	TaxPercent      float64    `json:"tax_percent" gorm:"type:decimal(5,2);default:10"`
	TaxAmount       float64    `json:"tax_amount" gorm:"type:decimal(18,2);default:0"`
	LineTotal       float64    `json:"line_total" gorm:"type:decimal(18,2);default:0"`
	Notes           string     `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (QuotationLineItem) TableName() string {
	return "quotation_line_items"
}

// CalculateLineTotal calculates the line total
func (li *QuotationLineItem) CalculateLineTotal() {
	subtotal := li.Quantity * li.UnitPrice

	// Apply discount
	if li.DiscountPercent > 0 {
		li.DiscountAmount = subtotal * li.DiscountPercent / 100
	}

	// Calculate tax
	taxableAmount := subtotal - li.DiscountAmount
	li.TaxAmount = taxableAmount * li.TaxPercent / 100

	// Line total (including tax)
	li.LineTotal = taxableAmount + li.TaxAmount
}
