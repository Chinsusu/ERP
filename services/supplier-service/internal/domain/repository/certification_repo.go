package repository

import (
	"context"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CertificationFilter represents filters for listing certifications
type CertificationFilter struct {
	SupplierID *uuid.UUID
	Type       string
	Status     string
	ExpiringIn *int // days until expiry
}

// CertificationRepository defines the interface for certification data access
type CertificationRepository interface {
	// CRUD operations
	Create(ctx context.Context, cert *entity.Certification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Certification, error)
	Update(ctx context.Context, cert *entity.Certification) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List operations
	GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Certification, error)
	GetExpiring(ctx context.Context, days int) ([]*entity.Certification, error)
	GetExpired(ctx context.Context) ([]*entity.Certification, error)
	
	// Batch operations for scheduled jobs
	UpdateExpiredStatuses(ctx context.Context, cutoffDate time.Time) (int64, error)
	UpdateExpiringStatuses(ctx context.Context, cutoffDate time.Time) (int64, error)
}
