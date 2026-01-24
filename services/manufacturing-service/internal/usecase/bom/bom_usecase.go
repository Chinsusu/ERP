package bom

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateBOMUseCase handles BOM creation
type CreateBOMUseCase struct {
	repo          repository.BOMRepository
	eventPub      *event.Publisher
	encryptionKey []byte
}

// NewCreateBOMUseCase creates a new CreateBOMUseCase
func NewCreateBOMUseCase(repo repository.BOMRepository, eventPub *event.Publisher, encryptionKey []byte) *CreateBOMUseCase {
	return &CreateBOMUseCase{
		repo:          repo,
		eventPub:      eventPub,
		encryptionKey: encryptionKey,
	}
}

// CreateBOMInput is the input for creating a BOM
type CreateBOMInput struct {
	BOMNumber            string
	ProductID            uuid.UUID
	Version              int
	Name                 string
	Description          string
	BatchSize            float64
	BatchUnitID          uuid.UUID
	ConfidentialityLevel entity.ConfidentialityLevel
	LaborCost            float64
	OverheadCost         float64
	FormulaDetails       *entity.FormulaDetails
	Items                []CreateBOMItemInput
	CreatedBy            uuid.UUID
}

// CreateBOMItemInput is input for a BOM line item
type CreateBOMItemInput struct {
	LineNumber      int
	MaterialID      uuid.UUID
	ItemType        entity.BOMItemType
	Quantity        float64
	UOMID           uuid.UUID
	QuantityMin     *float64
	QuantityMax     *float64
	IsCritical      bool
	ScrapPercentage float64
	UnitCost        float64
	Notes           string
}

// Execute creates a new BOM
func (uc *CreateBOMUseCase) Execute(ctx context.Context, input CreateBOMInput) (*entity.BOM, error) {
	// Encrypt formula details if provided
	var encryptedFormula []byte
	var err error
	if input.FormulaDetails != nil {
		encryptedFormula, err = entity.EncryptFormula(input.FormulaDetails, uc.encryptionKey)
		if err != nil {
			return nil, err
		}
	}

	bom := &entity.BOM{
		BOMNumber:            input.BOMNumber,
		ProductID:            input.ProductID,
		Version:              input.Version,
		Name:                 input.Name,
		Description:          input.Description,
		Status:               entity.BOMStatusDraft,
		BatchSize:            input.BatchSize,
		BatchUnitID:          input.BatchUnitID,
		FormulaDetails:       encryptedFormula,
		ConfidentialityLevel: input.ConfidentialityLevel,
		LaborCost:            input.LaborCost,
		OverheadCost:         input.OverheadCost,
		CreatedBy:            &input.CreatedBy,
		UpdatedBy:            &input.CreatedBy,
	}

	// Create line items
	var items []entity.BOMLineItem
	var materialCost float64
	for _, item := range input.Items {
		totalCost := item.Quantity * item.UnitCost
		materialCost += totalCost
		items = append(items, entity.BOMLineItem{
			LineNumber:      item.LineNumber,
			MaterialID:      item.MaterialID,
			ItemType:        item.ItemType,
			Quantity:        item.Quantity,
			UOMID:           item.UOMID,
			QuantityMin:     item.QuantityMin,
			QuantityMax:     item.QuantityMax,
			IsCritical:      item.IsCritical,
			ScrapPercentage: item.ScrapPercentage,
			UnitCost:        item.UnitCost,
			TotalCost:       totalCost,
			Notes:           item.Notes,
		})
	}
	bom.Items = items
	bom.MaterialCost = materialCost
	bom.TotalCost = materialCost + bom.LaborCost + bom.OverheadCost

	if err := uc.repo.Create(ctx, bom); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishBOMCreated(event.BOMEvent{
		BOMID:     bom.ID.String(),
		BOMNumber: bom.BOMNumber,
		ProductID: bom.ProductID.String(),
		Version:   bom.Version,
		Status:    string(bom.Status),
	})

	return bom, nil
}

// GetBOMUseCase handles getting a BOM
type GetBOMUseCase struct {
	repo          repository.BOMRepository
	encryptionKey []byte
}

// NewGetBOMUseCase creates a new GetBOMUseCase
func NewGetBOMUseCase(repo repository.BOMRepository, encryptionKey []byte) *GetBOMUseCase {
	return &GetBOMUseCase{
		repo:          repo,
		encryptionKey: encryptionKey,
	}
}

// BOMResponse is the response for getting a BOM
type BOMResponse struct {
	BOM            *entity.BOM
	FormulaDetails *entity.FormulaDetails // Decrypted, only if user has permission
	CanViewFormula bool
}

// Execute gets a BOM by ID
func (uc *GetBOMUseCase) Execute(ctx context.Context, id uuid.UUID, canViewFormula bool) (*BOMResponse, error) {
	bom, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, entity.ErrBOMNotFound
	}

	response := &BOMResponse{
		BOM:            bom,
		CanViewFormula: canViewFormula,
	}

	// Decrypt formula if user has permission
	if canViewFormula && len(bom.FormulaDetails) > 0 {
		formula, err := entity.DecryptFormula(bom.FormulaDetails, uc.encryptionKey)
		if err == nil {
			response.FormulaDetails = formula
		}
	}

	return response, nil
}

// ListBOMsUseCase handles listing BOMs
type ListBOMsUseCase struct {
	repo repository.BOMRepository
}

// NewListBOMsUseCase creates a new ListBOMsUseCase
func NewListBOMsUseCase(repo repository.BOMRepository) *ListBOMsUseCase {
	return &ListBOMsUseCase{repo: repo}
}

// Execute lists BOMs
func (uc *ListBOMsUseCase) Execute(ctx context.Context, filter repository.BOMFilter) ([]*entity.BOM, int64, error) {
	return uc.repo.List(ctx, filter)
}

// ApproveBOMUseCase handles BOM approval
type ApproveBOMUseCase struct {
	repo     repository.BOMRepository
	eventPub *event.Publisher
}

// NewApproveBOMUseCase creates a new ApproveBOMUseCase
func NewApproveBOMUseCase(repo repository.BOMRepository, eventPub *event.Publisher) *ApproveBOMUseCase {
	return &ApproveBOMUseCase{repo: repo, eventPub: eventPub}
}

// Execute approves a BOM
func (uc *ApproveBOMUseCase) Execute(ctx context.Context, bomID uuid.UUID, approverID uuid.UUID) (*entity.BOM, error) {
	bom, err := uc.repo.GetByID(ctx, bomID)
	if err != nil {
		return nil, entity.ErrBOMNotFound
	}

	// Submit if in draft
	if bom.IsDraft() {
		if err := bom.Submit(); err != nil {
			return nil, err
		}
	}

	if err := bom.Approve(approverID); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, bom); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishBOMApproved(event.BOMEvent{
		BOMID:     bom.ID.String(),
		BOMNumber: bom.BOMNumber,
		ProductID: bom.ProductID.String(),
		Version:   bom.Version,
		Status:    string(bom.Status),
	})

	return bom, nil
}

// GetActiveBOMUseCase gets the active BOM for a product
type GetActiveBOMUseCase struct {
	repo repository.BOMRepository
}

// NewGetActiveBOMUseCase creates a new GetActiveBOMUseCase
func NewGetActiveBOMUseCase(repo repository.BOMRepository) *GetActiveBOMUseCase {
	return &GetActiveBOMUseCase{repo: repo}
}

// Execute gets the active BOM for a product
func (uc *GetActiveBOMUseCase) Execute(ctx context.Context, productID uuid.UUID) (*entity.BOM, error) {
	return uc.repo.GetActiveBOMForProduct(ctx, productID)
}
