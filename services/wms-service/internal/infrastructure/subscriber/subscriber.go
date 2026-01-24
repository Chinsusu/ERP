package subscriber

import (
	"context"
	"encoding/json"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/usecase/grn"
	"github.com/erp-cosmetics/wms-service/internal/usecase/reservation"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// EventSubscriber handles incoming events from other services
type EventSubscriber struct {
	nc               *nats.Conn
	logger           *zap.Logger
	createGRNUC      *grn.CreateGRNUseCase
	reserveStockUC   *reservation.CreateReservationUseCase
	releaseReservationUC *reservation.ReleaseReservationUseCase
	subscriptions    []*nats.Subscription
}

// NewEventSubscriber creates a new event subscriber
func NewEventSubscriber(
	nc *nats.Conn,
	logger *zap.Logger,
	createGRNUC *grn.CreateGRNUseCase,
	reserveStockUC *reservation.CreateReservationUseCase,
	releaseReservationUC *reservation.ReleaseReservationUseCase,
) *EventSubscriber {
	return &EventSubscriber{
		nc:                   nc,
		logger:               logger,
		createGRNUC:          createGRNUC,
		reserveStockUC:       reserveStockUC,
		releaseReservationUC: releaseReservationUC,
	}
}

// Start begins listening for events
func (s *EventSubscriber) Start() error {
	if s.nc == nil {
		s.logger.Warn("NATS not connected, skipping event subscriptions")
		return nil
	}

	// Subscribe to procurement events
	sub1, err := s.nc.Subscribe("procurement.po.received", s.handlePOReceived)
	if err != nil {
		return err
	}
	s.subscriptions = append(s.subscriptions, sub1)

	// Subscribe to sales order events
	sub2, err := s.nc.Subscribe("sales.order.confirmed", s.handleSalesOrderConfirmed)
	if err != nil {
		return err
	}
	s.subscriptions = append(s.subscriptions, sub2)

	// Subscribe to sales order cancellation
	sub3, err := s.nc.Subscribe("sales.order.cancelled", s.handleSalesOrderCancelled)
	if err != nil {
		return err
	}
	s.subscriptions = append(s.subscriptions, sub3)

	// Subscribe to manufacturing events
	sub4, err := s.nc.Subscribe("manufacturing.wo.started", s.handleWorkOrderStarted)
	if err != nil {
		return err
	}
	s.subscriptions = append(s.subscriptions, sub4)

	s.logger.Info("Event subscriber started",
		zap.Int("subscriptions", len(s.subscriptions)),
	)

	return nil
}

// Stop stops all subscriptions
func (s *EventSubscriber) Stop() {
	for _, sub := range s.subscriptions {
		sub.Unsubscribe()
	}
	s.logger.Info("Event subscriber stopped")
}

// POReceivedEvent represents a PO received event from procurement
type POReceivedEvent struct {
	POID            uuid.UUID `json:"po_id"`
	PONumber        string    `json:"po_number"`
	SupplierID      uuid.UUID `json:"supplier_id"`
	SupplierName    string    `json:"supplier_name"`
	WarehouseID     uuid.UUID `json:"warehouse_id"`
	ReceivedBy      uuid.UUID `json:"received_by"`
	LineItems       []POLineItemEvent `json:"line_items"`
}

// POLineItemEvent represents a line item in PO received event
type POLineItemEvent struct {
	MaterialID        uuid.UUID `json:"material_id"`
	Quantity          float64   `json:"quantity"`
	UnitID            uuid.UUID `json:"unit_id"`
	SupplierLotNumber string    `json:"supplier_lot_number,omitempty"`
	ExpiryDate        string    `json:"expiry_date,omitempty"` // RFC3339
}

// handlePOReceived handles PO received events - creates GRN
func (s *EventSubscriber) handlePOReceived(msg *nats.Msg) {
	var event POReceivedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("Failed to unmarshal PO received event", zap.Error(err))
		return
	}

	s.logger.Info("Received PO received event",
		zap.String("po_number", event.PONumber),
		zap.String("po_id", event.POID.String()),
	)

	// Auto-create GRN from PO received
	// This would call createGRNUC with the event data
	// Implementation depends on business requirements
}

// SalesOrderEvent represents a sales order event
type SalesOrderEvent struct {
	OrderID     uuid.UUID          `json:"order_id"`
	OrderNumber string             `json:"order_number"`
	CustomerID  uuid.UUID          `json:"customer_id"`
	LineItems   []OrderLineItemEvent `json:"line_items"`
}

// OrderLineItemEvent represents an order line item
type OrderLineItemEvent struct {
	MaterialID uuid.UUID `json:"material_id"`
	Quantity   float64   `json:"quantity"`
	UnitID     uuid.UUID `json:"unit_id"`
}

// handleSalesOrderConfirmed handles sales order confirmed events - reserves stock
func (s *EventSubscriber) handleSalesOrderConfirmed(msg *nats.Msg) {
	var event SalesOrderEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("Failed to unmarshal sales order event", zap.Error(err))
		return
	}

	s.logger.Info("Received sales order confirmed event",
		zap.String("order_number", event.OrderNumber),
	)

	ctx := context.Background()

	// Reserve stock for each line item
	for _, item := range event.LineItems {
		input := &reservation.CreateReservationInput{
			MaterialID:      item.MaterialID,
			Quantity:        item.Quantity,
			UnitID:          item.UnitID,
			ReservationType: entity.ReservationTypeSalesOrder,
			ReferenceID:     event.OrderID,
			ReferenceNumber: event.OrderNumber,
			CreatedBy:       uuid.Nil, // System
		}

		_, err := s.reserveStockUC.Execute(ctx, input)
		if err != nil {
			s.logger.Error("Failed to reserve stock for sales order",
				zap.String("order_number", event.OrderNumber),
				zap.String("material_id", item.MaterialID.String()),
				zap.Error(err),
			)
			// Publish reservation failed event
			continue
		}
	}

	s.logger.Info("Stock reserved for sales order",
		zap.String("order_number", event.OrderNumber),
		zap.Int("items", len(event.LineItems)),
	)
}

// handleSalesOrderCancelled handles sales order cancelled - releases reservations
func (s *EventSubscriber) handleSalesOrderCancelled(msg *nats.Msg) {
	var event struct {
		OrderID       uuid.UUID   `json:"order_id"`
		OrderNumber   string      `json:"order_number"`
		ReservationIDs []uuid.UUID `json:"reservation_ids"`
	}

	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("Failed to unmarshal sales order cancelled event", zap.Error(err))
		return
	}

	s.logger.Info("Received sales order cancelled event",
		zap.String("order_number", event.OrderNumber),
	)

	ctx := context.Background()

	// Release reservations
	for _, reservationID := range event.ReservationIDs {
		err := s.releaseReservationUC.Execute(ctx, reservationID)
		if err != nil {
			s.logger.Error("Failed to release reservation",
				zap.String("reservation_id", reservationID.String()),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("Reservations released for cancelled order",
		zap.String("order_number", event.OrderNumber),
	)
}

// WorkOrderEvent represents a work order event from manufacturing
type WorkOrderEvent struct {
	WorkOrderID     uuid.UUID              `json:"work_order_id"`
	WorkOrderNumber string                 `json:"work_order_number"`
	ProductID       uuid.UUID              `json:"product_id"`
	Quantity        float64                `json:"quantity"`
	Materials       []WorkOrderMaterialEvent `json:"materials"`
}

// WorkOrderMaterialEvent represents a material requirement
type WorkOrderMaterialEvent struct {
	MaterialID uuid.UUID `json:"material_id"`
	Quantity   float64   `json:"quantity"`
	UnitID     uuid.UUID `json:"unit_id"`
}

// handleWorkOrderStarted handles work order started - reserves materials
func (s *EventSubscriber) handleWorkOrderStarted(msg *nats.Msg) {
	var event WorkOrderEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("Failed to unmarshal work order event", zap.Error(err))
		return
	}

	s.logger.Info("Received work order started event",
		zap.String("wo_number", event.WorkOrderNumber),
	)

	ctx := context.Background()

	// Reserve materials for production
	for _, mat := range event.Materials {
		input := &reservation.CreateReservationInput{
			MaterialID:      mat.MaterialID,
			Quantity:        mat.Quantity,
			UnitID:          mat.UnitID,
			ReservationType: entity.ReservationTypeWorkOrder,
			ReferenceID:     event.WorkOrderID,
			ReferenceNumber: event.WorkOrderNumber,
			CreatedBy:       uuid.Nil, // System
		}

		_, err := s.reserveStockUC.Execute(ctx, input)
		if err != nil {
			s.logger.Error("Failed to reserve material for work order",
				zap.String("wo_number", event.WorkOrderNumber),
				zap.String("material_id", mat.MaterialID.String()),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("Materials reserved for work order",
		zap.String("wo_number", event.WorkOrderNumber),
		zap.Int("materials", len(event.Materials)),
	)
}

// ReservationRepository interface for querying reservations
type ReservationRepository interface {
	GetByReferenceID(ctx context.Context, referenceID uuid.UUID) ([]*entity.StockReservation, error)
}
