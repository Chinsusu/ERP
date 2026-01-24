package supplier

import (
	"context"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/erp-cosmetics/supplier-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// UpdateSupplierRequest represents the request to update a supplier
type UpdateSupplierRequest struct {
	Name         string  `json:"name" validate:"omitempty,min=2,max=255"`
	LegalName    string  `json:"legal_name"`
	TaxCode      string  `json:"tax_code"`
	Email        string  `json:"email" validate:"omitempty,email"`
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

// UpdateSupplierUseCase handles updating a supplier
type UpdateSupplierUseCase struct {
	supplierRepo repository.SupplierRepository
}

// NewUpdateSupplierUseCase creates a new UpdateSupplierUseCase
func NewUpdateSupplierUseCase(supplierRepo repository.SupplierRepository) *UpdateSupplierUseCase {
	return &UpdateSupplierUseCase{supplierRepo: supplierRepo}
}

// Execute updates a supplier
func (uc *UpdateSupplierUseCase) Execute(ctx context.Context, id uuid.UUID, req *UpdateSupplierRequest) (*entity.Supplier, error) {
	supplier, err := uc.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.LegalName != "" {
		supplier.LegalName = req.LegalName
	}
	if req.TaxCode != "" {
		supplier.TaxCode = req.TaxCode
	}
	if req.Email != "" {
		supplier.Email = req.Email
	}
	if req.Phone != "" {
		supplier.Phone = req.Phone
	}
	if req.Fax != "" {
		supplier.Fax = req.Fax
	}
	if req.Website != "" {
		supplier.Website = req.Website
	}
	if req.PaymentTerms != "" {
		supplier.PaymentTerms = req.PaymentTerms
	}
	if req.Currency != "" {
		supplier.Currency = req.Currency
	}
	if req.CreditLimit > 0 {
		supplier.CreditLimit = req.CreditLimit
	}
	if req.BankName != "" {
		supplier.BankName = req.BankName
	}
	if req.BankAccount != "" {
		supplier.BankAccount = req.BankAccount
	}
	if req.BankBranch != "" {
		supplier.BankBranch = req.BankBranch
	}
	if req.Notes != "" {
		supplier.Notes = req.Notes
	}
	supplier.UpdatedAt = time.Now()

	if err := uc.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// ApproveSupplierUseCase handles approving a supplier
type ApproveSupplierUseCase struct {
	supplierRepo repository.SupplierRepository
	eventPub     *event.Publisher
}

// NewApproveSupplierUseCase creates a new ApproveSupplierUseCase
func NewApproveSupplierUseCase(
	supplierRepo repository.SupplierRepository,
	eventPub *event.Publisher,
) *ApproveSupplierUseCase {
	return &ApproveSupplierUseCase{
		supplierRepo: supplierRepo,
		eventPub:     eventPub,
	}
}

// Execute approves a supplier
func (uc *ApproveSupplierUseCase) Execute(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID, notes string) (*entity.Supplier, error) {
	supplier, err := uc.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	supplier.Approve(approvedBy)
	if notes != "" {
		supplier.Notes = notes
	}
	supplier.UpdatedAt = time.Now()

	if err := uc.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishSupplierApproved(ctx, &event.SupplierEvent{
		SupplierID:   supplier.ID.String(),
		SupplierCode: supplier.Code,
		Name:         supplier.Name,
		Status:       string(supplier.Status),
		ActionBy:     approvedBy.String(),
		ActionAt:     time.Now(),
	})

	return supplier, nil
}

// BlockSupplierUseCase handles blocking a supplier
type BlockSupplierUseCase struct {
	supplierRepo repository.SupplierRepository
	eventPub     *event.Publisher
}

// NewBlockSupplierUseCase creates a new BlockSupplierUseCase
func NewBlockSupplierUseCase(
	supplierRepo repository.SupplierRepository,
	eventPub *event.Publisher,
) *BlockSupplierUseCase {
	return &BlockSupplierUseCase{
		supplierRepo: supplierRepo,
		eventPub:     eventPub,
	}
}

// Execute blocks a supplier
func (uc *BlockSupplierUseCase) Execute(ctx context.Context, id uuid.UUID, blockedBy uuid.UUID, reason string) (*entity.Supplier, error) {
	supplier, err := uc.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	supplier.Block(blockedBy, reason)
	supplier.UpdatedAt = time.Now()

	if err := uc.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishSupplierBlocked(ctx, &event.SupplierEvent{
		SupplierID:   supplier.ID.String(),
		SupplierCode: supplier.Code,
		Name:         supplier.Name,
		Status:       string(supplier.Status),
		Reason:       reason,
		ActionBy:     blockedBy.String(),
		ActionAt:     time.Now(),
	})

	return supplier, nil
}
