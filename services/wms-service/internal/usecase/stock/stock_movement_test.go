package stock_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/testmocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStockMovement_Auditing(t *testing.T) {
	// Arrange
	ctx := context.Background()
	stockRepo := new(testmocks.MockStockRepository)
	
	materialID := uuid.New()
	lotID := uuid.New()
	locationID := uuid.New()

	movement := entity.NewStockMovementIn(
		materialID,
		lotID,
		locationID,
		uuid.New(),
		uuid.New(),
		100,
		entity.ReferenceTypeGRN,
		nil,
		"MOV-UT-001",
	)

	stockRepo.On("CreateMovement", ctx, movement).Return(nil)

	// Act
	err := stockRepo.CreateMovement(ctx, movement)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, entity.MovementTypeIn, movement.MovementType)
	assert.Equal(t, 100.0, movement.Quantity)
	stockRepo.AssertExpectations(t)
}
