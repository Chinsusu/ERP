package repository

import (
	"context"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/google/uuid"
)

// UnitRepository defines unit of measure data access methods
type UnitRepository interface {
	// Units of Measure
	Create(ctx context.Context, unit *entity.UnitOfMeasure) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UnitOfMeasure, error)
	GetByCode(ctx context.Context, code string) (*entity.UnitOfMeasure, error)
	Update(ctx context.Context, unit *entity.UnitOfMeasure) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *UnitFilter) ([]entity.UnitOfMeasure, int64, error)
	
	// Unit Conversions
	CreateConversion(ctx context.Context, conv *entity.UnitConversion) error
	GetConversion(ctx context.Context, fromUnitID, toUnitID uuid.UUID) (*entity.UnitConversion, error)
	ListConversions(ctx context.Context, unitID uuid.UUID) ([]entity.UnitConversion, error)
	Convert(ctx context.Context, value float64, fromUnitID, toUnitID uuid.UUID) (float64, error)
}

// UnitFilter for listing units
type UnitFilter struct {
	UoMType    entity.UoMType
	IsBaseUnit *bool
	Status     string
	Search     string
	Page       int
	PageSize   int
}
