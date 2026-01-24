package stock

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// GetStockUseCase handles stock queries
type GetStockUseCase struct {
	stockRepo repository.StockRepository
}

// NewGetStockUseCase creates a new use case
func NewGetStockUseCase(stockRepo repository.StockRepository) *GetStockUseCase {
	return &GetStockUseCase{stockRepo: stockRepo}
}

// Execute returns stock list
func (uc *GetStockUseCase) Execute(ctx context.Context, filter *repository.StockFilter) ([]*entity.Stock, int64, error) {
	return uc.stockRepo.List(ctx, filter)
}

// GetByMaterial returns stock by material
func (uc *GetStockUseCase) GetByMaterial(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error) {
	return uc.stockRepo.GetByMaterial(ctx, materialID)
}

// GetSummary returns stock summary for a material
func (uc *GetStockUseCase) GetSummary(ctx context.Context, materialID uuid.UUID) (*entity.StockSummary, error) {
	return uc.stockRepo.GetMaterialSummary(ctx, materialID)
}

// GetExpiringStock returns stock expiring within days
func (uc *GetStockUseCase) GetExpiringStock(ctx context.Context, days int) ([]*entity.Stock, error) {
	return uc.stockRepo.GetExpiringStock(ctx, days)
}

// GetLowStock returns materials with low stock
func (uc *GetStockUseCase) GetLowStock(ctx context.Context, threshold float64) ([]*entity.StockSummary, error) {
	return uc.stockRepo.GetLowStockMaterials(ctx, threshold)
}

// IssueStockFEFOUseCase handles issuing stock with FEFO logic
type IssueStockFEFOUseCase struct {
	stockRepo   repository.StockRepository
	eventPub    *event.Publisher
}

// NewIssueStockFEFOUseCase creates a new use case
func NewIssueStockFEFOUseCase(stockRepo repository.StockRepository, eventPub *event.Publisher) *IssueStockFEFOUseCase {
	return &IssueStockFEFOUseCase{
		stockRepo:  stockRepo,
		eventPub:   eventPub,
	}
}

// IssueStockInput represents input for issuing stock
type IssueStockInput struct {
	MaterialID      uuid.UUID
	Quantity        float64
	UnitID          uuid.UUID
	ReferenceType   entity.ReferenceType
	ReferenceID     *uuid.UUID
	ReferenceNumber string
	CreatedBy       uuid.UUID
}

// IssueStockOutput represents output from issuing stock
type IssueStockOutput struct {
	MovementNumber string
	LotsIssued     []entity.LotIssued
}

// Execute issues stock using FEFO logic
func (uc *IssueStockFEFOUseCase) Execute(ctx context.Context, input *IssueStockInput) (*IssueStockOutput, error) {
	// Issue using FEFO (First Expired First Out)
	lotsIssued, err := uc.stockRepo.IssueStockFEFO(ctx, input.MaterialID, input.Quantity, input.CreatedBy)
	if err != nil {
		return nil, err
	}

	// Get movement number
	movementNumber, err := uc.stockRepo.GetNextMovementNumber(ctx, entity.MovementTypeOut)
	if err != nil {
		return nil, err
	}

	// Create movement records for each lot
	for _, lot := range lotsIssued {
		movement := entity.NewStockMovementOut(
			input.MaterialID,
			&lot.LotID,
			&lot.LocationID,
			input.UnitID,
			input.CreatedBy,
			lot.Quantity,
			input.ReferenceType,
			input.ReferenceID,
			movementNumber,
		)
		if err := uc.stockRepo.CreateMovement(ctx, movement); err != nil {
			return nil, err
		}
	}

	// Publish event
	lotsUsed := make([]event.LotUsedInIssue, len(lotsIssued))
	for i, lot := range lotsIssued {
		lotsUsed[i] = event.LotUsedInIssue{
			LotID:      lot.LotID.String(),
			LotNumber:  lot.LotNumber,
			Quantity:   lot.Quantity,
			ExpiryDate: lot.ExpiryDate.Format("2006-01-02"),
		}
	}

	refID := ""
	if input.ReferenceID != nil {
		refID = input.ReferenceID.String()
	}
	
	uc.eventPub.PublishStockIssued(&event.StockIssuedEvent{
		MaterialID:    input.MaterialID.String(),
		Quantity:      input.Quantity,
		LotsUsed:      lotsUsed,
		ReferenceType: string(input.ReferenceType),
		ReferenceID:   refID,
	})

	return &IssueStockOutput{
		MovementNumber: movementNumber,
		LotsIssued:     lotsIssued,
	}, nil
}

// ReserveStockUseCase handles stock reservation
type ReserveStockUseCase struct {
	stockRepo repository.StockRepository
	eventPub  *event.Publisher
}

// NewReserveStockUseCase creates a new use case
func NewReserveStockUseCase(stockRepo repository.StockRepository, eventPub *event.Publisher) *ReserveStockUseCase {
	return &ReserveStockUseCase{
		stockRepo: stockRepo,
		eventPub:  eventPub,
	}
}

// ReserveStockInput represents input for reserving stock
type ReserveStockInput struct {
	MaterialID      uuid.UUID
	Quantity        float64
	UnitID          uuid.UUID
	ReservationType entity.ReservationType
	ReferenceID     uuid.UUID
	ReferenceNumber string
	CreatedBy       uuid.UUID
}

// Execute reserves stock
func (uc *ReserveStockUseCase) Execute(ctx context.Context, input *ReserveStockInput) (*entity.StockReservation, error) {
	reservation := &entity.StockReservation{
		MaterialID:      input.MaterialID,
		Quantity:        input.Quantity,
		UnitID:          input.UnitID,
		ReservationType: input.ReservationType,
		ReferenceID:     input.ReferenceID,
		ReferenceNumber: input.ReferenceNumber,
		Status:          entity.ReservationStatusActive,
		CreatedBy:       input.CreatedBy,
	}

	if err := uc.stockRepo.ReserveStock(ctx, input.MaterialID, input.Quantity, reservation); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishStockReserved(&event.StockReservedEvent{
		ReservationID:   reservation.ID.String(),
		MaterialID:      input.MaterialID.String(),
		Quantity:        input.Quantity,
		ReservationType: string(input.ReservationType),
		ReferenceID:     input.ReferenceID.String(),
	})

	return reservation, nil
}

// ReleaseReservationUseCase handles releasing reservations
type ReleaseReservationUseCase struct {
	stockRepo repository.StockRepository
}

// NewReleaseReservationUseCase creates a new use case
func NewReleaseReservationUseCase(stockRepo repository.StockRepository) *ReleaseReservationUseCase {
	return &ReleaseReservationUseCase{stockRepo: stockRepo}
}

// Execute releases a reservation
func (uc *ReleaseReservationUseCase) Execute(ctx context.Context, reservationID uuid.UUID) error {
	return uc.stockRepo.ReleaseReservation(ctx, reservationID)
}
