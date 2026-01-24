package repository

import (
	"context"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// SampleRequestFilter represents filter options for sample request listing
type SampleRequestFilter struct {
	KOLID      *uuid.UUID
	CampaignID *uuid.UUID
	Status     entity.SampleRequestStatus
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
}

// SampleRequestRepository defines methods for sample request operations
type SampleRequestRepository interface {
	Create(ctx context.Context, request *entity.SampleRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SampleRequest, error)
	GetByNumber(ctx context.Context, number string) (*entity.SampleRequest, error)
	List(ctx context.Context, filter *SampleRequestFilter) ([]*entity.SampleRequest, int64, error)
	Update(ctx context.Context, request *entity.SampleRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SampleRequestStatus) error
	Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) error
	Reject(ctx context.Context, id uuid.UUID, reason string) error
	
	// Line items
	AddItem(ctx context.Context, item *entity.SampleItem) error
	GetItems(ctx context.Context, requestID uuid.UUID) ([]*entity.SampleItem, error)
	UpdateItem(ctx context.Context, item *entity.SampleItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	
	// Code generation
	GenerateRequestNumber(ctx context.Context) (string, error)
	
	// By KOL
	GetByKOL(ctx context.Context, kolID uuid.UUID) ([]*entity.SampleRequest, error)
	
	// Pending approval
	GetPendingApproval(ctx context.Context) ([]*entity.SampleRequest, error)
}

// SampleShipmentFilter represents filter options for shipment listing
type SampleShipmentFilter struct {
	SampleRequestID *uuid.UUID
	Status          entity.ShipmentStatus
	Courier         string
	DateFrom        string
	DateTo          string
	Page            int
	Limit           int
}

// SampleShipmentRepository defines methods for sample shipment operations
type SampleShipmentRepository interface {
	Create(ctx context.Context, shipment *entity.SampleShipment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SampleShipment, error)
	GetByNumber(ctx context.Context, number string) (*entity.SampleShipment, error)
	GetByRequest(ctx context.Context, requestID uuid.UUID) (*entity.SampleShipment, error)
	List(ctx context.Context, filter *SampleShipmentFilter) ([]*entity.SampleShipment, int64, error)
	Update(ctx context.Context, shipment *entity.SampleShipment) error
	
	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ShipmentStatus) error
	MarkDelivered(ctx context.Context, id uuid.UUID, proofURL string) error
	
	// Code generation
	GenerateShipmentNumber(ctx context.Context) (string, error)
}

// KOLPostFilter represents filter options for KOL posts
type KOLPostFilter struct {
	KOLID           *uuid.UUID
	CampaignID      *uuid.UUID
	CollaborationID *uuid.UUID
	Platform        entity.Platform
	Verified        *bool
	DateFrom        string
	DateTo          string
	Page            int
	Limit           int
}

// KOLPostRepository defines methods for KOL post operations
type KOLPostRepository interface {
	Create(ctx context.Context, post *entity.KOLPost) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.KOLPost, error)
	List(ctx context.Context, filter *KOLPostFilter) ([]*entity.KOLPost, int64, error)
	Update(ctx context.Context, post *entity.KOLPost) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// By KOL
	GetByKOL(ctx context.Context, kolID uuid.UUID) ([]*entity.KOLPost, error)
	
	// Verify
	Verify(ctx context.Context, id uuid.UUID, verifierID uuid.UUID) error
	
	// Metrics update
	UpdateMetrics(ctx context.Context, id uuid.UUID, likes, comments, shares, views, reach int) error
}
