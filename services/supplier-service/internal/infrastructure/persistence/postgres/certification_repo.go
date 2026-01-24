package postgres

import (
	"context"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type certificationRepository struct {
	db *gorm.DB
}

// NewCertificationRepository creates a new certification repository
func NewCertificationRepository(db *gorm.DB) repository.CertificationRepository {
	return &certificationRepository{db: db}
}

func (r *certificationRepository) Create(ctx context.Context, cert *entity.Certification) error {
	// Auto-update status before saving
	cert.UpdateStatus()
	return r.db.WithContext(ctx).Create(cert).Error
}

func (r *certificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Certification, error) {
	var cert entity.Certification
	err := r.db.WithContext(ctx).First(&cert, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	cert.UpdateStatus()
	return &cert, nil
}

func (r *certificationRepository) Update(ctx context.Context, cert *entity.Certification) error {
	cert.UpdateStatus()
	return r.db.WithContext(ctx).Save(cert).Error
}

func (r *certificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Certification{}, "id = ?", id).Error
}

func (r *certificationRepository) GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Certification, error) {
	var certs []*entity.Certification
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("expiry_date").
		Find(&certs).Error
	if err != nil {
		return nil, err
	}
	
	// Update status for each cert
	for _, cert := range certs {
		cert.UpdateStatus()
	}
	
	return certs, nil
}

func (r *certificationRepository) GetExpiring(ctx context.Context, days int) ([]*entity.Certification, error) {
	var certs []*entity.Certification
	cutoffDate := time.Now().AddDate(0, 0, days)
	
	err := r.db.WithContext(ctx).
		Where("expiry_date <= ? AND expiry_date >= CURRENT_DATE", cutoffDate).
		Preload("Supplier", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, name")
		}).
		Order("expiry_date").
		Find(&certs).Error
	if err != nil {
		return nil, err
	}
	
	for _, cert := range certs {
		cert.UpdateStatus()
	}
	
	return certs, nil
}

func (r *certificationRepository) GetExpired(ctx context.Context) ([]*entity.Certification, error) {
	var certs []*entity.Certification
	
	err := r.db.WithContext(ctx).
		Where("expiry_date < CURRENT_DATE").
		Order("expiry_date DESC").
		Find(&certs).Error
	if err != nil {
		return nil, err
	}
	
	for _, cert := range certs {
		cert.UpdateStatus()
	}
	
	return certs, nil
}

func (r *certificationRepository) UpdateExpiredStatuses(ctx context.Context, cutoffDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&entity.Certification{}).
		Where("expiry_date < ? AND status != ?", cutoffDate, entity.CertStatusExpired).
		Update("status", entity.CertStatusExpired)
	return result.RowsAffected, result.Error
}

func (r *certificationRepository) UpdateExpiringStatuses(ctx context.Context, cutoffDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&entity.Certification{}).
		Where("expiry_date >= CURRENT_DATE AND expiry_date <= ? AND status = ?", 
			cutoffDate, entity.CertStatusValid).
		Update("status", entity.CertStatusExpiringSoon)
	return result.RowsAffected, result.Error
}
