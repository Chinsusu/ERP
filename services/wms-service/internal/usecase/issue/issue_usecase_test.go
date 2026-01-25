package issue_test

import (
	"context"
	"testing"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/testmocks"
	"github.com/erp-cosmetics/wms-service/internal/usecase/issue"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGoodsIssueUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	issueRepo := new(testmocks.MockGoodsIssueRepository)
	stockRepo := new(testmocks.MockStockRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := issue.NewCreateGoodsIssueUseCase(issueRepo, stockRepo, eventPub)

	materialID := uuid.New()
	warehouseID := uuid.New()
	userID := uuid.New()
	unitID := uuid.New()
	lotID := uuid.New()
	locationID := uuid.New()

	input := &issue.CreateGoodsIssueInput{
		IssueDate:   time.Now(),
		IssueType:   entity.IssueTypeSales,
		WarehouseID: warehouseID,
		IssuedBy:    userID,
		Items: []issue.CreateGoodsIssueItemInput{
			{
				MaterialID: materialID,
				Quantity:   50,
				UnitID:     unitID,
			},
		},
	}

	issueRepo.On("GetNextIssueNumber", ctx).Return("GI-2026-00001", nil)
	issueRepo.On("Create", ctx, mock.AnythingOfType("*entity.GoodsIssue")).Return(nil)
	issueRepo.On("CreateLineItem", ctx, mock.AnythingOfType("*entity.GILineItem")).Return(nil)
	issueRepo.On("Update", ctx, mock.AnythingOfType("*entity.GoodsIssue")).Return(nil)

	stockRepo.On("IssueStockFEFO", ctx, materialID, 50.0, userID).Return([]entity.LotIssued{
		{
			LotID:      lotID,
			LotNumber:  "LOT-001",
			Quantity:   50,
			ExpiryDate: time.Now().AddDate(1, 0, 0),
			LocationID: locationID,
		},
	}, nil)
	stockRepo.On("GetNextMovementNumber", ctx, entity.MovementTypeOut).Return("MOV-OUT-001", nil)
	stockRepo.On("CreateMovement", ctx, mock.AnythingOfType("*entity.StockMovement")).Return(nil)

	eventPub.On("PublishStockIssued", mock.AnythingOfType("*event.StockIssuedEvent")).Return(nil)

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "GI-2026-00001", output.IssueNumber)
	assert.Equal(t, string(entity.GoodsIssueStatusCompleted), output.Status)
	assert.Len(t, output.LineItems, 1)
	assert.Equal(t, 50.0, output.LineItems[0].IssuedQty)

	issueRepo.AssertExpectations(t)
	stockRepo.AssertExpectations(t)
	eventPub.AssertExpectations(t)
}

func TestCreateGoodsIssueUseCase_Execute_InsufficientStock(t *testing.T) {
	// Arrange
	ctx := context.Background()
	issueRepo := new(testmocks.MockGoodsIssueRepository)
	stockRepo := new(testmocks.MockStockRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := issue.NewCreateGoodsIssueUseCase(issueRepo, stockRepo, eventPub)

	materialID := uuid.New()

	input := &issue.CreateGoodsIssueInput{
		Items: []issue.CreateGoodsIssueItemInput{
			{
				MaterialID: materialID,
				Quantity:   100,
			},
		},
	}

	issueRepo.On("GetNextIssueNumber", ctx).Return("GI-2026-00002", nil)
	issueRepo.On("Create", ctx, mock.Anything).Return(nil)
	stockRepo.On("IssueStockFEFO", ctx, materialID, 100.0, mock.Anything).Return(nil, entity.ErrInsufficientStock)

	// Act
	output, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInsufficientStock, err)
}
