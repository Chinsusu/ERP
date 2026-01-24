package inventory

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
)

// CreateInventoryCountUseCase handles creating inventory count
type CreateInventoryCountUseCase struct {
	countRepo    repository.InventoryCountRepository
	stockRepo    repository.StockRepository
	locationRepo repository.LocationRepository
}

// NewCreateInventoryCountUseCase creates a new use case
func NewCreateInventoryCountUseCase(
	countRepo repository.InventoryCountRepository,
	stockRepo repository.StockRepository,
	locationRepo repository.LocationRepository,
) *CreateInventoryCountUseCase {
	return &CreateInventoryCountUseCase{
		countRepo:    countRepo,
		stockRepo:    stockRepo,
		locationRepo: locationRepo,
	}
}

// CreateInventoryCountInput represents input for creating inventory count
type CreateInventoryCountInput struct {
	CountDate   time.Time
	CountType   entity.InventoryCountType
	WarehouseID uuid.UUID
	ZoneID      *uuid.UUID
	Notes       string
	CreatedBy   uuid.UUID
}

// Execute creates an inventory count with line items from current stock
func (uc *CreateInventoryCountUseCase) Execute(ctx context.Context, input *CreateInventoryCountInput) (*entity.InventoryCount, error) {
	// Generate count number
	countNumber, err := uc.countRepo.GetNextCountNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Create inventory count
	count := &entity.InventoryCount{
		CountNumber: countNumber,
		CountDate:   input.CountDate,
		CountType:   input.CountType,
		WarehouseID: input.WarehouseID,
		ZoneID:      input.ZoneID,
		Status:      entity.InventoryCountStatusDraft,
		Notes:       input.Notes,
		CreatedBy:   input.CreatedBy,
	}

	if err := uc.countRepo.Create(ctx, count); err != nil {
		return nil, err
	}

	// Get stock for the warehouse/zone and create line items
	filter := &repository.StockFilter{
		WarehouseID: &input.WarehouseID,
	}
	if input.ZoneID != nil {
		filter.ZoneID = input.ZoneID
	}
	hasStock := true
	filter.HasStock = &hasStock

	stocks, _, err := uc.stockRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	var lineItems []*entity.InventoryCountLineItem
	for _, stock := range stocks {
		lineItem := &entity.InventoryCountLineItem{
			InventoryCountID: count.ID,
			LocationID:       stock.LocationID,
			MaterialID:       stock.MaterialID,
			LotID:            stock.LotID,
			UnitID:           stock.UnitID,
			SystemQty:        stock.Quantity,
			IsCounted:        false,
		}
		lineItems = append(lineItems, lineItem)
	}

	if len(lineItems) > 0 {
		if err := uc.countRepo.CreateLineItems(ctx, lineItems); err != nil {
			return nil, err
		}
	}

	return count, nil
}

// StartInventoryCountUseCase handles starting inventory count
type StartInventoryCountUseCase struct {
	countRepo repository.InventoryCountRepository
}

// NewStartInventoryCountUseCase creates a new use case
func NewStartInventoryCountUseCase(countRepo repository.InventoryCountRepository) *StartInventoryCountUseCase {
	return &StartInventoryCountUseCase{countRepo: countRepo}
}

// Execute starts the inventory count
func (uc *StartInventoryCountUseCase) Execute(ctx context.Context, countID uuid.UUID) (*entity.InventoryCount, error) {
	count, err := uc.countRepo.GetByID(ctx, countID)
	if err != nil {
		return nil, err
	}

	if !count.CanStart() {
		return nil, entity.ErrInvalidStatus
	}

	count.Start()
	if err := uc.countRepo.Update(ctx, count); err != nil {
		return nil, err
	}

	return count, nil
}

// RecordCountUseCase handles recording a count
type RecordCountUseCase struct {
	countRepo repository.InventoryCountRepository
}

// NewRecordCountUseCase creates a new use case
func NewRecordCountUseCase(countRepo repository.InventoryCountRepository) *RecordCountUseCase {
	return &RecordCountUseCase{countRepo: countRepo}
}

// RecordCountInput represents input for recording a count
type RecordCountInput struct {
	LineItemID uuid.UUID
	CountedQty float64
	CountedBy  uuid.UUID
	Notes      string
}

// Execute records the counted quantity
func (uc *RecordCountUseCase) Execute(ctx context.Context, input *RecordCountInput) (*entity.InventoryCountLineItem, error) {
	items, err := uc.countRepo.GetLineItemsByCountID(ctx, input.LineItemID)
	if err != nil {
		return nil, err
	}

	var lineItem *entity.InventoryCountLineItem
	for _, item := range items {
		if item.ID == input.LineItemID {
			lineItem = item
			break
		}
	}

	if lineItem == nil {
		return nil, entity.ErrNotFound
	}

	lineItem.RecordCount(input.CountedQty, input.CountedBy)
	lineItem.Notes = input.Notes

	if err := uc.countRepo.UpdateLineItem(ctx, lineItem); err != nil {
		return nil, err
	}

	return lineItem, nil
}

// CompleteInventoryCountUseCase handles completing inventory count
type CompleteInventoryCountUseCase struct {
	countRepo repository.InventoryCountRepository
	stockRepo repository.StockRepository
}

// NewCompleteInventoryCountUseCase creates a new use case
func NewCompleteInventoryCountUseCase(
	countRepo repository.InventoryCountRepository,
	stockRepo repository.StockRepository,
) *CompleteInventoryCountUseCase {
	return &CompleteInventoryCountUseCase{
		countRepo: countRepo,
		stockRepo: stockRepo,
	}
}

// CompleteInventoryCountInput represents input for completing count
type CompleteInventoryCountInput struct {
	CountID       uuid.UUID
	ApplyVariance bool // If true, adjust stock based on variance
	ApprovedBy    uuid.UUID
}

// Execute completes the inventory count
func (uc *CompleteInventoryCountUseCase) Execute(ctx context.Context, input *CompleteInventoryCountInput) (*entity.InventoryCount, error) {
	count, err := uc.countRepo.GetByID(ctx, input.CountID)
	if err != nil {
		return nil, err
	}

	if !count.CanComplete() {
		return nil, entity.ErrInvalidStatus
	}

	// Check all items are counted
	pending, err := uc.countRepo.GetPendingItems(ctx, count.ID)
	if err != nil {
		return nil, err
	}
	if len(pending) > 0 {
		return nil, entity.ErrPendingItems
	}

	// Apply variance if requested
	if input.ApplyVariance {
		varianceItems, err := uc.countRepo.GetVarianceItems(ctx, count.ID)
		if err != nil {
			return nil, err
		}

		for _, item := range varianceItems {
			// Get stock and adjust
			stock, err := uc.stockRepo.GetByLocationMaterialLot(ctx, item.LocationID, item.MaterialID, item.LotID)
			if err != nil {
				continue
			}

			movementNumber, _ := uc.stockRepo.GetNextMovementNumber(ctx, entity.MovementTypeAdjustment)
			movement := &entity.StockMovement{
				MovementNumber: movementNumber,
				MovementType:   entity.MovementTypeAdjustment,
				ReferenceType:  entity.ReferenceTypeAdjustment,
				ReferenceID:    &count.ID,
				MaterialID:     item.MaterialID,
				LotID:          item.LotID,
				FromLocationID: &item.LocationID,
				ToLocationID:   &item.LocationID,
				Quantity:       item.Variance,
				UnitID:         item.UnitID,
				Notes:          "Inventory Count Adjustment",
				CreatedBy:      input.ApprovedBy,
			}

			uc.stockRepo.AdjustStock(ctx, stock, item.Variance, movement)
		}
	}

	count.Complete(input.ApprovedBy)
	if err := uc.countRepo.Update(ctx, count); err != nil {
		return nil, err
	}

	return count, nil
}

// GetInventoryCountUseCase handles getting inventory count
type GetInventoryCountUseCase struct {
	countRepo repository.InventoryCountRepository
}

// NewGetInventoryCountUseCase creates a new use case
func NewGetInventoryCountUseCase(countRepo repository.InventoryCountRepository) *GetInventoryCountUseCase {
	return &GetInventoryCountUseCase{countRepo: countRepo}
}

// Execute gets inventory count by ID
func (uc *GetInventoryCountUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.InventoryCount, error) {
	return uc.countRepo.GetByID(ctx, id)
}

// ListInventoryCountsUseCase handles listing inventory counts
type ListInventoryCountsUseCase struct {
	countRepo repository.InventoryCountRepository
}

// NewListInventoryCountsUseCase creates a new use case
func NewListInventoryCountsUseCase(countRepo repository.InventoryCountRepository) *ListInventoryCountsUseCase {
	return &ListInventoryCountsUseCase{countRepo: countRepo}
}

// Execute lists inventory counts
func (uc *ListInventoryCountsUseCase) Execute(ctx context.Context, filter *repository.InventoryCountFilter) ([]*entity.InventoryCount, int64, error) {
	return uc.countRepo.List(ctx, filter)
}
