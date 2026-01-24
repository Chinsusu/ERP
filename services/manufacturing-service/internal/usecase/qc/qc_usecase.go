package qc

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// GetCheckpointsUseCase handles getting QC checkpoints
type GetCheckpointsUseCase struct {
	repo repository.QCRepository
}

// NewGetCheckpointsUseCase creates a new GetCheckpointsUseCase
func NewGetCheckpointsUseCase(repo repository.QCRepository) *GetCheckpointsUseCase {
	return &GetCheckpointsUseCase{repo: repo}
}

// Execute gets all active checkpoints
func (uc *GetCheckpointsUseCase) Execute(ctx context.Context) ([]*entity.QCCheckpoint, error) {
	return uc.repo.GetCheckpoints(ctx)
}

// CreateInspectionUseCase handles creating QC inspections
type CreateInspectionUseCase struct {
	repo     repository.QCRepository
	eventPub *event.Publisher
}

// NewCreateInspectionUseCase creates a new CreateInspectionUseCase
func NewCreateInspectionUseCase(repo repository.QCRepository, eventPub *event.Publisher) *CreateInspectionUseCase {
	return &CreateInspectionUseCase{repo: repo, eventPub: eventPub}
}

// CreateInspectionInput is the input for creating an inspection
type CreateInspectionInput struct {
	InspectionType    entity.CheckpointType
	CheckpointID      *uuid.UUID
	ReferenceType     entity.ReferenceType
	ReferenceID       uuid.UUID
	ProductID         *uuid.UUID
	MaterialID        *uuid.UUID
	LotID             *uuid.UUID
	LotNumber         string
	InspectedQuantity float64
	SampleSize        *int
	InspectorID       uuid.UUID
	InspectorName     string
	Items             []CreateInspectionItemInput
}

// CreateInspectionItemInput is input for an inspection item
type CreateInspectionItemInput struct {
	ItemNumber    int
	TestName      string
	TestMethod    string
	Specification string
	TargetValue   string
	MinValue      string
	MaxValue      string
	ActualValue   string
	UOM           string
	Result        entity.ItemResult
	Notes         string
}

// Execute creates a new inspection
func (uc *CreateInspectionUseCase) Execute(ctx context.Context, input CreateInspectionInput) (*entity.QCInspection, error) {
	// Generate inspection number
	inspNumber, err := uc.repo.GenerateInspectionNumber(ctx)
	if err != nil {
		return nil, err
	}

	inspection := &entity.QCInspection{
		InspectionNumber:  inspNumber,
		InspectionType:    input.InspectionType,
		CheckpointID:      input.CheckpointID,
		ReferenceType:     input.ReferenceType,
		ReferenceID:       input.ReferenceID,
		ProductID:         input.ProductID,
		MaterialID:        input.MaterialID,
		LotID:             input.LotID,
		LotNumber:         input.LotNumber,
		InspectedQuantity: input.InspectedQuantity,
		SampleSize:        input.SampleSize,
		Result:            entity.InspectionResultPending,
		InspectorID:       input.InspectorID,
		InspectorName:     input.InspectorName,
	}

	if err := uc.repo.CreateInspection(ctx, inspection); err != nil {
		return nil, err
	}

	// Create inspection items
	var items []*entity.QCInspectionItem
	for _, item := range input.Items {
		items = append(items, &entity.QCInspectionItem{
			InspectionID:  inspection.ID,
			ItemNumber:    item.ItemNumber,
			TestName:      item.TestName,
			TestMethod:    item.TestMethod,
			Specification: item.Specification,
			TargetValue:   item.TargetValue,
			MinValue:      item.MinValue,
			MaxValue:      item.MaxValue,
			ActualValue:   item.ActualValue,
			UOM:           item.UOM,
			Result:        item.Result,
			Notes:         item.Notes,
		})
	}
	if len(items) > 0 {
		if err := uc.repo.CreateInspectionItems(ctx, items); err != nil {
			return nil, err
		}
	}

	return inspection, nil
}

// GetInspectionUseCase handles getting an inspection
type GetInspectionUseCase struct {
	repo repository.QCRepository
}

// NewGetInspectionUseCase creates a new GetInspectionUseCase
func NewGetInspectionUseCase(repo repository.QCRepository) *GetInspectionUseCase {
	return &GetInspectionUseCase{repo: repo}
}

// Execute gets an inspection by ID
func (uc *GetInspectionUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.QCInspection, error) {
	return uc.repo.GetInspectionByID(ctx, id)
}

// ListInspectionsUseCase handles listing inspections
type ListInspectionsUseCase struct {
	repo repository.QCRepository
}

// NewListInspectionsUseCase creates a new ListInspectionsUseCase
func NewListInspectionsUseCase(repo repository.QCRepository) *ListInspectionsUseCase {
	return &ListInspectionsUseCase{repo: repo}
}

// Execute lists inspections
func (uc *ListInspectionsUseCase) Execute(ctx context.Context, filter repository.QCFilter) ([]*entity.QCInspection, int64, error) {
	return uc.repo.ListInspections(ctx, filter)
}

// ApproveInspectionUseCase handles approving/rejecting inspections
type ApproveInspectionUseCase struct {
	repo     repository.QCRepository
	eventPub *event.Publisher
}

// NewApproveInspectionUseCase creates a new ApproveInspectionUseCase
func NewApproveInspectionUseCase(repo repository.QCRepository, eventPub *event.Publisher) *ApproveInspectionUseCase {
	return &ApproveInspectionUseCase{repo: repo, eventPub: eventPub}
}

// ApproveInspectionInput is input for approving an inspection
type ApproveInspectionInput struct {
	InspectionID     uuid.UUID
	Result           entity.InspectionResult
	AcceptedQuantity *float64
	RejectedQuantity *float64
	ApproverID       uuid.UUID
	Notes            string
}

// Execute approves or rejects an inspection
func (uc *ApproveInspectionUseCase) Execute(ctx context.Context, input ApproveInspectionInput) (*entity.QCInspection, error) {
	inspection, err := uc.repo.GetInspectionByID(ctx, input.InspectionID)
	if err != nil {
		return nil, entity.ErrQCInspectionNotFound
	}

	if !inspection.CanBeApproved() {
		return nil, entity.ErrQCAlreadyApproved
	}

	inspection.AcceptedQuantity = input.AcceptedQuantity
	inspection.RejectedQuantity = input.RejectedQuantity
	inspection.Notes = input.Notes

	switch input.Result {
	case entity.InspectionResultPassed:
		inspection.Pass(input.ApproverID)
	case entity.InspectionResultFailed:
		inspection.Fail(input.ApproverID)
	case entity.InspectionResultConditional:
		inspection.ConditionalPass(input.ApproverID)
	}

	inspection.CalculateScore()

	if err := uc.repo.UpdateInspection(ctx, inspection); err != nil {
		return nil, err
	}

	// Publish event
	qcEvent := event.QCEvent{
		InspectionID:     inspection.ID.String(),
		InspectionNumber: inspection.InspectionNumber,
		InspectionType:   string(inspection.InspectionType),
		ReferenceType:    string(inspection.ReferenceType),
		ReferenceID:      inspection.ReferenceID.String(),
		Result:           string(inspection.Result),
	}
	if inspection.LotID != nil {
		qcEvent.LotID = inspection.LotID.String()
	}

	if input.Result == entity.InspectionResultPassed {
		uc.eventPub.PublishQCPassed(qcEvent)
	} else if input.Result == entity.InspectionResultFailed {
		uc.eventPub.PublishQCFailed(qcEvent)
	}

	return inspection, nil
}
