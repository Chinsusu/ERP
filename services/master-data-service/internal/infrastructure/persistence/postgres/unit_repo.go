package postgres

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type unitRepository struct {
	db *gorm.DB
}

// NewUnitRepository creates a new unit repository
func NewUnitRepository(db *gorm.DB) repository.UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) Create(ctx context.Context, unit *entity.UnitOfMeasure) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *unitRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UnitOfMeasure, error) {
	var unit entity.UnitOfMeasure
	err := r.db.WithContext(ctx).
		Preload("BaseUnit").
		First(&unit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) GetByCode(ctx context.Context, code string) (*entity.UnitOfMeasure, error) {
	var unit entity.UnitOfMeasure
	err := r.db.WithContext(ctx).
		First(&unit, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) Update(ctx context.Context, unit *entity.UnitOfMeasure) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

func (r *unitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.UnitOfMeasure{}, "id = ?", id).Error
}

func (r *unitRepository) List(ctx context.Context, filter *repository.UnitFilter) ([]entity.UnitOfMeasure, int64, error) {
	var units []entity.UnitOfMeasure
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.UnitOfMeasure{})

	// Apply filters
	if filter.UoMType != "" {
		query = query.Where("uom_type = ?", filter.UoMType)
	}
	if filter.IsBaseUnit != nil {
		query = query.Where("is_base_unit = ?", *filter.IsBaseUnit)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	query = query.Order("uom_type ASC, code ASC")

	if err := query.Preload("BaseUnit").Find(&units).Error; err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (r *unitRepository) CreateConversion(ctx context.Context, conv *entity.UnitConversion) error {
	return r.db.WithContext(ctx).Create(conv).Error
}

func (r *unitRepository) GetConversion(ctx context.Context, fromUnitID, toUnitID uuid.UUID) (*entity.UnitConversion, error) {
	var conv entity.UnitConversion
	err := r.db.WithContext(ctx).
		First(&conv, "from_unit_id = ? AND to_unit_id = ?", fromUnitID, toUnitID).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *unitRepository) ListConversions(ctx context.Context, unitID uuid.UUID) ([]entity.UnitConversion, error) {
	var conversions []entity.UnitConversion
	err := r.db.WithContext(ctx).
		Preload("FromUnit").
		Preload("ToUnit").
		Where("from_unit_id = ? OR to_unit_id = ?", unitID, unitID).
		Find(&conversions).Error
	return conversions, err
}

func (r *unitRepository) Convert(ctx context.Context, value float64, fromUnitID, toUnitID uuid.UUID) (float64, error) {
	// Same unit, no conversion needed
	if fromUnitID == toUnitID {
		return value, nil
	}

	// Try direct conversion
	conv, err := r.GetConversion(ctx, fromUnitID, toUnitID)
	if err == nil {
		return value * conv.ConversionFactor, nil
	}

	// Try via base units
	fromUnit, err := r.GetByID(ctx, fromUnitID)
	if err != nil {
		return 0, fmt.Errorf("from unit not found: %w", err)
	}
	
	toUnit, err := r.GetByID(ctx, toUnitID)
	if err != nil {
		return 0, fmt.Errorf("to unit not found: %w", err)
	}

	// Check same type
	if fromUnit.UoMType != toUnit.UoMType {
		return 0, fmt.Errorf("cannot convert between different unit types: %s to %s", fromUnit.UoMType, toUnit.UoMType)
	}

	// Convert via conversion factors
	// from -> base -> to
	// value * fromConversionFactor / toConversionFactor
	if fromUnit.ConversionFactor == 0 || toUnit.ConversionFactor == 0 {
		return 0, fmt.Errorf("invalid conversion factors")
	}

	baseValue := value * fromUnit.ConversionFactor
	result := baseValue / toUnit.ConversionFactor
	
	return result, nil
}
