package dto

import (
	"time"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// ===== BOM DTOs =====

// CreateBOMRequest is the request for creating a BOM
type CreateBOMRequest struct {
	BOMNumber            string                       `json:"bom_number" binding:"required"`
	ProductID            uuid.UUID                    `json:"product_id" binding:"required"`
	Version              int                          `json:"version"`
	Name                 string                       `json:"name" binding:"required"`
	Description          string                       `json:"description"`
	BatchSize            float64                      `json:"batch_size" binding:"required"`
	BatchUnitID          uuid.UUID                    `json:"batch_uom_id" binding:"required"`
	ConfidentialityLevel string                       `json:"confidentiality_level"`
	LaborCost            float64                      `json:"labor_cost"`
	OverheadCost         float64                      `json:"overhead_cost"`
	FormulaDetails       *entity.FormulaDetails       `json:"formula_details"`
	Items                []CreateBOMItemRequest       `json:"items" binding:"required"`
}

// CreateBOMItemRequest is the request for a BOM line item
type CreateBOMItemRequest struct {
	LineNumber      int       `json:"line_number"`
	MaterialID      uuid.UUID `json:"material_id" binding:"required"`
	ItemType        string    `json:"item_type"`
	Quantity        float64   `json:"quantity" binding:"required"`
	UOMID           uuid.UUID `json:"uom_id" binding:"required"`
	QuantityMin     *float64  `json:"quantity_min"`
	QuantityMax     *float64  `json:"quantity_max"`
	IsCritical      bool      `json:"is_critical"`
	ScrapPercentage float64   `json:"scrap_percentage"`
	UnitCost        float64   `json:"unit_cost"`
	Notes           string    `json:"notes"`
}

// BOMResponse is the response for a BOM
type BOMResponse struct {
	ID                   uuid.UUID               `json:"id"`
	BOMNumber            string                  `json:"bom_number"`
	ProductID            uuid.UUID               `json:"product_id"`
	Version              int                     `json:"version"`
	Name                 string                  `json:"name"`
	Description          string                  `json:"description"`
	Status               string                  `json:"status"`
	BatchSize            float64                 `json:"batch_size"`
	ConfidentialityLevel string                  `json:"confidentiality_level"`
	MaterialCost         float64                 `json:"material_cost"`
	LaborCost            float64                 `json:"labor_cost"`
	OverheadCost         float64                 `json:"overhead_cost"`
	TotalCost            float64                 `json:"total_cost"`
	EffectiveFrom        *time.Time              `json:"effective_from,omitempty"`
	EffectiveTo          *time.Time              `json:"effective_to,omitempty"`
	Items                []BOMItemResponse       `json:"items,omitempty"`
	FormulaDetails       *entity.FormulaDetails  `json:"formula_details,omitempty"`
	Message              string                  `json:"message,omitempty"`
	CreatedAt            time.Time               `json:"created_at"`
}

// BOMItemResponse is the response for a BOM line item
type BOMItemResponse struct {
	ID         uuid.UUID `json:"id"`
	LineNumber int       `json:"line_number"`
	MaterialID uuid.UUID `json:"material_id"`
	ItemType   string    `json:"item_type"`
	Quantity   float64   `json:"quantity"`
	UOMID      uuid.UUID `json:"uom_id"`
	IsCritical bool      `json:"is_critical"`
	UnitCost   float64   `json:"unit_cost"`
	TotalCost  float64   `json:"total_cost"`
	Notes      string    `json:"notes,omitempty"`
}

// ===== Work Order DTOs =====

// CreateWORequest is the request for creating a work order
type CreateWORequest struct {
	ProductID        uuid.UUID  `json:"product_id" binding:"required"`
	BOMID            uuid.UUID  `json:"bom_id" binding:"required"`
	PlannedQuantity  float64    `json:"planned_quantity" binding:"required"`
	UOMID            uuid.UUID  `json:"uom_id" binding:"required"`
	PlannedStartDate string     `json:"planned_start_date"`
	PlannedEndDate   string     `json:"planned_end_date"`
	BatchNumber      string     `json:"batch_number"`
	SalesOrderID     *uuid.UUID `json:"sales_order_id"`
	ProductionLine   string     `json:"production_line"`
	Shift            string     `json:"shift"`
	Priority         string     `json:"priority"`
	Notes            string     `json:"notes"`
}

// WOResponse is the response for a work order
type WOResponse struct {
	ID               uuid.UUID    `json:"id"`
	WONumber         string       `json:"wo_number"`
	WODate           time.Time    `json:"wo_date"`
	ProductID        uuid.UUID    `json:"product_id"`
	BOMID            uuid.UUID    `json:"bom_id"`
	Status           string       `json:"status"`
	Priority         string       `json:"priority"`
	PlannedQuantity  float64      `json:"planned_quantity"`
	ActualQuantity   *float64     `json:"actual_quantity,omitempty"`
	GoodQuantity     *float64     `json:"good_quantity,omitempty"`
	RejectedQuantity *float64     `json:"rejected_quantity,omitempty"`
	YieldPercentage  *float64     `json:"yield_percentage,omitempty"`
	BatchNumber      string       `json:"batch_number"`
	PlannedStartDate *time.Time   `json:"planned_start_date,omitempty"`
	PlannedEndDate   *time.Time   `json:"planned_end_date,omitempty"`
	ActualStartDate  *time.Time   `json:"actual_start_date,omitempty"`
	ActualEndDate    *time.Time   `json:"actual_end_date,omitempty"`
	ProductionLine   string       `json:"production_line,omitempty"`
	Notes            string       `json:"notes,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
}

// StartWORequest is the request for starting a work order
type StartWORequest struct {
	SupervisorID uuid.UUID `json:"supervisor_id" binding:"required"`
}

// CompleteWORequest is the request for completing a work order
type CompleteWORequest struct {
	ActualQuantity   float64 `json:"actual_quantity" binding:"required"`
	GoodQuantity     float64 `json:"good_quantity" binding:"required"`
	RejectedQuantity float64 `json:"rejected_quantity"`
	Notes            string  `json:"notes"`
}

// ===== QC DTOs =====

// CreateInspectionRequest is the request for creating a QC inspection
type CreateInspectionRequest struct {
	InspectionType    string                        `json:"inspection_type" binding:"required"`
	CheckpointID      *uuid.UUID                    `json:"checkpoint_id"`
	ReferenceType     string                        `json:"reference_type" binding:"required"`
	ReferenceID       uuid.UUID                     `json:"reference_id" binding:"required"`
	ProductID         *uuid.UUID                    `json:"product_id"`
	MaterialID        *uuid.UUID                    `json:"material_id"`
	LotID             *uuid.UUID                    `json:"lot_id"`
	LotNumber         string                        `json:"lot_number"`
	InspectedQuantity float64                       `json:"inspected_quantity" binding:"required"`
	SampleSize        *int                          `json:"sample_size"`
	InspectorName     string                        `json:"inspector_name"`
	Items             []CreateInspectionItemRequest `json:"items"`
}

// CreateInspectionItemRequest is the request for a QC inspection item
type CreateInspectionItemRequest struct {
	ItemNumber    int    `json:"item_number"`
	TestName      string `json:"test_name" binding:"required"`
	TestMethod    string `json:"test_method"`
	Specification string `json:"specification"`
	TargetValue   string `json:"target_value"`
	MinValue      string `json:"min_value"`
	MaxValue      string `json:"max_value"`
	ActualValue   string `json:"actual_value"`
	UOM           string `json:"uom"`
	Result        string `json:"result" binding:"required"`
	Notes         string `json:"notes"`
}

// ApproveInspectionRequest is the request for approving an inspection
type ApproveInspectionRequest struct {
	Result           string   `json:"result" binding:"required"`
	AcceptedQuantity *float64 `json:"accepted_quantity"`
	RejectedQuantity *float64 `json:"rejected_quantity"`
	Notes            string   `json:"notes"`
}

// ===== NCR DTOs =====

// CreateNCRRequest is the request for creating an NCR
type CreateNCRRequest struct {
	NCType           string     `json:"nc_type" binding:"required"`
	Severity         string     `json:"severity"`
	ReferenceType    string     `json:"reference_type"`
	ReferenceID      *uuid.UUID `json:"reference_id"`
	ProductID        *uuid.UUID `json:"product_id"`
	MaterialID       *uuid.UUID `json:"material_id"`
	LotID            *uuid.UUID `json:"lot_id"`
	LotNumber        string     `json:"lot_number"`
	Description      string     `json:"description" binding:"required"`
	QuantityAffected *float64   `json:"quantity_affected"`
	UOMID            *uuid.UUID `json:"uom_id"`
	ImmediateAction  string     `json:"immediate_action"`
}

// CloseNCRRequest is the request for closing an NCR
type CloseNCRRequest struct {
	RootCause        string  `json:"root_cause"`
	CorrectiveAction string  `json:"corrective_action"`
	PreventiveAction string  `json:"preventive_action"`
	Disposition      string  `json:"disposition"`
	DispositionQty   *float64 `json:"disposition_quantity"`
	ClosureNotes     string  `json:"closure_notes"`
}
