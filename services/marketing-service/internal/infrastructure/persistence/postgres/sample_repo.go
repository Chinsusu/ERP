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

// SampleRequestRepository implements repository.SampleRequestRepository
type SampleRequestRepository struct {
	db *gorm.DB
}

// NewSampleRequestRepository creates a new sample request repository
func NewSampleRequestRepository(db *gorm.DB) *SampleRequestRepository {
	return &SampleRequestRepository{db: db}
}

func (r *SampleRequestRepository) Create(ctx context.Context, request *entity.SampleRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *SampleRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SampleRequest, error) {
	var request entity.SampleRequest
	err := r.db.WithContext(ctx).Preload("KOL").Preload("Campaign").Preload("Items").First(&request, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *SampleRequestRepository) GetByNumber(ctx context.Context, number string) (*entity.SampleRequest, error) {
	var request entity.SampleRequest
	err := r.db.WithContext(ctx).Preload("KOL").Preload("Items").First(&request, "request_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *SampleRequestRepository) List(ctx context.Context, filter *repository.SampleRequestFilter) ([]*entity.SampleRequest, int64, error) {
	var requests []*entity.SampleRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SampleRequest{})

	if filter.KOLID != nil {
		query = query.Where("kol_id = ?", filter.KOLID)
	}
	if filter.CampaignID != nil {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != "" {
		query = query.Where("request_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("request_date <= ?", filter.DateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if filter.Page < 1 {
		offset = 0
	}

	err := query.Preload("KOL").Preload("Campaign").Offset(offset).Limit(filter.Limit).Order("created_at DESC").Find(&requests).Error
	return requests, total, err
}

func (r *SampleRequestRepository) Update(ctx context.Context, request *entity.SampleRequest) error {
	request.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(request).Error
}

func (r *SampleRequestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SampleRequest{}, "id = ?", id).Error
}

func (r *SampleRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SampleRequestStatus) error {
	return r.db.WithContext(ctx).Model(&entity.SampleRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *SampleRequestRepository) Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.SampleRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      entity.SampleStatusApproved,
			"approved_by": approverID,
			"approved_at": now,
			"updated_at":  now,
		}).Error
}

func (r *SampleRequestRepository) Reject(ctx context.Context, id uuid.UUID, reason string) error {
	return r.db.WithContext(ctx).Model(&entity.SampleRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           entity.SampleStatusRejected,
			"rejection_reason": reason,
			"updated_at":       time.Now(),
		}).Error
}

func (r *SampleRequestRepository) AddItem(ctx context.Context, item *entity.SampleItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *SampleRequestRepository) GetItems(ctx context.Context, requestID uuid.UUID) ([]*entity.SampleItem, error) {
	var items []*entity.SampleItem
	err := r.db.WithContext(ctx).Where("sample_request_id = ?", requestID).Order("line_number ASC").Find(&items).Error
	return items, err
}

func (r *SampleRequestRepository) UpdateItem(ctx context.Context, item *entity.SampleItem) error {
	item.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *SampleRequestRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SampleItem{}, "id = ?", id).Error
}

func (r *SampleRequestRepository) GenerateRequestNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).Model(&entity.SampleRequest{}).
		Where("request_number LIKE ?", fmt.Sprintf("SR-%d-%%", year)).
		Count(&count)
	return fmt.Sprintf("SR-%d-%04d", year, count+1), nil
}

func (r *SampleRequestRepository) GetByKOL(ctx context.Context, kolID uuid.UUID) ([]*entity.SampleRequest, error) {
	var requests []*entity.SampleRequest
	err := r.db.WithContext(ctx).Preload("Items").Where("kol_id = ?", kolID).Find(&requests).Error
	return requests, err
}

func (r *SampleRequestRepository) GetPendingApproval(ctx context.Context) ([]*entity.SampleRequest, error) {
	var requests []*entity.SampleRequest
	err := r.db.WithContext(ctx).Preload("KOL").Preload("Items").
		Where("status = ?", entity.SampleStatusPendingApproval).
		Order("created_at ASC").Find(&requests).Error
	return requests, err
}

// SampleShipmentRepository implements repository.SampleShipmentRepository
type SampleShipmentRepository struct {
	db *gorm.DB
}

// NewSampleShipmentRepository creates a new sample shipment repository
func NewSampleShipmentRepository(db *gorm.DB) *SampleShipmentRepository {
	return &SampleShipmentRepository{db: db}
}

func (r *SampleShipmentRepository) Create(ctx context.Context, shipment *entity.SampleShipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

func (r *SampleShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SampleShipment, error) {
	var shipment entity.SampleShipment
	err := r.db.WithContext(ctx).Preload("SampleRequest").First(&shipment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *SampleShipmentRepository) GetByNumber(ctx context.Context, number string) (*entity.SampleShipment, error) {
	var shipment entity.SampleShipment
	err := r.db.WithContext(ctx).First(&shipment, "shipment_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *SampleShipmentRepository) GetByRequest(ctx context.Context, requestID uuid.UUID) (*entity.SampleShipment, error) {
	var shipment entity.SampleShipment
	err := r.db.WithContext(ctx).First(&shipment, "sample_request_id = ?", requestID).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *SampleShipmentRepository) List(ctx context.Context, filter *repository.SampleShipmentFilter) ([]*entity.SampleShipment, int64, error) {
	var shipments []*entity.SampleShipment
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SampleShipment{})

	if filter.SampleRequestID != nil {
		query = query.Where("sample_request_id = ?", filter.SampleRequestID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Courier != "" {
		query = query.Where("courier = ?", filter.Courier)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if filter.Page < 1 {
		offset = 0
	}

	err := query.Preload("SampleRequest").Offset(offset).Limit(filter.Limit).Order("created_at DESC").Find(&shipments).Error
	return shipments, total, err
}

func (r *SampleShipmentRepository) Update(ctx context.Context, shipment *entity.SampleShipment) error {
	shipment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(shipment).Error
}

func (r *SampleShipmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ShipmentStatus) error {
	return r.db.WithContext(ctx).Model(&entity.SampleShipment{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *SampleShipmentRepository) MarkDelivered(ctx context.Context, id uuid.UUID, proofURL string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.SampleShipment{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            entity.ShipmentStatusDelivered,
			"actual_delivery":   now,
			"proof_of_delivery": proofURL,
			"updated_at":        now,
		}).Error
}

func (r *SampleShipmentRepository) GenerateShipmentNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int64
	r.db.WithContext(ctx).Model(&entity.SampleShipment{}).
		Where("shipment_number LIKE ?", fmt.Sprintf("SHP-%d-%%", year)).
		Count(&count)
	return fmt.Sprintf("SHP-%d-%04d", year, count+1), nil
}

// KOLPostRepository implements repository.KOLPostRepository
type KOLPostRepository struct {
	db *gorm.DB
}

// NewKOLPostRepository creates a new KOL post repository
func NewKOLPostRepository(db *gorm.DB) *KOLPostRepository {
	return &KOLPostRepository{db: db}
}

func (r *KOLPostRepository) Create(ctx context.Context, post *entity.KOLPost) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *KOLPostRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.KOLPost, error) {
	var post entity.KOLPost
	err := r.db.WithContext(ctx).Preload("KOL").First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *KOLPostRepository) List(ctx context.Context, filter *repository.KOLPostFilter) ([]*entity.KOLPost, int64, error) {
	var posts []*entity.KOLPost
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.KOLPost{})

	if filter.KOLID != nil {
		query = query.Where("kol_id = ?", filter.KOLID)
	}
	if filter.CampaignID != nil {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.CollaborationID != nil {
		query = query.Where("collaboration_id = ?", filter.CollaborationID)
	}
	if filter.Platform != "" {
		query = query.Where("platform = ?", filter.Platform)
	}
	if filter.Verified != nil {
		query = query.Where("verified = ?", *filter.Verified)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if filter.Page < 1 {
		offset = 0
	}

	err := query.Preload("KOL").Offset(offset).Limit(filter.Limit).Order("post_date DESC").Find(&posts).Error
	return posts, total, err
}

func (r *KOLPostRepository) Update(ctx context.Context, post *entity.KOLPost) error {
	post.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *KOLPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.KOLPost{}, "id = ?", id).Error
}

func (r *KOLPostRepository) GetByKOL(ctx context.Context, kolID uuid.UUID) ([]*entity.KOLPost, error) {
	var posts []*entity.KOLPost
	err := r.db.WithContext(ctx).Where("kol_id = ?", kolID).Order("post_date DESC").Find(&posts).Error
	return posts, err
}

func (r *KOLPostRepository) Verify(ctx context.Context, id uuid.UUID, verifierID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.KOLPost{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"verified":    true,
			"verified_by": verifierID,
			"verified_at": now,
			"updated_at":  now,
		}).Error
}

func (r *KOLPostRepository) UpdateMetrics(ctx context.Context, id uuid.UUID, likes, comments, shares, views, reach int) error {
	return r.db.WithContext(ctx).Model(&entity.KOLPost{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"likes":      likes,
			"comments":   comments,
			"shares":     shares,
			"views":      views,
			"reach":      reach,
			"updated_at": time.Now(),
		}).Error
}
