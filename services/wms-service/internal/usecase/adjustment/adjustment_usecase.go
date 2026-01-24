package adjustment

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
)

// AdjustmentType defines types of adjustments
type AdjustmentType string

const (
	AdjustmentTypeCycleCount AdjustmentType = "CYCLE_COUNT"
	AdjustmentTypeDamage     AdjustmentType = "DAMAGE"
	AdjustmentTypeExpiry     AdjustmentType = "EXPIRY"
	AdjustmentTypeCorrection AdjustmentType = "CORRECTION"
)

// CreateAdjustmentUseCase handles stock adjustments
type CreateAdjustmentUseCase struct {
	stockRepo repository.StockRepository
}

// NewCreateAdjustmentUseCase creates a new use case
func NewCreateAdjustmentUseCase(stockRepo repository.StockRepository) *CreateAdjustmentUseCase {
	return &CreateAdjustmentUseCase{stockRepo: stockRepo}
}

// CreateAdjustmentInput represents input for creating adjustment
type CreateAdjustmentInput struct {
	AdjustmentDate time.Time
	AdjustmentType AdjustmentType
	LocationID     uuid.UUID
	MaterialID     uuid.UUID
	LotID          *uuid.UUID
	UnitID         uuid.UUID
	SystemQty      float64
	ActualQty      float64
	Reason         string
	Notes          string
	AdjustedBy     uuid.UUID
}

// CreateAdjustmentOutput represents adjustment output
type CreateAdjustmentOutput struct {
	AdjustmentNumber string
	Variance         float64
	MovementNumber   string
}

// Execute creates a stock adjustment
func (uc *CreateAdjustmentUseCase) Execute(ctx context.Context, input *CreateAdjustmentInput) (*CreateAdjustmentOutput, error) {
	// Calculate variance
	variance := input.ActualQty - input.SystemQty

	// Get stock record
	stock, err := uc.stockRepo.GetByLocationMaterialLot(ctx, input.LocationID, input.MaterialID, input.LotID)
	if err != nil {
		return nil, err
	}

	// Create movement record
	movementNumber, err := uc.stockRepo.GetNextMovementNumber(ctx, entity.MovementTypeAdjustment)
	if err != nil {
		return nil, err
	}

	movement := &entity.StockMovement{
		MovementNumber:  movementNumber,
		MovementType:    entity.MovementTypeAdjustment,
		ReferenceType:   entity.ReferenceTypeAdjustment,
		MaterialID:      input.MaterialID,
		LotID:           input.LotID,
		FromLocationID:  &input.LocationID,
		ToLocationID:    &input.LocationID,
		Quantity:        variance, // Can be negative
		UnitID:          input.UnitID,
		Notes:           input.Notes,
		CreatedBy:       input.AdjustedBy,
	}

	// Adjust stock
	if err := uc.stockRepo.AdjustStock(ctx, stock, variance, movement); err != nil {
		return nil, err
	}

	return &CreateAdjustmentOutput{
		AdjustmentNumber: movementNumber,
		Variance:         variance,
		MovementNumber:   movementNumber,
	}, nil
}

// TransferStockUseCase handles stock transfers between locations
type TransferStockUseCase struct {
	stockRepo repository.StockRepository
}

// NewTransferStockUseCase creates a new use case
func NewTransferStockUseCase(stockRepo repository.StockRepository) *TransferStockUseCase {
	return &TransferStockUseCase{stockRepo: stockRepo}
}

// TransferStockInput represents input for transferring stock
type TransferStockInput struct {
	MaterialID     uuid.UUID
	LotID          *uuid.UUID
	FromLocationID uuid.UUID
	ToLocationID   uuid.UUID
	Quantity       float64
	UnitID         uuid.UUID
	Reason         string
	TransferredBy  uuid.UUID
}

// Execute transfers stock between locations
func (uc *TransferStockUseCase) Execute(ctx context.Context, input *TransferStockInput) (string, error) {
	// Get source stock
	fromStock, err := uc.stockRepo.GetByLocationMaterialLot(ctx, input.FromLocationID, input.MaterialID, input.LotID)
	if err != nil {
		return "", err
	}

	// Check availability
	if fromStock.Quantity-fromStock.ReservedQty < input.Quantity {
		return "", entity.ErrInsufficientStock
	}

	// Deduct from source
	fromStock.Quantity -= input.Quantity

	// Create destination stock
	toStock := &entity.Stock{
		WarehouseID: fromStock.WarehouseID,
		ZoneID:      fromStock.ZoneID, // Will be updated based on location
		LocationID:  input.ToLocationID,
		MaterialID:  input.MaterialID,
		LotID:       input.LotID,
		Quantity:    input.Quantity,
		UnitID:      input.UnitID,
	}

	// Create movement
	movementNumber, _ := uc.stockRepo.GetNextMovementNumber(ctx, entity.MovementTypeTransfer)
	movement := entity.NewStockMovementTransfer(
		input.MaterialID,
		input.LotID,
		input.FromLocationID,
		input.ToLocationID,
		input.UnitID,
		input.TransferredBy,
		input.Quantity,
		movementNumber,
	)
	movement.Notes = input.Reason

	// Execute transfer
	if err := uc.stockRepo.TransferStock(ctx, fromStock, toStock, movement); err != nil {
		return "", err
	}

	return movementNumber, nil
}
