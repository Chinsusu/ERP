package supplier

import (
	"context"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/erp-cosmetics/supplier-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateSupplierRequest represents the request to create a supplier
type CreateSupplierRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=255"`
	LegalName    string  `json:"legal_name"`
	TaxCode      string  `json:"tax_code"`
	SupplierType string  `json:"supplier_type" validate:"required,oneof=MANUFACTURER TRADER IMPORTER"`
	BusinessType string  `json:"business_type" validate:"required,oneof=DOMESTIC INTERNATIONAL"`
	Email        string  `json:"email" validate:"omitempty,email"`
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

// CreateSupplierUseCase handles creating a new supplier
type CreateSupplierUseCase struct {
	supplierRepo repository.SupplierRepository
	eventPub     *event.Publisher
}

// NewCreateSupplierUseCase creates a new CreateSupplierUseCase
func NewCreateSupplierUseCase(
	supplierRepo repository.SupplierRepository,
	eventPub *event.Publisher,
) *CreateSupplierUseCase {
	return &CreateSupplierUseCase{
		supplierRepo: supplierRepo,
		eventPub:     eventPub,
	}
}

// Execute creates a new supplier
func (uc *CreateSupplierUseCase) Execute(ctx context.Context, req *CreateSupplierRequest) (*entity.Supplier, error) {
	// Generate code
	code, err := uc.supplierRepo.GetNextCode(ctx)
	if err != nil {
		return nil, err
	}

	// Set defaults
	currency := req.Currency
	if currency == "" {
		currency = "VND"
	}
	paymentTerms := req.PaymentTerms
	if paymentTerms == "" {
		paymentTerms = "Net 30"
	}
	businessType := req.BusinessType
	if businessType == "" {
		businessType = "DOMESTIC"
	}

	supplier := &entity.Supplier{
		ID:           uuid.New(),
		Code:         code,
		Name:         req.Name,
		LegalName:    req.LegalName,
		TaxCode:      req.TaxCode,
		SupplierType: entity.SupplierType(req.SupplierType),
		BusinessType: entity.BusinessType(businessType),
		Email:        req.Email,
		Phone:        req.Phone,
		Website:      req.Website,
		PaymentTerms: paymentTerms,
		Currency:     currency,
		CreditLimit:  req.CreditLimit,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankBranch:   req.BankBranch,
		Notes:        req.Notes,
		Status:       entity.SupplierStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.supplierRepo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishSupplierCreated(ctx, &event.SupplierEvent{
		SupplierID:   supplier.ID.String(),
		SupplierCode: supplier.Code,
		Name:         supplier.Name,
		ActionAt:     time.Now(),
	})

	return supplier, nil
}

// GetSupplierUseCase handles getting a supplier by ID
type GetSupplierUseCase struct {
	supplierRepo repository.SupplierRepository
}

// NewGetSupplierUseCase creates a new GetSupplierUseCase
func NewGetSupplierUseCase(supplierRepo repository.SupplierRepository) *GetSupplierUseCase {
	return &GetSupplierUseCase{supplierRepo: supplierRepo}
}

// Execute gets a supplier by ID
func (uc *GetSupplierUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Supplier, error) {
	return uc.supplierRepo.GetByID(ctx, id)
}

// ListSuppliersUseCase handles listing suppliers
type ListSuppliersUseCase struct {
	supplierRepo repository.SupplierRepository
}

// NewListSuppliersUseCase creates a new ListSuppliersUseCase
func NewListSuppliersUseCase(supplierRepo repository.SupplierRepository) *ListSuppliersUseCase {
	return &ListSuppliersUseCase{supplierRepo: supplierRepo}
}

// Execute lists suppliers with filters
func (uc *ListSuppliersUseCase) Execute(ctx context.Context, filter *repository.SupplierFilter) ([]*entity.Supplier, int64, error) {
	return uc.supplierRepo.List(ctx, filter)
}
