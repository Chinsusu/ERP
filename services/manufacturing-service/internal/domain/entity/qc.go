package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CheckpointType represents QC checkpoint type
type CheckpointType string

const (
	CheckpointTypeIQC  CheckpointType = "IQC"  // Incoming QC
	CheckpointTypeIPQC CheckpointType = "IPQC" // In-Process QC
	CheckpointTypeFQC  CheckpointType = "FQC"  // Final QC
)

// InspectionResult represents QC result
type InspectionResult string

const (
	InspectionResultPending     InspectionResult = "PENDING"
	InspectionResultPassed      InspectionResult = "PASSED"
	InspectionResultFailed      InspectionResult = "FAILED"
	InspectionResultConditional InspectionResult = "CONDITIONAL"
)

// ItemResult represents individual test result
type ItemResult string

const (
	ItemResultPass ItemResult = "PASS"
	ItemResultFail ItemResult = "FAIL"
	ItemResultNA   ItemResult = "N/A"
)

// ReferenceType represents what is being inspected
type ReferenceType string

const (
	ReferenceTypeWorkOrder ReferenceType = "WORK_ORDER"
	ReferenceTypeGRN       ReferenceType = "GRN"
	ReferenceTypeLot       ReferenceType = "LOT"
)

// QCCheckpoint represents a QC checkpoint template
type QCCheckpoint struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code           string         `json:"code" gorm:"type:varchar(20);unique;not null"`
	Name           string         `json:"name" gorm:"type:varchar(100);not null"`
	Description    string         `json:"description" gorm:"type:text"`
	CheckpointType CheckpointType `json:"checkpoint_type" gorm:"type:varchar(20);not null"`
	AppliesTo      string         `json:"applies_to" gorm:"type:varchar(20);default:'ALL'"` // ALL, MATERIAL, PRODUCT
	TestItems      json.RawMessage `json:"test_items" gorm:"type:jsonb;not null"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (QCCheckpoint) TableName() string {
	return "qc_checkpoints"
}

// TestItemTemplate represents a test item in a checkpoint
type TestItemTemplate struct {
	Name          string   `json:"name"`
	Method        string   `json:"method"`
	Specification string   `json:"specification"`
	Type          string   `json:"type"` // PASS_FAIL, NUMERIC
	Min           *float64 `json:"min,omitempty"`
	Max           *float64 `json:"max,omitempty"`
	Unit          string   `json:"unit,omitempty"`
}

// QCInspection represents an actual QC inspection
type QCInspection struct {
	ID                uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InspectionNumber  string           `json:"inspection_number" gorm:"type:varchar(30);unique;not null"`
	InspectionDate    time.Time        `json:"inspection_date" gorm:"not null"`
	InspectionType    CheckpointType   `json:"inspection_type" gorm:"type:varchar(20);not null"`
	CheckpointID      *uuid.UUID       `json:"checkpoint_id" gorm:"type:uuid"`
	ReferenceType     ReferenceType    `json:"reference_type" gorm:"type:varchar(30);not null"`
	ReferenceID       uuid.UUID        `json:"reference_id" gorm:"type:uuid;not null"`
	ProductID         *uuid.UUID       `json:"product_id" gorm:"type:uuid"`
	MaterialID        *uuid.UUID       `json:"material_id" gorm:"type:uuid"`
	LotID             *uuid.UUID       `json:"lot_id" gorm:"type:uuid"`
	LotNumber         string           `json:"lot_number" gorm:"type:varchar(50)"`
	InspectedQuantity float64          `json:"inspected_quantity" gorm:"type:decimal(15,4);not null"`
	AcceptedQuantity  *float64         `json:"accepted_quantity" gorm:"type:decimal(15,4)"`
	RejectedQuantity  *float64         `json:"rejected_quantity" gorm:"type:decimal(15,4)"`
	SampleSize        *int             `json:"sample_size"`
	Result            InspectionResult `json:"result" gorm:"type:varchar(20);default:'PENDING'"`
	OverallScore      *float64         `json:"overall_score" gorm:"type:decimal(5,2)"`
	InspectorID       uuid.UUID        `json:"inspector_id" gorm:"type:uuid;not null"`
	InspectorName     string           `json:"inspector_name" gorm:"type:varchar(100)"`
	ApprovedBy        *uuid.UUID       `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt        *time.Time       `json:"approved_at"`
	TestResults       json.RawMessage  `json:"test_results" gorm:"type:jsonb"`
	Notes             string           `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time        `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Associations
	Items []QCInspectionItem `json:"items,omitempty" gorm:"foreignKey:InspectionID"`
}

// TableName returns the table name
func (QCInspection) TableName() string {
	return "qc_inspections"
}

// QCInspectionItem represents an individual test in an inspection
type QCInspectionItem struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InspectionID  uuid.UUID  `json:"inspection_id" gorm:"type:uuid;not null"`
	ItemNumber    int        `json:"item_number" gorm:"not null"`
	TestName      string     `json:"test_name" gorm:"type:varchar(100);not null"`
	TestMethod    string     `json:"test_method" gorm:"type:varchar(100)"`
	Specification string     `json:"specification" gorm:"type:varchar(200)"`
	TargetValue   string     `json:"target_value" gorm:"type:varchar(100)"`
	MinValue      string     `json:"min_value" gorm:"type:varchar(100)"`
	MaxValue      string     `json:"max_value" gorm:"type:varchar(100)"`
	ActualValue   string     `json:"actual_value" gorm:"type:varchar(100)"`
	UOM           string     `json:"uom" gorm:"type:varchar(20)"`
	Result        ItemResult `json:"result" gorm:"type:varchar(20);not null"`
	Notes         string     `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (QCInspectionItem) TableName() string {
	return "qc_inspection_items"
}

// QC business methods

// IsPending returns true if inspection is pending
func (q *QCInspection) IsPending() bool {
	return q.Result == InspectionResultPending
}

// CanBeApproved returns true if inspection can be approved
func (q *QCInspection) CanBeApproved() bool {
	return q.Result == InspectionResultPending
}

// Pass marks the inspection as passed
func (q *QCInspection) Pass(approverID uuid.UUID) {
	q.Result = InspectionResultPassed
	q.ApprovedBy = &approverID
	now := time.Now()
	q.ApprovedAt = &now
	q.UpdatedAt = now
}

// Fail marks the inspection as failed
func (q *QCInspection) Fail(approverID uuid.UUID) {
	q.Result = InspectionResultFailed
	q.ApprovedBy = &approverID
	now := time.Now()
	q.ApprovedAt = &now
	q.UpdatedAt = now
}

// ConditionalPass marks the inspection as conditionally passed
func (q *QCInspection) ConditionalPass(approverID uuid.UUID) {
	q.Result = InspectionResultConditional
	q.ApprovedBy = &approverID
	now := time.Now()
	q.ApprovedAt = &now
	q.UpdatedAt = now
}

// CalculateScore calculates the overall score based on items
func (q *QCInspection) CalculateScore() {
	if len(q.Items) == 0 {
		return
	}
	
	passCount := 0
	totalCount := 0
	
	for _, item := range q.Items {
		if item.Result != ItemResultNA {
			totalCount++
			if item.Result == ItemResultPass {
				passCount++
			}
		}
	}
	
	if totalCount > 0 {
		score := float64(passCount) / float64(totalCount) * 100
		q.OverallScore = &score
	}
}
