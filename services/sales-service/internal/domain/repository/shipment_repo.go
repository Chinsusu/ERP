package repository

import (
	"context"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/google/uuid"
)

// ShipmentFilter defines filter options for shipments
type ShipmentFilter struct {
	SalesOrderID *uuid.UUID
	Status       entity.ShipmentStatus
	DateFrom     string
	DateTo       string
	Carrier      string
	Page         int
	Limit        int
}

// ShipmentRepository defines shipment repository interface
type ShipmentRepository interface {
	// CRUD
	Create(ctx context.Context, shipment *entity.Shipment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Shipment, error)
	GetByNumber(ctx context.Context, number string) (*entity.Shipment, error)
	Update(ctx context.Context, shipment *entity.Shipment) error
	List(ctx context.Context, filter *ShipmentFilter) ([]*entity.Shipment, int64, error)

	// Number generation
	GetNextShipmentNumber(ctx context.Context) (string, error)

	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ShipmentStatus) error

	// By sales order
	GetBySalesOrder(ctx context.Context, salesOrderID uuid.UUID) ([]*entity.Shipment, error)

	// By tracking number
	GetByTrackingNumber(ctx context.Context, trackingNumber string) (*entity.Shipment, error)
}

// ReturnFilter defines filter options for returns
type ReturnFilter struct {
	SalesOrderID *uuid.UUID
	Status       entity.ReturnStatus
	ReturnType   entity.ReturnType
	DateFrom     string
	DateTo       string
	Page         int
	Limit        int
}

// ReturnRepository defines return repository interface
type ReturnRepository interface {
	// CRUD
	Create(ctx context.Context, ret *entity.Return) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Return, error)
	GetByNumber(ctx context.Context, number string) (*entity.Return, error)
	Update(ctx context.Context, ret *entity.Return) error
	List(ctx context.Context, filter *ReturnFilter) ([]*entity.Return, int64, error)

	// Number generation
	GetNextReturnNumber(ctx context.Context) (string, error)

	// Line items
	CreateLineItem(ctx context.Context, item *entity.ReturnLineItem) error
	GetLineItems(ctx context.Context, returnID uuid.UUID) ([]*entity.ReturnLineItem, error)

	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ReturnStatus) error

	// By sales order
	GetBySalesOrder(ctx context.Context, salesOrderID uuid.UUID) ([]*entity.Return, error)
}
