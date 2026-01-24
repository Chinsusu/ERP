package repository

import (
	"context"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CampaignFilter represents filter options for campaign listing
type CampaignFilter struct {
	Search       string
	CampaignType entity.CampaignType
	Status       entity.CampaignStatus
	DateFrom     string
	DateTo       string
	Page         int
	Limit        int
}

// CampaignRepository defines methods for campaign operations
type CampaignRepository interface {
	Create(ctx context.Context, campaign *entity.Campaign) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Campaign, error)
	GetByCode(ctx context.Context, code string) (*entity.Campaign, error)
	List(ctx context.Context, filter *CampaignFilter) ([]*entity.Campaign, int64, error)
	Update(ctx context.Context, campaign *entity.Campaign) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.CampaignStatus) error
	
	// Performance updates
	UpdatePerformance(ctx context.Context, id uuid.UUID, impressions, reach, engagement, conversions int, revenue float64) error
	IncrementSpent(ctx context.Context, id uuid.UUID, amount float64) error
	
	// Code generation
	GenerateCampaignCode(ctx context.Context, prefix string) (string, error)
	
	// Active campaigns
	GetActiveCampaigns(ctx context.Context) ([]*entity.Campaign, error)
}

// KOLCollaborationFilter represents filter options for collaborations
type KOLCollaborationFilter struct {
	CampaignID *uuid.UUID
	KOLID      *uuid.UUID
	Status     entity.CollaborationStatus
	Page       int
	Limit      int
}

// KOLCollaborationRepository defines methods for collaboration operations
type KOLCollaborationRepository interface {
	Create(ctx context.Context, collab *entity.KOLCollaboration) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.KOLCollaboration, error)
	List(ctx context.Context, filter *KOLCollaborationFilter) ([]*entity.KOLCollaboration, int64, error)
	Update(ctx context.Context, collab *entity.KOLCollaboration) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Status
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.CollaborationStatus) error
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus, amount float64) error
	
	// Performance
	UpdatePerformance(ctx context.Context, id uuid.UUID, impressions, engagement, reach int) error
	IncrementPostCount(ctx context.Context, id uuid.UUID) error
	
	// By KOL
	GetByKOL(ctx context.Context, kolID uuid.UUID) ([]*entity.KOLCollaboration, error)
	GetByKOLAndCampaign(ctx context.Context, kolID, campaignID uuid.UUID) (*entity.KOLCollaboration, error)
	
	// Code generation
	GenerateCode(ctx context.Context) (string, error)
}
