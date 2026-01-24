package repository

import (
	"context"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CustomerFilter defines filter options for customers
type CustomerFilter struct {
	Search          string
	CustomerType    entity.CustomerType
	CustomerGroupID *uuid.UUID
	Status          entity.CustomerStatus
	Page            int
	Limit           int
}

// CustomerRepository defines customer repository interface
type CustomerRepository interface {
	// Customer CRUD
	Create(ctx context.Context, customer *entity.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
	GetByCode(ctx context.Context, code string) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *CustomerFilter) ([]*entity.Customer, int64, error)

	// Customer code generation
	GetNextCustomerCode(ctx context.Context) (string, error)

	// Credit management
	UpdateBalance(ctx context.Context, customerID uuid.UUID, amount float64) error
	GetAvailableCredit(ctx context.Context, customerID uuid.UUID) (float64, error)

	// Addresses
	CreateAddress(ctx context.Context, address *entity.CustomerAddress) error
	GetAddresses(ctx context.Context, customerID uuid.UUID) ([]*entity.CustomerAddress, error)
	UpdateAddress(ctx context.Context, address *entity.CustomerAddress) error
	DeleteAddress(ctx context.Context, addressID uuid.UUID) error
	GetDefaultAddress(ctx context.Context, customerID uuid.UUID, addressType entity.AddressType) (*entity.CustomerAddress, error)

	// Contacts
	CreateContact(ctx context.Context, contact *entity.CustomerContact) error
	GetContacts(ctx context.Context, customerID uuid.UUID) ([]*entity.CustomerContact, error)
	UpdateContact(ctx context.Context, contact *entity.CustomerContact) error
	DeleteContact(ctx context.Context, contactID uuid.UUID) error
	GetPrimaryContact(ctx context.Context, customerID uuid.UUID) (*entity.CustomerContact, error)
}

// CustomerGroupRepository defines customer group repository interface
type CustomerGroupRepository interface {
	Create(ctx context.Context, group *entity.CustomerGroup) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CustomerGroup, error)
	GetByCode(ctx context.Context, code string) (*entity.CustomerGroup, error)
	Update(ctx context.Context, group *entity.CustomerGroup) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, activeOnly bool) ([]*entity.CustomerGroup, error)
}
