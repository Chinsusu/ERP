package dto

import (
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CreateSupplierRequest DTO
type CreateSupplierRequest struct {
	Name         string  `json:"name" binding:"required,min=2,max=255"`
	LegalName    string  `json:"legal_name"`
	TaxCode      string  `json:"tax_code"`
	SupplierType string  `json:"supplier_type" binding:"required,oneof=MANUFACTURER TRADER IMPORTER"`
	BusinessType string  `json:"business_type" binding:"omitempty,oneof=DOMESTIC INTERNATIONAL"`
	Email        string  `json:"email" binding:"omitempty,email"`
	Phone        string  `json:"phone"`
	Website      string  `json:"website"`
	PaymentTerms string  `json:"payment_terms"`
	Currency     string  `json:"currency"`
	CreditLimit  float64 `json:"credit_limit"`
	BankName     string  `json:"bank_name"`
	BankAccount  string  `json:"bank_account"`
	BankBranch   string  `json:"bank_branch"`
	Notes        string  `json:"notes"`
}

// UpdateSupplierRequest DTO
type UpdateSupplierRequest struct {
	Name         string  `json:"name" binding:"omitempty,min=2,max=255"`
	LegalName    string  `json:"legal_name"`
	TaxCode      string  `json:"tax_code"`
	Email        string  `json:"email" binding:"omitempty,email"`
	Phone        string  `json:"phone"`
	Fax          string  `json:"fax"`
	Website      string  `json:"website"`
	PaymentTerms string  `json:"payment_terms"`
	Currency     string  `json:"currency"`
	CreditLimit  float64 `json:"credit_limit"`
	BankName     string  `json:"bank_name"`
	BankAccount  string  `json:"bank_account"`
	BankBranch   string  `json:"bank_branch"`
	Notes        string  `json:"notes"`
}

// ApproveRequest DTO
type ApproveRequest struct {
	Notes string `json:"notes"`
}

// BlockRequest DTO
type BlockRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// SupplierResponse DTO
type SupplierResponse struct {
	ID             uuid.UUID             `json:"id"`
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	LegalName      string                `json:"legal_name,omitempty"`
	TaxCode        string                `json:"tax_code,omitempty"`
	SupplierType   string                `json:"supplier_type"`
	BusinessType   string                `json:"business_type"`
	Email          string                `json:"email,omitempty"`
	Phone          string                `json:"phone,omitempty"`
	Website        string                `json:"website,omitempty"`
	PaymentTerms   string                `json:"payment_terms"`
	Currency       string                `json:"currency"`
	CreditLimit    float64               `json:"credit_limit"`
	OverallRating  float64               `json:"overall_rating"`
	Status         string                `json:"status"`
	BlockedReason  string                `json:"blocked_reason,omitempty"`
	HasValidGMP    bool                  `json:"has_valid_gmp"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Addresses      []AddressResponse     `json:"addresses,omitempty"`
	Contacts       []ContactResponse     `json:"contacts,omitempty"`
	Certifications []CertificationResponse `json:"certifications,omitempty"`
}

// ToSupplierResponse converts entity to response DTO
func ToSupplierResponse(s *entity.Supplier) *SupplierResponse {
	resp := &SupplierResponse{
		ID:            s.ID,
		Code:          s.Code,
		Name:          s.Name,
		LegalName:     s.LegalName,
		TaxCode:       s.TaxCode,
		SupplierType:  string(s.SupplierType),
		BusinessType:  string(s.BusinessType),
		Email:         s.Email,
		Phone:         s.Phone,
		Website:       s.Website,
		PaymentTerms:  s.PaymentTerms,
		Currency:      s.Currency,
		CreditLimit:   s.CreditLimit,
		OverallRating: s.OverallRating,
		Status:        string(s.Status),
		BlockedReason: s.BlockedReason,
		HasValidGMP:   s.HasValidGMP(),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}

	// Convert addresses
	for _, a := range s.Addresses {
		resp.Addresses = append(resp.Addresses, *ToAddressResponse(&a))
	}

	// Convert contacts
	for _, c := range s.Contacts {
		resp.Contacts = append(resp.Contacts, *ToContactResponse(&c))
	}

	// Convert certifications
	for _, c := range s.Certifications {
		resp.Certifications = append(resp.Certifications, *ToCertificationResponse(&c))
	}

	return resp
}

// SupplierListResponse for list endpoint
type SupplierListResponse struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	SupplierType  string    `json:"supplier_type"`
	BusinessType  string    `json:"business_type"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Status        string    `json:"status"`
	OverallRating float64   `json:"overall_rating"`
	HasValidGMP   bool      `json:"has_valid_gmp"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToSupplierListResponse converts entity to list response DTO
func ToSupplierListResponse(s *entity.Supplier) *SupplierListResponse {
	return &SupplierListResponse{
		ID:            s.ID,
		Code:          s.Code,
		Name:          s.Name,
		SupplierType:  string(s.SupplierType),
		BusinessType:  string(s.BusinessType),
		Email:         s.Email,
		Phone:         s.Phone,
		Status:        string(s.Status),
		OverallRating: s.OverallRating,
		HasValidGMP:   s.HasValidGMP(),
		CreatedAt:     s.CreatedAt,
	}
}
