package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KOLTierRepository implements repository.KOLTierRepository
type KOLTierRepository struct {
	db *gorm.DB
}

// NewKOLTierRepository creates a new KOL tier repository
func NewKOLTierRepository(db *gorm.DB) *KOLTierRepository {
	return &KOLTierRepository{db: db}
}

func (r *KOLTierRepository) Create(ctx context.Context, tier *entity.KOLTier) error {
	return r.db.WithContext(ctx).Create(tier).Error
}

func (r *KOLTierRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.KOLTier, error) {
	var tier entity.KOLTier
	err := r.db.WithContext(ctx).First(&tier, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *KOLTierRepository) GetByCode(ctx context.Context, code string) (*entity.KOLTier, error) {
	var tier entity.KOLTier
	err := r.db.WithContext(ctx).First(&tier, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *KOLTierRepository) List(ctx context.Context, activeOnly bool) ([]*entity.KOLTier, error) {
	var tiers []*entity.KOLTier
	query := r.db.WithContext(ctx).Order("priority ASC")
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&tiers).Error
	return tiers, err
}

func (r *KOLTierRepository) Update(ctx context.Context, tier *entity.KOLTier) error {
	tier.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(tier).Error
}

func (r *KOLTierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.KOLTier{}, "id = ?", id).Error
}

// KOLRepository implements repository.KOLRepository
type KOLRepository struct {
	db *gorm.DB
}

// NewKOLRepository creates a new KOL repository
func NewKOLRepository(db *gorm.DB) *KOLRepository {
	return &KOLRepository{db: db}
}

func (r *KOLRepository) Create(ctx context.Context, kol *entity.KOL) error {
	return r.db.WithContext(ctx).Create(kol).Error
}

func (r *KOLRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.KOL, error) {
	var kol entity.KOL
	err := r.db.WithContext(ctx).Preload("Tier").First(&kol, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &kol, nil
}

func (r *KOLRepository) GetByCode(ctx context.Context, code string) (*entity.KOL, error) {
	var kol entity.KOL
	err := r.db.WithContext(ctx).Preload("Tier").First(&kol, "kol_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &kol, nil
}

func (r *KOLRepository) List(ctx context.Context, filter *repository.KOLFilter) ([]*entity.KOL, int64, error) {
	var kols []*entity.KOL
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.KOL{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR kol_code ILIKE ? OR instagram_handle ILIKE ?", search, search, search)
	}
	if filter.TierID != nil {
		query = query.Where("tier_id = ?", filter.TierID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Niche != "" {
		query = query.Where("niche = ?", filter.Niche)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if filter.Page < 1 {
		offset = 0
	}

	err := query.Preload("Tier").Offset(offset).Limit(filter.Limit).Order("name ASC").Find(&kols).Error
	return kols, total, err
}

func (r *KOLRepository) Update(ctx context.Context, kol *entity.KOL) error {
	kol.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(kol).Error
}

func (r *KOLRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.KOL{}, "id = ?", id).Error
}

func (r *KOLRepository) IncrementPostCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.KOL{}).Where("id = ?", id).
		UpdateColumn("total_posts", gorm.Expr("total_posts + 1")).Error
}

func (r *KOLRepository) IncrementSampleCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.KOL{}).Where("id = ?", id).
		UpdateColumn("total_samples_received", gorm.Expr("total_samples_received + 1")).Error
}

func (r *KOLRepository) IncrementCollaborationCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.KOL{}).Where("id = ?", id).
		UpdateColumn("total_collaborations", gorm.Expr("total_collaborations + 1")).Error
}

func (r *KOLRepository) UpdateLastCollaborationDate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.KOL{}).Where("id = ?", id).
		Update("last_collaboration_date", time.Now()).Error
}

func (r *KOLRepository) GenerateKOLCode(ctx context.Context) (string, error) {
	var count int64
	r.db.WithContext(ctx).Model(&entity.KOL{}).Count(&count)
	return fmt.Sprintf("KOL-%04d", count+1), nil
}
