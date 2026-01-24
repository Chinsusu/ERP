package repository

import (
	"context"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/google/uuid"
)

// SalesOrderFilter defines filter options for sales orders
type SalesOrderFilter struct {
	CustomerID    *uuid.UUID
	Status        entity.SOStatus
	PaymentStatus entity.PaymentStatus
	DateFrom      string
	DateTo        string
	Page          int
	Limit         int
}

// SalesOrderRepository defines sales order repository interface
type SalesOrderRepository interface {
	// CRUD
	Create(ctx context.Context, order *entity.SalesOrder) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SalesOrder, error)
	GetByNumber(ctx context.Context, number string) (*entity.SalesOrder, error)
	Update(ctx context.Context, order *entity.SalesOrder) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *SalesOrderFilter) ([]*entity.SalesOrder, int64, error)

	// Number generation
	GetNextSONumber(ctx context.Context) (string, error)

	// Line items
	CreateLineItem(ctx context.Context, item *entity.SOLineItem) error
	UpdateLineItem(ctx context.Context, item *entity.SOLineItem) error
	DeleteLineItem(ctx context.Context, itemID uuid.UUID) error
	GetLineItems(ctx context.Context, orderID uuid.UUID) ([]*entity.SOLineItem, error)

	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SOStatus) error
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error

	// Shipped quantity
	UpdateShippedQuantity(ctx context.Context, lineItemID uuid.UUID, shippedQty float64) error
	UpdateLineItemReservation(ctx context.Context, lineItemID uuid.UUID, reservationID uuid.UUID) error

	// By customer
	GetByCustomer(ctx context.Context, customerID uuid.UUID, limit int) ([]*entity.SalesOrder, error)
	GetPendingOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]*entity.SalesOrder, error)
}
