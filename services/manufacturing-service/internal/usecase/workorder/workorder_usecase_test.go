package workorder_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/testmocks"
	"github.com/erp-cosmetics/manufacturing-service/internal/testutils"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/workorder"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWOUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	woRepo := new(testmocks.MockWorkOrderRepository)
	bomRepo := new(testmocks.MockBOMRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := workorder.NewCreateWOUseCase(woRepo, bomRepo, eventPub)

	bom := testutils.NewBOMBuilder().
		WithItems([]entity.BOMLineItem{
			{ID: uuid.New(), MaterialID: uuid.New(), Quantity: 10, UnitCost: 5},
		}).
		Build()

	input := workorder.CreateWOInput{
		ProductID:       bom.ProductID,
		BOMID:           bom.ID,
		PlannedQuantity: 200, // 2x batch size
		PlannedStartDate: "2026-01-26",
		CreatedBy:       uuid.New(),
	}

	bomRepo.On("GetByID", ctx, bom.ID).Return(bom, nil)
	woRepo.On("GenerateWONumber", ctx).Return("WO-2026-0001", nil)
	woRepo.On("Create", ctx, mock.AnythingOfType("*entity.WorkOrder")).Return(nil)
	woRepo.On("CreateLineItems", ctx, mock.Anything).Return(nil)
	eventPub.On("PublishWOCreated", mock.Anything).Return(nil)

	// Act
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, entity.WOStatusPlanned, res.Status)
	assert.Equal(t, 200.0, res.PlannedQuantity)
	
	woRepo.AssertExpectations(t)
	bomRepo.AssertExpectations(t)
	eventPub.AssertExpectations(t)
}

func TestStartWOUseCase_Execute_TriggersWMS(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := new(testmocks.MockWorkOrderRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := workorder.NewStartWOUseCase(repo, eventPub)
	woID := uuid.New()
	targetWO := testutils.NewWorkOrderBuilder().WithStatus(entity.WOStatusReleased).Build()
	
	repo.On("GetByID", ctx, woID).Return(targetWO, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)
	eventPub.On("PublishWOStarted", mock.Anything).Return(nil)

	// Act
	res, err := uc.Execute(ctx, woID, uuid.New())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, entity.WOStatusInProgress, res.Status)
	eventPub.AssertCalled(t, "PublishWOStarted", mock.Anything)
}
