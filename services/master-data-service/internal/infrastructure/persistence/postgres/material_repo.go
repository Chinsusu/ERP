package postgres

import (
	"context"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type materialRepository struct {
	db *gorm.DB
}

// NewMaterialRepository creates a new material repository
func NewMaterialRepository(db *gorm.DB) repository.MaterialRepository {
	return &materialRepository{db: db}
}

func (r *materialRepository) Create(ctx context.Context, material *entity.Material) error {
	return r.db.WithContext(ctx).Create(material).Error
}

func (r *materialRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Material, error) {
	var material entity.Material
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("BaseUnit").
		Preload("PurchaseUnit").
		Preload("StockUnit").
		Preload("Specifications").
		First(&material, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *materialRepository) GetByCode(ctx context.Context, code string) (*entity.Material, error) {
	var material entity.Material
	err := r.db.WithContext(ctx).
		First(&material, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *materialRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Material, error) {
	var materials []entity.Material
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("BaseUnit").
		Where("id IN ?", ids).
		Find(&materials).Error
	return materials, err
}

func (r *materialRepository) Update(ctx context.Context, material *entity.Material) error {
	return r.db.WithContext(ctx).Save(material).Error
}

func (r *materialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Material{}, "id = ?", id).Error
}

func (r *materialRepository) List(ctx context.Context, filter *repository.MaterialFilter) ([]entity.Material, int64, error) {
	var materials []entity.Material
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Material{})

	query = r.applyMaterialFilters(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	query = query.Order("code ASC")

	if err := query.
		Preload("Category").
		Preload("BaseUnit").
		Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

func (r *materialRepository) Search(ctx context.Context, queryStr string, filter *repository.MaterialFilter) ([]entity.Material, int64, error) {
	var materials []entity.Material
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Material{})

	// Apply search
	if queryStr != "" {
		search := "%" + queryStr + "%"
		query = query.Where(
			"name ILIKE ? OR name_en ILIKE ? OR code ILIKE ? OR inci_name ILIKE ? OR cas_number ILIKE ?",
			search, search, search, search, search,
		)
	}

	query = r.applyMaterialFilters(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if filter != nil && filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	query = query.Order("name ASC")

	if err := query.
		Preload("Category").
		Preload("BaseUnit").
		Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

func (r *materialRepository) applyMaterialFilters(query *gorm.DB, filter *repository.MaterialFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.MaterialType != "" {
		query = query.Where("material_type = ?", filter.MaterialType)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.StorageCondition != "" {
		query = query.Where("storage_condition = ?", filter.StorageCondition)
	}
	if filter.IsOrganic != nil {
		query = query.Where("is_organic = ?", *filter.IsOrganic)
	}
	if filter.IsNatural != nil {
		query = query.Where("is_natural = ?", *filter.IsNatural)
	}
	if filter.IsAllergen != nil {
		query = query.Where("is_allergen = ?", *filter.IsAllergen)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	return query
}

func (r *materialRepository) GetNextSequence(ctx context.Context, materialType entity.MaterialType) (int, error) {
	var count int64
	prefix := "MAT"
	switch materialType {
	case entity.MaterialTypeRaw:
		prefix = "RM"
	case entity.MaterialTypePackaging:
		prefix = "PKG"
	case entity.MaterialTypeConsumable:
		prefix = "CON"
	case entity.MaterialTypeSemiFinished:
		prefix = "SF"
	}
	
	err := r.db.WithContext(ctx).
		Model(&entity.Material{}).
		Where("code LIKE ?", prefix+"-%").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count) + 1, nil
}

func (r *materialRepository) AddSpecification(ctx context.Context, spec *entity.MaterialSpecification) error {
	return r.db.WithContext(ctx).Create(spec).Error
}

func (r *materialRepository) GetSpecifications(ctx context.Context, materialID uuid.UUID) ([]entity.MaterialSpecification, error) {
	var specs []entity.MaterialSpecification
	err := r.db.WithContext(ctx).
		Where("material_id = ?", materialID).
		Order("created_at ASC").
		Find(&specs).Error
	return specs, err
}

func (r *materialRepository) UpdateSpecification(ctx context.Context, spec *entity.MaterialSpecification) error {
	return r.db.WithContext(ctx).Save(spec).Error
}

func (r *materialRepository) DeleteSpecification(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.MaterialSpecification{}, "id = ?", id).Error
}
