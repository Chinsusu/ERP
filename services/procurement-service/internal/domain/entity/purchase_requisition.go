package entity

import (
	"time"

	"github.com/google/uuid"
)

// PRPriority represents the priority of a PR
type PRPriority string

const (
	PRPriorityLow    PRPriority = "LOW"
	PRPriorityNormal PRPriority = "NORMAL"
	PRPriorityHigh   PRPriority = "HIGH"
	PRPriorityUrgent PRPriority = "URGENT"
)

// PRStatus represents the status of a PR
type PRStatus string

const (
	PRStatusDraft           PRStatus = "DRAFT"
	PRStatusSubmitted       PRStatus = "SUBMITTED"
	PRStatusPendingApproval PRStatus = "PENDING_APPROVAL"
	PRStatusApproved        PRStatus = "APPROVED"
	PRStatusRejected        PRStatus = "REJECTED"
	PRStatusConvertedToPO   PRStatus = "CONVERTED_TO_PO"
	PRStatusCancelled       PRStatus = "CANCELLED"
)

// ApprovalLevel represents the approval level based on amount
type ApprovalLevel string

const (
	ApprovalLevelAuto               ApprovalLevel = "AUTO"
	ApprovalLevelDeptManager        ApprovalLevel = "DEPARTMENT_MANAGER"
	ApprovalLevelProcurementManager ApprovalLevel = "PROCUREMENT_MANAGER"
	ApprovalLevelCFO                ApprovalLevel = "CFO"
)

// PurchaseRequisition represents a PR entity
type PurchaseRequisition struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PRNumber        string        `json:"pr_number" gorm:"column:pr_number;type:varchar(20);unique;not null"`
	PRDate          time.Time     `json:"pr_date" gorm:"type:date;not null;default:CURRENT_DATE"`
	RequiredDate    time.Time     `json:"required_date" gorm:"type:date;not null"`
	Priority        PRPriority    `json:"priority" gorm:"type:varchar(20);not null;default:'NORMAL'"`
	Status          PRStatus      `json:"status" gorm:"type:varchar(30);not null;default:'DRAFT'"`
	RequesterID     uuid.UUID     `json:"requester_id" gorm:"type:uuid;not null"`
	DepartmentID    *uuid.UUID    `json:"department_id" gorm:"type:uuid"`
	Justification   string        `json:"justification" gorm:"type:text"`
	TotalAmount     float64       `json:"total_amount" gorm:"type:decimal(15,2);not null;default:0"`
	Currency        string        `json:"currency" gorm:"type:varchar(3);not null;default:'VND'"`
	ApprovalLevel   ApprovalLevel `json:"approval_level" gorm:"type:varchar(30)"`
	Notes           string        `json:"notes" gorm:"type:text"`
	SubmittedAt     *time.Time    `json:"submitted_at"`
	ApprovedAt      *time.Time    `json:"approved_at"`
	ApprovedBy      *uuid.UUID    `json:"approved_by" gorm:"type:uuid"`
	RejectedAt      *time.Time    `json:"rejected_at"`
	RejectedBy      *uuid.UUID    `json:"rejected_by" gorm:"type:uuid"`
	RejectionReason string        `json:"rejection_reason" gorm:"type:text"`
	POID            *uuid.UUID    `json:"po_id" gorm:"column:po_id;type:uuid"`
	CreatedAt       time.Time     `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time     `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       *time.Time    `json:"deleted_at" gorm:"index"`

	// Relations
	LineItems []PRLineItem  `json:"line_items,omitempty" gorm:"foreignKey:PRID"`
	Approvals []PRApproval  `json:"approvals,omitempty" gorm:"foreignKey:PRID"`
}

// TableName returns the table name
func (PurchaseRequisition) TableName() string {
	return "purchase_requisitions"
}

// CalculateTotalAmount calculates total from line items
func (pr *PurchaseRequisition) CalculateTotalAmount() {
	total := 0.0
	for _, item := range pr.LineItems {
		total += item.LineTotal
	}
	pr.TotalAmount = total
}

// DetermineApprovalLevel determines approval level based on amount
func (pr *PurchaseRequisition) DetermineApprovalLevel() ApprovalLevel {
	amount := pr.TotalAmount
	switch {
	case amount < 10000000: // < 10M VND
		return ApprovalLevelAuto
	case amount < 50000000: // 10M - 50M
		return ApprovalLevelDeptManager
	case amount < 200000000: // 50M - 200M
		return ApprovalLevelProcurementManager
	default: // > 200M
		return ApprovalLevelCFO
	}
}

// Submit submits PR for approval
func (pr *PurchaseRequisition) Submit() {
	now := time.Now()
	pr.Status = PRStatusSubmitted
	pr.SubmittedAt = &now
	pr.ApprovalLevel = pr.DetermineApprovalLevel()
	
	// Auto-approve if under threshold
	if pr.ApprovalLevel == ApprovalLevelAuto {
		pr.Status = PRStatusApproved
		pr.ApprovedAt = &now
	} else {
		pr.Status = PRStatusPendingApproval
	}
}

// Approve approves the PR
func (pr *PurchaseRequisition) Approve(approverID uuid.UUID) {
	now := time.Now()
	pr.Status = PRStatusApproved
	pr.ApprovedAt = &now
	pr.ApprovedBy = &approverID
}

// Reject rejects the PR
func (pr *PurchaseRequisition) Reject(rejectedBy uuid.UUID, reason string) {
	now := time.Now()
	pr.Status = PRStatusRejected
	pr.RejectedAt = &now
	pr.RejectedBy = &rejectedBy
	pr.RejectionReason = reason
}

// MarkConvertedToPO marks PR as converted to PO
func (pr *PurchaseRequisition) MarkConvertedToPO(poID uuid.UUID) {
	pr.Status = PRStatusConvertedToPO
	pr.POID = &poID
}

// CanSubmit returns true if PR can be submitted
func (pr *PurchaseRequisition) CanSubmit() bool {
	return pr.Status == PRStatusDraft && len(pr.LineItems) > 0
}

// CanApprove returns true if PR can be approved
func (pr *PurchaseRequisition) CanApprove() bool {
	return pr.Status == PRStatusPendingApproval
}

// CanConvert returns true if PR can be converted to PO
func (pr *PurchaseRequisition) CanConvert() bool {
	return pr.Status == PRStatusApproved
}
