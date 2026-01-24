package entity

import (
	"time"

	"github.com/google/uuid"
)

// SOStatus represents sales order status
type SOStatus string

const (
	SOStatusDraft            SOStatus = "DRAFT"
	SOStatusConfirmed        SOStatus = "CONFIRMED"
	SOStatusProcessing       SOStatus = "PROCESSING"
	SOStatusPartiallyShipped SOStatus = "PARTIALLY_SHIPPED"
	SOStatusShipped          SOStatus = "SHIPPED"
	SOStatusDelivered        SOStatus = "DELIVERED"
	SOStatusCancelled        SOStatus = "CANCELLED"
)

// PaymentMethod represents payment method
type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "CASH"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodCredit       PaymentMethod = "CREDIT"
	PaymentMethodCOD          PaymentMethod = "COD"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPartial PaymentStatus = "PARTIAL"
	PaymentStatusPaid    PaymentStatus = "PAID"
)

// SalesOrder represents a sales order
type SalesOrder struct {
	ID                 uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SONumber           string        `json:"so_number" gorm:"type:varchar(20);unique;not null"`
	CustomerID         uuid.UUID     `json:"customer_id" gorm:"type:uuid;not null"`
	Customer           *Customer     `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	QuotationID        *uuid.UUID    `json:"quotation_id" gorm:"type:uuid"`
	Quotation          *Quotation    `json:"quotation,omitempty" gorm:"foreignKey:QuotationID"`
	SODate             time.Time     `json:"so_date" gorm:"type:date;not null"`
	DeliveryDate       *time.Time    `json:"delivery_date" gorm:"type:date"`
	DeliveryAddress    string        `json:"delivery_address" gorm:"type:text"`
	BillingAddress     string        `json:"billing_address" gorm:"type:text"`
	Subtotal           float64       `json:"subtotal" gorm:"type:decimal(18,2);default:0"`
	DiscountPercent    float64       `json:"discount_percent" gorm:"type:decimal(5,2);default:0"`
	DiscountAmount     float64       `json:"discount_amount" gorm:"type:decimal(18,2);default:0"`
	TaxPercent         float64       `json:"tax_percent" gorm:"type:decimal(5,2);default:10"`
	TaxAmount          float64       `json:"tax_amount" gorm:"type:decimal(18,2);default:0"`
	TotalAmount        float64       `json:"total_amount" gorm:"type:decimal(18,2);default:0"`
	Status             SOStatus      `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	PaymentMethod      PaymentMethod `json:"payment_method" gorm:"type:varchar(20);default:'BANK_TRANSFER'"`
	PaymentStatus      PaymentStatus `json:"payment_status" gorm:"type:varchar(20);default:'PENDING'"`
	ConfirmedAt        *time.Time    `json:"confirmed_at" gorm:"type:timestamp"`
	ConfirmedBy        *uuid.UUID    `json:"confirmed_by" gorm:"type:uuid"`
	ShippedAt          *time.Time    `json:"shipped_at" gorm:"type:timestamp"`
	DeliveredAt        *time.Time    `json:"delivered_at" gorm:"type:timestamp"`
	CancelledAt        *time.Time    `json:"cancelled_at" gorm:"type:timestamp"`
	CancelledBy        *uuid.UUID    `json:"cancelled_by" gorm:"type:uuid"`
	CancellationReason string        `json:"cancellation_reason" gorm:"type:text"`
	Notes              string        `json:"notes" gorm:"type:text"`
	CreatedBy          *uuid.UUID    `json:"created_by" gorm:"type:uuid"`
	UpdatedBy          *uuid.UUID    `json:"updated_by" gorm:"type:uuid"`
	CreatedAt          time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time     `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	LineItems []SOLineItem `json:"line_items,omitempty" gorm:"foreignKey:SalesOrderID"`
}

func (SalesOrder) TableName() string {
	return "sales_orders"
}

// CanBeConfirmed checks if order can be confirmed
func (so *SalesOrder) CanBeConfirmed() bool {
	return so.Status == SOStatusDraft && len(so.LineItems) > 0
}

// CanBeCancelled checks if order can be cancelled
func (so *SalesOrder) CanBeCancelled() bool {
	return so.Status == SOStatusDraft || so.Status == SOStatusConfirmed
}

// CanBeShipped checks if order can be shipped
func (so *SalesOrder) CanBeShipped() bool {
	return so.Status == SOStatusConfirmed || so.Status == SOStatusProcessing || so.Status == SOStatusPartiallyShipped
}

// Confirm confirms the sales order
func (so *SalesOrder) Confirm(userID uuid.UUID) {
	now := time.Now()
	so.Status = SOStatusConfirmed
	so.ConfirmedAt = &now
	so.ConfirmedBy = &userID
	so.UpdatedAt = now
}

// StartProcessing marks order as processing
func (so *SalesOrder) StartProcessing() {
	so.Status = SOStatusProcessing
	so.UpdatedAt = time.Now()
}

// MarkPartiallyShipped marks order as partially shipped
func (so *SalesOrder) MarkPartiallyShipped() {
	now := time.Now()
	so.Status = SOStatusPartiallyShipped
	if so.ShippedAt == nil {
		so.ShippedAt = &now
	}
	so.UpdatedAt = now
}

// MarkShipped marks order as fully shipped
func (so *SalesOrder) MarkShipped() {
	now := time.Now()
	so.Status = SOStatusShipped
	so.ShippedAt = &now
	so.UpdatedAt = now
}

// MarkDelivered marks order as delivered
func (so *SalesOrder) MarkDelivered() {
	now := time.Now()
	so.Status = SOStatusDelivered
	so.DeliveredAt = &now
	so.UpdatedAt = now
}

// Cancel cancels the sales order
func (so *SalesOrder) Cancel(userID uuid.UUID, reason string) {
	now := time.Now()
	so.Status = SOStatusCancelled
	so.CancelledAt = &now
	so.CancelledBy = &userID
	so.CancellationReason = reason
	so.UpdatedAt = now
}

// CalculateTotals calculates subtotal, tax, and total
func (so *SalesOrder) CalculateTotals() {
	so.Subtotal = 0
	for _, item := range so.LineItems {
		so.Subtotal += item.LineTotal
	}

	// Apply discount
	if so.DiscountPercent > 0 {
		so.DiscountAmount = so.Subtotal * so.DiscountPercent / 100
	}

	// Calculate tax
	taxableAmount := so.Subtotal - so.DiscountAmount
	so.TaxAmount = taxableAmount * so.TaxPercent / 100

	// Calculate total
	so.TotalAmount = taxableAmount + so.TaxAmount
}

// GetShippedPercentage returns the percentage of items shipped
func (so *SalesOrder) GetShippedPercentage() float64 {
	if len(so.LineItems) == 0 {
		return 0
	}

	totalQty := 0.0
	shippedQty := 0.0
	for _, item := range so.LineItems {
		totalQty += item.Quantity
		shippedQty += item.ShippedQuantity
	}

	if totalQty == 0 {
		return 0
	}
	return (shippedQty / totalQty) * 100
}

// IsFullyShipped checks if all items are shipped
func (so *SalesOrder) IsFullyShipped() bool {
	for _, item := range so.LineItems {
		if item.ShippedQuantity < item.Quantity {
			return false
		}
	}
	return true
}

// SOLineItem represents a line item in sales order
type SOLineItem struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SalesOrderID    uuid.UUID  `json:"sales_order_id" gorm:"type:uuid;not null"`
	LineNumber      int        `json:"line_number" gorm:"not null"`
	ProductID       uuid.UUID  `json:"product_id" gorm:"type:uuid;not null"`
	ProductCode     string     `json:"product_code" gorm:"type:varchar(50)"`
	ProductName     string     `json:"product_name" gorm:"type:varchar(200)"`
	Quantity        float64    `json:"quantity" gorm:"type:decimal(18,3);not null"`
	ShippedQuantity float64    `json:"shipped_quantity" gorm:"type:decimal(18,3);default:0"`
	UomID           *uuid.UUID `json:"uom_id" gorm:"type:uuid"`
	UnitPrice       float64    `json:"unit_price" gorm:"type:decimal(18,2);not null"`
	DiscountPercent float64    `json:"discount_percent" gorm:"type:decimal(5,2);default:0"`
	DiscountAmount  float64    `json:"discount_amount" gorm:"type:decimal(18,2);default:0"`
	TaxPercent      float64    `json:"tax_percent" gorm:"type:decimal(5,2);default:10"`
	TaxAmount       float64    `json:"tax_amount" gorm:"type:decimal(18,2);default:0"`
	LineTotal       float64    `json:"line_total" gorm:"type:decimal(18,2);default:0"`
	ReservationID   *uuid.UUID `json:"reservation_id" gorm:"type:uuid"`
	Notes           string     `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (SOLineItem) TableName() string {
	return "so_line_items"
}

// CalculateLineTotal calculates the line total
func (li *SOLineItem) CalculateLineTotal() {
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

// GetRemainingQuantity returns quantity not yet shipped
func (li *SOLineItem) GetRemainingQuantity() float64 {
	return li.Quantity - li.ShippedQuantity
}

// IsFullyShipped checks if item is fully shipped
func (li *SOLineItem) IsFullyShipped() bool {
	return li.ShippedQuantity >= li.Quantity
}

// AddShippedQuantity adds to shipped quantity
func (li *SOLineItem) AddShippedQuantity(qty float64) {
	li.ShippedQuantity += qty
	if li.ShippedQuantity > li.Quantity {
		li.ShippedQuantity = li.Quantity
	}
	li.UpdatedAt = time.Now()
}
