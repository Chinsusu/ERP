package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// WOStatus represents work order status
type WOStatus string

const (
	WOStatusPlanned    WOStatus = "PLANNED"
	WOStatusReleased   WOStatus = "RELEASED"
	WOStatusInProgress WOStatus = "IN_PROGRESS"
	WOStatusQCPending  WOStatus = "QC_PENDING"
	WOStatusCompleted  WOStatus = "COMPLETED"
	WOStatusCancelled  WOStatus = "CANCELLED"
)

// WOPriority represents work order priority
type WOPriority string

const (
	WOPriorityLow    WOPriority = "LOW"
	WOPriorityNormal WOPriority = "NORMAL"
	WOPriorityHigh   WOPriority = "HIGH"
	WOPriorityUrgent WOPriority = "URGENT"
)

// WorkOrder represents a manufacturing work order
type WorkOrder struct {
	ID                uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WONumber          string      `json:"wo_number" gorm:"type:varchar(30);unique;not null"`
	WODate            time.Time   `json:"wo_date" gorm:"type:date;not null"`
	ProductID         uuid.UUID   `json:"product_id" gorm:"type:uuid;not null"`
	BOMID             uuid.UUID   `json:"bom_id" gorm:"type:uuid;not null"`
	Status            WOStatus    `json:"status" gorm:"type:varchar(30);default:'PLANNED'"`
	Priority          WOPriority  `json:"priority" gorm:"type:varchar(20);default:'NORMAL'"`
	PlannedQuantity   float64     `json:"planned_quantity" gorm:"type:decimal(15,4);not null"`
	UOMID             uuid.UUID   `json:"uom_id" gorm:"type:uuid;not null"`
	PlannedStartDate  *time.Time  `json:"planned_start_date"`
	PlannedEndDate    *time.Time  `json:"planned_end_date"`
	ActualStartDate   *time.Time  `json:"actual_start_date"`
	ActualEndDate     *time.Time  `json:"actual_end_date"`
	ActualQuantity    *float64    `json:"actual_quantity" gorm:"type:decimal(15,4)"`
	GoodQuantity      *float64    `json:"good_quantity" gorm:"type:decimal(15,4)"`
	RejectedQuantity  *float64    `json:"rejected_quantity" gorm:"type:decimal(15,4)"`
	YieldPercentage   *float64    `json:"yield_percentage" gorm:"type:decimal(5,2)"`
	BatchNumber       string      `json:"batch_number" gorm:"type:varchar(50)"`
	OutputLotID       *uuid.UUID  `json:"output_lot_id" gorm:"type:uuid"`
	SalesOrderID      *uuid.UUID  `json:"sales_order_id" gorm:"type:uuid"`
	ProductionLine    string      `json:"production_line" gorm:"type:varchar(50)"`
	Shift             string      `json:"shift" gorm:"type:varchar(20)"`
	SupervisorID      *uuid.UUID  `json:"supervisor_id" gorm:"type:uuid"`
	Notes             string      `json:"notes" gorm:"type:text"`
	CreatedBy         *uuid.UUID  `json:"created_by" gorm:"type:uuid"`
	UpdatedBy         *uuid.UUID  `json:"updated_by" gorm:"type:uuid"`
	CreatedAt         time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Associations
	Items          []WOLineItem       `json:"items,omitempty" gorm:"foreignKey:WorkOrderID"`
	MaterialIssues []WOMaterialIssue  `json:"material_issues,omitempty" gorm:"foreignKey:WorkOrderID"`
	BOM            *BOM               `json:"bom,omitempty" gorm:"foreignKey:BOMID"`
}

// TableName returns the table name
func (WorkOrder) TableName() string {
	return "work_orders"
}

// WOLineItem represents a planned material for work order
type WOLineItem struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkOrderID    uuid.UUID  `json:"work_order_id" gorm:"type:uuid;not null"`
	BOMLineItemID  *uuid.UUID `json:"bom_line_item_id" gorm:"type:uuid"`
	LineNumber     int        `json:"line_number" gorm:"not null"`
	MaterialID     uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	PlannedQuantity float64   `json:"planned_quantity" gorm:"type:decimal(15,4);not null"`
	IssuedQuantity float64    `json:"issued_quantity" gorm:"type:decimal(15,4);default:0"`
	UOMID          uuid.UUID  `json:"uom_id" gorm:"type:uuid;not null"`
	IsCritical     bool       `json:"is_critical" gorm:"default:false"`
	Notes          string     `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (WOLineItem) TableName() string {
	return "wo_line_items"
}

// WOMaterialIssue represents an actual material issue to a work order
type WOMaterialIssue struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkOrderID    uuid.UUID  `json:"work_order_id" gorm:"type:uuid;not null"`
	WOLineItemID   *uuid.UUID `json:"wo_line_item_id" gorm:"type:uuid"`
	IssueNumber    string     `json:"issue_number" gorm:"type:varchar(30);not null"`
	IssueDate      time.Time  `json:"issue_date" gorm:"default:CURRENT_TIMESTAMP"`
	MaterialID     uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	LotID          uuid.UUID  `json:"lot_id" gorm:"type:uuid;not null"`
	LotNumber      string     `json:"lot_number" gorm:"type:varchar(50);not null"`
	Quantity       float64    `json:"quantity" gorm:"type:decimal(15,4);not null"`
	UOMID          uuid.UUID  `json:"uom_id" gorm:"type:uuid;not null"`
	WMSMovementID  *uuid.UUID `json:"wms_movement_id" gorm:"type:uuid"`
	IssuedBy       *uuid.UUID `json:"issued_by" gorm:"type:uuid"`
	Notes          string     `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (WOMaterialIssue) TableName() string {
	return "wo_material_issues"
}

// Work Order business methods

// CanBeReleased returns true if WO can be released
func (w *WorkOrder) CanBeReleased() bool {
	return w.Status == WOStatusPlanned
}

// CanBeStarted returns true if WO can be started
func (w *WorkOrder) CanBeStarted() bool {
	return w.Status == WOStatusReleased
}

// CanBeCompleted returns true if WO can be completed
func (w *WorkOrder) CanBeCompleted() bool {
	return w.Status == WOStatusInProgress || w.Status == WOStatusQCPending
}

// CanBeCancelled returns true if WO can be cancelled
func (w *WorkOrder) CanBeCancelled() bool {
	return w.Status == WOStatusPlanned || w.Status == WOStatusReleased
}

// Release releases the work order
func (w *WorkOrder) Release() error {
	if !w.CanBeReleased() {
		return errors.New("work order cannot be released from current status")
	}
	w.Status = WOStatusReleased
	w.UpdatedAt = time.Now()
	return nil
}

// Start starts the work order
func (w *WorkOrder) Start(supervisorID uuid.UUID) error {
	if !w.CanBeStarted() {
		return errors.New("work order cannot be started from current status")
	}
	w.Status = WOStatusInProgress
	w.SupervisorID = &supervisorID
	now := time.Now()
	w.ActualStartDate = &now
	w.UpdatedAt = now
	return nil
}

// SendToQC sends the work order to QC
func (w *WorkOrder) SendToQC() error {
	if w.Status != WOStatusInProgress {
		return errors.New("only in-progress work orders can be sent to QC")
	}
	w.Status = WOStatusQCPending
	w.UpdatedAt = time.Now()
	return nil
}

// Complete completes the work order
func (w *WorkOrder) Complete(actualQty, goodQty, rejectedQty float64) error {
	if !w.CanBeCompleted() {
		return errors.New("work order cannot be completed from current status")
	}
	w.Status = WOStatusCompleted
	w.ActualQuantity = &actualQty
	w.GoodQuantity = &goodQty
	w.RejectedQuantity = &rejectedQty
	
	// Calculate yield
	if w.PlannedQuantity > 0 {
		yield := (goodQty / w.PlannedQuantity) * 100
		w.YieldPercentage = &yield
	}
	
	now := time.Now()
	w.ActualEndDate = &now
	w.UpdatedAt = now
	return nil
}

// Cancel cancels the work order
func (w *WorkOrder) Cancel() error {
	if !w.CanBeCancelled() {
		return errors.New("work order cannot be cancelled from current status")
	}
	w.Status = WOStatusCancelled
	w.UpdatedAt = time.Now()
	return nil
}

// GetOutstandingQuantity returns the remaining quantity to be issued
func (item *WOLineItem) GetOutstandingQuantity() float64 {
	return item.PlannedQuantity - item.IssuedQuantity
}

// IsFullyIssued returns true if all planned quantity has been issued
func (item *WOLineItem) IsFullyIssued() bool {
	return item.IssuedQuantity >= item.PlannedQuantity
}
