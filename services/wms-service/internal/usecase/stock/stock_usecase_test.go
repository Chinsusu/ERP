package stock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/testmocks"
	"github.com/erp-cosmetics/wms-service/internal/testutils"
	"github.com/erp-cosmetics/wms-service/internal/usecase/stock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReserveStockUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := stock.NewReserveStockUseCase(stockRepo, eventPub)

	materialID := uuid.New()
	req := &stock.ReserveStockInput{
		MaterialID:      materialID,
		Quantity:        50,
		ReservationType: entity.ReservationTypeProduction,
		ReferenceID:     uuid.New(),
		ReferenceNumber: "WO-001",
		CreatedBy:       uuid.New(),
	}

	stockRepo.On("ReserveStock", ctx, materialID, 50.0, mock.AnythingOfType("*entity.StockReservation")).Return(nil)
	eventPub.On("PublishStockReserved", mock.Anything).Return(nil)

	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, entity.ReservationStatusActive, res.Status)
	stockRepo.AssertExpectations(t)
}

func TestReserveStockUseCase_Execute_InsufficientStock(t *testing.T) {
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	
	uc := stock.NewReserveStockUseCase(stockRepo, nil)

	materialID := uuid.New()
	req := &stock.ReserveStockInput{
		MaterialID: materialID,
		Quantity:   1000,
	}

	stockRepo.On("ReserveStock", ctx, materialID, 1000.0, mock.Anything).Return(errors.New("insufficient stock"))

	res, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestReleaseReservationUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	uc := stock.NewReleaseReservationUseCase(stockRepo)

	reservationID := uuid.New()
	stockRepo.On("ReleaseReservation", ctx, reservationID).Return(nil)

	err := uc.Execute(ctx, reservationID)

	assert.NoError(t, err)
	stockRepo.AssertExpectations(t)
}

func TestGetStockUseCase_GetExpiringStock(t *testing.T) {
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	uc := stock.NewGetStockUseCase(stockRepo)

	materialID := uuid.New()
	lot := testutils.NewLotBuilder().WithMaterialID(materialID).Build()
	stocks := []*entity.Stock{
		testutils.NewStockBuilder().WithLot(lot).Build(),
	}

	stockRepo.On("GetExpiringStock", ctx, 30).Return(stocks, nil)

	res, err := uc.GetExpiringStock(ctx, 30)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, lot.ID, *res[0].LotID)
}
