package repository

import (
	"context"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/google/uuid"
)

// SupplierFilter represents filters for listing suppliers
type SupplierFilter struct {
	SupplierType string
	BusinessType string
	Status       string
	Country      string
	HasGMP       *bool
	MinRating    *float64
	Search       string
	Page         int
	Limit        int
}

// SupplierRepository defines the interface for supplier data access
type SupplierRepository interface {
	// CRUD operations
	Create(ctx context.Context, supplier *entity.Supplier) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Supplier, error)
	GetByCode(ctx context.Context, code string) (*entity.Supplier, error)
	Update(ctx context.Context, supplier *entity.Supplier) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List operations
	List(ctx context.Context, filter *SupplierFilter) ([]*entity.Supplier, int64, error)
	
	// Status operations
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SupplierStatus) error
	
	// Code generation
	GetNextCode(ctx context.Context) (string, error)
	
	// Rating operations
	UpdateRating(ctx context.Context, id uuid.UUID, quality, delivery, service, overall float64) error
}

// AddressRepository defines the interface for address data access
type AddressRepository interface {
	Create(ctx context.Context, address *entity.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Address, error)
	GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Address, error)
	Update(ctx context.Context, address *entity.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, supplierID, addressID uuid.UUID) error
}

// ContactRepository defines the interface for contact data access
type ContactRepository interface {
	Create(ctx context.Context, contact *entity.Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Contact, error)
	GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Contact, error)
	Update(ctx context.Context, contact *entity.Contact) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, supplierID, contactID uuid.UUID) error
}
