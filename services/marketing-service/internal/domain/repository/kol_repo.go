package repository

import (
	"context"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// KOLTierRepository defines methods for KOL tier operations
type KOLTierRepository interface {
	Create(ctx context.Context, tier *entity.KOLTier) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.KOLTier, error)
	GetByCode(ctx context.Context, code string) (*entity.KOLTier, error)
	List(ctx context.Context, activeOnly bool) ([]*entity.KOLTier, error)
	Update(ctx context.Context, tier *entity.KOLTier) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// KOLFilter represents filter options for KOL listing
type KOLFilter struct {
	Search    string
	TierID    *uuid.UUID
	Category  entity.KOLCategory
	Niche     string
	Status    entity.KOLStatus
	Platform  entity.Platform
	Page      int
	Limit     int
}

// KOLRepository defines methods for KOL operations
type KOLRepository interface {
	Create(ctx context.Context, kol *entity.KOL) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.KOL, error)
	GetByCode(ctx context.Context, code string) (*entity.KOL, error)
	List(ctx context.Context, filter *KOLFilter) ([]*entity.KOL, int64, error)
	Update(ctx context.Context, kol *entity.KOL) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Stats
	IncrementPostCount(ctx context.Context, id uuid.UUID) error
	IncrementSampleCount(ctx context.Context, id uuid.UUID) error
	IncrementCollaborationCount(ctx context.Context, id uuid.UUID) error
	UpdateLastCollaborationDate(ctx context.Context, id uuid.UUID) error
	
	// Code generation
	GenerateKOLCode(ctx context.Context) (string, error)
}
