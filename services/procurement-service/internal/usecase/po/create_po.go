package po

import (
	"context"
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreatePOFromPRRequest represents request to convert PR to PO
type CreatePOFromPRRequest struct {
	PRID                 uuid.UUID  `json:"-"`
	SupplierID           uuid.UUID  `json:"supplier_id" validate:"required"`
	SupplierCode         string     `json:"supplier_code"`
	SupplierName         string     `json:"supplier_name"`
	DeliveryAddress      string     `json:"delivery_address"`
	DeliveryTerms        string     `json:"delivery_terms"`
	PaymentTerms         string     `json:"payment_terms"`
	ExpectedDeliveryDate string     `json:"expected_delivery_date"`
	Notes                string     `json:"notes"`
	CreatedBy            uuid.UUID  `json:"-"`
}

// CreatePOFromPRUseCase handles converting PR to PO
type CreatePOFromPRUseCase struct {
	prRepo   repository.PRRepository
	poRepo   repository.PORepository
	eventPub *event.Publisher
}

// NewCreatePOFromPRUseCase creates a new use case
func NewCreatePOFromPRUseCase(
	prRepo repository.PRRepository,
	poRepo repository.PORepository,
	eventPub *event.Publisher,
) *CreatePOFromPRUseCase {
	return &CreatePOFromPRUseCase{
		prRepo:   prRepo,
		poRepo:   poRepo,
		eventPub: eventPub,
	}
}

// Execute converts an approved PR to PO
func (uc *CreatePOFromPRUseCase) Execute(ctx context.Context, req *CreatePOFromPRRequest) (*entity.PurchaseOrder, error) {
	// Get PR
	pr, err := uc.prRepo.GetByID(ctx, req.PRID)
	if err != nil {
		return nil, err
	}

	if !pr.CanConvert() {
		return nil, entity.ErrPRNotApproved
	}

	// Generate PO number
	poNumber, err := uc.poRepo.GetNextPONumber(ctx)
	if err != nil {
		return nil, err
	}

	// Parse expected date
	var expectedDate *time.Time
	if req.ExpectedDeliveryDate != "" {
		if d, err := time.Parse("2006-01-02", req.ExpectedDeliveryDate); err == nil {
			expectedDate = &d
		}
	}

	deliveryTerms := req.DeliveryTerms
	if deliveryTerms == "" {
		deliveryTerms = "EXW"
	}
	paymentTerms := req.PaymentTerms
	if paymentTerms == "" {
		paymentTerms = "Net 30"
	}

	po := &entity.PurchaseOrder{
		ID:                   uuid.New(),
		PONumber:             poNumber,
		PODate:               time.Now(),
		PRID:                 &pr.ID,
		SupplierID:           req.SupplierID,
		SupplierCode:         req.SupplierCode,
		SupplierName:         req.SupplierName,
		Status:               entity.POStatusDraft,
		DeliveryAddress:      req.DeliveryAddress,
		DeliveryTerms:        deliveryTerms,
		PaymentTerms:         paymentTerms,
		ExpectedDeliveryDate: expectedDate,
		Currency:             pr.Currency,
		Notes:                req.Notes,
		CreatedBy:            req.CreatedBy,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Create PO
	if err := uc.poRepo.Create(ctx, po); err != nil {
		return nil, err
	}

	// Create PO line items from PR line items
	for i, prLine := range pr.LineItems {
		poLine := &entity.POLineItem{
			ID:             uuid.New(),
			POID:           po.ID,
			PRLineItemID:   &prLine.ID,
			LineNumber:     i + 1,
			MaterialID:     prLine.MaterialID,
			MaterialCode:   prLine.MaterialCode,
			MaterialName:   prLine.MaterialName,
			Quantity:       prLine.Quantity,
			ReceivedQty:    0,
			PendingQty:     prLine.Quantity,
			UOMCode:        prLine.UOMCode,
			UnitPrice:      prLine.UnitPrice,
			Currency:       prLine.Currency,
			Specifications: prLine.Specifications,
			Status:         entity.POLineStatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		poLine.CalculateLineTotal()

		if err := uc.poRepo.CreateLineItem(ctx, poLine); err != nil {
			return nil, err
		}
		po.LineItems = append(po.LineItems, *poLine)
	}

	// Calculate totals
	po.CalculateTotals()
	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, err
	}

	// Mark PR as converted
	pr.MarkConvertedToPO(po.ID)
	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishPOCreated(ctx, &event.POEvent{
		POID:         po.ID.String(),
		PONumber:     po.PONumber,
		Status:       string(po.Status),
		SupplierID:   po.SupplierID.String(),
		SupplierName: po.SupplierName,
		TotalAmount:  po.GrandTotal,
		Timestamp:    time.Now(),
	})

	return po, nil
}

// GetPOUseCase handles getting a PO
type GetPOUseCase struct {
	poRepo repository.PORepository
}

// NewGetPOUseCase creates a new use case
func NewGetPOUseCase(poRepo repository.PORepository) *GetPOUseCase {
	return &GetPOUseCase{poRepo: poRepo}
}

// Execute gets a PO by ID
func (uc *GetPOUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.PurchaseOrder, error) {
	return uc.poRepo.GetByID(ctx, id)
}

// ListPOsUseCase handles listing POs
type ListPOsUseCase struct {
	poRepo repository.PORepository
}

// NewListPOsUseCase creates a new use case
func NewListPOsUseCase(poRepo repository.PORepository) *ListPOsUseCase {
	return &ListPOsUseCase{poRepo: poRepo}
}

// Execute lists POs with filters
func (uc *ListPOsUseCase) Execute(ctx context.Context, filter *repository.POFilter) ([]*entity.PurchaseOrder, int64, error) {
	return uc.poRepo.List(ctx, filter)
}
