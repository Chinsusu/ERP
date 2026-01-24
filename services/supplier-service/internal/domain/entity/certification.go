package entity

import (
	"time"

	"github.com/google/uuid"
)

// CertificationType represents the type of certification
type CertificationType string

const (
	CertTypeGMP      CertificationType = "GMP"
	CertTypeISO9001  CertificationType = "ISO9001"
	CertTypeISO22716 CertificationType = "ISO22716"
	CertTypeOrganic  CertificationType = "ORGANIC"
	CertTypeEcocert  CertificationType = "ECOCERT"
	CertTypeHalal    CertificationType = "HALAL"
	CertTypeCosmos   CertificationType = "COSMOS"
	CertTypeOther    CertificationType = "OTHER"
)

// CertStatus represents the status of a certification
type CertStatus string

const (
	CertStatusValid       CertStatus = "VALID"
	CertStatusExpiringSoon CertStatus = "EXPIRING_SOON"
	CertStatusExpired     CertStatus = "EXPIRED"
)

// Certification represents a supplier certification (GMP, ISO, Organic, etc.)
type Certification struct {
	ID              uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SupplierID      uuid.UUID         `json:"supplier_id" gorm:"type:uuid;not null"`
	Type            CertificationType `json:"certification_type" gorm:"column:certification_type;type:varchar(20);not null"`
	CertNumber      string            `json:"certificate_number" gorm:"column:certificate_number;type:varchar(100);not null"`
	IssuingBody     string            `json:"issuing_body" gorm:"type:varchar(255);not null"`
	IssueDate       time.Time         `json:"issue_date" gorm:"type:date;not null"`
	ExpiryDate      time.Time         `json:"expiry_date" gorm:"type:date;not null"`
	DocumentURL     string            `json:"document_url" gorm:"type:varchar(500)"`
	Status          CertStatus        `json:"status" gorm:"type:varchar(20);not null;default:'VALID'"`
	VerifiedBy      *uuid.UUID        `json:"verified_by" gorm:"type:uuid"`
	VerifiedAt      *time.Time        `json:"verified_at"`
	Notes           string            `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time         `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time         `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Computed field (not stored in DB)
	DaysUntilExpiry int `json:"days_until_expiry" gorm:"-"`
}

// TableName returns the table name for GORM
func (Certification) TableName() string {
	return "supplier_certifications"
}

// UpdateStatus updates the certification status based on expiry date
func (c *Certification) UpdateStatus() {
	c.DaysUntilExpiry = int(time.Until(c.ExpiryDate).Hours() / 24)
	
	switch {
	case c.DaysUntilExpiry < 0:
		c.Status = CertStatusExpired
	case c.DaysUntilExpiry <= 90:
		c.Status = CertStatusExpiringSoon
	default:
		c.Status = CertStatusValid
	}
}

// IsExpired returns true if the certification is expired
func (c *Certification) IsExpired() bool {
	return c.ExpiryDate.Before(time.Now())
}

// IsExpiringSoon returns true if expiring within N days
func (c *Certification) IsExpiringSoon(days int) bool {
	return time.Until(c.ExpiryDate).Hours()/24 <= float64(days)
}

// CalculateDaysUntilExpiry calculates days until expiry
func (c *Certification) CalculateDaysUntilExpiry() int {
	return int(time.Until(c.ExpiryDate).Hours() / 24)
}
