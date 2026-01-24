package entity

import (
	"time"

	"github.com/google/uuid"
)

// ReturnStatus represents return status
type ReturnStatus string

const (
	ReturnStatusPending   ReturnStatus = "PENDING"
	ReturnStatusApproved  ReturnStatus = "APPROVED"
	ReturnStatusReceived  ReturnStatus = "RECEIVED"
	ReturnStatusInspected ReturnStatus = "INSPECTED"
	ReturnStatusCompleted ReturnStatus = "COMPLETED"
	ReturnStatusRejected  ReturnStatus = "REJECTED"
)

// ReturnType represents return type
type ReturnType string

const (
	ReturnTypeRefund   ReturnType = "REFUND"
	ReturnTypeExchange ReturnType = "EXCHANGE"
	ReturnTypeCredit   ReturnType = "CREDIT"
)

// ItemCondition represents returned item condition
type ItemCondition string

const (
	ItemConditionGood      ItemCondition = "GOOD"
	ItemConditionDamaged   ItemCondition = "DAMAGED"
	ItemConditionDefective ItemCondition = "DEFECTIVE"
	ItemConditionExpired   ItemCondition = "EXPIRED"
)

// ItemAction represents action for returned item
type ItemAction string

const (
	ItemActionRefund   ItemAction = "REFUND"
	ItemActionExchange ItemAction = "EXCHANGE"
	ItemActionCredit   ItemAction = "CREDIT"
	ItemActionDispose  ItemAction = "DISPOSE"
)

// Return represents a sales return
type Return struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReturnNumber string        `json:"return_number" gorm:"type:varchar(20);unique;not null"`
	SalesOrderID uuid.UUID     `json:"sales_order_id" gorm:"type:uuid;not null"`
	SalesOrder   *SalesOrder   `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	ShipmentID   *uuid.UUID    `json:"shipment_id" gorm:"type:uuid"`
	Shipment     *Shipment     `json:"shipment,omitempty" gorm:"foreignKey:ShipmentID"`
	ReturnDate   time.Time     `json:"return_date" gorm:"type:date;not null"`
	ReturnReason string        `json:"return_reason" gorm:"type:varchar(100)"`
	ReturnType   ReturnType    `json:"return_type" gorm:"type:varchar(20);default:'REFUND'"`
	Status       ReturnStatus  `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	Subtotal     float64       `json:"subtotal" gorm:"type:decimal(18,2);default:0"`
	RefundAmount float64       `json:"refund_amount" gorm:"type:decimal(18,2);default:0"`
	Notes        string        `json:"notes" gorm:"type:text"`
	ApprovedBy   *uuid.UUID    `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt   *time.Time    `json:"approved_at" gorm:"type:timestamp"`
	CompletedAt  *time.Time    `json:"completed_at" gorm:"type:timestamp"`
	CreatedBy    *uuid.UUID    `json:"created_by" gorm:"type:uuid"`
	UpdatedBy    *uuid.UUID    `json:"updated_by" gorm:"type:uuid"`
	CreatedAt    time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	LineItems []ReturnLineItem `json:"line_items,omitempty" gorm:"foreignKey:ReturnID"`
}

func (Return) TableName() string {
	return "returns"
}

// Approve approves the return
func (r *Return) Approve(userID uuid.UUID) {
	now := time.Now()
	r.Status = ReturnStatusApproved
	r.ApprovedBy = &userID
	r.ApprovedAt = &now
	r.UpdatedAt = now
}

// MarkReceived marks return as received
func (r *Return) MarkReceived() {
	r.Status = ReturnStatusReceived
	r.UpdatedAt = time.Now()
}

// MarkInspected marks return as inspected
func (r *Return) MarkInspected() {
	r.Status = ReturnStatusInspected
	r.UpdatedAt = time.Now()
}

// Complete completes the return
func (r *Return) Complete() {
	now := time.Now()
	r.Status = ReturnStatusCompleted
	r.CompletedAt = &now
	r.UpdatedAt = now
}

// Reject rejects the return
func (r *Return) Reject() {
	r.Status = ReturnStatusRejected
	r.UpdatedAt = time.Now()
}

// CalculateTotals calculates return totals
func (r *Return) CalculateTotals() {
	r.Subtotal = 0
	for _, item := range r.LineItems {
		r.Subtotal += item.Quantity * item.UnitPrice
	}
	r.RefundAmount = r.Subtotal
}

// ReturnLineItem represents a line item in return
type ReturnLineItem struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReturnID     uuid.UUID     `json:"return_id" gorm:"type:uuid;not null"`
	SOLineItemID *uuid.UUID    `json:"so_line_item_id" gorm:"type:uuid"`
	ProductID    uuid.UUID     `json:"product_id" gorm:"type:uuid;not null"`
	ProductCode  string        `json:"product_code" gorm:"type:varchar(50)"`
	ProductName  string        `json:"product_name" gorm:"type:varchar(200)"`
	Quantity     float64       `json:"quantity" gorm:"type:decimal(18,3);not null"`
	UnitPrice    float64       `json:"unit_price" gorm:"type:decimal(18,2);not null"`
	Reason       string        `json:"reason" gorm:"type:varchar(200)"`
	Condition    ItemCondition `json:"condition" gorm:"type:varchar(20);default:'GOOD'"`
	Action       ItemAction    `json:"action" gorm:"type:varchar(20);default:'REFUND'"`
	CreatedAt    time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (ReturnLineItem) TableName() string {
	return "return_line_items"
}
