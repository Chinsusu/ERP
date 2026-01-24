package repository

import (
	"context"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/google/uuid"
)

// PRFilter represents filters for listing PRs
type PRFilter struct {
	Status      string
	RequesterID *uuid.UUID
	Priority    string
	DateFrom    string
	DateTo      string
	Search      string
	Page        int
	Limit       int
}

// PRRepository defines the interface for PR data access
type PRRepository interface {
	// CRUD
	Create(ctx context.Context, pr *entity.PurchaseRequisition) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PurchaseRequisition, error)
	GetByNumber(ctx context.Context, prNumber string) (*entity.PurchaseRequisition, error)
	Update(ctx context.Context, pr *entity.PurchaseRequisition) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List
	List(ctx context.Context, filter *PRFilter) ([]*entity.PurchaseRequisition, int64, error)
	
	// Code generation
	GetNextPRNumber(ctx context.Context) (string, error)
	
	// Line items
	CreateLineItem(ctx context.Context, item *entity.PRLineItem) error
	UpdateLineItem(ctx context.Context, item *entity.PRLineItem) error
	DeleteLineItem(ctx context.Context, id uuid.UUID) error
	GetLineItemsByPRID(ctx context.Context, prID uuid.UUID) ([]*entity.PRLineItem, error)
	
	// Approvals
	CreateApproval(ctx context.Context, approval *entity.PRApproval) error
	GetApprovalsByPRID(ctx context.Context, prID uuid.UUID) ([]*entity.PRApproval, error)
}
