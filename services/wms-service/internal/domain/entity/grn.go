package entity

import (
	"time"

	"github.com/google/uuid"
)

// GRNStatus represents GRN status
type GRNStatus string

const (
	GRNStatusDraft      GRNStatus = "DRAFT"
	GRNStatusInProgress GRNStatus = "IN_PROGRESS"
	GRNStatusCompleted  GRNStatus = "COMPLETED"
	GRNStatusCancelled  GRNStatus = "CANCELLED"
)

// GRN represents a Goods Receipt Note
type GRN struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GRNNumber          string     `json:"grn_number" gorm:"type:varchar(30);unique;not null"` // GRN-YYYY-XXXX
	GRNDate            time.Time  `json:"grn_date" gorm:"type:date;not null"`
	POID               *uuid.UUID `json:"po_id" gorm:"type:uuid"`
	PONumber           string     `json:"po_number" gorm:"type:varchar(30)"`
	SupplierID         *uuid.UUID `json:"supplier_id" gorm:"type:uuid"`
	WarehouseID        uuid.UUID  `json:"warehouse_id" gorm:"type:uuid;not null"`
	DeliveryNoteNumber string     `json:"delivery_note_number" gorm:"type:varchar(50)"`
	VehicleNumber      string     `json:"vehicle_number" gorm:"type:varchar(20)"`
	Status             GRNStatus  `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	QCStatus           QCStatus   `json:"qc_status" gorm:"type:varchar(20);default:'PENDING'"`
	QCNotes            string     `json:"qc_notes" gorm:"type:text"`
	Notes              string     `json:"notes" gorm:"type:text"`
	ReceivedBy         *uuid.UUID `json:"received_by" gorm:"type:uuid"`
	CompletedAt        *time.Time `json:"completed_at"`
	CreatedAt          time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Warehouse *Warehouse    `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	LineItems []GRNLineItem `json:"line_items,omitempty" gorm:"foreignKey:GRNID"`
}

// TableName returns the table name
func (GRN) TableName() string {
	return "grns"
}

// CanComplete returns true if GRN can be completed
func (g *GRN) CanComplete() bool {
	return g.Status == GRNStatusDraft || g.Status == GRNStatusInProgress
}

// Complete completes the GRN
func (g *GRN) Complete(qcStatus QCStatus, qcNotes string) {
	now := time.Now()
	g.Status = GRNStatusCompleted
	g.QCStatus = qcStatus
	g.QCNotes = qcNotes
	g.CompletedAt = &now
	g.UpdatedAt = now
}

// Cancel cancels the GRN
func (g *GRN) Cancel() {
	g.Status = GRNStatusCancelled
	g.UpdatedAt = time.Now()
}

// StartProcessing starts processing the GRN
func (g *GRN) StartProcessing() {
	g.Status = GRNStatusInProgress
	g.UpdatedAt = time.Now()
}

// GRNLineItem represents a line item in GRN
type GRNLineItem struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GRNID             uuid.UUID  `json:"grn_id" gorm:"type:uuid;not null"`
	LineNumber        int        `json:"line_number" gorm:"not null"`
	POLineItemID      *uuid.UUID `json:"po_line_item_id" gorm:"type:uuid"`
	MaterialID        uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	ExpectedQty       *float64   `json:"expected_qty" gorm:"type:decimal(15,4)"`
	ReceivedQty       float64    `json:"received_qty" gorm:"type:decimal(15,4);not null"`
	AcceptedQty       *float64   `json:"accepted_qty" gorm:"type:decimal(15,4)"`
	RejectedQty       float64    `json:"rejected_qty" gorm:"type:decimal(15,4);default:0"`
	UnitID            uuid.UUID  `json:"unit_id" gorm:"type:uuid;not null"`
	LotID             *uuid.UUID `json:"lot_id" gorm:"type:uuid"`
	SupplierLotNumber string     `json:"supplier_lot_number" gorm:"type:varchar(50)"`
	ManufacturedDate  *time.Time `json:"manufactured_date" gorm:"type:date"`
	ExpiryDate        time.Time  `json:"expiry_date" gorm:"type:date;not null"`
	LocationID        *uuid.UUID `json:"location_id" gorm:"type:uuid"`
	QCStatus          QCStatus   `json:"qc_status" gorm:"type:varchar(20);default:'PENDING'"`
	QCNotes           string     `json:"qc_notes" gorm:"type:text"`
	CreatedAt         time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Lot      *Lot      `json:"lot,omitempty" gorm:"foreignKey:LotID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName returns the table name
func (GRNLineItem) TableName() string {
	return "grn_line_items"
}

// PassQC passes QC for the line item
func (i *GRNLineItem) PassQC(acceptedQty float64) {
	i.QCStatus = QCStatusPassed
	i.AcceptedQty = &acceptedQty
	i.RejectedQty = i.ReceivedQty - acceptedQty
}

// FailQC fails QC for the line item
func (i *GRNLineItem) FailQC(notes string) {
	i.QCStatus = QCStatusFailed
	rejected := i.ReceivedQty
	i.RejectedQty = rejected
	zero := 0.0
	i.AcceptedQty = &zero
	i.QCNotes = notes
}
