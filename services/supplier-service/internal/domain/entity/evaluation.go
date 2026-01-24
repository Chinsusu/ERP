package entity

import (
	"time"

	"github.com/google/uuid"
)

// EvaluationStatus represents the status of an evaluation
type EvaluationStatus string

const (
	EvaluationStatusDraft     EvaluationStatus = "DRAFT"
	EvaluationStatusSubmitted EvaluationStatus = "SUBMITTED"
	EvaluationStatusApproved  EvaluationStatus = "APPROVED"
)

// Evaluation represents a supplier performance evaluation
type Evaluation struct {
	ID                   uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SupplierID           uuid.UUID        `json:"supplier_id" gorm:"type:uuid;not null"`
	EvaluationDate       time.Time        `json:"evaluation_date" gorm:"type:date;not null"`
	EvaluationPeriod     string           `json:"evaluation_period" gorm:"type:varchar(20);not null"`
	QualityScore         float64          `json:"quality_score" gorm:"type:decimal(3,2);not null"`
	DeliveryScore        float64          `json:"delivery_score" gorm:"type:decimal(3,2);not null"`
	PriceScore           float64          `json:"price_score" gorm:"type:decimal(3,2);not null"`
	ServiceScore         float64          `json:"service_score" gorm:"type:decimal(3,2);not null"`
	DocumentationScore   float64          `json:"documentation_score" gorm:"type:decimal(3,2);not null"`
	OverallScore         float64          `json:"overall_score" gorm:"type:decimal(3,2);not null"`
	OnTimeDeliveryRate   float64          `json:"on_time_delivery_rate" gorm:"type:decimal(5,2)"`
	QualityAcceptanceRate float64         `json:"quality_acceptance_rate" gorm:"type:decimal(5,2)"`
	LeadTimeAdherence    float64          `json:"lead_time_adherence" gorm:"type:decimal(5,2)"`
	Strengths            string           `json:"strengths" gorm:"type:text"`
	Weaknesses           string           `json:"weaknesses" gorm:"type:text"`
	ActionItems          string           `json:"action_items" gorm:"type:text"`
	EvaluatedBy          uuid.UUID        `json:"evaluated_by" gorm:"type:uuid;not null"`
	Status               EvaluationStatus `json:"status" gorm:"type:varchar(20);not null;default:'DRAFT'"`
	ApprovedBy           *uuid.UUID       `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt           *time.Time       `json:"approved_at"`
	CreatedAt            time.Time        `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time        `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for GORM
func (Evaluation) TableName() string {
	return "supplier_evaluations"
}

// CalculateOverallScore calculates the overall score from component scores
func (e *Evaluation) CalculateOverallScore() {
	e.OverallScore = (e.QualityScore + e.DeliveryScore + e.PriceScore + e.ServiceScore + e.DocumentationScore) / 5
}

// Submit changes status to submitted
func (e *Evaluation) Submit() {
	e.Status = EvaluationStatusSubmitted
}

// Approve changes status to approved
func (e *Evaluation) Approve(approvedBy uuid.UUID) {
	now := time.Now()
	e.Status = EvaluationStatusApproved
	e.ApprovedBy = &approvedBy
	e.ApprovedAt = &now
}
