package customer

import (
	"context"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateCustomerInput represents input for creating customer
type CreateCustomerInput struct {
	Name            string
	TaxCode         string
	CustomerType    entity.CustomerType
	CustomerGroupID *uuid.UUID
	Email           string
	Phone           string
	Website         string
	PaymentTerms    string
	CreditLimit     float64
	Currency        string
	Notes           string
	CreatedBy       *uuid.UUID
}

// CreateCustomerUseCase handles customer creation
type CreateCustomerUseCase struct {
	customerRepo repository.CustomerRepository
	eventPub     *event.Publisher
}

// NewCreateCustomerUseCase creates a new use case
func NewCreateCustomerUseCase(repo repository.CustomerRepository, eventPub *event.Publisher) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		customerRepo: repo,
		eventPub:     eventPub,
	}
}

// Execute creates a new customer
func (uc *CreateCustomerUseCase) Execute(ctx context.Context, input *CreateCustomerInput) (*entity.Customer, error) {
	// Generate customer code
	code, err := uc.customerRepo.GetNextCustomerCode(ctx)
	if err != nil {
		return nil, err
	}

	customer := &entity.Customer{
		CustomerCode:    code,
		Name:            input.Name,
		TaxCode:         input.TaxCode,
		CustomerType:    input.CustomerType,
		CustomerGroupID: input.CustomerGroupID,
		Email:           input.Email,
		Phone:           input.Phone,
		Website:         input.Website,
		PaymentTerms:    input.PaymentTerms,
		CreditLimit:     input.CreditLimit,
		Currency:        input.Currency,
		Notes:           input.Notes,
		Status:          entity.CustomerStatusActive,
		CreatedBy:       input.CreatedBy,
	}

	if customer.Currency == "" {
		customer.Currency = "VND"
	}
	if customer.PaymentTerms == "" {
		customer.PaymentTerms = "Net 30"
	}
	if customer.CustomerType == "" {
		customer.CustomerType = entity.CustomerTypeRetail
	}

	if err := uc.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	// Publish event
	if uc.eventPub != nil {
		uc.eventPub.PublishCustomerCreated(&event.CustomerCreatedEvent{
			CustomerID:   customer.ID.String(),
			CustomerCode: customer.CustomerCode,
			Name:         customer.Name,
			CustomerType: string(customer.CustomerType),
			Email:        customer.Email,
		})
	}

	return customer, nil
}

// GetCustomerUseCase handles getting customer
type GetCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// NewGetCustomerUseCase creates a new use case
func NewGetCustomerUseCase(repo repository.CustomerRepository) *GetCustomerUseCase {
	return &GetCustomerUseCase{customerRepo: repo}
}

// Execute gets a customer by ID
func (uc *GetCustomerUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	return uc.customerRepo.GetByID(ctx, id)
}

// ListCustomersUseCase handles listing customers
type ListCustomersUseCase struct {
	customerRepo repository.CustomerRepository
}

// NewListCustomersUseCase creates a new use case
func NewListCustomersUseCase(repo repository.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{customerRepo: repo}
}

// Execute lists customers with filters
func (uc *ListCustomersUseCase) Execute(ctx context.Context, filter *repository.CustomerFilter) ([]*entity.Customer, int64, error) {
	return uc.customerRepo.List(ctx, filter)
}

// UpdateCustomerInput represents input for updating customer
type UpdateCustomerInput struct {
	Name            string
	TaxCode         string
	CustomerType    entity.CustomerType
	CustomerGroupID *uuid.UUID
	Email           string
	Phone           string
	Website         string
	PaymentTerms    string
	CreditLimit     float64
	Currency        string
	Status          entity.CustomerStatus
	Notes           string
	UpdatedBy       *uuid.UUID
}

// UpdateCustomerUseCase handles updating customer
type UpdateCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// NewUpdateCustomerUseCase creates a new use case
func NewUpdateCustomerUseCase(repo repository.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{customerRepo: repo}
}

// Execute updates a customer
func (uc *UpdateCustomerUseCase) Execute(ctx context.Context, id uuid.UUID, input *UpdateCustomerInput) (*entity.Customer, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	customer.Name = input.Name
	customer.TaxCode = input.TaxCode
	customer.CustomerType = input.CustomerType
	customer.CustomerGroupID = input.CustomerGroupID
	customer.Email = input.Email
	customer.Phone = input.Phone
	customer.Website = input.Website
	customer.PaymentTerms = input.PaymentTerms
	customer.CreditLimit = input.CreditLimit
	customer.Currency = input.Currency
	customer.Status = input.Status
	customer.Notes = input.Notes
	customer.UpdatedBy = input.UpdatedBy

	if err := uc.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// DeleteCustomerUseCase handles deleting customer
type DeleteCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// NewDeleteCustomerUseCase creates a new use case
func NewDeleteCustomerUseCase(repo repository.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{customerRepo: repo}
}

// Execute deletes a customer
func (uc *DeleteCustomerUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.customerRepo.Delete(ctx, id)
}

// CheckCreditUseCase handles credit limit checking
type CheckCreditUseCase struct {
	customerRepo repository.CustomerRepository
}

// NewCheckCreditUseCase creates a new use case
func NewCheckCreditUseCase(repo repository.CustomerRepository) *CheckCreditUseCase {
	return &CheckCreditUseCase{customerRepo: repo}
}

// CreditCheckResult represents credit check result
type CreditCheckResult struct {
	CustomerID        uuid.UUID
	CreditLimit       float64
	CurrentBalance    float64
	AvailableCredit   float64
	RequestedAmount   float64
	WithinLimit       bool
}

// Execute checks if customer can place order of given amount
func (uc *CheckCreditUseCase) Execute(ctx context.Context, customerID uuid.UUID, orderAmount float64) (*CreditCheckResult, error) {
	customer, err := uc.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	availableCredit := customer.CreditLimit - customer.CurrentBalance

	return &CreditCheckResult{
		CustomerID:      customerID,
		CreditLimit:     customer.CreditLimit,
		CurrentBalance:  customer.CurrentBalance,
		AvailableCredit: availableCredit,
		RequestedAmount: orderAmount,
		WithinLimit:     availableCredit >= orderAmount,
	}, nil
}
