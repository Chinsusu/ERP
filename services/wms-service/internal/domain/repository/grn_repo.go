package repository

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/google/uuid"
)

// GRNRepository defines GRN repository interface
type GRNRepository interface {
	Create(ctx context.Context, grn *entity.GRN) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.GRN, error)
	GetByNumber(ctx context.Context, grnNumber string) (*entity.GRN, error)
	GetByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.GRN, error)
	List(ctx context.Context, filter *GRNFilter) ([]*entity.GRN, int64, error)
	Update(ctx context.Context, grn *entity.GRN) error
	
	// Line items
	CreateLineItem(ctx context.Context, item *entity.GRNLineItem) error
	GetLineItemsByGRNID(ctx context.Context, grnID uuid.UUID) ([]*entity.GRNLineItem, error)
	UpdateLineItem(ctx context.Context, item *entity.GRNLineItem) error
	
	// Number generation
	GetNextGRNNumber(ctx context.Context) (string, error)
}

// GRNFilter defines filter options for GRNs
type GRNFilter struct {
	WarehouseID *uuid.UUID
	SupplierID  *uuid.UUID
	POID        *uuid.UUID
	Status      string
	QCStatus    string
	Search      string
	Page        int
	Limit       int
}

// GoodsIssueRepository defines goods issue repository interface
type GoodsIssueRepository interface {
	Create(ctx context.Context, issue *entity.GoodsIssue) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.GoodsIssue, error)
	GetByNumber(ctx context.Context, issueNumber string) (*entity.GoodsIssue, error)
	List(ctx context.Context, filter *GoodsIssueFilter) ([]*entity.GoodsIssue, int64, error)
	Update(ctx context.Context, issue *entity.GoodsIssue) error
	
	// Line items
	CreateLineItem(ctx context.Context, item *entity.GILineItem) error
	CreateLineItems(ctx context.Context, items []*entity.GILineItem) error
	GetLineItemsByIssueID(ctx context.Context, issueID uuid.UUID) ([]*entity.GILineItem, error)
	
	// Number generation
	GetNextIssueNumber(ctx context.Context) (string, error)
}

// GoodsIssueFilter defines filter options for goods issues
type GoodsIssueFilter struct {
	WarehouseID   *uuid.UUID
	IssueType     string
	ReferenceType string
	Status        string
	Search        string
	Page          int
	Limit         int
}
