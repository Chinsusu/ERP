package material

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/erp-cosmetics/master-data-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// UseCase handles material business logic
type UseCase struct {
	repo           repository.MaterialRepository
	categoryRepo   repository.CategoryRepository
	unitRepo       repository.UnitRepository
	publisher      *event.Publisher
	autoGenCodes   bool
}

// NewUseCase creates a new material use case
func NewUseCase(
	repo repository.MaterialRepository,
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

// CreateRequest represents material creation request
type CreateRequest struct {
	Code             string
	Name             string
	NameEN           string
	Description      string
	MaterialType     entity.MaterialType
	CategoryID       *uuid.UUID
	BaseUnitID       uuid.UUID
	PurchaseUnitID   *uuid.UUID
	StockUnitID      *uuid.UUID
	INCIName         string
	CASNumber        string
	IsAllergen       bool
	AllergenInfo     string
	IsOrganic        bool
	IsNatural        bool
	IsVegan          bool
	OriginCountry    string
	StorageCondition entity.StorageCondition
	MinTemp          *float64
	MaxTemp          *float64
	ShelfLifeDays    int
	IsHazardous      bool
	LeadTimeDays     int
	MinOrderQty      float64
	ReorderPoint     float64
	SafetyStock      float64
	StandardCost     float64
	Currency         string
}

// Create creates a new material
func (uc *UseCase) Create(ctx context.Context, req *CreateRequest) (*entity.Material, error) {
	// Generate code if auto-generation is enabled
	code := req.Code
	if code == "" && uc.autoGenCodes {
		seq, err := uc.repo.GetNextSequence(ctx, req.MaterialType)
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}
		prefix := "MAT"
		switch req.MaterialType {
		case entity.MaterialTypeRaw:
			prefix = "RM"
		case entity.MaterialTypePackaging:
			prefix = "PKG"
		case entity.MaterialTypeConsumable:
			prefix = "CON"
		case entity.MaterialTypeSemiFinished:
			prefix = "SF"
		}
		code = fmt.Sprintf("%s-%04d", prefix, seq)
	}

	// Check if code already exists
	existing, _ := uc.repo.GetByCode(ctx, code)
	if existing != nil {
		return nil, fmt.Errorf("material with code %s already exists", code)
	}

	// Validate base unit exists
	_, err := uc.unitRepo.GetByID(ctx, req.BaseUnitID)
	if err != nil {
		return nil, fmt.Errorf("base unit not found: %w", err)
	}

	// Default values
	if req.StorageCondition == "" {
		req.StorageCondition = entity.StorageAmbient
	}
	if req.ShelfLifeDays == 0 {
		req.ShelfLifeDays = 365
	}
	if req.LeadTimeDays == 0 {
		req.LeadTimeDays = 14
	}
	if req.Currency == "" {
		req.Currency = "VND"
	}

	material := &entity.Material{
		Code:                code,
		Name:                req.Name,
		NameEN:              req.NameEN,
		Description:         req.Description,
		MaterialType:        req.MaterialType,
		CategoryID:          req.CategoryID,
		BaseUnitID:          req.BaseUnitID,
		PurchaseUnitID:      req.PurchaseUnitID,
		StockUnitID:         req.StockUnitID,
		INCIName:            req.INCIName,
		CASNumber:           req.CASNumber,
		IsAllergen:          req.IsAllergen,
		AllergenInfo:        req.AllergenInfo,
		IsOrganic:           req.IsOrganic,
		IsNatural:           req.IsNatural,
		IsVegan:             req.IsVegan,
		OriginCountry:       req.OriginCountry,
		StorageCondition:    req.StorageCondition,
		MinTemp:             req.MinTemp,
		MaxTemp:             req.MaxTemp,
		ShelfLifeDays:       req.ShelfLifeDays,
		IsHazardous:         req.IsHazardous,
		LeadTimeDays:        req.LeadTimeDays,
		MinOrderQty:         req.MinOrderQty,
		ReorderPoint:        req.ReorderPoint,
		SafetyStock:         req.SafetyStock,
		StandardCost:        req.StandardCost,
		Currency:            req.Currency,
		Status:              "active",
	}

	if err := material.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, material); err != nil {
		return nil, fmt.Errorf("failed to create material: %w", err)
	}

	// Publish event
	uc.publisher.Publish(event.EventMaterialCreated, event.MaterialCreatedEvent{
		MaterialID:   material.ID.String(),
		MaterialCode: material.Code,
		Name:         material.Name,
		MaterialType: string(material.MaterialType),
	})

	return material, nil
}

// GetByID retrieves a material by ID
func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Material, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByIDs retrieves materials by IDs
func (uc *UseCase) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Material, error) {
	return uc.repo.GetByIDs(ctx, ids)
}

// Update updates an existing material
func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*entity.Material, error) {
	material, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("material not found: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		material.Name = name
	}
	if desc, ok := updates["description"].(string); ok {
		material.Description = desc
	}
	if inci, ok := updates["inci_name"].(string); ok {
		material.INCIName = inci
	}
	if cas, ok := updates["cas_number"].(string); ok {
		material.CASNumber = cas
	}
	if status, ok := updates["status"].(string); ok {
		material.Status = status
	}
	// Add more fields as needed...

	if err := uc.repo.Update(ctx, material); err != nil {
		return nil, fmt.Errorf("failed to update material: %w", err)
	}

	return material, nil
}

// Delete soft deletes a material
func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// List lists materials with filters
func (uc *UseCase) List(ctx context.Context, filter *repository.MaterialFilter) ([]entity.Material, int64, error) {
	return uc.repo.List(ctx, filter)
}

// Search searches materials
func (uc *UseCase) Search(ctx context.Context, query string, filter *repository.MaterialFilter) ([]entity.Material, int64, error) {
	return uc.repo.Search(ctx, query, filter)
}

// AddSpecification adds a specification to a material
func (uc *UseCase) AddSpecification(ctx context.Context, materialID uuid.UUID, spec *entity.MaterialSpecification) error {
	spec.MaterialID = materialID
	return uc.repo.AddSpecification(ctx, spec)
}

// GetSpecifications gets specifications for a material
func (uc *UseCase) GetSpecifications(ctx context.Context, materialID uuid.UUID) ([]entity.MaterialSpecification, error) {
	return uc.repo.GetSpecifications(ctx, materialID)
}
