package issue

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// EventPublisher defines event publishing interface for issue
type EventPublisher interface {
	PublishStockIssued(event *event.StockIssuedEvent) error
}

// CreateGoodsIssueUseCase handles goods issue creation with FEFO
type CreateGoodsIssueUseCase struct {
	issueRepo repository.GoodsIssueRepository
	stockRepo repository.StockRepository
	eventPub  EventPublisher
}

// NewCreateGoodsIssueUseCase creates a new use case
func NewCreateGoodsIssueUseCase(
	issueRepo repository.GoodsIssueRepository,
	stockRepo repository.StockRepository,
	eventPub EventPublisher,
) *CreateGoodsIssueUseCase {
	return &CreateGoodsIssueUseCase{
		issueRepo: issueRepo,
		stockRepo: stockRepo,
		eventPub:  eventPub,
	}
}

// CreateGoodsIssueInput represents input for creating goods issue
type CreateGoodsIssueInput struct {
	IssueDate       time.Time
	IssueType       entity.IssueType
	ReferenceType   entity.ReferenceType
	ReferenceID     *uuid.UUID
	ReferenceNumber string
	WarehouseID     uuid.UUID
	Notes           string
	IssuedBy        uuid.UUID
	Items           []CreateGoodsIssueItemInput
}

// CreateGoodsIssueItemInput represents an item to issue
type CreateGoodsIssueItemInput struct {
	MaterialID uuid.UUID
	Quantity   float64
	UnitID     uuid.UUID
}

// CreateGoodsIssueOutput represents output from goods issue
type CreateGoodsIssueOutput struct {
	ID          uuid.UUID
	IssueNumber string
	Status      string
	LineItems   []GoodsIssueLineItemOutput
}

// GoodsIssueLineItemOutput represents issued line item with FEFO details
type GoodsIssueLineItemOutput struct {
	LineNumber  int
	MaterialID  uuid.UUID
	IssuedQty   float64
	LotsUsed    []entity.LotIssued
}

// Execute creates goods issue using FEFO logic
func (uc *CreateGoodsIssueUseCase) Execute(ctx context.Context, input *CreateGoodsIssueInput) (*CreateGoodsIssueOutput, error) {
	// Generate issue number
	issueNumber, err := uc.issueRepo.GetNextIssueNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Create goods issue
	issue := &entity.GoodsIssue{
		IssueNumber:     issueNumber,
		IssueDate:       input.IssueDate,
		IssueType:       input.IssueType,
		ReferenceType:   input.ReferenceType,
		ReferenceID:     input.ReferenceID,
		ReferenceNumber: input.ReferenceNumber,
		WarehouseID:     input.WarehouseID,
		Status:          entity.GoodsIssueStatusDraft,
		Notes:           input.Notes,
		IssuedBy:        &input.IssuedBy,
	}

	if err := uc.issueRepo.Create(ctx, issue); err != nil {
		return nil, err
	}

	output := &CreateGoodsIssueOutput{
		ID:          issue.ID,
		IssueNumber: issueNumber,
		Status:      string(issue.Status),
	}

	// Process each item with FEFO
	for i, item := range input.Items {
		// Issue using FEFO logic
		lotsIssued, err := uc.stockRepo.IssueStockFEFO(ctx, item.MaterialID, item.Quantity, input.IssuedBy)
		if err != nil {
			return nil, err
		}

		// Create line item for each lot used
		lineNumber := i + 1
		for _, lotIssued := range lotsIssued {
			lineItem := &entity.GILineItem{
				GoodsIssueID: issue.ID,
				LineNumber:   lineNumber,
				MaterialID:   item.MaterialID,
				RequestedQty: item.Quantity,
				IssuedQty:    lotIssued.Quantity,
				UnitID:       item.UnitID,
				LotID:        &lotIssued.LotID,
				LocationID:   &lotIssued.LocationID,
			}
			if err := uc.issueRepo.CreateLineItem(ctx, lineItem); err != nil {
				return nil, err
			}
		}

		// Create movement records
		movementNumber, _ := uc.stockRepo.GetNextMovementNumber(ctx, entity.MovementTypeOut)
		for _, lotIssued := range lotsIssued {
			movement := entity.NewStockMovementOut(
				item.MaterialID,
				&lotIssued.LotID,
				&lotIssued.LocationID,
				item.UnitID,
				input.IssuedBy,
				lotIssued.Quantity,
				entity.ReferenceTypeGI,
				&issue.ID,
				movementNumber,
			)
			uc.stockRepo.CreateMovement(ctx, movement)
		}

		output.LineItems = append(output.LineItems, GoodsIssueLineItemOutput{
			LineNumber: lineNumber,
			MaterialID: item.MaterialID,
			IssuedQty:  item.Quantity,
			LotsUsed:   lotsIssued,
		})

		// Publish event for each material issued
		lotsUsed := make([]event.LotUsedInIssue, len(lotsIssued))
		for j, lot := range lotsIssued {
			lotsUsed[j] = event.LotUsedInIssue{
				LotID:      lot.LotID.String(),
				LotNumber:  lot.LotNumber,
				Quantity:   lot.Quantity,
				ExpiryDate: lot.ExpiryDate.Format("2006-01-02"),
			}
		}

		refID := ""
		if input.ReferenceID != nil {
			refID = input.ReferenceID.String()
		}
		uc.eventPub.PublishStockIssued(&event.StockIssuedEvent{
			MaterialID:    item.MaterialID.String(),
			Quantity:      item.Quantity,
			LotsUsed:      lotsUsed,
			ReferenceType: string(input.ReferenceType),
			ReferenceID:   refID,
		})
	}

	// Mark issue as completed
	issue.Complete()
	uc.issueRepo.Update(ctx, issue)
	output.Status = string(entity.GoodsIssueStatusCompleted)

	return output, nil
}

// GetGoodsIssueUseCase handles getting goods issue
type GetGoodsIssueUseCase struct {
	issueRepo repository.GoodsIssueRepository
}

// NewGetGoodsIssueUseCase creates a new use case
func NewGetGoodsIssueUseCase(issueRepo repository.GoodsIssueRepository) *GetGoodsIssueUseCase {
	return &GetGoodsIssueUseCase{issueRepo: issueRepo}
}

// Execute gets goods issue by ID
func (uc *GetGoodsIssueUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.GoodsIssue, error) {
	return uc.issueRepo.GetByID(ctx, id)
}

// ListGoodsIssuesUseCase handles listing goods issues
type ListGoodsIssuesUseCase struct {
	issueRepo repository.GoodsIssueRepository
}

// NewListGoodsIssuesUseCase creates a new use case
func NewListGoodsIssuesUseCase(issueRepo repository.GoodsIssueRepository) *ListGoodsIssuesUseCase {
	return &ListGoodsIssuesUseCase{issueRepo: issueRepo}
}

// Execute lists goods issues
func (uc *ListGoodsIssuesUseCase) Execute(ctx context.Context, filter *repository.GoodsIssueFilter) ([]*entity.GoodsIssue, int64, error) {
	return uc.issueRepo.List(ctx, filter)
}
