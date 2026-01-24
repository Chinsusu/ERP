package grpc

import (
	"context"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/usecase/po"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProcurementServer implements the gRPC Procurement service
type ProcurementServer struct {
	poRepo          repository.PORepository
	updateReceivedUC *po.UpdateReceivedQtyUseCase
	logger          *zap.Logger
}

// NewProcurementServer creates a new gRPC server
func NewProcurementServer(
	poRepo repository.PORepository,
	updateReceivedUC *po.UpdateReceivedQtyUseCase,
	logger *zap.Logger,
) *ProcurementServer {
	return &ProcurementServer{
		poRepo:          poRepo,
		updateReceivedUC: updateReceivedUC,
		logger:          logger,
	}
}

// GetPO returns a purchase order by ID
func (s *ProcurementServer) GetPO(ctx context.Context, poID string) (*entity.PurchaseOrder, error) {
	id, err := uuid.Parse(poID)
	if err != nil {
		return nil, err
	}
	return s.poRepo.GetByID(ctx, id)
}

// GetPOsBySupplier returns all POs for a supplier
func (s *ProcurementServer) GetPOsBySupplier(ctx context.Context, supplierID string, status string) ([]*entity.PurchaseOrder, error) {
	id, err := uuid.Parse(supplierID)
	if err != nil {
		return nil, err
	}
	
	pos, err := s.poRepo.GetBySupplierID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Filter by status if provided
	if status != "" {
		filtered := make([]*entity.PurchaseOrder, 0)
		for _, po := range pos {
			if string(po.Status) == status {
				filtered = append(filtered, po)
			}
		}
		return filtered, nil
	}
	
	return pos, nil
}

// UpdatePOReceivedQty updates received quantity for a PO line item
func (s *ProcurementServer) UpdatePOReceivedQty(
	ctx context.Context,
	poLineItemID string,
	receivedQty float64,
	grnID string,
	grnNumber string,
) (*UpdateReceivedQtyResult, error) {
	lineItemUUID, err := uuid.Parse(poLineItemID)
	if err != nil {
		return nil, err
	}
	
	var grnUUID *uuid.UUID
	if grnID != "" {
		if parsed, err := uuid.Parse(grnID); err == nil {
			grnUUID = &parsed
		}
	}
	
	updatedPO, err := s.updateReceivedUC.Execute(ctx, lineItemUUID, receivedQty, grnUUID, grnNumber)
	if err != nil {
		return nil, err
	}
	
	// Find the updated line item
	var lineItem *entity.POLineItem
	for i := range updatedPO.LineItems {
		if updatedPO.LineItems[i].ID == lineItemUUID {
			lineItem = &updatedPO.LineItems[i]
			break
		}
	}
	
	result := &UpdateReceivedQtyResult{
		Success:        true,
		Message:        "Received quantity updated successfully",
		POID:           updatedPO.ID.String(),
		POStatus:       string(updatedPO.Status),
	}
	
	if lineItem != nil {
		result.LineItemStatus = string(lineItem.Status)
		result.TotalReceived = lineItem.ReceivedQty
		result.PendingQty = lineItem.PendingQty
	}
	
	s.logger.Info("PO received quantity updated via gRPC",
		zap.String("po_id", updatedPO.ID.String()),
		zap.String("line_item_id", poLineItemID),
		zap.Float64("received_qty", receivedQty),
	)
	
	return result, nil
}

// UpdateReceivedQtyResult represents the result of updating received qty
type UpdateReceivedQtyResult struct {
	Success        bool
	Message        string
	POID           string
	POStatus       string
	LineItemStatus string
	TotalReceived  float64
	PendingQty     float64
}

// GetPendingPOs returns POs pending receipt
func (s *ProcurementServer) GetPendingPOs(ctx context.Context, page, limit int) ([]*entity.PurchaseOrder, int64, error) {
	filter := &repository.POFilter{
		Status: string(entity.POStatusConfirmed),
		Page:   page,
		Limit:  limit,
	}
	return s.poRepo.List(ctx, filter)
}
