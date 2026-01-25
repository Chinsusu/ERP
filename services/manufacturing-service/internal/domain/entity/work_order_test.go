package entity_test

import (
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWorkOrder_Lifecycle(t *testing.T) {
	wo := &entity.WorkOrder{
		Status:          entity.WOStatusPlanned,
		PlannedQuantity: 100,
	}

	t.Run("Planned to Released", func(t *testing.T) {
		err := wo.Release()
		assert.NoError(t, err)
		assert.Equal(t, entity.WOStatusReleased, wo.Status)
	})

	t.Run("Released to In Progress", func(t *testing.T) {
		supervisorID := uuid.New()
		err := wo.Start(supervisorID)
		assert.NoError(t, err)
		assert.Equal(t, entity.WOStatusInProgress, wo.Status)
		assert.Equal(t, &supervisorID, wo.SupervisorID)
		assert.NotNil(t, wo.ActualStartDate)
	})

	t.Run("In Progress to Completed with Yield", func(t *testing.T) {
		err := wo.Complete(100, 95, 5)
		assert.NoError(t, err)
		assert.Equal(t, entity.WOStatusCompleted, wo.Status)
		assert.Equal(t, 95.0, *wo.YieldPercentage)
		assert.NotNil(t, wo.ActualEndDate)
	})

	t.Run("Start from Planned Fails", func(t *testing.T) {
		wo.Status = entity.WOStatusPlanned
		err := wo.Start(uuid.New())
		assert.Error(t, err)
	})
}

func TestWorkOrder_Cancellation(t *testing.T) {
	t.Run("Cancel Planned Success", func(t *testing.T) {
		wo := &entity.WorkOrder{Status: entity.WOStatusPlanned}
		err := wo.Cancel()
		assert.NoError(t, err)
		assert.Equal(t, entity.WOStatusCancelled, wo.Status)
	})

	t.Run("Cancel InProgress Fails", func(t *testing.T) {
		wo := &entity.WorkOrder{Status: entity.WOStatusInProgress}
		err := wo.Cancel()
		assert.Error(t, err)
	})
}
