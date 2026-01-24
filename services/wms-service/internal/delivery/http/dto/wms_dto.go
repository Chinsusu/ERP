package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateGRNRequest represents request to create GRN
type CreateGRNRequest struct {
	GRNDate            string              `json:"grn_date" binding:"required"`
	POID               *uuid.UUID          `json:"po_id"`
	PONumber           string              `json:"po_number"`
	SupplierID         *uuid.UUID          `json:"supplier_id"`
	WarehouseID        uuid.UUID           `json:"warehouse_id" binding:"required"`
	DeliveryNoteNumber string              `json:"delivery_note_number"`
	VehicleNumber      string              `json:"vehicle_number"`
	Notes              string              `json:"notes"`
	Items              []CreateGRNItemRequest `json:"items" binding:"required,dive"`
}

// CreateGRNItemRequest represents GRN line item request
type CreateGRNItemRequest struct {
	POLineItemID      *uuid.UUID `json:"po_line_item_id"`
	MaterialID        uuid.UUID  `json:"material_id" binding:"required"`
	ExpectedQty       *float64   `json:"expected_qty"`
	ReceivedQty       float64    `json:"received_qty" binding:"required,gt=0"`
	UnitID            uuid.UUID  `json:"unit_id" binding:"required"`
	SupplierLotNumber string     `json:"supplier_lot_number"`
	ManufacturedDate  *string    `json:"manufactured_date"`
	ExpiryDate        string     `json:"expiry_date" binding:"required"`
	LocationID        *uuid.UUID `json:"location_id"`
}

// CompleteGRNRequest represents request to complete GRN
type CompleteGRNRequest struct {
	QCStatus string `json:"qc_status" binding:"required,oneof=PASSED FAILED"`
	QCNotes  string `json:"qc_notes"`
}

// GRNResponse represents GRN response
type GRNResponse struct {
	ID                 uuid.UUID              `json:"id"`
	GRNNumber          string                 `json:"grn_number"`
	GRNDate            string                 `json:"grn_date"`
	PONumber           string                 `json:"po_number,omitempty"`
	WarehouseID        uuid.UUID              `json:"warehouse_id"`
	WarehouseName      string                 `json:"warehouse_name,omitempty"`
	DeliveryNoteNumber string                 `json:"delivery_note_number,omitempty"`
	VehicleNumber      string                 `json:"vehicle_number,omitempty"`
	Status             string                 `json:"status"`
	QCStatus           string                 `json:"qc_status"`
	QCNotes            string                 `json:"qc_notes,omitempty"`
	Notes              string                 `json:"notes,omitempty"`
	CompletedAt        *time.Time             `json:"completed_at,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	LineItems          []GRNLineItemResponse  `json:"line_items,omitempty"`
}

// GRNLineItemResponse represents GRN line item response
type GRNLineItemResponse struct {
	ID                uuid.UUID  `json:"id"`
	LineNumber        int        `json:"line_number"`
	MaterialID        uuid.UUID  `json:"material_id"`
	ExpectedQty       *float64   `json:"expected_qty,omitempty"`
	ReceivedQty       float64    `json:"received_qty"`
	AcceptedQty       *float64   `json:"accepted_qty,omitempty"`
	RejectedQty       float64    `json:"rejected_qty"`
	LotNumber         string     `json:"lot_number,omitempty"`
	SupplierLotNumber string     `json:"supplier_lot_number,omitempty"`
	ExpiryDate        string     `json:"expiry_date"`
	LocationCode      string     `json:"location_code,omitempty"`
	QCStatus          string     `json:"qc_status"`
}

// IssueStockRequest represents request to issue stock
type IssueStockRequest struct {
	IssueDate       string                `json:"issue_date" binding:"required"`
	IssueType       string                `json:"issue_type" binding:"required,oneof=PRODUCTION SALES SAMPLE SCRAP ADJUSTMENT"`
	ReferenceID     *uuid.UUID            `json:"reference_id"`
	ReferenceNumber string                `json:"reference_number"`
	WarehouseID     uuid.UUID             `json:"warehouse_id" binding:"required"`
	Notes           string                `json:"notes"`
	Items           []IssueStockItemRequest `json:"items" binding:"required,dive"`
}

// IssueStockItemRequest represents issue stock item request
type IssueStockItemRequest struct {
	MaterialID uuid.UUID `json:"material_id" binding:"required"`
	Quantity   float64   `json:"quantity" binding:"required,gt=0"`
	UnitID     uuid.UUID `json:"unit_id" binding:"required"`
}

// IssueStockResponse represents issue stock response
type IssueStockResponse struct {
	IssueNumber string            `json:"issue_number"`
	LotsIssued  []LotIssuedResponse `json:"lots_issued"`
}

// LotIssuedResponse represents lot issued in response
type LotIssuedResponse struct {
	LotNumber  string  `json:"lot_number"`
	Quantity   float64 `json:"quantity"`
	ExpiryDate string  `json:"expiry_date"`
}

// ReserveStockRequest represents request to reserve stock
type ReserveStockRequest struct {
	MaterialID      uuid.UUID `json:"material_id" binding:"required"`
	Quantity        float64   `json:"quantity" binding:"required,gt=0"`
	UnitID          uuid.UUID `json:"unit_id" binding:"required"`
	ReservationType string    `json:"reservation_type" binding:"required,oneof=SALES_ORDER WORK_ORDER TRANSFER"`
	ReferenceID     uuid.UUID `json:"reference_id" binding:"required"`
	ReferenceNumber string    `json:"reference_number"`
	ExpiresAt       *string   `json:"expires_at"`
}

// ReserveStockResponse represents reserve stock response
type ReserveStockResponse struct {
	ReservationID    uuid.UUID `json:"reservation_id"`
	ReservedQuantity float64   `json:"reserved_quantity"`
}

// TransferStockRequest represents request to transfer stock
type TransferStockRequest struct {
	MaterialID     uuid.UUID  `json:"material_id" binding:"required"`
	LotID          *uuid.UUID `json:"lot_id"`
	FromLocationID uuid.UUID  `json:"from_location_id" binding:"required"`
	ToLocationID   uuid.UUID  `json:"to_location_id" binding:"required"`
	Quantity       float64    `json:"quantity" binding:"required,gt=0"`
	UnitID         uuid.UUID  `json:"unit_id" binding:"required"`
	Reason         string     `json:"reason"`
}

// AdjustmentRequest represents stock adjustment request
type AdjustmentRequest struct {
	AdjustmentDate string     `json:"adjustment_date" binding:"required"`
	AdjustmentType string     `json:"adjustment_type" binding:"required,oneof=CYCLE_COUNT DAMAGE EXPIRY CORRECTION"`
	LocationID     uuid.UUID  `json:"location_id" binding:"required"`
	MaterialID     uuid.UUID  `json:"material_id" binding:"required"`
	LotID          *uuid.UUID `json:"lot_id"`
	SystemQty      float64    `json:"system_qty" binding:"required"`
	ActualQty      float64    `json:"actual_qty" binding:"required"`
	UnitID         uuid.UUID  `json:"unit_id" binding:"required"`
	Reason         string     `json:"reason"`
	Notes          string     `json:"notes"`
}
