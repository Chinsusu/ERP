package entity

import (
	"time"

	"github.com/google/uuid"
)

// SupplierType represents the type of supplier
type SupplierType string

const (
	SupplierTypeManufacturer SupplierType = "MANUFACTURER"
	SupplierTypeTrader       SupplierType = "TRADER"
	SupplierTypeImporter     SupplierType = "IMPORTER"
)

// BusinessType represents domestic or international supplier
type BusinessType string

const (
	BusinessTypeDomestic      BusinessType = "DOMESTIC"
	BusinessTypeInternational BusinessType = "INTERNATIONAL"
)

// SupplierStatus represents the status of a supplier
type SupplierStatus string

const (
	SupplierStatusPending  SupplierStatus = "PENDING"
	SupplierStatusApproved SupplierStatus = "APPROVED"
	SupplierStatusBlocked  SupplierStatus = "BLOCKED"
	SupplierStatusInactive SupplierStatus = "INACTIVE"
)

// Supplier represents the supplier entity
type Supplier struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code            string         `json:"code" gorm:"type:varchar(20);unique;not null"`
	Name            string         `json:"name" gorm:"type:varchar(255);not null"`
	LegalName       string         `json:"legal_name" gorm:"type:varchar(255)"`
	TaxCode         string         `json:"tax_code" gorm:"type:varchar(50)"`
	SupplierType    SupplierType   `json:"supplier_type" gorm:"type:varchar(20);not null"`
	BusinessType    BusinessType   `json:"business_type" gorm:"type:varchar(20);not null;default:'DOMESTIC'"`
	Email           string         `json:"email" gorm:"type:varchar(255)"`
	Phone           string         `json:"phone" gorm:"type:varchar(50)"`
	Fax             string         `json:"fax" gorm:"type:varchar(50)"`
	Website         string         `json:"website" gorm:"type:varchar(255)"`
	PaymentTerms    string         `json:"payment_terms" gorm:"type:varchar(50);default:'Net 30'"`
	Currency        string         `json:"currency" gorm:"type:varchar(3);default:'VND'"`
	CreditLimit     float64        `json:"credit_limit" gorm:"type:decimal(15,2);default:0"`
	BankName        string         `json:"bank_name" gorm:"type:varchar(255)"`
	BankAccount     string         `json:"bank_account" gorm:"type:varchar(100)"`
	BankBranch      string         `json:"bank_branch" gorm:"type:varchar(255)"`
	QualityRating   float64        `json:"quality_rating" gorm:"type:decimal(3,2);default:0"`
	DeliveryRating  float64        `json:"delivery_rating" gorm:"type:decimal(3,2);default:0"`
	ServiceRating   float64        `json:"service_rating" gorm:"type:decimal(3,2);default:0"`
	OverallRating   float64        `json:"overall_rating" gorm:"type:decimal(3,2);default:0"`
	Status          SupplierStatus `json:"status" gorm:"type:varchar(20);not null;default:'PENDING'"`
	BlockedReason   string         `json:"blocked_reason" gorm:"type:text"`
	BlockedBy       *uuid.UUID     `json:"blocked_by" gorm:"type:uuid"`
	BlockedAt       *time.Time     `json:"blocked_at"`
	ApprovedBy      *uuid.UUID     `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt      *time.Time     `json:"approved_at"`
	Notes           string         `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       *time.Time     `json:"deleted_at" gorm:"index"`

	// Relations
	Addresses      []Address       `json:"addresses,omitempty" gorm:"foreignKey:SupplierID"`
	Contacts       []Contact       `json:"contacts,omitempty" gorm:"foreignKey:SupplierID"`
	Certifications []Certification `json:"certifications,omitempty" gorm:"foreignKey:SupplierID"`
	Evaluations    []Evaluation    `json:"evaluations,omitempty" gorm:"foreignKey:SupplierID"`
}

// TableName returns the table name for GORM
func (Supplier) TableName() string {
	return "suppliers"
}

// IsActive returns true if supplier is approved
func (s *Supplier) IsActive() bool {
	return s.Status == SupplierStatusApproved
}

// CanCreatePO returns true if supplier can receive purchase orders
func (s *Supplier) CanCreatePO() bool {
	return s.Status == SupplierStatusApproved
}

// UpdateRating recalculates overall rating from component ratings
func (s *Supplier) UpdateRating() {
	if s.QualityRating > 0 || s.DeliveryRating > 0 || s.ServiceRating > 0 {
		count := 0.0
		sum := 0.0
		if s.QualityRating > 0 {
			sum += s.QualityRating
			count++
		}
		if s.DeliveryRating > 0 {
			sum += s.DeliveryRating
			count++
		}
		if s.ServiceRating > 0 {
			sum += s.ServiceRating
			count++
		}
		if count > 0 {
			s.OverallRating = sum / count
		}
	}
}

// Approve marks the supplier as approved
func (s *Supplier) Approve(approvedBy uuid.UUID) {
	now := time.Now()
	s.Status = SupplierStatusApproved
	s.ApprovedBy = &approvedBy
	s.ApprovedAt = &now
}

// Block marks the supplier as blocked
func (s *Supplier) Block(blockedBy uuid.UUID, reason string) {
	now := time.Now()
	s.Status = SupplierStatusBlocked
	s.BlockedBy = &blockedBy
	s.BlockedAt = &now
	s.BlockedReason = reason
}

// HasValidGMP returns true if supplier has a valid GMP certificate
func (s *Supplier) HasValidGMP() bool {
	for _, cert := range s.Certifications {
		if cert.Type == CertTypeGMP && cert.Status == CertStatusValid {
			return true
		}
	}
	return false
}

// GetValidCertTypes returns list of valid certification types
func (s *Supplier) GetValidCertTypes() []CertificationType {
	var types []CertificationType
	for _, cert := range s.Certifications {
		if cert.Status == CertStatusValid {
			types = append(types, cert.Type)
		}
	}
	return types
}
