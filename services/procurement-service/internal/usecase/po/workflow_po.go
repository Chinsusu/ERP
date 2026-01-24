package po

import (
	"context"
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// ConfirmPOUseCase handles confirming a PO
type ConfirmPOUseCase struct {
	poRepo   repository.PORepository
	eventPub *event.Publisher
}

// NewConfirmPOUseCase creates a new use case
func NewConfirmPOUseCase(poRepo repository.PORepository, eventPub *event.Publisher) *ConfirmPOUseCase {
	return &ConfirmPOUseCase{poRepo: poRepo, eventPub: eventPub}
}

// Execute confirms a PO
func (uc *ConfirmPOUseCase) Execute(ctx context.Context, id uuid.UUID, confirmedBy uuid.UUID) (*entity.PurchaseOrder, error) {
	po, err := uc.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !po.CanConfirm() {
		return nil, entity.ErrInvalidPOStatus
	}

	po.Confirm(confirmedBy)
	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, err
	}

	// Publish event - WMS subscribes to this
	expectedDate := ""
	if po.ExpectedDeliveryDate != nil {
		expectedDate = po.ExpectedDeliveryDate.Format("2006-01-02")
	}
	uc.eventPub.PublishPOConfirmed(ctx, &event.POEvent{
		POID:                 po.ID.String(),
		PONumber:             po.PONumber,
		Status:               string(po.Status),
		SupplierID:           po.SupplierID.String(),
		SupplierName:         po.SupplierName,
		TotalAmount:          po.GrandTotal,
		ExpectedDeliveryDate: expectedDate,
		Timestamp:            time.Now(),
	})

	return po, nil
}

// CancelPOUseCase handles cancelling a PO
type CancelPOUseCase struct {
	poRepo   repository.PORepository
	eventPub *event.Publisher
}

// NewCancelPOUseCase creates a new use case
func NewCancelPOUseCase(poRepo repository.PORepository, eventPub *event.Publisher) *CancelPOUseCase {
	return &CancelPOUseCase{poRepo: poRepo, eventPub: eventPub}
}

// Execute cancels a PO
func (uc *CancelPOUseCase) Execute(ctx context.Context, id uuid.UUID, cancelledBy uuid.UUID, reason string) (*entity.PurchaseOrder, error) {
	po, err := uc.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !po.CanCancel() {
		return nil, entity.ErrInvalidPOStatus
	}

	po.Cancel(cancelledBy, reason)
	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, err
	}

	return po, nil
}

// UpdateReceivedQtyUseCase handles updating received quantity (called by WMS)
type UpdateReceivedQtyUseCase struct {
	poRepo   repository.PORepository
	eventPub *event.Publisher
}

// NewUpdateReceivedQtyUseCase creates a new use case
func NewUpdateReceivedQtyUseCase(poRepo repository.PORepository, eventPub *event.Publisher) *UpdateReceivedQtyUseCase {
	return &UpdateReceivedQtyUseCase{poRepo: poRepo, eventPub: eventPub}
}

// Execute updates received qty for a PO line item
func (uc *UpdateReceivedQtyUseCase) Execute(ctx context.Context, poLineItemID uuid.UUID, receivedQty float64, grnID *uuid.UUID, grnNumber string) (*entity.PurchaseOrder, error) {
	// Get line item
	lineItem, err := uc.poRepo.GetLineItemByID(ctx, poLineItemID)
	if err != nil {
		return nil, err
	}

	// Get PO
	po, err := uc.poRepo.GetByID(ctx, lineItem.POID)
	if err != nil {
		return nil, err
	}

	// Update line item
	lineItem.UpdateReceivedQty(receivedQty)
	if err := uc.poRepo.UpdateLineItem(ctx, lineItem); err != nil {
		return nil, err
	}

	// Create receipt record
	receipt := &entity.POReceipt{
		ID:           uuid.New(),
		POID:         po.ID,
		POLineItemID: lineItem.ID,
		GRNID:        grnID,
		GRNNumber:    grnNumber,
		ReceivedQty:  receivedQty,
		ReceivedDate: time.Now(),
		QCStatus:     entity.QCStatusPending,
		CreatedAt:    time.Now(),
	}
	if err := uc.poRepo.CreateReceipt(ctx, receipt); err != nil {
		return nil, err
	}

	// Refresh line items and update PO status
	lineItemPtrs, _ := uc.poRepo.GetLineItemsByPOID(ctx, po.ID)
	po.LineItems = make([]entity.POLineItem, len(lineItemPtrs))
	for i, item := range lineItemPtrs {
		po.LineItems[i] = *item
	}
	po.UpdateReceivedStatus()
	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishPOReceived(ctx, &event.POEvent{
		POID:        po.ID.String(),
		PONumber:    po.PONumber,
		Status:      string(po.Status),
		SupplierID:  po.SupplierID.String(),
		TotalAmount: po.GrandTotal,
		Timestamp:   time.Now(),
	})

	return po, nil
}

// ClosePOUseCase handles closing a PO
type ClosePOUseCase struct {
	poRepo   repository.PORepository
	eventPub *event.Publisher
}

// NewClosePOUseCase creates a new use case
func NewClosePOUseCase(poRepo repository.PORepository, eventPub *event.Publisher) *ClosePOUseCase {
	return &ClosePOUseCase{poRepo: poRepo, eventPub: eventPub}
}

// Execute closes a fully received PO
func (uc *ClosePOUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.PurchaseOrder, error) {
	po, err := uc.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !po.CanClose() {
		return nil, entity.ErrInvalidPOStatus
	}

	po.Close()
	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishPOClosed(ctx, &event.POEvent{
		POID:       po.ID.String(),
		PONumber:   po.PONumber,
		Status:     string(po.Status),
		SupplierID: po.SupplierID.String(),
		Timestamp:  time.Now(),
	})

	return po, nil
}

// GetPOReceiptsUseCase handles getting receipts for a PO
type GetPOReceiptsUseCase struct {
	poRepo repository.PORepository
}

// NewGetPOReceiptsUseCase creates a new use case
func NewGetPOReceiptsUseCase(poRepo repository.PORepository) *GetPOReceiptsUseCase {
	return &GetPOReceiptsUseCase{poRepo: poRepo}
}

// Execute gets receipts for a PO
func (uc *GetPOReceiptsUseCase) Execute(ctx context.Context, poID uuid.UUID) ([]*entity.POReceipt, error) {
	return uc.poRepo.GetReceiptsByPOID(ctx, poID)
}
