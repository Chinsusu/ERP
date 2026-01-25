package bom_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/testmocks"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/bom"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBOMUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := new(testmocks.MockBOMRepository)
	eventPub := new(testmocks.MockEventPublisher)
	key := []byte("thisis32bytekeyforaesgcmtesting!")
	
	uc := bom.NewCreateBOMUseCase(repo, eventPub, key)

	input := bom.CreateBOMInput{
		BOMNumber:   "BOM-001",
		ProductID:   uuid.New(),
		Name:        "Face Cream",
		BatchSize:   100,
		FormulaDetails: &entity.FormulaDetails{
			Notes: "Sensitive detail",
		},
		CreatedBy: uuid.New(),
	}

	repo.On("Create", ctx, mock.AnythingOfType("*entity.BOM")).Return(nil).Run(func(args mock.Arguments) {
		b := args.Get(1).(*entity.BOM)
		assert.NotEmpty(t, b.FormulaDetails) // Must be encrypted
	})
	eventPub.On("PublishBOMCreated", mock.Anything).Return(nil)

	// Act
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	repo.AssertExpectations(t)
}

func TestGetBOMUseCase_Execute_Permissions(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := new(testmocks.MockBOMRepository)
	key := []byte("thisis32bytekeyforaesgcmtesting!")
	
	uc := bom.NewGetBOMUseCase(repo, key)
	bomID := uuid.New()
	
	formula := &entity.FormulaDetails{Notes: "Confidential"}
	encrypted, _ := entity.EncryptFormula(formula, key)
	
	targetBOM := &entity.BOM{
		ID:             bomID,
		FormulaDetails: encrypted,
	}

	repo.On("GetByID", ctx, bomID).Return(targetBOM, nil)

	t.Run("View with permission shows formula", func(t *testing.T) {
		res, err := uc.Execute(ctx, bomID, true)
		assert.NoError(t, err)
		assert.NotNil(t, res.FormulaDetails)
		assert.Equal(t, "Confidential", res.FormulaDetails.Notes)
	})

	t.Run("View without permission hides formula", func(t *testing.T) {
		res, err := uc.Execute(ctx, bomID, false)
		assert.NoError(t, err)
		assert.Nil(t, res.FormulaDetails)
	})
}
