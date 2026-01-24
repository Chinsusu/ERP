package entity

import (
	"time"

	"github.com/google/uuid"
)

// BatchTraceability links material lots to product lots for traceability
type BatchTraceability struct {
	ID                  uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorkOrderID         uuid.UUID  `json:"work_order_id" gorm:"type:uuid;not null"`
	WOMaterialIssueID   *uuid.UUID `json:"wo_material_issue_id" gorm:"type:uuid"`
	MaterialID          uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	MaterialLotID       uuid.UUID  `json:"material_lot_id" gorm:"type:uuid;not null"`
	MaterialLotNumber   string     `json:"material_lot_number" gorm:"type:varchar(50);not null"`
	MaterialQuantity    float64    `json:"material_quantity" gorm:"type:decimal(15,4);not null"`
	MaterialUOMID       uuid.UUID  `json:"material_uom_id" gorm:"type:uuid;not null"`
	SupplierLotNumber   string     `json:"supplier_lot_number" gorm:"type:varchar(100)"`
	ProductID           uuid.UUID  `json:"product_id" gorm:"type:uuid;not null"`
	ProductLotID        *uuid.UUID `json:"product_lot_id" gorm:"type:uuid"`
	ProductLotNumber    string     `json:"product_lot_number" gorm:"type:varchar(50)"`
	TraceDate           time.Time  `json:"trace_date" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt           time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (BatchTraceability) TableName() string {
	return "batch_traceability"
}

// BackwardTraceResult represents the result of a backward trace
type BackwardTraceResult struct {
	FinishedLot   FinishedLotInfo    `json:"finished_lot"`
	WorkOrder     WorkOrderInfo      `json:"work_order"`
	MaterialsUsed []MaterialTraceInfo `json:"materials_used"`
}

// ForwardTraceResult represents the result of a forward trace
type ForwardTraceResult struct {
	MaterialLot     MaterialLotInfo    `json:"material_lot"`
	UsedInProducts  []ProductTraceInfo `json:"used_in_products"`
	TotalQtyUsed    float64            `json:"total_quantity_used"`
	RemainingQty    float64            `json:"remaining_quantity"`
}

// FinishedLotInfo contains finished product lot info
type FinishedLotInfo struct {
	LotNumber        string    `json:"lot_number"`
	ProductCode      string    `json:"product_code"`
	ProductName      string    `json:"product_name"`
	Quantity         float64   `json:"quantity"`
	ManufacturedDate time.Time `json:"manufactured_date"`
}

// WorkOrderInfo contains work order summary
type WorkOrderInfo struct {
	WONumber   string `json:"wo_number"`
	Supervisor string `json:"supervisor"`
}

// MaterialTraceInfo contains material trace info
type MaterialTraceInfo struct {
	MaterialCode    string  `json:"material_code"`
	MaterialName    string  `json:"material_name"`
	LotNumber       string  `json:"lot_number"`
	Quantity        float64 `json:"quantity"`
	UOM             string  `json:"uom"`
	Supplier        string  `json:"supplier"`
	SupplierLot     string  `json:"supplier_lot"`
}

// MaterialLotInfo contains material lot info for forward trace
type MaterialLotInfo struct {
	LotNumber    string `json:"lot_number"`
	MaterialCode string `json:"material_code"`
	MaterialName string `json:"material_name"`
	Supplier     string `json:"supplier"`
}

// ProductTraceInfo contains product trace info for forward trace
type ProductTraceInfo struct {
	ProductLot     string    `json:"product_lot"`
	ProductName    string    `json:"product_name"`
	WONumber       string    `json:"wo_number"`
	QuantityUsed   float64   `json:"quantity_used"`
	ProductionDate time.Time `json:"production_date"`
}
