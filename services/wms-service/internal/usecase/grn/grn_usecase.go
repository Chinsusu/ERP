package grn

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateGRNUseCase handles GRN creation
type CreateGRNUseCase struct {
	grnRepo      repository.GRNRepository
	lotRepo      repository.LotRepository
	stockRepo    repository.StockRepository
	zoneRepo     repository.ZoneRepository
	locationRepo repository.LocationRepository
	eventPub     *event.Publisher
}

// NewCreateGRNUseCase creates a new use case
func NewCreateGRNUseCase(
	grnRepo repository.GRNRepository,
	lotRepo repository.LotRepository,
	stockRepo repository.StockRepository,
	zoneRepo repository.ZoneRepository,
	locationRepo repository.LocationRepository,
	eventPub *event.Publisher,
) *CreateGRNUseCase {
	return &CreateGRNUseCase{
		grnRepo:      grnRepo,
		lotRepo:      lotRepo,
		stockRepo:    stockRepo,
		zoneRepo:     zoneRepo,
		locationRepo: locationRepo,
		eventPub:     eventPub,
	}
}

// CreateGRNInput represents input for creating GRN
type CreateGRNInput struct {
	GRNDate            time.Time
	POID               *uuid.UUID
	PONumber           string
	SupplierID         *uuid.UUID
	WarehouseID        uuid.UUID
	DeliveryNoteNumber string
	VehicleNumber      string
	Notes              string
	ReceivedBy         uuid.UUID
	Items              []CreateGRNItemInput
}

// CreateGRNItemInput represents input for GRN line item
type CreateGRNItemInput struct {
	POLineItemID      *uuid.UUID
	MaterialID        uuid.UUID
	ExpectedQty       *float64
	ReceivedQty       float64
	UnitID            uuid.UUID
	SupplierLotNumber string
	ManufacturedDate  *time.Time
	ExpiryDate        time.Time
	LocationID        *uuid.UUID
}

// Execute creates a GRN
func (uc *CreateGRNUseCase) Execute(ctx context.Context, input *CreateGRNInput) (*entity.GRN, error) {
	// Generate GRN number
	grnNumber, err := uc.grnRepo.GetNextGRNNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Get quarantine zone for initial placement
	quarantineZone, _ := uc.zoneRepo.GetQuarantineZone(ctx, input.WarehouseID)
	var defaultLocationID *uuid.UUID
	if quarantineZone != nil {
		locations, _ := uc.locationRepo.GetByZoneID(ctx, quarantineZone.ID)
		if len(locations) > 0 {
			defaultLocationID = &locations[0].ID
		}
	}

	// Create GRN
	grn := &entity.GRN{
		GRNNumber:          grnNumber,
		GRNDate:            input.GRNDate,
		POID:               input.POID,
		PONumber:           input.PONumber,
		SupplierID:         input.SupplierID,
		WarehouseID:        input.WarehouseID,
		DeliveryNoteNumber: input.DeliveryNoteNumber,
		VehicleNumber:      input.VehicleNumber,
		Status:             entity.GRNStatusDraft,
		QCStatus:           entity.QCStatusPending,
		Notes:              input.Notes,
		ReceivedBy:         &input.ReceivedBy,
	}

	if err := uc.grnRepo.Create(ctx, grn); err != nil {
		return nil, err
	}

	// Create line items with lots
	for i, item := range input.Items {
		// Generate lot number
		lotNumber, err := uc.lotRepo.GetNextLotNumber(ctx)
		if err != nil {
			return nil, err
		}

		// Create lot
		lot := &entity.Lot{
			LotNumber:         lotNumber,
			MaterialID:        item.MaterialID,
			SupplierID:        input.SupplierID,
			SupplierLotNumber: item.SupplierLotNumber,
			ManufacturedDate:  item.ManufacturedDate,
			ExpiryDate:        item.ExpiryDate,
			ReceivedDate:      input.GRNDate,
			GRNID:             &grn.ID,
			QCStatus:          entity.QCStatusPending,
			Status:            entity.LotStatusAvailable,
		}

		if err := uc.lotRepo.Create(ctx, lot); err != nil {
			return nil, err
		}

		// Determine location
		locationID := item.LocationID
		if locationID == nil && defaultLocationID != nil {
			locationID = defaultLocationID
		}

		// Create GRN line item
		lineItem := &entity.GRNLineItem{
			GRNID:             grn.ID,
			LineNumber:        i + 1,
			POLineItemID:      item.POLineItemID,
			MaterialID:        item.MaterialID,
			ExpectedQty:       item.ExpectedQty,
			ReceivedQty:       item.ReceivedQty,
			UnitID:            item.UnitID,
			LotID:             &lot.ID,
			SupplierLotNumber: item.SupplierLotNumber,
			ManufacturedDate:  item.ManufacturedDate,
			ExpiryDate:        item.ExpiryDate,
			LocationID:        locationID,
			QCStatus:          entity.QCStatusPending,
		}

		if err := uc.grnRepo.CreateLineItem(ctx, lineItem); err != nil {
			return nil, err
		}
	}

	// Publish event
	poID := ""
	if input.POID != nil {
		poID = input.POID.String()
	}
	uc.eventPub.PublishGRNCreated(&event.GRNCreatedEvent{
		GRNID:     grn.ID.String(),
		GRNNumber: grnNumber,
		POID:      poID,
	})

	return grn, nil
}

// CompleteGRNUseCase handles completing GRN after QC
type CompleteGRNUseCase struct {
	grnRepo   repository.GRNRepository
	lotRepo   repository.LotRepository
	stockRepo repository.StockRepository
	zoneRepo  repository.ZoneRepository
	eventPub  *event.Publisher
}

// NewCompleteGRNUseCase creates a new use case
func NewCompleteGRNUseCase(
	grnRepo repository.GRNRepository,
	lotRepo repository.LotRepository,
	stockRepo repository.StockRepository,
	zoneRepo repository.ZoneRepository,
	eventPub *event.Publisher,
) *CompleteGRNUseCase {
	return &CompleteGRNUseCase{
		grnRepo:   grnRepo,
		lotRepo:   lotRepo,
		stockRepo: stockRepo,
		zoneRepo:  zoneRepo,
		eventPub:  eventPub,
	}
}

// CompleteGRNInput represents input for completing GRN
type CompleteGRNInput struct {
	GRNID    uuid.UUID
	QCStatus entity.QCStatus
	QCNotes  string
}

// Execute completes the GRN after QC
func (uc *CompleteGRNUseCase) Execute(ctx context.Context, input *CompleteGRNInput) (*entity.GRN, error) {
	// Get GRN
	grn, err := uc.grnRepo.GetByID(ctx, input.GRNID)
	if err != nil {
		return nil, err
	}

	if !grn.CanComplete() {
		return nil, entity.ErrAlreadyCompleted
	}

	// Get line items
	items, err := uc.grnRepo.GetLineItemsByGRNID(ctx, grn.ID)
	if err != nil {
		return nil, err
	}

	// Process each line item
	eventItems := make([]event.GRNCompletedEventItem, 0)
	for _, item := range items {
		// Update lot QC status
		if item.LotID != nil {
			lot, err := uc.lotRepo.GetByID(ctx, *item.LotID)
			if err != nil {
				return nil, err
			}

			if input.QCStatus == entity.QCStatusPassed {
				lot.PassQC()
				item.PassQC(item.ReceivedQty)
			} else {
				lot.FailQC()
				item.FailQC(input.QCNotes)
			}

			if err := uc.lotRepo.Update(ctx, lot); err != nil {
				return nil, err
			}

			// Create stock if QC passed
			if input.QCStatus == entity.QCStatusPassed && item.LocationID != nil {
				location, _ := uc.stockRepo.GetByID(ctx, *item.LocationID)
				
				stock := &entity.Stock{
					WarehouseID: grn.WarehouseID,
					ZoneID:      location.ZoneID,
					LocationID:  *item.LocationID,
					MaterialID:  item.MaterialID,
					LotID:       item.LotID,
					Quantity:    item.ReceivedQty,
					UnitID:      item.UnitID,
				}

				movementNumber, _ := uc.stockRepo.GetNextMovementNumber(ctx, entity.MovementTypeIn)
				movement := entity.NewStockMovementIn(
					item.MaterialID,
					*item.LotID,
					*item.LocationID,
					item.UnitID,
					*grn.ReceivedBy,
					item.ReceivedQty,
					entity.ReferenceTypeGRN,
					&grn.ID,
					movementNumber,
				)

				if err := uc.stockRepo.ReceiveStock(ctx, stock, movement); err != nil {
					return nil, err
				}

				// Publish stock received event
				uc.eventPub.PublishStockReceived(&event.StockReceivedEvent{
					MaterialID:  item.MaterialID.String(),
					LotID:       item.LotID.String(),
					Quantity:    item.ReceivedQty,
					LocationID:  item.LocationID.String(),
					WarehouseID: grn.WarehouseID.String(),
				})
			}

			// Add to event items
			acceptedQty := 0.0
			if item.AcceptedQty != nil {
				acceptedQty = *item.AcceptedQty
			}
			eventItems = append(eventItems, event.GRNCompletedEventItem{
				MaterialID:  item.MaterialID.String(),
				LotID:       item.LotID.String(),
				LotNumber:   lot.LotNumber,
				ReceivedQty: item.ReceivedQty,
				AcceptedQty: acceptedQty,
			})
		}

		// Update line item
		if err := uc.grnRepo.UpdateLineItem(ctx, item); err != nil {
			return nil, err
		}
	}

	// Complete GRN
	grn.Complete(input.QCStatus, input.QCNotes)
	if err := uc.grnRepo.Update(ctx, grn); err != nil {
		return nil, err
	}

	// Publish GRN completed event (Procurement will update PO)
	poID := ""
	if grn.POID != nil {
		poID = grn.POID.String()
	}
	uc.eventPub.PublishGRNCompleted(&event.GRNCompletedEvent{
		GRNID:       grn.ID.String(),
		GRNNumber:   grn.GRNNumber,
		POID:        poID,
		WarehouseID: grn.WarehouseID.String(),
		Items:       eventItems,
	})

	return grn, nil
}

// GetGRNUseCase handles getting GRN
type GetGRNUseCase struct {
	grnRepo repository.GRNRepository
}

// NewGetGRNUseCase creates a new use case
func NewGetGRNUseCase(grnRepo repository.GRNRepository) *GetGRNUseCase {
	return &GetGRNUseCase{grnRepo: grnRepo}
}

// Execute gets a GRN by ID
func (uc *GetGRNUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.GRN, error) {
	return uc.grnRepo.GetByID(ctx, id)
}

// ListGRNsUseCase handles listing GRNs
type ListGRNsUseCase struct {
	grnRepo repository.GRNRepository
}

// NewListGRNsUseCase creates a new use case
func NewListGRNsUseCase(grnRepo repository.GRNRepository) *ListGRNsUseCase {
	return &ListGRNsUseCase{grnRepo: grnRepo}
}

// Execute lists GRNs
func (uc *ListGRNsUseCase) Execute(ctx context.Context, filter *repository.GRNFilter) ([]*entity.GRN, int64, error) {
	return uc.grnRepo.List(ctx, filter)
}
