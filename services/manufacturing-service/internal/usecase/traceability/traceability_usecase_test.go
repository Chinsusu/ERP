package traceability_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/testmocks"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/traceability"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTraceabilityUseCase_Trace_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := new(testmocks.MockTraceabilityRepository)
	woRepo := new(testmocks.MockWorkOrderRepository)
	
	ucBackward := traceability.NewTraceBackwardUseCase(repo, woRepo)
	ucForward := traceability.NewTraceForwardUseCase(repo)

	productLotID := uuid.New()
	matLotID1 := uuid.New()

	backwardTraces := []*entity.BatchTraceability{
		{
			MaterialLotID:     matLotID1,
			MaterialQuantity:  50,
			MaterialLotNumber: "MAT-LOT-A",
			WorkOrderID:       uuid.New(),
			ProductLotNumber:  "PROD-LOT-001",
		},
	}

	forwardTraces := []*entity.BatchTraceability{
		{
			ProductLotID:      &productLotID,
			MaterialQuantity:  50,
			ProductLotNumber:  "PROD-LOT-001",
			MaterialLotNumber: "MAT-LOT-A",
		},
	}

	repo.On("GetByProductLot", ctx, productLotID).Return(backwardTraces, nil)
	repo.On("GetByMaterialLot", ctx, matLotID1).Return(forwardTraces, nil)
	woRepo.On("GetByID", ctx, backwardTraces[0].WorkOrderID).Return(&entity.WorkOrder{WONumber: "WO-001"}, nil)

	t.Run("Backward Trace Success", func(t *testing.T) {
		res, err := ucBackward.Execute(ctx, productLotID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "MAT-LOT-A", res.MaterialsUsed[0].LotNumber)
	})

	t.Run("Forward Trace Success", func(t *testing.T) {
		res, err := ucForward.Execute(ctx, matLotID1)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "PROD-LOT-001", res.UsedInProducts[0].ProductLot)
	})
}
