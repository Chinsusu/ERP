package repository

import (
	"context"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/google/uuid"
)

// QuotationFilter defines filter options for quotations
type QuotationFilter struct {
	CustomerID  *uuid.UUID
	Status      entity.QuotationStatus
	DateFrom    string
	DateTo      string
	ExpiringSoon bool
	Page        int
	Limit       int
}

// QuotationRepository defines quotation repository interface
type QuotationRepository interface {
	// CRUD
	Create(ctx context.Context, quotation *entity.Quotation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Quotation, error)
	GetByNumber(ctx context.Context, number string) (*entity.Quotation, error)
	Update(ctx context.Context, quotation *entity.Quotation) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *QuotationFilter) ([]*entity.Quotation, int64, error)

	// Number generation
	GetNextQuotationNumber(ctx context.Context) (string, error)

	// Line items
	CreateLineItem(ctx context.Context, item *entity.QuotationLineItem) error
	UpdateLineItem(ctx context.Context, item *entity.QuotationLineItem) error
	DeleteLineItem(ctx context.Context, itemID uuid.UUID) error
	GetLineItems(ctx context.Context, quotationID uuid.UUID) ([]*entity.QuotationLineItem, error)

	// Status updates
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.QuotationStatus) error

	// Expiring quotations
	GetExpiringQuotations(ctx context.Context, days int) ([]*entity.Quotation, error)
	MarkExpiredQuotations(ctx context.Context) (int64, error)
}
