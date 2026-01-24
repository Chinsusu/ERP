package product

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/erp-cosmetics/master-data-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// UseCase handles product business logic
type UseCase struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
	unitRepo     repository.UnitRepository
	publisher    *event.Publisher
	autoGenCodes bool
}

// NewUseCase creates a new product use case
func NewUseCase(
	repo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	unitRepo repository.UnitRepository,
	publisher *event.Publisher,
	autoGenCodes bool,
) *UseCase {
	return &UseCase{
		repo:         repo,
		categoryRepo: categoryRepo,
		unitRepo:     unitRepo,
		publisher:    publisher,
		autoGenCodes: autoGenCodes,
	}
}

// CreateRequest represents product creation request
type CreateRequest struct {
	Code                   string
	SKU                    string
	Barcode                string
	Name                   string
	NameEN                 string
	Description            string
	CategoryID             *uuid.UUID
	ProductLine            string
	Brand                  string
	Volume                 *float64
	VolumeUnit             string
	CosmeticLicenseNumber  string
	IngredientsSummary     string
	TargetSkinType         string
	UsageInstructions      string
	PackagingType          string
	StandardCost           float64
	StandardPrice          float64
	RecommendedRetailPrice *float64
	Currency               string
	BaseUnitID             uuid.UUID
	SalesUnitID            *uuid.UUID
	ShelfLifeMonths        int
}

// Create creates a new product
func (uc *UseCase) Create(ctx context.Context, req *CreateRequest) (*entity.Product, error) {
	// Generate code if auto-generation is enabled
	code := req.Code
	if code == "" && uc.autoGenCodes && req.CategoryID != nil {
		category, err := uc.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err == nil {
			seq, err := uc.repo.GetNextSequence(ctx, category.Code)
			if err == nil {
				code = fmt.Sprintf("FG-%s-%04d", category.Code, seq)
			}
		}
		if code == "" {
			seq, _ := uc.repo.GetNextSequence(ctx, "PROD")
			code = fmt.Sprintf("FG-PROD-%04d", seq)
		}
	}

	// Check if code already exists
	existing, _ := uc.repo.GetByCode(ctx, code)
	if existing != nil {
		return nil, fmt.Errorf("product with code %s already exists", code)
	}

	// Check if SKU already exists
	existing, _ = uc.repo.GetBySKU(ctx, req.SKU)
	if existing != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
	}

	// Validate base unit exists
	_, err := uc.unitRepo.GetByID(ctx, req.BaseUnitID)
	if err != nil {
		return nil, fmt.Errorf("base unit not found: %w", err)
	}

	// Default values
	if req.Currency == "" {
		req.Currency = "VND"
	}
	if req.ShelfLifeMonths == 0 {
		req.ShelfLifeMonths = 24
	}

	product := &entity.Product{
		Code:                   code,
		SKU:                    req.SKU,
		Barcode:                req.Barcode,
		Name:                   req.Name,
		NameEN:                 req.NameEN,
		Description:            req.Description,
		CategoryID:             req.CategoryID,
		ProductLine:            req.ProductLine,
		Brand:                  req.Brand,
		Volume:                 req.Volume,
		VolumeUnit:             req.VolumeUnit,
		CosmeticLicenseNumber:  req.CosmeticLicenseNumber,
		IngredientsSummary:     req.IngredientsSummary,
		TargetSkinType:         req.TargetSkinType,
		UsageInstructions:      req.UsageInstructions,
		PackagingType:          req.PackagingType,
		StandardCost:           req.StandardCost,
		StandardPrice:          req.StandardPrice,
		RecommendedRetailPrice: req.RecommendedRetailPrice,
		Currency:               req.Currency,
		BaseUnitID:             req.BaseUnitID,
		SalesUnitID:            req.SalesUnitID,
		ShelfLifeMonths:        req.ShelfLifeMonths,
		Status:                 "active",
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Publish event
	uc.publisher.Publish(event.EventProductCreated, event.ProductCreatedEvent{
		ProductID:   product.ID.String(),
		ProductCode: product.Code,
		SKU:         product.SKU,
		Name:        product.Name,
	})

	return product, nil
}

// GetByID retrieves a product by ID
func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByIDs retrieves products by IDs
func (uc *UseCase) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Product, error) {
	return uc.repo.GetByIDs(ctx, ids)
}

// Update updates an existing product
func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*entity.Product, error) {
	product, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		product.Name = name
	}
	if desc, ok := updates["description"].(string); ok {
		product.Description = desc
	}
	if status, ok := updates["status"].(string); ok {
		product.Status = status
	}
	// Add more fields as needed...

	if err := uc.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// Delete soft deletes a product
func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// List lists products with filters
func (uc *UseCase) List(ctx context.Context, filter *repository.ProductFilter) ([]entity.Product, int64, error) {
	return uc.repo.List(ctx, filter)
}

// Search searches products
func (uc *UseCase) Search(ctx context.Context, query string, filter *repository.ProductFilter) ([]entity.Product, int64, error) {
	return uc.repo.Search(ctx, query, filter)
}

// GetByCategory gets products by category
func (uc *UseCase) GetByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]entity.Product, int64, error) {
	return uc.repo.GetByCategory(ctx, categoryID, page, pageSize)
}

// AddImage adds an image to a product
func (uc *UseCase) AddImage(ctx context.Context, productID uuid.UUID, image *entity.ProductImage) error {
	image.ProductID = productID
	return uc.repo.AddImage(ctx, image)
}

// GetImages gets images for a product
func (uc *UseCase) GetImages(ctx context.Context, productID uuid.UUID) ([]entity.ProductImage, error) {
	return uc.repo.GetImages(ctx, productID)
}

// SetPrimaryImage sets the primary image for a product
func (uc *UseCase) SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error {
	return uc.repo.SetPrimaryImage(ctx, productID, imageID)
}
