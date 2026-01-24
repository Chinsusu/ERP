package grpc

import "google.golang.org/protobuf/types/known/timestamppb"

// Proto message types - simplified version without protoc generation
// In production, these would be generated from wms.proto

// CheckStockRequest represents stock availability check request
type CheckStockRequest struct {
	MaterialId        string  `json:"material_id"`
	RequestedQuantity float64 `json:"requested_quantity"`
}

// StockAvailabilityResponse represents stock availability response
type StockAvailabilityResponse struct {
	MaterialId        string  `json:"material_id"`
	TotalQuantity     float64 `json:"total_quantity"`
	ReservedQuantity  float64 `json:"reserved_quantity"`
	AvailableQuantity float64 `json:"available_quantity"`
	IsAvailable       bool    `json:"is_available"`
	ShortageQuantity  float64 `json:"shortage_quantity"`
}

// ReserveStockRequest represents reserve stock request
type ReserveStockRequest struct {
	MaterialId      string  `json:"material_id"`
	Quantity        float64 `json:"quantity"`
	UnitId          string  `json:"unit_id"`
	ReservationType string  `json:"reservation_type"`
	ReferenceId     string  `json:"reference_id"`
	ReferenceNumber string  `json:"reference_number"`
}

// ReserveStockResponse represents reserve stock response
type ReserveStockResponse struct {
	ReservationId    string  `json:"reservation_id"`
	ReservedQuantity float64 `json:"reserved_quantity"`
	Success          bool    `json:"success"`
	Message          string  `json:"message"`
}

// ReleaseReservationRequest represents release reservation request
type ReleaseReservationRequest struct {
	ReservationId string `json:"reservation_id"`
}

// ReleaseReservationResponse represents release reservation response
type ReleaseReservationResponse struct {
	Success bool `json:"success"`
}

// MaterialQuantity represents material quantity
type MaterialQuantity struct {
	MaterialId string  `json:"material_id"`
	Quantity   float64 `json:"quantity"`
	UnitId     string  `json:"unit_id"`
}

// IssueStockRequest represents issue stock request
type IssueStockRequest struct {
	WarehouseId     string              `json:"warehouse_id"`
	ReferenceType   string              `json:"reference_type"`
	ReferenceId     string              `json:"reference_id"`
	ReferenceNumber string              `json:"reference_number"`
	Items           []*MaterialQuantity `json:"items"`
	IssuedBy        string              `json:"issued_by"`
}

// LotIssued represents lot issued in response
type LotIssued struct {
	LotId      string                     `json:"lot_id"`
	LotNumber  string                     `json:"lot_number"`
	Quantity   float64                    `json:"quantity"`
	ExpiryDate *timestamppb.Timestamp `json:"expiry_date"`
	LocationId string                     `json:"location_id"`
}

// IssueLineItem represents issue line item
type IssueLineItem struct {
	MaterialId     string       `json:"material_id"`
	IssuedQuantity float64      `json:"issued_quantity"`
	LotsUsed       []*LotIssued `json:"lots_used"`
}

// IssueStockResponse represents issue stock response
type IssueStockResponse struct {
	IssueNumber string           `json:"issue_number"`
	LineItems   []*IssueLineItem `json:"line_items"`
	Success     bool             `json:"success"`
	Message     string           `json:"message"`
}

// GetLotRequest represents get lot request
type GetLotRequest struct {
	LotId string `json:"lot_id"`
}

// LotInfo represents lot info
type LotInfo struct {
	LotId             string                     `json:"lot_id"`
	LotNumber         string                     `json:"lot_number"`
	MaterialId        string                     `json:"material_id"`
	ExpiryDate        *timestamppb.Timestamp `json:"expiry_date"`
	QcStatus          string                     `json:"qc_status"`
	Status            string                     `json:"status"`
	AvailableQuantity float64                    `json:"available_quantity"`
}

// LotInfoResponse represents lot info response
type LotInfoResponse struct {
	Lot *LotInfo `json:"lot"`
}

// GetLotsByMaterialRequest represents get lots by material request
type GetLotsByMaterialRequest struct {
	MaterialId    string `json:"material_id"`
	AvailableOnly bool   `json:"available_only"`
}

// LotsResponse represents lots response
type LotsResponse struct {
	Lots []*LotInfo `json:"lots"`
}

// ReceiveLineItem represents receive line item
type ReceiveLineItem struct {
	MaterialId        string                     `json:"material_id"`
	Quantity          float64                    `json:"quantity"`
	UnitId            string                     `json:"unit_id"`
	SupplierLotNumber string                     `json:"supplier_lot_number"`
	ExpiryDate        *timestamppb.Timestamp `json:"expiry_date"`
	LocationId        string                     `json:"location_id"`
}

// ReceiveStockRequest represents receive stock request
type ReceiveStockRequest struct {
	GrnId       string             `json:"grn_id"`
	PoId        string             `json:"po_id"`
	WarehouseId string             `json:"warehouse_id"`
	Items       []*ReceiveLineItem `json:"items"`
	ReceivedBy  string             `json:"received_by"`
}

// ReceiveStockResponse represents receive stock response
type ReceiveStockResponse struct {
	GrnNumber  string   `json:"grn_number"`
	LotNumbers []string `json:"lot_numbers"`
	Success    bool     `json:"success"`
}
