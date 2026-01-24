package sales_order

import (
	"context"
	"errors"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

var (
	ErrOrderNotFound       = errors.New("sales order not found")
	ErrOrderCannotConfirm  = errors.New("order cannot be confirmed")
	ErrOrderCannotCancel   = errors.New("order cannot be cancelled")
	ErrOrderCannotShip     = errors.New("order cannot be shipped")
	ErrInsufficientCredit  = errors.New("insufficient credit limit")
)

// CreateOrderInput represents input for creating sales order
type CreateOrderInput struct {
	CustomerID      uuid.UUID
	QuotationID     *uuid.UUID
	SODate          time.Time
	DeliveryDate    *time.Time
	DeliveryAddress string
	BillingAddress  string
	DiscountPercent float64
	TaxPercent      float64
	PaymentMethod   entity.PaymentMethod
	Notes           string
	Items           []OrderItemInput
	CreatedBy       *uuid.UUID
}

// OrderItemInput represents a line item input
type OrderItemInput struct {
	ProductID       uuid.UUID
	ProductCode     string
	ProductName     string
	Quantity        float64
	UomID           *uuid.UUID
	UnitPrice       float64
	DiscountPercent float64
	TaxPercent      float64
	Notes           string
}

// CreateOrderUseCase handles sales order creation
type CreateOrderUseCase struct {
	orderRepo    repository.SalesOrderRepository
	customerRepo repository.CustomerRepository
	eventPub     *event.Publisher
}

// NewCreateOrderUseCase creates a new use case
func NewCreateOrderUseCase(
	orderRepo repository.SalesOrderRepository,
	customerRepo repository.CustomerRepository,
	eventPub *event.Publisher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		eventPub:     eventPub,
	}
}

// Execute creates a new sales order
func (uc *CreateOrderUseCase) Execute(ctx context.Context, input *CreateOrderInput) (*entity.SalesOrder, error) {
	// Verify customer exists
	_, err := uc.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return nil, err
	}

	// Generate SO number
	soNumber, err := uc.orderRepo.GetNextSONumber(ctx)
	if err != nil {
		return nil, err
	}

	order := &entity.SalesOrder{
		SONumber:        soNumber,
		CustomerID:      input.CustomerID,
		QuotationID:     input.QuotationID,
		SODate:          input.SODate,
		DeliveryDate:    input.DeliveryDate,
		DeliveryAddress: input.DeliveryAddress,
		BillingAddress:  input.BillingAddress,
		DiscountPercent: input.DiscountPercent,
		TaxPercent:      input.TaxPercent,
		PaymentMethod:   input.PaymentMethod,
		PaymentStatus:   entity.PaymentStatusPending,
		Status:          entity.SOStatusDraft,
		Notes:           input.Notes,
		CreatedBy:       input.CreatedBy,
	}

	// Create line items
	for i, item := range input.Items {
		lineItem := entity.SOLineItem{
			LineNumber:      i + 1,
			ProductID:       item.ProductID,
			ProductCode:     item.ProductCode,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			UomID:           item.UomID,
			UnitPrice:       item.UnitPrice,
			DiscountPercent: item.DiscountPercent,
			TaxPercent:      item.TaxPercent,
			Notes:           item.Notes,
		}
		lineItem.CalculateLineTotal()
		order.LineItems = append(order.LineItems, lineItem)
	}

	// Calculate totals
	order.CalculateTotals()

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Publish event
	if uc.eventPub != nil {
		items := make([]event.OrderLineItem, len(order.LineItems))
		for i, item := range order.LineItems {
			items[i] = event.OrderLineItem{
				ProductID:   item.ProductID.String(),
				ProductCode: item.ProductCode,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			}
		}
		uc.eventPub.PublishOrderCreated(&event.OrderCreatedEvent{
			SOID:        order.ID.String(),
			SONumber:    order.SONumber,
			CustomerID:  order.CustomerID.String(),
			TotalAmount: order.TotalAmount,
			Items:       items,
		})
	}

	return order, nil
}

// GetOrderUseCase handles getting sales order
type GetOrderUseCase struct {
	orderRepo repository.SalesOrderRepository
}

// NewGetOrderUseCase creates a new use case
func NewGetOrderUseCase(repo repository.SalesOrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: repo}
}

// Execute gets a sales order by ID
func (uc *GetOrderUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.SalesOrder, error) {
	return uc.orderRepo.GetByID(ctx, id)
}

// ListOrdersUseCase handles listing sales orders
type ListOrdersUseCase struct {
	orderRepo repository.SalesOrderRepository
}

// NewListOrdersUseCase creates a new use case
func NewListOrdersUseCase(repo repository.SalesOrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{orderRepo: repo}
}

// Execute lists sales orders with filters
func (uc *ListOrdersUseCase) Execute(ctx context.Context, filter *repository.SalesOrderFilter) ([]*entity.SalesOrder, int64, error) {
	return uc.orderRepo.List(ctx, filter)
}

// ConfirmOrderUseCase handles confirming sales order
type ConfirmOrderUseCase struct {
	orderRepo       repository.SalesOrderRepository
	customerRepo    repository.CustomerRepository
	eventPub        *event.Publisher
	enableCreditCheck bool
}

// NewConfirmOrderUseCase creates a new use case
func NewConfirmOrderUseCase(
	orderRepo repository.SalesOrderRepository,
	customerRepo repository.CustomerRepository,
	eventPub *event.Publisher,
	enableCreditCheck bool,
) *ConfirmOrderUseCase {
	return &ConfirmOrderUseCase{
		orderRepo:       orderRepo,
		customerRepo:    customerRepo,
		eventPub:        eventPub,
		enableCreditCheck: enableCreditCheck,
	}
}

// Execute confirms a sales order
func (uc *ConfirmOrderUseCase) Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*entity.SalesOrder, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !order.CanBeConfirmed() {
		return nil, ErrOrderCannotConfirm
	}

	// Check credit limit
	if uc.enableCreditCheck {
		customer, err := uc.customerRepo.GetByID(ctx, order.CustomerID)
		if err != nil {
			return nil, err
		}

		if !customer.CanPlaceOrder(order.TotalAmount) {
			return nil, ErrInsufficientCredit
		}

		// Update customer balance
		if err := uc.customerRepo.UpdateBalance(ctx, customer.ID, order.TotalAmount); err != nil {
			return nil, err
		}
	}

	// Confirm order
	order.Confirm(userID)

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish event -> WMS will reserve stock
	if uc.eventPub != nil {
		deliveryDate := ""
		if order.DeliveryDate != nil {
			deliveryDate = order.DeliveryDate.Format("2006-01-02")
		}
		
		items := make([]event.OrderLineItem, len(order.LineItems))
		for i, item := range order.LineItems {
			items[i] = event.OrderLineItem{
				ProductID:   item.ProductID.String(),
				ProductCode: item.ProductCode,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			}
		}
		
		uc.eventPub.PublishOrderConfirmed(&event.OrderConfirmedEvent{
			SOID:            order.ID.String(),
			SONumber:        order.SONumber,
			CustomerID:      order.CustomerID.String(),
			DeliveryDate:    deliveryDate,
			DeliveryAddress: order.DeliveryAddress,
			TotalAmount:     order.TotalAmount,
			Items:           items,
		})
	}

	return order, nil
}

// CancelOrderUseCase handles cancelling sales order
type CancelOrderUseCase struct {
	orderRepo    repository.SalesOrderRepository
	customerRepo repository.CustomerRepository
	eventPub     *event.Publisher
}

// NewCancelOrderUseCase creates a new use case
func NewCancelOrderUseCase(
	orderRepo repository.SalesOrderRepository,
	customerRepo repository.CustomerRepository,
	eventPub *event.Publisher,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		eventPub:     eventPub,
	}
}

// Execute cancels a sales order
func (uc *CancelOrderUseCase) Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, reason string) (*entity.SalesOrder, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !order.CanBeCancelled() {
		return nil, ErrOrderCannotCancel
	}

	// If order was confirmed, release the credit hold
	if order.Status == entity.SOStatusConfirmed {
		if err := uc.customerRepo.UpdateBalance(ctx, order.CustomerID, -order.TotalAmount); err != nil {
			return nil, err
		}
	}

	order.Cancel(userID, reason)

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish event -> WMS will release reservations
	if uc.eventPub != nil {
		uc.eventPub.PublishOrderCancelled(&event.OrderCancelledEvent{
			SOID:        order.ID.String(),
			SONumber:    order.SONumber,
			CustomerID:  order.CustomerID.String(),
			Reason:      reason,
			CancelledBy: userID.String(),
		})
	}

	return order, nil
}

// ShipOrderUseCase handles shipping sales order
type ShipOrderUseCase struct {
	orderRepo repository.SalesOrderRepository
	eventPub  *event.Publisher
}

// NewShipOrderUseCase creates a new use case
func NewShipOrderUseCase(orderRepo repository.SalesOrderRepository, eventPub *event.Publisher) *ShipOrderUseCase {
	return &ShipOrderUseCase{
		orderRepo: orderRepo,
		eventPub:  eventPub,
	}
}

// Execute marks order as shipped
func (uc *ShipOrderUseCase) Execute(ctx context.Context, orderID uuid.UUID) (*entity.SalesOrder, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !order.CanBeShipped() {
		return nil, ErrOrderCannotShip
	}

	order.MarkShipped()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// DeliverOrderUseCase handles marking order as delivered
type DeliverOrderUseCase struct {
	orderRepo repository.SalesOrderRepository
	eventPub  *event.Publisher
}

// NewDeliverOrderUseCase creates a new use case
func NewDeliverOrderUseCase(orderRepo repository.SalesOrderRepository, eventPub *event.Publisher) *DeliverOrderUseCase {
	return &DeliverOrderUseCase{
		orderRepo: orderRepo,
		eventPub:  eventPub,
	}
}

// Execute marks order as delivered
func (uc *DeliverOrderUseCase) Execute(ctx context.Context, orderID uuid.UUID) (*entity.SalesOrder, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order.MarkDelivered()

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish event
	if uc.eventPub != nil {
		uc.eventPub.PublishOrderDelivered(&event.OrderDeliveredEvent{
			SOID:        order.ID.String(),
			SONumber:    order.SONumber,
			CustomerID:  order.CustomerID.String(),
			DeliveredAt: order.DeliveredAt.Format(time.RFC3339),
		})
	}

	return order, nil
}
