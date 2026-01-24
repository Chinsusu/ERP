package sample

import (
	"context"
	"time"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	"github.com/erp-cosmetics/marketing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// SampleItemInput represents input for sample item
type SampleItemInput struct {
	ProductID   uuid.UUID
	ProductCode string
	ProductName string
	Quantity    int
	UnitValue   float64
	Notes       string
}

// CreateSampleRequestInput represents input for creating sample request
type CreateSampleRequestInput struct {
	KOLID            uuid.UUID
	CampaignID       *uuid.UUID
	RequestReason    string
	DeliveryAddress  string
	RecipientName    string
	RecipientPhone   string
	ExpectedPostDate *time.Time
	ExpectedReach    int
	Items            []SampleItemInput
	Notes            string
}

// CreateSampleRequestUseCase handles sample request creation
type CreateSampleRequestUseCase struct {
	repo      repository.SampleRequestRepository
	kolRepo   repository.KOLRepository
	publisher *event.Publisher
}

// NewCreateSampleRequestUseCase creates a new use case
func NewCreateSampleRequestUseCase(repo repository.SampleRequestRepository, kolRepo repository.KOLRepository, publisher *event.Publisher) *CreateSampleRequestUseCase {
	return &CreateSampleRequestUseCase{repo: repo, kolRepo: kolRepo, publisher: publisher}
}

// Execute creates a new sample request
func (uc *CreateSampleRequestUseCase) Execute(ctx context.Context, input *CreateSampleRequestInput) (*entity.SampleRequest, error) {
	// Verify KOL exists
	kol, err := uc.kolRepo.GetByID(ctx, input.KOLID)
	if err != nil {
		return nil, err
	}

	number, err := uc.repo.GenerateRequestNumber(ctx)
	if err != nil {
		return nil, err
	}

	request := &entity.SampleRequest{
		RequestNumber:    number,
		KOLID:            input.KOLID,
		CampaignID:       input.CampaignID,
		RequestDate:      time.Now(),
		RequestReason:    input.RequestReason,
		DeliveryAddress:  input.DeliveryAddress,
		RecipientName:    input.RecipientName,
		RecipientPhone:   input.RecipientPhone,
		ExpectedPostDate: input.ExpectedPostDate,
		ExpectedReach:    input.ExpectedReach,
		Notes:            input.Notes,
		Status:           entity.SampleStatusDraft,
	}

	// Use KOL address if not provided
	if request.DeliveryAddress == "" {
		request.DeliveryAddress = kol.AddressLine1
	}
	if request.RecipientName == "" {
		request.RecipientName = kol.Name
	}
	if request.RecipientPhone == "" {
		request.RecipientPhone = kol.Phone
	}

	if err := uc.repo.Create(ctx, request); err != nil {
		return nil, err
	}

	// Add items
	var totalValue float64
	for i, itemInput := range input.Items {
		item := &entity.SampleItem{
			SampleRequestID: request.ID,
			LineNumber:      i + 1,
			ProductID:       itemInput.ProductID,
			ProductCode:     itemInput.ProductCode,
			ProductName:     itemInput.ProductName,
			Quantity:        itemInput.Quantity,
			UnitValue:       itemInput.UnitValue,
			TotalValue:      float64(itemInput.Quantity) * itemInput.UnitValue,
			Notes:           itemInput.Notes,
		}
		if err := uc.repo.AddItem(ctx, item); err != nil {
			return nil, err
		}
		totalValue += item.TotalValue
		request.Items = append(request.Items, *item)
	}

	// Update totals
	request.TotalItems = len(input.Items)
	request.TotalValue = totalValue
	if err := uc.repo.Update(ctx, request); err != nil {
		return nil, err
	}

	if uc.publisher != nil {
		uc.publisher.PublishSampleRequestCreated(ctx, map[string]interface{}{
			"request_id":     request.ID,
			"request_number": request.RequestNumber,
			"kol_id":         request.KOLID,
			"total_value":    request.TotalValue,
		})
	}

	return request, nil
}

// GetSampleRequestUseCase handles getting a sample request
type GetSampleRequestUseCase struct {
	repo repository.SampleRequestRepository
}

// NewGetSampleRequestUseCase creates a new use case
func NewGetSampleRequestUseCase(repo repository.SampleRequestRepository) *GetSampleRequestUseCase {
	return &GetSampleRequestUseCase{repo: repo}
}

// Execute gets a sample request by ID
func (uc *GetSampleRequestUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.SampleRequest, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListSampleRequestsUseCase handles listing sample requests
type ListSampleRequestsUseCase struct {
	repo repository.SampleRequestRepository
}

// NewListSampleRequestsUseCase creates a new use case
func NewListSampleRequestsUseCase(repo repository.SampleRequestRepository) *ListSampleRequestsUseCase {
	return &ListSampleRequestsUseCase{repo: repo}
}

// Execute lists sample requests with filter
func (uc *ListSampleRequestsUseCase) Execute(ctx context.Context, filter *repository.SampleRequestFilter) ([]*entity.SampleRequest, int64, error) {
	return uc.repo.List(ctx, filter)
}

// ApproveSampleRequestUseCase handles approving a sample request
type ApproveSampleRequestUseCase struct {
	repo      repository.SampleRequestRepository
	kolRepo   repository.KOLRepository
	publisher *event.Publisher
}

// NewApproveSampleRequestUseCase creates a new use case
func NewApproveSampleRequestUseCase(repo repository.SampleRequestRepository, kolRepo repository.KOLRepository, publisher *event.Publisher) *ApproveSampleRequestUseCase {
	return &ApproveSampleRequestUseCase{repo: repo, kolRepo: kolRepo, publisher: publisher}
}

// Execute approves a sample request
func (uc *ApproveSampleRequestUseCase) Execute(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*entity.SampleRequest, error) {
	request, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !request.CanBeApproved() {
		return nil, ErrRequestCannotBeApproved
	}

	request.Approve(approverID)

	if err := uc.repo.Update(ctx, request); err != nil {
		return nil, err
	}

	// Increment KOL sample count
	uc.kolRepo.IncrementSampleCount(ctx, request.KOLID)

	if uc.publisher != nil {
		uc.publisher.PublishSampleRequestApproved(ctx, map[string]interface{}{
			"request_id":     request.ID,
			"request_number": request.RequestNumber,
			"kol_id":         request.KOLID,
			"approved_by":    approverID,
			"items":          request.Items,
		})
	}

	return request, nil
}

// RejectSampleRequestUseCase handles rejecting a sample request
type RejectSampleRequestUseCase struct {
	repo repository.SampleRequestRepository
}

// NewRejectSampleRequestUseCase creates a new use case
func NewRejectSampleRequestUseCase(repo repository.SampleRequestRepository) *RejectSampleRequestUseCase {
	return &RejectSampleRequestUseCase{repo: repo}
}

// Execute rejects a sample request
func (uc *RejectSampleRequestUseCase) Execute(ctx context.Context, id uuid.UUID, reason string) (*entity.SampleRequest, error) {
	request, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !request.CanBeApproved() {
		return nil, ErrRequestCannotBeApproved
	}

	request.Reject(reason)

	if err := uc.repo.Update(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

// ShipSampleInput represents input for shipping a sample
type ShipSampleInput struct {
	RequestID         uuid.UUID
	Courier           string
	TrackingNumber    string
	EstimatedDelivery *time.Time
}

// ShipSampleUseCase handles shipping a sample
type ShipSampleUseCase struct {
	requestRepo  repository.SampleRequestRepository
	shipmentRepo repository.SampleShipmentRepository
	publisher    *event.Publisher
}

// NewShipSampleUseCase creates a new use case
func NewShipSampleUseCase(requestRepo repository.SampleRequestRepository, shipmentRepo repository.SampleShipmentRepository, publisher *event.Publisher) *ShipSampleUseCase {
	return &ShipSampleUseCase{requestRepo: requestRepo, shipmentRepo: shipmentRepo, publisher: publisher}
}

// Execute ships a sample
func (uc *ShipSampleUseCase) Execute(ctx context.Context, input *ShipSampleInput) (*entity.SampleShipment, error) {
	request, err := uc.requestRepo.GetByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	if request.Status != entity.SampleStatusApproved {
		return nil, ErrRequestMustBeApproved
	}

	number, err := uc.shipmentRepo.GenerateShipmentNumber(ctx)
	if err != nil {
		return nil, err
	}

	shipment := &entity.SampleShipment{
		ShipmentNumber:    number,
		SampleRequestID:   input.RequestID,
		ShipmentDate:      time.Now(),
		Courier:           input.Courier,
		TrackingNumber:    input.TrackingNumber,
		RecipientName:     request.RecipientName,
		RecipientPhone:    request.RecipientPhone,
		DeliveryAddress:   request.DeliveryAddress,
		EstimatedDelivery: input.EstimatedDelivery,
		Status:            entity.ShipmentStatusShipped,
	}

	if err := uc.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, err
	}

	// Update request status
	request.MarkShipped()
	if err := uc.requestRepo.Update(ctx, request); err != nil {
		return nil, err
	}

	if uc.publisher != nil {
		uc.publisher.PublishSampleShipped(ctx, map[string]interface{}{
			"shipment_id":     shipment.ID,
			"shipment_number": shipment.ShipmentNumber,
			"request_id":      input.RequestID,
			"tracking_number": input.TrackingNumber,
			"courier":         input.Courier,
		})
	}

	return shipment, nil
}

// Custom errors
var (
	ErrRequestCannotBeApproved = &SampleError{Message: "request cannot be approved"}
	ErrRequestMustBeApproved   = &SampleError{Message: "request must be approved before shipping"}
)

// SampleError represents a sample-related error
type SampleError struct {
	Message string
}

func (e *SampleError) Error() string {
	return e.Message
}
