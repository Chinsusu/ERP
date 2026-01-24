package repository

import (
	"context"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/google/uuid"
)

// ProductRepository defines product data access methods
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetByCode(ctx context.Context, code string) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *ProductFilter) ([]entity.Product, int64, error)
	Search(ctx context.Context, query string, filter *ProductFilter) ([]entity.Product, int64, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]entity.Product, int64, error)
	GetNextSequence(ctx context.Context, categoryCode string) (int, error)
	
	// Images
	AddImage(ctx context.Context, image *entity.ProductImage) error
	GetImages(ctx context.Context, productID uuid.UUID) ([]entity.ProductImage, error)
	UpdateImage(ctx context.Context, image *entity.ProductImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
	SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error
}

// ProductFilter for listing and searching products
type ProductFilter struct {
	CategoryID     *uuid.UUID
	ProductLine    string
	Brand          string
	Status         string
	LicenseExpiring int // Number of days until expiry
	Search         string
	Page           int
	PageSize       int
}
