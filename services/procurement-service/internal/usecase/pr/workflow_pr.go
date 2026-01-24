package pr

import (
	"context"
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// SubmitPRUseCase handles submitting a PR for approval
type SubmitPRUseCase struct {
	prRepo   repository.PRRepository
	eventPub *event.Publisher
}

// NewSubmitPRUseCase creates a new use case
func NewSubmitPRUseCase(prRepo repository.PRRepository, eventPub *event.Publisher) *SubmitPRUseCase {
	return &SubmitPRUseCase{prRepo: prRepo, eventPub: eventPub}
}

// Execute submits a PR for approval
func (uc *SubmitPRUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.PurchaseRequisition, error) {
	pr, err := uc.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !pr.CanSubmit() {
		return nil, entity.ErrInvalidPRStatus
	}

	pr.Submit()
	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishPRSubmitted(ctx, &event.PREvent{
		PRID:        pr.ID.String(),
		PRNumber:    pr.PRNumber,
		Status:      string(pr.Status),
		TotalAmount: pr.TotalAmount,
		Timestamp:   time.Now(),
	})

	return pr, nil
}

// ApprovePRUseCase handles approving a PR
type ApprovePRUseCase struct {
	prRepo   repository.PRRepository
	eventPub *event.Publisher
}

// NewApprovePRUseCase creates a new use case
func NewApprovePRUseCase(prRepo repository.PRRepository, eventPub *event.Publisher) *ApprovePRUseCase {
	return &ApprovePRUseCase{prRepo: prRepo, eventPub: eventPub}
}

// Execute approves a PR
func (uc *ApprovePRUseCase) Execute(ctx context.Context, id uuid.UUID, approverID uuid.UUID, notes string) (*entity.PurchaseRequisition, error) {
	pr, err := uc.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !pr.CanApprove() {
		return nil, entity.ErrInvalidPRStatus
	}

	pr.Approve(approverID)
	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	// Create approval record
	approval := &entity.PRApproval{
		ID:            uuid.New(),
		PRID:          pr.ID,
		ApproverID:    approverID,
		ApprovalLevel: pr.ApprovalLevel,
		Action:        "APPROVED",
		Notes:         notes,
		CreatedAt:     time.Now(),
	}
	uc.prRepo.CreateApproval(ctx, approval)

	// Publish event
	uc.eventPub.PublishPRApproved(ctx, &event.PREvent{
		PRID:        pr.ID.String(),
		PRNumber:    pr.PRNumber,
		Status:      string(pr.Status),
		TotalAmount: pr.TotalAmount,
		ApprovedBy:  approverID.String(),
		Timestamp:   time.Now(),
	})

	return pr, nil
}

// RejectPRUseCase handles rejecting a PR
type RejectPRUseCase struct {
	prRepo   repository.PRRepository
	eventPub *event.Publisher
}

// NewRejectPRUseCase creates a new use case
func NewRejectPRUseCase(prRepo repository.PRRepository, eventPub *event.Publisher) *RejectPRUseCase {
	return &RejectPRUseCase{prRepo: prRepo, eventPub: eventPub}
}

// Execute rejects a PR
func (uc *RejectPRUseCase) Execute(ctx context.Context, id uuid.UUID, rejectedBy uuid.UUID, reason string) (*entity.PurchaseRequisition, error) {
	pr, err := uc.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !pr.CanApprove() {
		return nil, entity.ErrInvalidPRStatus
	}

	pr.Reject(rejectedBy, reason)
	if err := uc.prRepo.Update(ctx, pr); err != nil {
		return nil, err
	}

	// Create approval record
	approval := &entity.PRApproval{
		ID:            uuid.New(),
		PRID:          pr.ID,
		ApproverID:    rejectedBy,
		ApprovalLevel: pr.ApprovalLevel,
		Action:        "REJECTED",
		Notes:         reason,
		CreatedAt:     time.Now(),
	}
	uc.prRepo.CreateApproval(ctx, approval)

	// Publish event
	uc.eventPub.PublishPRRejected(ctx, &event.PREvent{
		PRID:       pr.ID.String(),
		PRNumber:   pr.PRNumber,
		Status:     string(pr.Status),
		RejectedBy: rejectedBy.String(),
		Reason:     reason,
		Timestamp:  time.Now(),
	})

	return pr, nil
}
