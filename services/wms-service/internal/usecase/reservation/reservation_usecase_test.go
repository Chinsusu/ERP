package reservation_test

import (
	"context"
	"testing"
	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/testmocks"
	"github.com/erp-cosmetics/wms-service/internal/usecase/reservation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateReservationUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := reservation.NewCreateReservationUseCase(stockRepo, eventPub)

	materialID := uuid.New()
	refID := uuid.New()

	input := &reservation.CreateReservationInput{
		MaterialID:      materialID,
		Quantity:        50,
		ReservationType: entity.ReservationTypeSalesOrder,
		ReferenceID:     refID,
	}

	stockRepo.On("ReserveStock", ctx, materialID, 50.0, mock.AnythingOfType("*entity.StockReservation")).Return(nil)
	eventPub.On("PublishStockReserved", mock.AnythingOfType("*event.StockReservedEvent")).Return(nil)

	// Act
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, entity.ReservationStatusActive, res.Status)
	assert.Equal(t, 50.0, res.Quantity)

	stockRepo.AssertExpectations(t)
	eventPub.AssertExpectations(t)
}

func TestCreateReservationUseCase_Execute_InsufficientStock(t *testing.T) {
	// Arrange
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := reservation.NewCreateReservationUseCase(stockRepo, eventPub)

	materialID := uuid.New()

	input := &reservation.CreateReservationInput{
		MaterialID: materialID,
		Quantity:   1000,
	}

	stockRepo.On("ReserveStock", ctx, materialID, 1000.0, mock.Anything).Return(entity.ErrInsufficientStock)

	// Act
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, entity.ErrInsufficientStock, err)
}

func TestReleaseReservationUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	uc := reservation.NewReleaseReservationUseCase(stockRepo)

	resID := uuid.New()
	stockRepo.On("ReleaseReservation", ctx, resID).Return(nil)

	// Act
	err := uc.Execute(ctx, resID)

	// Assert
	assert.NoError(t, err)
	stockRepo.AssertExpectations(t)
}
