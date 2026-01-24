package repository

import (
	"context"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/google/uuid"
)

// POFilter represents filters for listing POs
type POFilter struct {
	Status     string
	SupplierID *uuid.UUID
	PRID       *uuid.UUID
	DateFrom   string
	DateTo     string
	Search     string
	Page       int
	Limit      int
}

// PORepository defines the interface for PO data access
type PORepository interface {
	// CRUD
	Create(ctx context.Context, po *entity.PurchaseOrder) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PurchaseOrder, error)
	GetByNumber(ctx context.Context, poNumber string) (*entity.PurchaseOrder, error)
	Update(ctx context.Context, po *entity.PurchaseOrder) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List
	List(ctx context.Context, filter *POFilter) ([]*entity.PurchaseOrder, int64, error)
	GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.PurchaseOrder, error)
	
	// Code generation
	GetNextPONumber(ctx context.Context) (string, error)
	
	// Line items
	CreateLineItem(ctx context.Context, item *entity.POLineItem) error
	UpdateLineItem(ctx context.Context, item *entity.POLineItem) error
	DeleteLineItem(ctx context.Context, id uuid.UUID) error
	GetLineItemsByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.POLineItem, error)
	GetLineItemByID(ctx context.Context, id uuid.UUID) (*entity.POLineItem, error)
	
	// Amendments
	CreateAmendment(ctx context.Context, amendment *entity.POAmendment) error
	GetAmendmentsByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.POAmendment, error)
	
	// Receipts
	CreateReceipt(ctx context.Context, receipt *entity.POReceipt) error
	GetReceiptsByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.POReceipt, error)
	GetReceiptsByLineItemID(ctx context.Context, lineItemID uuid.UUID) ([]*entity.POReceipt, error)
}
