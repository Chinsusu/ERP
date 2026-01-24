package repository

import (
	"context"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CategoryRepository defines category data access methods
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetByCode(ctx context.Context, code string) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *CategoryFilter) ([]entity.Category, int64, error)
	GetTree(ctx context.Context, categoryType entity.CategoryType) ([]entity.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]entity.Category, error)
}

// CategoryFilter for listing categories
type CategoryFilter struct {
	CategoryType entity.CategoryType
	ParentID     *uuid.UUID
	Status       string
	Search       string
	Page         int
	PageSize     int
}
