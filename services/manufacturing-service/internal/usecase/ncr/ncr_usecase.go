package ncr

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateNCRUseCase handles NCR creation
type CreateNCRUseCase struct {
	repo     repository.NCRRepository
	eventPub *event.Publisher
}

// NewCreateNCRUseCase creates a new CreateNCRUseCase
func NewCreateNCRUseCase(repo repository.NCRRepository, eventPub *event.Publisher) *CreateNCRUseCase {
	return &CreateNCRUseCase{repo: repo, eventPub: eventPub}
}

// CreateNCRInput is the input for creating an NCR
type CreateNCRInput struct {
	NCType           entity.NCType
	Severity         entity.NCRSeverity
	ReferenceType    string
	ReferenceID      *uuid.UUID
	ProductID        *uuid.UUID
	MaterialID       *uuid.UUID
	LotID            *uuid.UUID
	LotNumber        string
	Description      string
	QuantityAffected *float64
	UOMID            *uuid.UUID
	ImmediateAction  string
	CreatedBy        uuid.UUID
}

// Execute creates a new NCR
func (uc *CreateNCRUseCase) Execute(ctx context.Context, input CreateNCRInput) (*entity.NCR, error) {
	ncrNumber, err := uc.repo.GenerateNCRNumber(ctx)
	if err != nil {
		return nil, err
	}

	ncr := &entity.NCR{
		NCRNumber:        ncrNumber,
		NCType:           input.NCType,
		Severity:         input.Severity,
		Status:           entity.NCRStatusOpen,
		ReferenceType:    input.ReferenceType,
		ReferenceID:      input.ReferenceID,
		ProductID:        input.ProductID,
		MaterialID:       input.MaterialID,
		LotID:            input.LotID,
		LotNumber:        input.LotNumber,
		Description:      input.Description,
		QuantityAffected: input.QuantityAffected,
		UOMID:            input.UOMID,
		ImmediateAction:  input.ImmediateAction,
		CreatedBy:        &input.CreatedBy,
		UpdatedBy:        &input.CreatedBy,
	}

	if err := uc.repo.Create(ctx, ncr); err != nil {
		return nil, err
	}

	// Publish event - triggers notification
	ncrEvent := event.NCREvent{
		NCRID:       ncr.ID.String(),
		NCRNumber:   ncr.NCRNumber,
		NCType:      string(ncr.NCType),
		Severity:    string(ncr.Severity),
		Description: ncr.Description,
	}
	if ncr.LotID != nil {
		ncrEvent.LotID = ncr.LotID.String()
	}
	uc.eventPub.PublishNCRCreated(ncrEvent)

	return ncr, nil
}

// GetNCRUseCase handles getting an NCR
type GetNCRUseCase struct {
	repo repository.NCRRepository
}

// NewGetNCRUseCase creates a new GetNCRUseCase
func NewGetNCRUseCase(repo repository.NCRRepository) *GetNCRUseCase {
	return &GetNCRUseCase{repo: repo}
}

// Execute gets an NCR by ID
func (uc *GetNCRUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.NCR, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListNCRsUseCase handles listing NCRs
type ListNCRsUseCase struct {
	repo repository.NCRRepository
}

// NewListNCRsUseCase creates a new ListNCRsUseCase
func NewListNCRsUseCase(repo repository.NCRRepository) *ListNCRsUseCase {
	return &ListNCRsUseCase{repo: repo}
}

// Execute lists NCRs
func (uc *ListNCRsUseCase) Execute(ctx context.Context, filter repository.NCRFilter) ([]*entity.NCR, int64, error) {
	return uc.repo.List(ctx, filter)
}

// CloseNCRUseCase handles closing an NCR
type CloseNCRUseCase struct {
	repo repository.NCRRepository
}

// NewCloseNCRUseCase creates a new CloseNCRUseCase
func NewCloseNCRUseCase(repo repository.NCRRepository) *CloseNCRUseCase {
	return &CloseNCRUseCase{repo: repo}
}

// CloseNCRInput is input for closing an NCR
type CloseNCRInput struct {
	NCRID            uuid.UUID
	RootCause        string
	CorrectiveAction string
	PreventiveAction string
	Disposition      *entity.Disposition
	DispositionQty   *float64
	ClosureNotes     string
	ClosedBy         uuid.UUID
}

// Execute closes an NCR
func (uc *CloseNCRUseCase) Execute(ctx context.Context, input CloseNCRInput) (*entity.NCR, error) {
	ncr, err := uc.repo.GetByID(ctx, input.NCRID)
	if err != nil {
		return nil, entity.ErrNCRNotFound
	}

	if !ncr.CanBeClosed() {
		return nil, entity.ErrNCRAlreadyClosed
	}

	ncr.RootCause = input.RootCause
	ncr.CorrectiveAction = input.CorrectiveAction
	ncr.PreventiveAction = input.PreventiveAction
	ncr.Disposition = input.Disposition
	ncr.DispositionQuantity = input.DispositionQty

	if err := ncr.Close(input.ClosedBy, input.ClosureNotes); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, ncr); err != nil {
		return nil, err
	}

	return ncr, nil
}
