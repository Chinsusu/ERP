package subscriber

import (
	"context"
	"encoding/json"

	"github.com/erp-cosmetics/procurement-service/internal/usecase/po"
	"github.com/erp-cosmetics/shared/pkg/nats"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EventSubscriber handles incoming NATS events
type EventSubscriber struct {
	client           *nats.Client
	updateReceivedUC *po.UpdateReceivedQtyUseCase
	logger           *zap.Logger
}

// NewEventSubscriber creates a new event subscriber
func NewEventSubscriber(
	client *nats.Client,
	updateReceivedUC *po.UpdateReceivedQtyUseCase,
	logger *zap.Logger,
) *EventSubscriber {
	return &EventSubscriber{
		client:           client,
		updateReceivedUC: updateReceivedUC,
		logger:           logger,
	}
}

// GRNCompletedEvent represents the wms.grn.completed event payload
type GRNCompletedEvent struct {
	GRNID          string  `json:"grn_id"`
	GRNNumber      string  `json:"grn_number"`
	POID           string  `json:"po_id"`
	POLineItemID   string  `json:"po_line_item_id"`
	MaterialID     string  `json:"material_id"`
	ReceivedQty    float64 `json:"received_qty"`
	AcceptedQty    float64 `json:"accepted_qty"`
	RejectedQty    float64 `json:"rejected_qty"`
	BatchNumber    string  `json:"batch_number"`
	LotNumber      string  `json:"lot_number"`
	ExpiryDate     string  `json:"expiry_date"`
	StorageLocation string `json:"storage_location"`
	ReceivedBy     string  `json:"received_by"`
	ReceivedDate   string  `json:"received_date"`
}

// SupplierBlockedEvent represents the supplier.blocked event payload
type SupplierBlockedEvent struct {
	SupplierID   string `json:"supplier_id"`
	SupplierCode string `json:"supplier_code"`
	SupplierName string `json:"supplier_name"`
	Reason       string `json:"reason"`
	BlockedBy    string `json:"blocked_by"`
}

// Subscribe starts listening to relevant events
func (s *EventSubscriber) Subscribe(ctx context.Context) error {
	// Subscribe to WMS GRN completed events
	if err := s.subscribeToGRNCompleted(ctx); err != nil {
		return err
	}

	// Subscribe to Supplier blocked events
	if err := s.subscribeToSupplierBlocked(ctx); err != nil {
		return err
	}

	s.logger.Info("Event subscribers started",
		zap.Strings("topics", []string{"wms.grn.completed", "supplier.blocked"}),
	)

	return nil
}

func (s *EventSubscriber) subscribeToGRNCompleted(ctx context.Context) error {
	_, err := s.client.Subscribe("wms.grn.completed", "procurement-service", func(msg []byte) error {
		var event GRNCompletedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			s.logger.Error("Failed to unmarshal GRN completed event", zap.Error(err))
			return err
		}

		s.logger.Info("Received GRN completed event",
			zap.String("grn_number", event.GRNNumber),
			zap.String("po_line_item_id", event.POLineItemID),
			zap.Float64("received_qty", event.ReceivedQty),
		)

		// Parse UUIDs
		lineItemID, err := uuid.Parse(event.POLineItemID)
		if err != nil {
			s.logger.Error("Invalid PO line item ID", zap.String("id", event.POLineItemID), zap.Error(err))
			return err
		}

		var grnID *uuid.UUID
		if event.GRNID != "" {
			if parsed, err := uuid.Parse(event.GRNID); err == nil {
				grnID = &parsed
			}
		}

		// Update received quantity
		_, err = s.updateReceivedUC.Execute(ctx, lineItemID, event.AcceptedQty, grnID, event.GRNNumber)
		if err != nil {
			s.logger.Error("Failed to update PO received quantity",
				zap.String("po_line_item_id", event.POLineItemID),
				zap.Error(err),
			)
			return err
		}

		s.logger.Info("PO received quantity updated from GRN",
			zap.String("grn_number", event.GRNNumber),
			zap.String("po_line_item_id", event.POLineItemID),
			zap.Float64("accepted_qty", event.AcceptedQty),
		)
		return nil
	})
	return err
}

func (s *EventSubscriber) subscribeToSupplierBlocked(ctx context.Context) error {
	_, err := s.client.Subscribe("supplier.blocked", "procurement-service", func(msg []byte) error {
		var event SupplierBlockedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			s.logger.Error("Failed to unmarshal supplier blocked event", zap.Error(err))
			return err
		}

		s.logger.Warn("Supplier blocked - check pending POs",
			zap.String("supplier_id", event.SupplierID),
			zap.String("supplier_name", event.SupplierName),
			zap.String("reason", event.Reason),
		)

		// TODO: Implement logic to:
		// 1. Find all confirmed POs for this supplier
		// 2. Send notification to procurement team
		// 3. Optionally flag POs for review
		return nil
	})
	return err
}
