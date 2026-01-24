package reservation

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateReservationUseCase handles stock reservation
type CreateReservationUseCase struct {
	stockRepo repository.StockRepository
	eventPub  *event.Publisher
}

// NewCreateReservationUseCase creates a new use case
func NewCreateReservationUseCase(stockRepo repository.StockRepository, eventPub *event.Publisher) *CreateReservationUseCase {
	return &CreateReservationUseCase{
		stockRepo: stockRepo,
		eventPub:  eventPub,
	}
}

// CreateReservationInput represents input for creating reservation
type CreateReservationInput struct {
	MaterialID      uuid.UUID
	Quantity        float64
	UnitID          uuid.UUID
	ReservationType entity.ReservationType
	ReferenceID     uuid.UUID
	ReferenceNumber string
	ExpiresAt       *time.Time
	CreatedBy       uuid.UUID
}

// Execute creates a stock reservation
func (uc *CreateReservationUseCase) Execute(ctx context.Context, input *CreateReservationInput) (*entity.StockReservation, error) {
	reservation := &entity.StockReservation{
		MaterialID:      input.MaterialID,
		Quantity:        input.Quantity,
		UnitID:          input.UnitID,
		ReservationType: input.ReservationType,
		ReferenceID:     input.ReferenceID,
		ReferenceNumber: input.ReferenceNumber,
		Status:          entity.ReservationStatusActive,
		ExpiresAt:       input.ExpiresAt,
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

// CheckAvailabilityUseCase checks stock availability
type CheckAvailabilityUseCase struct {
	stockRepo repository.StockRepository
}

// NewCheckAvailabilityUseCase creates a new use case
func NewCheckAvailabilityUseCase(stockRepo repository.StockRepository) *CheckAvailabilityUseCase {
	return &CheckAvailabilityUseCase{stockRepo: stockRepo}
}

// AvailabilityResult represents stock availability result
type AvailabilityResult struct {
	MaterialID     uuid.UUID
	TotalQuantity  float64
	ReservedQty    float64
	AvailableQty   float64
	IsAvailable    bool
	RequestedQty   float64
	ShortageQty    float64
	ExpiringInDays int
}

// Execute checks if material is available
func (uc *CheckAvailabilityUseCase) Execute(ctx context.Context, materialID uuid.UUID, requestedQty float64) (*AvailabilityResult, error) {
	summary, err := uc.stockRepo.GetMaterialSummary(ctx, materialID)
	if err != nil {
		return nil, err
	}

	result := &AvailabilityResult{
		MaterialID:    materialID,
		TotalQuantity: summary.TotalQuantity,
		ReservedQty:   summary.TotalReserved,
		AvailableQty:  summary.TotalAvailable,
		RequestedQty:  requestedQty,
		IsAvailable:   summary.TotalAvailable >= requestedQty,
	}

	if !result.IsAvailable {
		result.ShortageQty = requestedQty - summary.TotalAvailable
	}

	return result, nil
}
