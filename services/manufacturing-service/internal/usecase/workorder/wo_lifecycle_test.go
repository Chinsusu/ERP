package workorder_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/testmocks"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/workorder"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkOrder_Lifecycle_And_Yield(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := new(testmocks.MockWorkOrderRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	woID := uuid.New()
	userID := uuid.New()
	
	// Initial state: PLANNED
	wo := &entity.WorkOrder{
		ID:              woID,
		Status:          entity.WOStatusPlanned,
		PlannedQuantity: 100,
	}

	repo.On("GetByID", ctx, woID).Return(wo, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)
	eventPub.On("PublishWOReleased", mock.Anything).Return(nil)
	eventPub.On("PublishWOStarted", mock.Anything).Return(nil)
	eventPub.On("PublishWOCompleted", mock.Anything).Return(nil)

	// 1. RELEASE
	releaseUC := workorder.NewReleaseWOUseCase(repo, eventPub)
	res, err := releaseUC.Execute(ctx, woID, userID)
	assert.NoError(t, err)
	assert.Equal(t, entity.WOStatusReleased, res.Status)

	// 2. START
	startUC := workorder.NewStartWOUseCase(repo, eventPub)
	res, err = startUC.Execute(ctx, woID, userID)
	assert.NoError(t, err)
	assert.Equal(t, entity.WOStatusInProgress, res.Status)
	assert.NotNil(t, res.ActualStartDate)

	// 3. COMPLETE with YIELD
	completeUC := workorder.NewCompleteWOUseCase(repo, eventPub)
	res, err = completeUC.Execute(ctx, workorder.CompleteWOInput{
		WorkOrderID:      woID,
		ActualQuantity:   105,
		GoodQuantity:     98,
		RejectedQuantity: 7,
		CompletedBy:      userID,
	})
	
	assert.NoError(t, err)
	assert.Equal(t, entity.WOStatusCompleted, res.Status)
	assert.Equal(t, 98.0, *res.YieldPercentage) // 98/100 * 100
	assert.NotNil(t, res.ActualEndDate)
}

func TestWorkOrder_Cancel_ReleaseReservations(t *testing.T) {
	ctx := context.Background()
	repo := new(testmocks.MockWorkOrderRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := workorder.NewCancelWOUseCase(repo, eventPub)
	woID := uuid.New()
	
	wo := &entity.WorkOrder{
		ID:     woID,
		Status: entity.WOStatusReleased,
	}

	repo.On("GetByID", ctx, woID).Return(wo, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)
	eventPub.On("PublishWOCancelled", mock.Anything).Return(nil)

	res, err := uc.Execute(ctx, woID, uuid.New())

	assert.NoError(t, err)
	assert.Equal(t, entity.WOStatusCancelled, res.Status)
	eventPub.AssertCalled(t, "PublishWOCancelled", mock.Anything)
}
