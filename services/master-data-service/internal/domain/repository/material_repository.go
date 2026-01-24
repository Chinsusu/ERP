package repository

import (
	"context"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/google/uuid"
)

// MaterialRepository defines material data access methods
type MaterialRepository interface {
	Create(ctx context.Context, material *entity.Material) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Material, error)
	GetByCode(ctx context.Context, code string) (*entity.Material, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Material, error)
	Update(ctx context.Context, material *entity.Material) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *MaterialFilter) ([]entity.Material, int64, error)
	Search(ctx context.Context, query string, filter *MaterialFilter) ([]entity.Material, int64, error)
	GetNextSequence(ctx context.Context, materialType entity.MaterialType) (int, error)
	
	// Specifications
	AddSpecification(ctx context.Context, spec *entity.MaterialSpecification) error
	GetSpecifications(ctx context.Context, materialID uuid.UUID) ([]entity.MaterialSpecification, error)
	UpdateSpecification(ctx context.Context, spec *entity.MaterialSpecification) error
	DeleteSpecification(ctx context.Context, id uuid.UUID) error
}

// MaterialFilter for listing and searching materials
type MaterialFilter struct {
	MaterialType     entity.MaterialType
	CategoryID       *uuid.UUID
	StorageCondition entity.StorageCondition
	IsOrganic        *bool
	IsNatural        *bool
	IsAllergen       *bool
	Status           string
	Search           string
	Page             int
	PageSize         int
}
