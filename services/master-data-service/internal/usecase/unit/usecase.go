package unit

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase handles unit of measure business logic
type UseCase struct {
	repo repository.UnitRepository
}

// NewUseCase creates a new unit use case
func NewUseCase(repo repository.UnitRepository) *UseCase {
	return &UseCase{repo: repo}
}

// CreateRequest represents unit creation request
type CreateRequest struct {
	Code             string
	Name             string
	NameEN           string
	Symbol           string
	UoMType          entity.UoMType
	IsBaseUnit       bool
	BaseUnitID       *uuid.UUID
	ConversionFactor float64
	Description      string
}

// Create creates a new unit of measure
func (uc *UseCase) Create(ctx context.Context, req *CreateRequest) (*entity.UnitOfMeasure, error) {
	// Check if code already exists
	existing, _ := uc.repo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("unit with code %s already exists", req.Code)
	}

	// Default conversion factor to 1
	if req.ConversionFactor == 0 {
		req.ConversionFactor = 1
	}

	unit := &entity.UnitOfMeasure{
		Code:             req.Code,
		Name:             req.Name,
		NameEN:           req.NameEN,
		Symbol:           req.Symbol,
		UoMType:          req.UoMType,
		IsBaseUnit:       req.IsBaseUnit,
		BaseUnitID:       req.BaseUnitID,
		ConversionFactor: req.ConversionFactor,
		Description:      req.Description,
		Status:           "active",
	}

	if err := unit.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}

	return unit, nil
}

// GetByID retrieves a unit by ID
func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.UnitOfMeasure, error) {
	return uc.repo.GetByID(ctx, id)
}

// Update updates an existing unit
func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, name, description string) (*entity.UnitOfMeasure, error) {
	unit, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unit not found: %w", err)
	}

	if name != "" {
		unit.Name = name
	}
	if description != "" {
		unit.Description = description
	}

	if err := uc.repo.Update(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to update unit: %w", err)
	}

	return unit, nil
}

// Delete soft deletes a unit
func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// List lists units with filters
func (uc *UseCase) List(ctx context.Context, filter *repository.UnitFilter) ([]entity.UnitOfMeasure, int64, error) {
	return uc.repo.List(ctx, filter)
}

// ConvertRequest represents unit conversion request
type ConvertRequest struct {
	Value      float64
	FromUnitID uuid.UUID
	ToUnitID   uuid.UUID
}

// ConvertResponse represents unit conversion response
type ConvertResponse struct {
	OriginalValue   float64 `json:"original_value"`
	OriginalUoM     string  `json:"original_uom"`
	ConvertedValue  float64 `json:"converted_value"`
	ConvertedUoM    string  `json:"converted_uom"`
}

// Convert converts a value from one unit to another
func (uc *UseCase) Convert(ctx context.Context, req *ConvertRequest) (*ConvertResponse, error) {
	fromUnit, err := uc.repo.GetByID(ctx, req.FromUnitID)
	if err != nil {
		return nil, fmt.Errorf("from unit not found: %w", err)
	}

	toUnit, err := uc.repo.GetByID(ctx, req.ToUnitID)
	if err != nil {
		return nil, fmt.Errorf("to unit not found: %w", err)
	}

	convertedValue, err := uc.repo.Convert(ctx, req.Value, req.FromUnitID, req.ToUnitID)
	if err != nil {
		return nil, err
	}

	return &ConvertResponse{
		OriginalValue:  req.Value,
		OriginalUoM:    fromUnit.Code,
		ConvertedValue: convertedValue,
		ConvertedUoM:   toUnit.Code,
	}, nil
}

// GetConversions gets all conversions for a unit
func (uc *UseCase) GetConversions(ctx context.Context, unitID uuid.UUID) ([]entity.UnitConversion, error) {
	return uc.repo.ListConversions(ctx, unitID)
}
