package shipment

import (
	"context"
	"errors"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

var (
	ErrShipmentNotFound     = errors.New("shipment not found")
	ErrShipmentCannotShip   = errors.New("shipment cannot be shipped")
	ErrShipmentCannotDeliver = errors.New("shipment cannot be marked as delivered")
)

// CreateShipmentInput represents input for creating shipment
type CreateShipmentInput struct {
	SalesOrderID          uuid.UUID
	Carrier               string
	TrackingNumber        string
	ShippingMethod        string
	ShippingCost          float64
	RecipientName         string
	RecipientPhone        string
	DeliveryAddress       string
	Notes                 string
	CreatedBy             *uuid.UUID
}

// CreateShipmentUseCase handles shipment creation
type CreateShipmentUseCase struct {
	shipmentRepo repository.ShipmentRepository
	orderRepo    repository.SalesOrderRepository
	eventPub     *event.Publisher
}

// NewCreateShipmentUseCase creates a new use case
func NewCreateShipmentUseCase(
	shipmentRepo repository.ShipmentRepository,
	orderRepo repository.SalesOrderRepository,
	eventPub *event.Publisher,
) *CreateShipmentUseCase {
	return &CreateShipmentUseCase{
		shipmentRepo: shipmentRepo,
		orderRepo:    orderRepo,
		eventPub:     eventPub,
	}
}

// Execute creates a new shipment
func (uc *CreateShipmentUseCase) Execute(ctx context.Context, input *CreateShipmentInput) (*entity.Shipment, error) {
	// Verify order exists and can be shipped
	order, err := uc.orderRepo.GetByID(ctx, input.SalesOrderID)
	if err != nil {
		return nil, err
	}

	if !order.CanBeShipped() {
		return nil, ErrShipmentCannotShip
	}

	// Generate shipment number
	number, err := uc.shipmentRepo.GetNextShipmentNumber(ctx)
	if err != nil {
		return nil, err
	}

	shipment := &entity.Shipment{
		ShipmentNumber:  number,
		SalesOrderID:    input.SalesOrderID,
		Carrier:         input.Carrier,
		TrackingNumber:  input.TrackingNumber,
		ShippingMethod:  input.ShippingMethod,
		ShippingCost:    input.ShippingCost,
		RecipientName:   input.RecipientName,
		RecipientPhone:  input.RecipientPhone,
		DeliveryAddress: input.DeliveryAddress,
		Status:          entity.ShipmentStatusPending,
		Notes:           input.Notes,
		CreatedBy:       input.CreatedBy,
	}

	if err := uc.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

// GetShipmentUseCase handles getting shipment
type GetShipmentUseCase struct {
	shipmentRepo repository.ShipmentRepository
}

// NewGetShipmentUseCase creates a new use case
func NewGetShipmentUseCase(repo repository.ShipmentRepository) *GetShipmentUseCase {
	return &GetShipmentUseCase{shipmentRepo: repo}
}

// Execute gets a shipment by ID
func (uc *GetShipmentUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Shipment, error) {
	return uc.shipmentRepo.GetByID(ctx, id)
}

// ListShipmentsUseCase handles listing shipments
type ListShipmentsUseCase struct {
	shipmentRepo repository.ShipmentRepository
}

// NewListShipmentsUseCase creates a new use case
func NewListShipmentsUseCase(repo repository.ShipmentRepository) *ListShipmentsUseCase {
	return &ListShipmentsUseCase{shipmentRepo: repo}
}

// Execute lists shipments with filters
func (uc *ListShipmentsUseCase) Execute(ctx context.Context, filter *repository.ShipmentFilter) ([]*entity.Shipment, int64, error) {
	return uc.shipmentRepo.List(ctx, filter)
}

// ShipShipmentUseCase handles shipping a shipment
type ShipShipmentUseCase struct {
	shipmentRepo repository.ShipmentRepository
	orderRepo    repository.SalesOrderRepository
	eventPub     *event.Publisher
}

// NewShipShipmentUseCase creates a new use case
func NewShipShipmentUseCase(
	shipmentRepo repository.ShipmentRepository,
	orderRepo repository.SalesOrderRepository,
	eventPub *event.Publisher,
) *ShipShipmentUseCase {
	return &ShipShipmentUseCase{
		shipmentRepo: shipmentRepo,
		orderRepo:    orderRepo,
		eventPub:     eventPub,
	}
}

// ShipInput represents input for shipping
type ShipInput struct {
	ShipmentID     uuid.UUID
	Carrier        string
	TrackingNumber string
}

// Execute ships a shipment
func (uc *ShipShipmentUseCase) Execute(ctx context.Context, input *ShipInput) (*entity.Shipment, error) {
	shipment, err := uc.shipmentRepo.GetByID(ctx, input.ShipmentID)
	if err != nil {
		return nil, err
	}

	if !shipment.CanBeShipped() {
		return nil, ErrShipmentCannotShip
	}

	shipment.Ship(input.Carrier, input.TrackingNumber)

	if err := uc.shipmentRepo.Update(ctx, shipment); err != nil {
		return nil, err
	}

	// Update order status
	order, _ := uc.orderRepo.GetByID(ctx, shipment.SalesOrderID)
	if order != nil {
		order.MarkShipped()
		uc.orderRepo.Update(ctx, order)

		// Publish event
		if uc.eventPub != nil {
			uc.eventPub.PublishOrderShipped(&event.OrderShippedEvent{
				SOID:           order.ID.String(),
				SONumber:       order.SONumber,
				CustomerID:     order.CustomerID.String(),
				ShipmentID:     shipment.ID.String(),
				ShipmentNumber: shipment.ShipmentNumber,
				Carrier:        shipment.Carrier,
				TrackingNumber: shipment.TrackingNumber,
			})
		}
	}

	return shipment, nil
}

// DeliverShipmentUseCase handles marking shipment as delivered
type DeliverShipmentUseCase struct {
	shipmentRepo repository.ShipmentRepository
	orderRepo    repository.SalesOrderRepository
	eventPub     *event.Publisher
}

// NewDeliverShipmentUseCase creates a new use case
func NewDeliverShipmentUseCase(
	shipmentRepo repository.ShipmentRepository,
	orderRepo repository.SalesOrderRepository,
	eventPub *event.Publisher,
) *DeliverShipmentUseCase {
	return &DeliverShipmentUseCase{
		shipmentRepo: shipmentRepo,
		orderRepo:    orderRepo,
		eventPub:     eventPub,
	}
}

// Execute marks shipment as delivered
func (uc *DeliverShipmentUseCase) Execute(ctx context.Context, shipmentID uuid.UUID) (*entity.Shipment, error) {
	shipment, err := uc.shipmentRepo.GetByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	if !shipment.CanBeDelivered() {
		return nil, ErrShipmentCannotDeliver
	}

	shipment.MarkDelivered()

	if err := uc.shipmentRepo.Update(ctx, shipment); err != nil {
		return nil, err
	}

	// Update order status
	order, _ := uc.orderRepo.GetByID(ctx, shipment.SalesOrderID)
	if order != nil {
		order.MarkDelivered()
		uc.orderRepo.Update(ctx, order)

		// Publish event
		if uc.eventPub != nil {
			uc.eventPub.PublishOrderDelivered(&event.OrderDeliveredEvent{
				SOID:        order.ID.String(),
				SONumber:    order.SONumber,
				CustomerID:  order.CustomerID.String(),
				DeliveredAt: shipment.ActualDeliveryDate.Format("2006-01-02"),
			})
		}
	}

	return shipment, nil
}
