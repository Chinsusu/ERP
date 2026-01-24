package entity

import (
	"time"

	"github.com/google/uuid"
)

// NCType represents non-conformance type
type NCType string

const (
	NCTypeMaterial  NCType = "MATERIAL"
	NCTypeProcess   NCType = "PROCESS"
	NCTypeProduct   NCType = "PRODUCT"
	NCTypeEquipment NCType = "EQUIPMENT"
)

// NCRSeverity represents NCR severity level
type NCRSeverity string

const (
	NCRSeverityLow      NCRSeverity = "LOW"
	NCRSeverityMedium   NCRSeverity = "MEDIUM"
	NCRSeverityHigh     NCRSeverity = "HIGH"
	NCRSeverityCritical NCRSeverity = "CRITICAL"
)

// NCRStatus represents NCR status
type NCRStatus string

const (
	NCRStatusOpen             NCRStatus = "OPEN"
	NCRStatusInvestigation    NCRStatus = "INVESTIGATION"
	NCRStatusCorrectiveAction NCRStatus = "CORRECTIVE_ACTION"
	NCRStatusClosed           NCRStatus = "CLOSED"
)

// Disposition represents NCR disposition
type Disposition string

const (
	DispositionUseAsIs          Disposition = "USE_AS_IS"
	DispositionRework           Disposition = "REWORK"
	DispositionScrap            Disposition = "SCRAP"
	DispositionReturnToSupplier Disposition = "RETURN_TO_SUPPLIER"
)

// NCR represents a Non-Conformance Report
type NCR struct {
	ID                  uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NCRNumber           string      `json:"ncr_number" gorm:"type:varchar(30);unique;not null"`
	NCRDate             time.Time   `json:"ncr_date" gorm:"type:date;not null"`
	NCType              NCType      `json:"nc_type" gorm:"type:varchar(30);not null"`
	Severity            NCRSeverity `json:"severity" gorm:"type:varchar(20);default:'MEDIUM'"`
	Status              NCRStatus   `json:"status" gorm:"type:varchar(30);default:'OPEN'"`
	ReferenceType       string      `json:"reference_type" gorm:"type:varchar(30)"`
	ReferenceID         *uuid.UUID  `json:"reference_id" gorm:"type:uuid"`
	ProductID           *uuid.UUID  `json:"product_id" gorm:"type:uuid"`
	MaterialID          *uuid.UUID  `json:"material_id" gorm:"type:uuid"`
	LotID               *uuid.UUID  `json:"lot_id" gorm:"type:uuid"`
	LotNumber           string      `json:"lot_number" gorm:"type:varchar(50)"`
	Description         string      `json:"description" gorm:"type:text;not null"`
	QuantityAffected    *float64    `json:"quantity_affected" gorm:"type:decimal(15,4)"`
	UOMID               *uuid.UUID  `json:"uom_id" gorm:"type:uuid"`
	RootCause           string      `json:"root_cause" gorm:"type:text"`
	InvestigationDate   *time.Time  `json:"investigation_date"`
	InvestigatedBy      *uuid.UUID  `json:"investigated_by" gorm:"type:uuid"`
	ImmediateAction     string      `json:"immediate_action" gorm:"type:text"`
	CorrectiveAction    string      `json:"corrective_action" gorm:"type:text"`
	PreventiveAction    string      `json:"preventive_action" gorm:"type:text"`
	Disposition         *Disposition `json:"disposition" gorm:"type:varchar(30)"`
	DispositionQuantity *float64    `json:"disposition_quantity" gorm:"type:decimal(15,4)"`
	DispositionDate     *time.Time  `json:"disposition_date"`
	DispositionBy       *uuid.UUID  `json:"disposition_by" gorm:"type:uuid"`
	ClosedAt            *time.Time  `json:"closed_at"`
	ClosedBy            *uuid.UUID  `json:"closed_by" gorm:"type:uuid"`
	ClosureNotes        string      `json:"closure_notes" gorm:"type:text"`
	CreatedBy           *uuid.UUID  `json:"created_by" gorm:"type:uuid"`
	UpdatedBy           *uuid.UUID  `json:"updated_by" gorm:"type:uuid"`
	CreatedAt           time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (NCR) TableName() string {
	return "ncrs"
}

// NCR business methods

// IsOpen returns true if NCR is open
func (n *NCR) IsOpen() bool {
	return n.Status == NCRStatusOpen
}

// CanBeClosed returns true if NCR can be closed
func (n *NCR) CanBeClosed() bool {
	return n.Status != NCRStatusClosed
}

// StartInvestigation moves NCR to investigation phase
func (n *NCR) StartInvestigation(investigatorID uuid.UUID) {
	n.Status = NCRStatusInvestigation
	n.InvestigatedBy = &investigatorID
	now := time.Now()
	n.InvestigationDate = &now
	n.UpdatedAt = now
}

// RecordRootCause records the root cause
func (n *NCR) RecordRootCause(rootCause string) {
	n.RootCause = rootCause
	n.Status = NCRStatusCorrectiveAction
	n.UpdatedAt = time.Now()
}

// SetDisposition sets the disposition for the NCR
func (n *NCR) SetDisposition(disposition Disposition, quantity float64, disposedBy uuid.UUID) {
	n.Disposition = &disposition
	n.DispositionQuantity = &quantity
	n.DispositionBy = &disposedBy
	now := time.Now()
	n.DispositionDate = &now
	n.UpdatedAt = now
}

// Close closes the NCR
func (n *NCR) Close(closedBy uuid.UUID, notes string) error {
	if !n.CanBeClosed() {
		return nil
	}
	n.Status = NCRStatusClosed
	n.ClosedBy = &closedBy
	n.ClosureNotes = notes
	now := time.Now()
	n.ClosedAt = &now
	n.UpdatedAt = now
	return nil
}
