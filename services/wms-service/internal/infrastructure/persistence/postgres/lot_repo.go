package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type lotRepository struct {
	db *gorm.DB
}

// NewLotRepository creates a new lot repository
func NewLotRepository(db *gorm.DB) repository.LotRepository {
	return &lotRepository{db: db}
}

func (r *lotRepository) Create(ctx context.Context, lot *entity.Lot) error {
	return r.db.WithContext(ctx).Create(lot).Error
}

func (r *lotRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lot, error) {
	var lot entity.Lot
	err := r.db.WithContext(ctx).First(&lot, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lot, nil
}

func (r *lotRepository) GetByLotNumber(ctx context.Context, lotNumber string) (*entity.Lot, error) {
	var lot entity.Lot
	err := r.db.WithContext(ctx).First(&lot, "lot_number = ?", lotNumber).Error
	if err != nil {
		return nil, err
	}
	return &lot, nil
}

func (r *lotRepository) List(ctx context.Context, filter *repository.LotFilter) ([]*entity.Lot, int64, error) {
	var lots []*entity.Lot
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Lot{})

	if filter.MaterialID != nil {
		query = query.Where("material_id = ?", *filter.MaterialID)
	}
	if filter.SupplierID != nil {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.QCStatus != "" {
		query = query.Where("qc_status = ?", filter.QCStatus)
	}
	if filter.ExpiryBefore != nil {
		query = query.Where("expiry_date <= ?", *filter.ExpiryBefore)
	}
	if filter.ExpiryAfter != nil {
		query = query.Where("expiry_date >= ?", *filter.ExpiryAfter)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("lot_number ILIKE ? OR supplier_lot_number ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(20)
	}
	if filter.Page > 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if err := query.Order("expiry_date ASC").Find(&lots).Error; err != nil {
		return nil, 0, err
	}

	return lots, total, nil
}

func (r *lotRepository) Update(ctx context.Context, lot *entity.Lot) error {
	lot.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(lot).Error
}

// GetAvailableLots returns lots sorted by expiry date (FEFO)
func (r *lotRepository) GetAvailableLots(ctx context.Context, materialID uuid.UUID) ([]*entity.Lot, error) {
	var lots []*entity.Lot
	err := r.db.WithContext(ctx).
		Where("material_id = ?", materialID).
		Where("status = ?", entity.LotStatusAvailable).
		Where("qc_status = ?", entity.QCStatusPassed).
		Where("expiry_date > ?", time.Now()).
		Order("expiry_date ASC"). // FEFO: earliest expiry first
		Find(&lots).Error
	return lots, err
}

// GetExpiringLots returns lots expiring within the specified days
func (r *lotRepository) GetExpiringLots(ctx context.Context, days int) ([]*entity.Lot, error) {
	var lots []*entity.Lot
	threshold := time.Now().AddDate(0, 0, days)
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.LotStatusAvailable).
		Where("expiry_date <= ?", threshold).
		Where("expiry_date > ?", time.Now()).
		Order("expiry_date ASC").
		Find(&lots).Error
	return lots, err
}

// GetExpiredLots returns all expired lots that are still marked as available
func (r *lotRepository) GetExpiredLots(ctx context.Context) ([]*entity.Lot, error) {
	var lots []*entity.Lot
	err := r.db.WithContext(ctx).
		Where("status = ?", entity.LotStatusAvailable).
		Where("expiry_date <= ?", time.Now()).
		Find(&lots).Error
	return lots, err
}

// GetNextLotNumber generates the next lot number
func (r *lotRepository) GetNextLotNumber(ctx context.Context) (string, error) {
	var count int64
	yearMonth := time.Now().Format("200601")

	r.db.WithContext(ctx).
		Model(&entity.Lot{}).
		Where("lot_number LIKE ?", fmt.Sprintf("LOT-%s-%%", yearMonth)).
		Count(&count)

	return fmt.Sprintf("LOT-%s-%04d", yearMonth, count+1), nil
}

// MarkExpired marks lots as expired
func (r *lotRepository) MarkExpired(ctx context.Context, lotIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Lot{}).
		Where("id IN ?", lotIDs).
		Updates(map[string]interface{}{
			"status":     entity.LotStatusExpired,
			"updated_at": time.Now(),
		}).Error
}
