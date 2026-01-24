package category

import (
	"context"
	"fmt"

	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	"github.com/erp-cosmetics/master-data-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// UseCase handles category business logic
type UseCase struct {
	repo      repository.CategoryRepository
	publisher *event.Publisher
}

// NewUseCase creates a new category use case
func NewUseCase(repo repository.CategoryRepository, publisher *event.Publisher) *UseCase {
	return &UseCase{
		repo:      repo,
		publisher: publisher,
	}
}

// CreateRequest represents category creation request
type CreateRequest struct {
	Code         string
	Name         string
	NameEN       string
	Description  string
	CategoryType entity.CategoryType
	ParentID     *uuid.UUID
	SortOrder    int
}

// Create creates a new category
func (uc *UseCase) Create(ctx context.Context, req *CreateRequest) (*entity.Category, error) {
	// Check if code already exists
	existing, _ := uc.repo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("category with code %s already exists", req.Code)
	}

	category := &entity.Category{
		Code:         req.Code,
		Name:         req.Name,
		NameEN:       req.NameEN,
		Description:  req.Description,
		CategoryType: req.CategoryType,
		ParentID:     req.ParentID,
		SortOrder:    req.SortOrder,
		Status:       "active",
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Publish event
	uc.publisher.Publish(event.EventCategoryCreated, event.CategoryCreatedEvent{
		CategoryID:   category.ID.String(),
		Code:         category.Code,
		Name:         category.Name,
		CategoryType: string(category.CategoryType),
	})

	return category, nil
}

// GetByID retrieves a category by ID
func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return uc.repo.GetByID(ctx, id)
}

// UpdateRequest represents category update request
type UpdateRequest struct {
	ID          uuid.UUID
	Name        string
	NameEN      string
	Description string
	SortOrder   int
	Status      string
}

// Update updates an existing category
func (uc *UseCase) Update(ctx context.Context, req *UpdateRequest) (*entity.Category, error) {
	category, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.NameEN != "" {
		category.NameEN = req.NameEN
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.SortOrder > 0 {
		category.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		category.Status = req.Status
	}

	if err := uc.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// Delete soft deletes a category
func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// List lists categories with filters
func (uc *UseCase) List(ctx context.Context, filter *repository.CategoryFilter) ([]entity.Category, int64, error) {
	return uc.repo.List(ctx, filter)
}

// GetTree retrieves categories as a hierarchical tree
func (uc *UseCase) GetTree(ctx context.Context, categoryType entity.CategoryType) ([]entity.Category, error) {
	return uc.repo.GetTree(ctx, categoryType)
}
