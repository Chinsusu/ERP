package pr

import (
	"context"
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreatePRRequest represents request to create PR
type CreatePRRequest struct {
	RequiredDate  string           `json:"required_date" validate:"required"`
	Priority      string           `json:"priority"`
	Justification string           `json:"justification"`
	Notes         string           `json:"notes"`
	RequesterID   uuid.UUID        `json:"-"`
	DepartmentID  *uuid.UUID       `json:"-"`
	Items         []PRLineItemRequest `json:"items" validate:"required,min=1"`
}

// PRLineItemRequest represents a line item in create request
type PRLineItemRequest struct {
	MaterialID   uuid.UUID `json:"material_id" validate:"required"`
	MaterialCode string    `json:"material_code"`
	MaterialName string    `json:"material_name"`
	Quantity     float64   `json:"quantity" validate:"required,gt=0"`
	UOMCode      string    `json:"uom_code"`
	UnitPrice    float64   `json:"unit_price"`
	Specifications string  `json:"specifications"`
}

// CreatePRUseCase handles creating a PR
type CreatePRUseCase struct {
	prRepo   repository.PRRepository
	eventPub *event.Publisher
}

// NewCreatePRUseCase creates a new use case
func NewCreatePRUseCase(prRepo repository.PRRepository, eventPub *event.Publisher) *CreatePRUseCase {
	return &CreatePRUseCase{prRepo: prRepo, eventPub: eventPub}
}

// Execute creates a new PR
func (uc *CreatePRUseCase) Execute(ctx context.Context, req *CreatePRRequest) (*entity.PurchaseRequisition, error) {
	// Generate PR number
	prNumber, err := uc.prRepo.GetNextPRNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Parse required date
	requiredDate, err := time.Parse("2006-01-02", req.RequiredDate)
	if err != nil {
		requiredDate = time.Now().AddDate(0, 0, 14) // Default 2 weeks
	}

	priority := entity.PRPriority(req.Priority)
	if priority == "" {
		priority = entity.PRPriorityNormal
	}

	pr := &entity.PurchaseRequisition{
		ID:            uuid.New(),
		PRNumber:      prNumber,
		PRDate:        time.Now(),
		RequiredDate:  requiredDate,
		Priority:      priority,
		Status:        entity.PRStatusDraft,
		RequesterID:   req.RequesterID,
		DepartmentID:  req.DepartmentID,
		Justification: req.Justification,
		Notes:         req.Notes,
		Currency:      "VND",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create PR first
	if err := uc.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	// Create line items
	for i, item := range req.Items {
		lineItem := &entity.PRLineItem{
			ID:             uuid.New(),
			PRID:           pr.ID,
			LineNumber:     i + 1,
			MaterialID:     item.MaterialID,
			MaterialCode:   item.MaterialCode,
			MaterialName:   item.MaterialName,
			Quantity:       item.Quantity,
			UOMCode:        item.UOMCode,
			UnitPrice:      item.UnitPrice,
			Currency:       "VND",
			Specifications: item.Specifications,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		lineItem.CalculateLineTotal()
		
		if err := uc.prRepo.CreateLineItem(ctx, lineItem); err != nil {
			return nil, err
		}
		pr.LineItems = append(pr.LineItems, *lineItem)
	}

	// Calculate total
	pr.CalculateTotalAmount()
	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishPRCreated(ctx, &event.PREvent{
		PRID:        pr.ID.String(),
		PRNumber:    pr.PRNumber,
		Status:      string(pr.Status),
		TotalAmount: pr.TotalAmount,
		RequesterID: pr.RequesterID.String(),
		Timestamp:   time.Now(),
	})

	return pr, nil
}

// GetPRUseCase handles getting a PR
type GetPRUseCase struct {
	prRepo repository.PRRepository
}

// NewGetPRUseCase creates a new use case
func NewGetPRUseCase(prRepo repository.PRRepository) *GetPRUseCase {
	return &GetPRUseCase{prRepo: prRepo}
}

// Execute gets a PR by ID
func (uc *GetPRUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.PurchaseRequisition, error) {
	return uc.prRepo.GetByID(ctx, id)
}

// ListPRsUseCase handles listing PRs
type ListPRsUseCase struct {
	prRepo repository.PRRepository
}

// NewListPRsUseCase creates a new use case
func NewListPRsUseCase(prRepo repository.PRRepository) *ListPRsUseCase {
	return &ListPRsUseCase{prRepo: prRepo}
}

// Execute lists PRs with filters
func (uc *ListPRsUseCase) Execute(ctx context.Context, filter *repository.PRFilter) ([]*entity.PurchaseRequisition, int64, error) {
	return uc.prRepo.List(ctx, filter)
}
