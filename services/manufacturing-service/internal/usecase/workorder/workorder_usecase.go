package workorder

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// EventPublisher defines event publishing interface for work order
type EventPublisher interface {
	PublishWOCreated(event event.WOEvent) error
	PublishWOReleased(event event.WOEvent) error
	PublishWOStarted(event event.WOEvent) error
	PublishWOCompleted(event event.WOCompletedEvent) error
}

// CreateWOUseCase handles work order creation
type CreateWOUseCase struct {
	woRepo   repository.WorkOrderRepository
	bomRepo  repository.BOMRepository
	eventPub EventPublisher
}

// NewCreateWOUseCase creates a new CreateWOUseCase
func NewCreateWOUseCase(woRepo repository.WorkOrderRepository, bomRepo repository.BOMRepository, eventPub EventPublisher) *CreateWOUseCase {
	return &CreateWOUseCase{
		woRepo:   woRepo,
		bomRepo:  bomRepo,
		eventPub: eventPub,
	}
}

// CreateWOInput is the input for creating a work order
type CreateWOInput struct {
	ProductID        uuid.UUID
	BOMID            uuid.UUID
	PlannedQuantity  float64
	UOMID            uuid.UUID
	PlannedStartDate string
	PlannedEndDate   string
	BatchNumber      string
	SalesOrderID     *uuid.UUID
	ProductionLine   string
	Shift            string
	Priority         entity.WOPriority
	Notes            string
	CreatedBy        uuid.UUID
}

// Execute creates a new work order
func (uc *CreateWOUseCase) Execute(ctx context.Context, input CreateWOInput) (*entity.WorkOrder, error) {
	// Validate BOM exists and is approved
	bom, err := uc.bomRepo.GetByID(ctx, input.BOMID)
	if err != nil {
		return nil, entity.ErrBOMNotFound
	}
	if !bom.IsActive() {
		return nil, entity.ErrBOMNotPendingApproval
	}

	// Generate WO number
	woNumber, err := uc.woRepo.GenerateWONumber(ctx)
	if err != nil {
		return nil, err
	}

	wo := &entity.WorkOrder{
		WONumber:        woNumber,
		ProductID:       input.ProductID,
		BOMID:           input.BOMID,
		Status:          entity.WOStatusPlanned,
		Priority:        input.Priority,
		PlannedQuantity: input.PlannedQuantity,
		UOMID:           input.UOMID,
		BatchNumber:     input.BatchNumber,
		SalesOrderID:    input.SalesOrderID,
		ProductionLine:  input.ProductionLine,
		Shift:           input.Shift,
		Notes:           input.Notes,
		CreatedBy:       &input.CreatedBy,
		UpdatedBy:       &input.CreatedBy,
	}

	if err := uc.woRepo.Create(ctx, wo); err != nil {
		return nil, err
	}

	// Create line items from BOM
	var woItems []*entity.WOLineItem
	ratio := input.PlannedQuantity / bom.BatchSize
	for i, bomItem := range bom.Items {
		woItems = append(woItems, &entity.WOLineItem{
			WorkOrderID:     wo.ID,
			BOMLineItemID:   &bomItem.ID,
			LineNumber:      i + 1,
			MaterialID:      bomItem.MaterialID,
			PlannedQuantity: bomItem.Quantity * ratio,
			IssuedQuantity:  0,
			UOMID:           bomItem.UOMID,
			IsCritical:      bomItem.IsCritical,
		})
	}
	if err := uc.woRepo.CreateLineItems(ctx, woItems); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishWOCreated(event.WOEvent{
		WOID:            wo.ID.String(),
		WONumber:        wo.WONumber,
		ProductID:       wo.ProductID.String(),
		BOMID:           wo.BOMID.String(),
		BatchNumber:     wo.BatchNumber,
		PlannedQuantity: wo.PlannedQuantity,
		Status:          string(wo.Status),
	})

	return wo, nil
}

// GetWOUseCase handles getting a work order
type GetWOUseCase struct {
	repo repository.WorkOrderRepository
}

// NewGetWOUseCase creates a new GetWOUseCase
func NewGetWOUseCase(repo repository.WorkOrderRepository) *GetWOUseCase {
	return &GetWOUseCase{repo: repo}
}

// Execute gets a work order by ID
func (uc *GetWOUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.WorkOrder, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListWOsUseCase handles listing work orders
type ListWOsUseCase struct {
	repo repository.WorkOrderRepository
}

// NewListWOsUseCase creates a new ListWOsUseCase
func NewListWOsUseCase(repo repository.WorkOrderRepository) *ListWOsUseCase {
	return &ListWOsUseCase{repo: repo}
}

// Execute lists work orders
func (uc *ListWOsUseCase) Execute(ctx context.Context, filter repository.WOFilter) ([]*entity.WorkOrder, int64, error) {
	return uc.repo.List(ctx, filter)
}

// ReleaseWOUseCase handles releasing a work order
type ReleaseWOUseCase struct {
	repo     repository.WorkOrderRepository
	eventPub EventPublisher
}

// NewReleaseWOUseCase creates a new ReleaseWOUseCase
func NewReleaseWOUseCase(repo repository.WorkOrderRepository, eventPub EventPublisher) *ReleaseWOUseCase {
	return &ReleaseWOUseCase{repo: repo, eventPub: eventPub}
}

// Execute releases a work order
func (uc *ReleaseWOUseCase) Execute(ctx context.Context, woID uuid.UUID, updatedBy uuid.UUID) (*entity.WorkOrder, error) {
	wo, err := uc.repo.GetByID(ctx, woID)
	if err != nil {
		return nil, entity.ErrWONotFound
	}

	if err := wo.Release(); err != nil {
		return nil, entity.ErrWOCannotRelease
	}
	wo.UpdatedBy = &updatedBy

	if err := uc.repo.Update(ctx, wo); err != nil {
		return nil, err
	}

	uc.eventPub.PublishWOReleased(event.WOEvent{
		WOID:            wo.ID.String(),
		WONumber:        wo.WONumber,
		ProductID:       wo.ProductID.String(),
		BOMID:           wo.BOMID.String(),
		BatchNumber:     wo.BatchNumber,
		PlannedQuantity: wo.PlannedQuantity,
		Status:          string(wo.Status),
	})

	return wo, nil
}

// StartWOUseCase handles starting a work order
type StartWOUseCase struct {
	repo     repository.WorkOrderRepository
	eventPub EventPublisher
}

// NewStartWOUseCase creates a new StartWOUseCase
func NewStartWOUseCase(repo repository.WorkOrderRepository, eventPub EventPublisher) *StartWOUseCase {
	return &StartWOUseCase{repo: repo, eventPub: eventPub}
}

// Execute starts a work order
func (uc *StartWOUseCase) Execute(ctx context.Context, woID uuid.UUID, supervisorID uuid.UUID) (*entity.WorkOrder, error) {
	wo, err := uc.repo.GetByID(ctx, woID)
	if err != nil {
		return nil, entity.ErrWONotFound
	}

	if err := wo.Start(supervisorID); err != nil {
		return nil, entity.ErrWOCannotStart
	}

	if err := uc.repo.Update(ctx, wo); err != nil {
		return nil, err
	}

	// This event triggers WMS to reserve materials
	uc.eventPub.PublishWOStarted(event.WOEvent{
		WOID:            wo.ID.String(),
		WONumber:        wo.WONumber,
		ProductID:       wo.ProductID.String(),
		BOMID:           wo.BOMID.String(),
		BatchNumber:     wo.BatchNumber,
		PlannedQuantity: wo.PlannedQuantity,
		Status:          string(wo.Status),
	})

	return wo, nil
}

// CompleteWOUseCase handles completing a work order
type CompleteWOUseCase struct {
	repo      repository.WorkOrderRepository
	traceRepo repository.TraceabilityRepository
	eventPub  EventPublisher
}

// NewCompleteWOUseCase creates a new CompleteWOUseCase
func NewCompleteWOUseCase(repo repository.WorkOrderRepository, traceRepo repository.TraceabilityRepository, eventPub EventPublisher) *CompleteWOUseCase {
	return &CompleteWOUseCase{
		repo:      repo,
		traceRepo: traceRepo,
		eventPub:  eventPub,
	}
}

// CompleteWOInput is the input for completing a work order
type CompleteWOInput struct {
	WOID             uuid.UUID
	ActualQuantity   float64
	GoodQuantity     float64
	RejectedQuantity float64
	Notes            string
	UpdatedBy        uuid.UUID
}

// Execute completes a work order
func (uc *CompleteWOUseCase) Execute(ctx context.Context, input CompleteWOInput) (*entity.WorkOrder, error) {
	wo, err := uc.repo.GetByID(ctx, input.WOID)
	if err != nil {
		return nil, entity.ErrWONotFound
	}

	if err := wo.Complete(input.ActualQuantity, input.GoodQuantity, input.RejectedQuantity); err != nil {
		return nil, entity.ErrWOCannotComplete
	}
	wo.Notes = input.Notes
	wo.UpdatedBy = &input.UpdatedBy

	if err := uc.repo.Update(ctx, wo); err != nil {
		return nil, err
	}

	// This event triggers WMS to receive finished goods
	uc.eventPub.PublishWOCompleted(event.WOCompletedEvent{
		WOID:         wo.ID.String(),
		WONumber:     wo.WONumber,
		ProductID:    wo.ProductID.String(),
		BatchNumber:  wo.BatchNumber,
		GoodQuantity: input.GoodQuantity,
	})

	return wo, nil
}
