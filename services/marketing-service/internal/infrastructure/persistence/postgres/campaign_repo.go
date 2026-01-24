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

// CampaignRepository implements repository.CampaignRepository
type CampaignRepository struct {
	db *gorm.DB
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

func (r *CampaignRepository) Create(ctx context.Context, campaign *entity.Campaign) error {
	return r.db.WithContext(ctx).Create(campaign).Error
}

func (r *CampaignRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Campaign, error) {
	var campaign entity.Campaign
	err := r.db.WithContext(ctx).First(&campaign, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *CampaignRepository) GetByCode(ctx context.Context, code string) (*entity.Campaign, error) {
	var campaign entity.Campaign
	err := r.db.WithContext(ctx).First(&campaign, "campaign_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *CampaignRepository) List(ctx context.Context, filter *repository.CampaignFilter) ([]*entity.Campaign, int64, error) {
	var campaigns []*entity.Campaign
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Campaign{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR campaign_code ILIKE ?", search, search)
	}
	if filter.CampaignType != "" {
		query = query.Where("campaign_type = ?", filter.CampaignType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		query = query.Where("start_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("end_date <= ?", filter.DateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if filter.Page < 1 {
		offset = 0
	}

	err := query.Offset(offset).Limit(filter.Limit).Order("created_at DESC").Find(&campaigns).Error
	return campaigns, total, err
}

func (r *CampaignRepository) Update(ctx context.Context, campaign *entity.Campaign) error {
	campaign.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(campaign).Error
}

func (r *CampaignRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Campaign{}, "id = ?", id).Error
}

func (r *CampaignRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.CampaignStatus) error {
	return r.db.WithContext(ctx).Model(&entity.Campaign{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *CampaignRepository) UpdatePerformance(ctx context.Context, id uuid.UUID, impressions, reach, engagement, conversions int, revenue float64) error {
	return r.db.WithContext(ctx).Model(&entity.Campaign{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"impressions":       impressions,
			"reach":             reach,
			"engagement":        engagement,
			"conversions":       conversions,
			"revenue_generated": revenue,
			"updated_at":        time.Now(),
		}).Error
}

func (r *CampaignRepository) IncrementSpent(ctx context.Context, id uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&entity.Campaign{}).Where("id = ?", id).
		UpdateColumn("spent", gorm.Expr("spent + ?", amount)).Error
}

func (r *CampaignRepository) GenerateCampaignCode(ctx context.Context, prefix string) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).Model(&entity.Campaign{}).
		Where("campaign_code LIKE ?", fmt.Sprintf("CAMP-%s-%d-%%", prefix, year)).
		Count(&count)
	return fmt.Sprintf("CAMP-%s-%d-%04d", prefix, year, count+1), nil
}

func (r *CampaignRepository) GetActiveCampaigns(ctx context.Context) ([]*entity.Campaign, error) {
	var campaigns []*entity.Campaign
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("status = ? AND start_date <= ? AND end_date >= ?", entity.CampaignStatusActive, now, now).
		Find(&campaigns).Error
	return campaigns, err
}

// KOLCollaborationRepository implements repository.KOLCollaborationRepository
type KOLCollaborationRepository struct {
	db *gorm.DB
}

// NewKOLCollaborationRepository creates a new collaboration repository
func NewKOLCollaborationRepository(db *gorm.DB) *KOLCollaborationRepository {
	return &KOLCollaborationRepository{db: db}
}

func (r *KOLCollaborationRepository) Create(ctx context.Context, collab *entity.KOLCollaboration) error {
	return r.db.WithContext(ctx).Create(collab).Error
}

func (r *KOLCollaborationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.KOLCollaboration, error) {
	var collab entity.KOLCollaboration
	err := r.db.WithContext(ctx).Preload("Campaign").Preload("KOL").First(&collab, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &collab, nil
}

func (r *KOLCollaborationRepository) List(ctx context.Context, filter *repository.KOLCollaborationFilter) ([]*entity.KOLCollaboration, int64, error) {
	var collabs []*entity.KOLCollaboration
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.KOLCollaboration{})

	if filter.CampaignID != nil {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.KOLID != nil {
		query = query.Where("kol_id = ?", filter.KOLID)
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

	err := query.Preload("Campaign").Preload("KOL").Offset(offset).Limit(filter.Limit).Order("created_at DESC").Find(&collabs).Error
	return collabs, total, err
}

func (r *KOLCollaborationRepository) Update(ctx context.Context, collab *entity.KOLCollaboration) error {
	collab.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(collab).Error
}

func (r *KOLCollaborationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.KOLCollaboration{}, "id = ?", id).Error
}

func (r *KOLCollaborationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.CollaborationStatus) error {
	return r.db.WithContext(ctx).Model(&entity.KOLCollaboration{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *KOLCollaborationRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus, amount float64) error {
	return r.db.WithContext(ctx).Model(&entity.KOLCollaboration{}).Where("id = ?", id).
		Updates(map[string]interface{}{"payment_status": status, "paid_amount": amount, "updated_at": time.Now()}).Error
}

func (r *KOLCollaborationRepository) UpdatePerformance(ctx context.Context, id uuid.UUID, impressions, engagement, reach int) error {
	return r.db.WithContext(ctx).Model(&entity.KOLCollaboration{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_impressions": impressions,
			"total_engagement":  engagement,
			"total_reach":       reach,
			"updated_at":        time.Now(),
		}).Error
}

func (r *KOLCollaborationRepository) IncrementPostCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.KOLCollaboration{}).Where("id = ?", id).
		UpdateColumn("actual_posts", gorm.Expr("actual_posts + 1")).Error
}

func (r *KOLCollaborationRepository) GetByKOL(ctx context.Context, kolID uuid.UUID) ([]*entity.KOLCollaboration, error) {
	var collabs []*entity.KOLCollaboration
	err := r.db.WithContext(ctx).Preload("Campaign").Where("kol_id = ?", kolID).Find(&collabs).Error
	return collabs, err
}

func (r *KOLCollaborationRepository) GetByKOLAndCampaign(ctx context.Context, kolID, campaignID uuid.UUID) (*entity.KOLCollaboration, error) {
	var collab entity.KOLCollaboration
	err := r.db.WithContext(ctx).First(&collab, "kol_id = ? AND campaign_id = ?", kolID, campaignID).Error
	if err != nil {
		return nil, err
	}
	return &collab, nil
}

func (r *KOLCollaborationRepository) GenerateCode(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).Model(&entity.KOLCollaboration{}).
		Where("collaboration_code LIKE ?", fmt.Sprintf("COLLAB-%d-%%", year)).
		Count(&count)
	return fmt.Sprintf("COLLAB-%d-%04d", year, count+1), nil
}
