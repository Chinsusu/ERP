package postgres

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	// If parent exists, calculate path and level
	if category.ParentID != nil {
		parent, err := r.GetByID(ctx, *category.ParentID)
		if err != nil {
			return fmt.Errorf("parent category not found: %w", err)
		}
		category.Level = parent.Level + 1
		category.GeneratePath(parent.Path)
	} else {
		category.Level = 0
		category.GeneratePath("")
	}
	
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByCode(ctx context.Context, code string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).
		First(&category, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) List(ctx context.Context, filter *repository.CategoryFilter) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Category{})

	// Apply filters
	if filter.CategoryType != "" {
		query = query.Where("category_type = ?", filter.CategoryType)
	}
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	} else if filter.ParentID == nil && filter.Search == "" {
		// Only show root categories by default
		query = query.Where("parent_id IS NULL")
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

	// Order by sort_order
	query = query.Order("sort_order ASC, name ASC")

	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) GetTree(ctx context.Context, categoryType entity.CategoryType) ([]entity.Category, error) {
	var rootCategories []entity.Category

	query := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Where("deleted_at IS NULL").
		Order("sort_order ASC, name ASC")

	if categoryType != "" {
		query = query.Where("category_type = ?", categoryType)
	}

	if err := query.Find(&rootCategories).Error; err != nil {
		return nil, err
	}

	// Load all children recursively
	for i := range rootCategories {
		if err := r.loadChildren(ctx, &rootCategories[i]); err != nil {
			return nil, err
		}
	}

	return rootCategories, nil
}

func (r *categoryRepository) loadChildren(ctx context.Context, parent *entity.Category) error {
	var children []entity.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parent.ID).
		Where("deleted_at IS NULL").
		Order("sort_order ASC, name ASC").
		Find(&children).Error
	if err != nil {
		return err
	}

	parent.Children = children

	// Recursively load children
	for i := range parent.Children {
		if err := r.loadChildren(ctx, &parent.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

func (r *categoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]entity.Category, error) {
	var children []entity.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Where("deleted_at IS NULL").
		Order("sort_order ASC, name ASC").
		Find(&children).Error
	return children, err
}
