package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("BaseUnit").
		Preload("SalesUnit").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByCode(ctx context.Context, code string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).
		First(&product, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).
		First(&product, "sku = ?", sku).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("BaseUnit").
		Where("id IN ?", ids).
		Find(&products).Error
	return products, err
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id).Error
}

func (r *productRepository) List(ctx context.Context, filter *repository.ProductFilter) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	query = r.applyProductFilters(query, filter)

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
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Search(ctx context.Context, queryStr string, filter *repository.ProductFilter) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	// Apply search
	if queryStr != "" {
		search := "%" + queryStr + "%"
		query = query.Where(
			"name ILIKE ? OR name_en ILIKE ? OR code ILIKE ? OR sku ILIKE ? OR barcode ILIKE ?",
			search, search, search, search, search,
		)
	}

	query = r.applyProductFilters(query, filter)

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
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) applyProductFilters(query *gorm.DB, filter *repository.ProductFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.ProductLine != "" {
		query = query.Where("product_line = ?", filter.ProductLine)
	}
	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.LicenseExpiring > 0 {
		expiryDate := time.Now().AddDate(0, 0, filter.LicenseExpiring)
		query = query.Where("license_expiry_date IS NOT NULL AND license_expiry_date <= ?", expiryDate)
	}

	return query
}

func (r *productRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]entity.Product, int64, error) {
	filter := &repository.ProductFilter{
		CategoryID: &categoryID,
		Page:       page,
		PageSize:   pageSize,
	}
	return r.List(ctx, filter)
}

func (r *productRepository) GetNextSequence(ctx context.Context, categoryCode string) (int, error) {
	var count int64
	pattern := fmt.Sprintf("FG-%s-%%", categoryCode)
	
	err := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("code LIKE ?", pattern).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count) + 1, nil
}

func (r *productRepository) AddImage(ctx context.Context, image *entity.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *productRepository) GetImages(ctx context.Context, productID uuid.UUID) ([]entity.ProductImage, error) {
	var images []entity.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

func (r *productRepository) UpdateImage(ctx context.Context, image *entity.ProductImage) error {
	return r.db.WithContext(ctx).Save(image).Error
}

func (r *productRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductImage{}, "id = ?", id).Error
}

func (r *productRepository) SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Reset all images to non-primary
		if err := tx.Model(&entity.ProductImage{}).
			Where("product_id = ?", productID).
			Update("is_primary", false).Error; err != nil {
			return err
		}

		// Set the specified image as primary
		return tx.Model(&entity.ProductImage{}).
			Where("id = ?", imageID).
			Update("is_primary", true).Error
	})
}
