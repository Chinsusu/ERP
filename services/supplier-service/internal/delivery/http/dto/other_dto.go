package dto

import (
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CreateAddressRequest DTO
type CreateAddressRequest struct {
	AddressType  string `json:"address_type" binding:"required,oneof=BILLING SHIPPING FACTORY OFFICE WAREHOUSE"`
	AddressLine1 string `json:"address_line1" binding:"required"`
	AddressLine2 string `json:"address_line2"`
	Ward         string `json:"ward"`
	District     string `json:"district"`
	City         string `json:"city" binding:"required"`
	Province     string `json:"province"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	IsPrimary    bool   `json:"is_primary"`
}

// AddressResponse DTO
type AddressResponse struct {
	ID           uuid.UUID `json:"id"`
	AddressType  string    `json:"address_type"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2,omitempty"`
	Ward         string    `json:"ward,omitempty"`
	District     string    `json:"district,omitempty"`
	City         string    `json:"city"`
	Province     string    `json:"province,omitempty"`
	Country      string    `json:"country"`
	PostalCode   string    `json:"postal_code,omitempty"`
	IsPrimary    bool      `json:"is_primary"`
}

// ToAddressResponse converts entity to response DTO
func ToAddressResponse(a *entity.Address) *AddressResponse {
	return &AddressResponse{
		ID:           a.ID,
		AddressType:  string(a.AddressType),
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		Ward:         a.Ward,
		District:     a.District,
		City:         a.City,
		Province:     a.Province,
		Country:      a.Country,
		PostalCode:   a.PostalCode,
		IsPrimary:    a.IsPrimary,
	}
}

// CreateContactRequest DTO
type CreateContactRequest struct {
	ContactType string `json:"contact_type" binding:"required,oneof=PRIMARY SALES TECHNICAL QUALITY ACCOUNTING LOGISTICS"`
	FullName    string `json:"full_name" binding:"required"`
	Position    string `json:"position"`
	Department  string `json:"department"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	IsPrimary   bool   `json:"is_primary"`
}

// ContactResponse DTO
type ContactResponse struct {
	ID          uuid.UUID `json:"id"`
	ContactType string    `json:"contact_type"`
	FullName    string    `json:"full_name"`
	Position    string    `json:"position,omitempty"`
	Department  string    `json:"department,omitempty"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Mobile      string    `json:"mobile,omitempty"`
	IsPrimary   bool      `json:"is_primary"`
}

// ToContactResponse converts entity to response DTO
func ToContactResponse(c *entity.Contact) *ContactResponse {
	return &ContactResponse{
		ID:          c.ID,
		ContactType: string(c.ContactType),
		FullName:    c.FullName,
		Position:    c.Position,
		Department:  c.Department,
		Email:       c.Email,
		Phone:       c.Phone,
		Mobile:      c.Mobile,
		IsPrimary:   c.IsPrimary,
	}
}

// CreateCertificationRequest DTO
type CreateCertificationRequest struct {
	Type        string `json:"certification_type" binding:"required,oneof=GMP ISO9001 ISO22716 ORGANIC ECOCERT HALAL COSMOS OTHER"`
	CertNumber  string `json:"certificate_number" binding:"required"`
	IssuingBody string `json:"issuing_body" binding:"required"`
	IssueDate   string `json:"issue_date" binding:"required"`
	ExpiryDate  string `json:"expiry_date" binding:"required"`
	DocumentURL string `json:"document_url"`
	Notes       string `json:"notes"`
}

// CertificationResponse DTO
type CertificationResponse struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"certification_type"`
	CertNumber      string    `json:"certificate_number"`
	IssuingBody     string    `json:"issuing_body"`
	IssueDate       time.Time `json:"issue_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
	Status          string    `json:"status"`
	DaysUntilExpiry int       `json:"days_until_expiry"`
	DocumentURL     string    `json:"document_url,omitempty"`
}

// ToCertificationResponse converts entity to response DTO
func ToCertificationResponse(c *entity.Certification) *CertificationResponse {
	c.UpdateStatus()
	return &CertificationResponse{
		ID:              c.ID,
		Type:            string(c.Type),
		CertNumber:      c.CertNumber,
		IssuingBody:     c.IssuingBody,
		IssueDate:       c.IssueDate,
		ExpiryDate:      c.ExpiryDate,
		Status:          string(c.Status),
		DaysUntilExpiry: c.DaysUntilExpiry,
		DocumentURL:     c.DocumentURL,
	}
}

// CreateEvaluationRequest DTO
type CreateEvaluationRequest struct {
	EvaluationDate        string  `json:"evaluation_date" binding:"required"`
	EvaluationPeriod      string  `json:"evaluation_period" binding:"required"`
	QualityScore          float64 `json:"quality_score" binding:"required,min=1,max=5"`
	DeliveryScore         float64 `json:"delivery_score" binding:"required,min=1,max=5"`
	PriceScore            float64 `json:"price_score" binding:"required,min=1,max=5"`
	ServiceScore          float64 `json:"service_score" binding:"required,min=1,max=5"`
	DocumentationScore    float64 `json:"documentation_score" binding:"required,min=1,max=5"`
	OnTimeDeliveryRate    float64 `json:"on_time_delivery_rate"`
	QualityAcceptanceRate float64 `json:"quality_acceptance_rate"`
	LeadTimeAdherence     float64 `json:"lead_time_adherence"`
	Strengths             string  `json:"strengths"`
	Weaknesses            string  `json:"weaknesses"`
	ActionItems           string  `json:"action_items"`
}

// EvaluationResponse DTO
type EvaluationResponse struct {
	ID                    uuid.UUID `json:"id"`
	EvaluationDate        time.Time `json:"evaluation_date"`
	EvaluationPeriod      string    `json:"evaluation_period"`
	QualityScore          float64   `json:"quality_score"`
	DeliveryScore         float64   `json:"delivery_score"`
	PriceScore            float64   `json:"price_score"`
	ServiceScore          float64   `json:"service_score"`
	DocumentationScore    float64   `json:"documentation_score"`
	OverallScore          float64   `json:"overall_score"`
	OnTimeDeliveryRate    float64   `json:"on_time_delivery_rate"`
	QualityAcceptanceRate float64   `json:"quality_acceptance_rate"`
	Strengths             string    `json:"strengths,omitempty"`
	Weaknesses            string    `json:"weaknesses,omitempty"`
	ActionItems           string    `json:"action_items,omitempty"`
	Status                string    `json:"status"`
}

// ToEvaluationResponse converts entity to response DTO
func ToEvaluationResponse(e *entity.Evaluation) *EvaluationResponse {
	return &EvaluationResponse{
		ID:                    e.ID,
		EvaluationDate:        e.EvaluationDate,
		EvaluationPeriod:      e.EvaluationPeriod,
		QualityScore:          e.QualityScore,
		DeliveryScore:         e.DeliveryScore,
		PriceScore:            e.PriceScore,
		ServiceScore:          e.ServiceScore,
		DocumentationScore:    e.DocumentationScore,
		OverallScore:          e.OverallScore,
		OnTimeDeliveryRate:    e.OnTimeDeliveryRate,
		QualityAcceptanceRate: e.QualityAcceptanceRate,
		Strengths:             e.Strengths,
		Weaknesses:            e.Weaknesses,
		ActionItems:           e.ActionItems,
		Status:                string(e.Status),
	}
}
