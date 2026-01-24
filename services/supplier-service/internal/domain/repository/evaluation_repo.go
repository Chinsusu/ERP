package repository

import (
	"context"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/google/uuid"
)

// EvaluationFilter represents filters for listing evaluations
type EvaluationFilter struct {
	SupplierID *uuid.UUID
	Period     string
	Status     string
	Page       int
	Limit      int
}

// EvaluationRepository defines the interface for evaluation data access
type EvaluationRepository interface {
	// CRUD operations
	Create(ctx context.Context, eval *entity.Evaluation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Evaluation, error)
	Update(ctx context.Context, eval *entity.Evaluation) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List operations
	GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Evaluation, error)
	List(ctx context.Context, filter *EvaluationFilter) ([]*entity.Evaluation, int64, error)
	
	// Aggregations
	GetLatestBySupplierID(ctx context.Context, supplierID uuid.UUID) (*entity.Evaluation, error)
	GetAverageScoresBySupplierID(ctx context.Context, supplierID uuid.UUID) (quality, delivery, service float64, err error)
}
