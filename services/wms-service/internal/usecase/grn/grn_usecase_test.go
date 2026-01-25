package grn_test

import (
	"context"
	"testing"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/testmocks"
	"github.com/erp-cosmetics/wms-service/internal/usecase/grn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGRNUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	grnRepo := new(testmocks.MockGRNRepository)
	lotRepo := new(testmocks.MockLotRepository)
	stockRepo := new(testmocks.MockStockRepository)
	zoneRepo := new(testmocks.MockZoneRepository)
	locationRepo := new(testmocks.MockLocationRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := grn.NewCreateGRNUseCase(grnRepo, lotRepo, stockRepo, zoneRepo, locationRepo, eventPub)

	materialID := uuid.New()
	warehouseID := uuid.New()
	unitID := uuid.New()

	input := &grn.CreateGRNInput{
		GRNDate:     time.Now(),
		WarehouseID: warehouseID,
		ReceivedBy:  uuid.New(),
		Items: []grn.CreateGRNItemInput{
			{
				MaterialID: materialID,
				ReceivedQty: 100,
				UnitID:      unitID,
				ExpiryDate:  time.Now().AddDate(1, 0, 0),
			},
		},
	}

	grnRepo.On("GetNextGRNNumber", ctx).Return("GRN-2026-00001", nil)
	zoneRepo.On("GetQuarantineZone", ctx, warehouseID).Return(&entity.Zone{ID: uuid.New()}, nil)
	locationRepo.On("GetByZoneID", ctx, mock.Anything).Return([]*entity.Location{{ID: uuid.New()}}, nil)
	grnRepo.On("Create", ctx, mock.AnythingOfType("*entity.GRN")).Return(nil)
	lotRepo.On("GetNextLotNumber", ctx).Return("LOT-2026-00001", nil)
	lotRepo.On("Create", ctx, mock.AnythingOfType("*entity.Lot")).Return(nil)
	grnRepo.On("CreateLineItem", ctx, mock.AnythingOfType("*entity.GRNLineItem")).Return(nil)
	eventPub.On("PublishGRNCreated", mock.AnythingOfType("*event.GRNCreatedEvent")).Return(nil)

	// Act
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "GRN-2026-00001", res.GRNNumber)
	assert.Equal(t, entity.GRNStatusDraft, res.Status)

	grnRepo.AssertExpectations(t)
	lotRepo.AssertExpectations(t)
	eventPub.AssertExpectations(t)
}

func TestCompleteGRNUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	grnRepo := new(testmocks.MockGRNRepository)
	lotRepo := new(testmocks.MockLotRepository)
	stockRepo := new(testmocks.MockStockRepository)
	zoneRepo := new(testmocks.MockZoneRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := grn.NewCompleteGRNUseCase(grnRepo, lotRepo, stockRepo, zoneRepo, eventPub)

	grnID := uuid.New()
	materialID := uuid.New()
	lotID := uuid.New()
	locationID := uuid.New()
	userID := uuid.New()

	targetGRN := &entity.GRN{
		ID:         grnID,
		Status:     entity.GRNStatusDraft,
		ReceivedBy: &userID,
		WarehouseID: uuid.New(),
	}

	lineItems := []*entity.GRNLineItem{
		{
			GRNID:       grnID,
			MaterialID:  materialID,
			LotID:       &lotID,
			ReceivedQty: 100,
			LocationID:  &locationID,
		},
	}

	grnRepo.On("GetByID", ctx, grnID).Return(targetGRN, nil)
	grnRepo.On("GetLineItemsByGRNID", ctx, grnID).Return(lineItems, nil)
	lotRepo.On("GetByID", ctx, lotID).Return(&entity.Lot{ID: lotID, LotNumber: "LOT-001"}, nil)
	lotRepo.On("Update", ctx, mock.Anything).Return(nil)
	grnRepo.On("UpdateLineItem", ctx, mock.Anything).Return(nil)
	
	// Mock location for stock creation
	stockRepo.On("GetByID", ctx, locationID).Return(&entity.Stock{LocationID: locationID, ZoneID: uuid.New()}, nil)
	stockRepo.On("GetNextMovementNumber", ctx, entity.MovementTypeIn).Return("MOV-IN-001", nil)
	stockRepo.On("ReceiveStock", ctx, mock.Anything, mock.Anything).Return(nil)
	
	grnRepo.On("Update", ctx, mock.Anything).Return(nil)
	
	eventPub.On("PublishStockReceived", mock.Anything).Return(nil)
	eventPub.On("PublishGRNCompleted", mock.Anything).Return(nil)

	// Act
	res, err := uc.Execute(ctx, &grn.CompleteGRNInput{
		GRNID:    grnID,
		QCStatus: entity.QCStatusPassed,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, entity.GRNStatusCompleted, res.Status)

	grnRepo.AssertExpectations(t)
	lotRepo.AssertExpectations(t)
	stockRepo.AssertExpectations(t)
	eventPub.AssertExpectations(t)
}

func TestCompleteGRNUseCase_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	grnRepo := new(testmocks.MockGRNRepository)
	uc := grn.NewCompleteGRNUseCase(grnRepo, nil, nil, nil, nil)

	grnID := uuid.New()
	grnRepo.On("GetByID", ctx, grnID).Return(nil, errors.New("not found"))

	res, err := uc.Execute(ctx, &grn.CompleteGRNInput{GRNID: grnID})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "not found")
}

func TestCompleteGRNUseCase_Execute_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	grnRepo := new(testmocks.MockGRNRepository)
	uc := grn.NewCompleteGRNUseCase(grnRepo, nil, nil, nil, nil)

	grnID := uuid.New()
	targetGRN := &entity.GRN{
		ID:     grnID,
		Status: entity.GRNStatusCompleted,
	}

	grnRepo.On("GetByID", ctx, grnID).Return(targetGRN, nil)

	res, err := uc.Execute(ctx, &grn.CompleteGRNInput{GRNID: grnID})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "already completed")
}
