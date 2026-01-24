package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type qcRepository struct {
	db *gorm.DB
}

// NewQCRepository creates a new QC repository
func NewQCRepository(db *gorm.DB) repository.QCRepository {
	return &qcRepository{db: db}
}

// Checkpoints
func (r *qcRepository) GetCheckpoints(ctx context.Context) ([]*entity.QCCheckpoint, error) {
	var checkpoints []*entity.QCCheckpoint
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("checkpoint_type, code").
		Find(&checkpoints).Error
	return checkpoints, err
}

func (r *qcRepository) GetCheckpointByID(ctx context.Context, id uuid.UUID) (*entity.QCCheckpoint, error) {
	var checkpoint entity.QCCheckpoint
	err := r.db.WithContext(ctx).First(&checkpoint, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &checkpoint, nil
}

func (r *qcRepository) GetCheckpointsByType(ctx context.Context, cpType entity.CheckpointType) ([]*entity.QCCheckpoint, error) {
	var checkpoints []*entity.QCCheckpoint
	err := r.db.WithContext(ctx).
		Where("checkpoint_type = ? AND is_active = ?", cpType, true).
		Find(&checkpoints).Error
	return checkpoints, err
}

// Inspections
func (r *qcRepository) CreateInspection(ctx context.Context, inspection *entity.QCInspection) error {
	return r.db.WithContext(ctx).Create(inspection).Error
}

func (r *qcRepository) GetInspectionByID(ctx context.Context, id uuid.UUID) (*entity.QCInspection, error) {
	var inspection entity.QCInspection
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&inspection, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &inspection, nil
}

func (r *qcRepository) GetInspectionByNumber(ctx context.Context, number string) (*entity.QCInspection, error) {
	var inspection entity.QCInspection
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&inspection, "inspection_number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &inspection, nil
}

func (r *qcRepository) ListInspections(ctx context.Context, filter repository.QCFilter) ([]*entity.QCInspection, int64, error) {
	var inspections []*entity.QCInspection
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.QCInspection{})

	if filter.InspectionType != nil {
		query = query.Where("inspection_type = ?", *filter.InspectionType)
	}
	if filter.Result != nil {
		query = query.Where("result = ?", *filter.Result)
	}
	if filter.ReferenceType != nil {
		query = query.Where("reference_type = ?", *filter.ReferenceType)
	}
	if filter.ReferenceID != nil {
		query = query.Where("reference_id = ?", *filter.ReferenceID)
	}

	query.Count(&total)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Order("created_at DESC").Find(&inspections).Error
	return inspections, total, err
}

func (r *qcRepository) UpdateInspection(ctx context.Context, inspection *entity.QCInspection) error {
	return r.db.WithContext(ctx).Save(inspection).Error
}

// Inspection items
func (r *qcRepository) CreateInspectionItems(ctx context.Context, items []*entity.QCInspectionItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *qcRepository) GetInspectionItems(ctx context.Context, inspectionID uuid.UUID) ([]*entity.QCInspectionItem, error) {
	var items []*entity.QCInspectionItem
	err := r.db.WithContext(ctx).
		Where("inspection_id = ?", inspectionID).
		Order("item_number ASC").
		Find(&items).Error
	return items, err
}

func (r *qcRepository) GenerateInspectionNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	r.db.WithContext(ctx).Model(&entity.QCInspection{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count)
	return fmt.Sprintf("QC-%d-%04d", year, count+1), nil
}

// NCR Repository
type ncrRepository struct {
	db *gorm.DB
}

// NewNCRRepository creates a new NCR repository
func NewNCRRepository(db *gorm.DB) repository.NCRRepository {
	return &ncrRepository{db: db}
}

func (r *ncrRepository) Create(ctx context.Context, ncr *entity.NCR) error {
	return r.db.WithContext(ctx).Create(ncr).Error
}

func (r *ncrRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.NCR, error) {
	var ncr entity.NCR
	err := r.db.WithContext(ctx).First(&ncr, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ncr, nil
}

func (r *ncrRepository) GetByNumber(ctx context.Context, ncrNumber string) (*entity.NCR, error) {
	var ncr entity.NCR
	err := r.db.WithContext(ctx).First(&ncr, "ncr_number = ?", ncrNumber).Error
	if err != nil {
		return nil, err
	}
	return &ncr, nil
}

func (r *ncrRepository) List(ctx context.Context, filter repository.NCRFilter) ([]*entity.NCR, int64, error) {
	var ncrs []*entity.NCR
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.NCR{})

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}
	if filter.NCType != nil {
		query = query.Where("nc_type = ?", *filter.NCType)
	}

	query.Count(&total)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Order("created_at DESC").Find(&ncrs).Error
	return ncrs, total, err
}

func (r *ncrRepository) Update(ctx context.Context, ncr *entity.NCR) error {
	return r.db.WithContext(ctx).Save(ncr).Error
}

func (r *ncrRepository) GenerateNCRNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	r.db.WithContext(ctx).Model(&entity.NCR{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count)
	return fmt.Sprintf("NCR-%d-%04d", year, count+1), nil
}

// Traceability Repository
type traceabilityRepository struct {
	db *gorm.DB
}

// NewTraceabilityRepository creates a new traceability repository
func NewTraceabilityRepository(db *gorm.DB) repository.TraceabilityRepository {
	return &traceabilityRepository{db: db}
}

func (r *traceabilityRepository) Create(ctx context.Context, trace *entity.BatchTraceability) error {
	return r.db.WithContext(ctx).Create(trace).Error
}

func (r *traceabilityRepository) CreateBatch(ctx context.Context, traces []*entity.BatchTraceability) error {
	return r.db.WithContext(ctx).Create(&traces).Error
}

func (r *traceabilityRepository) GetByProductLot(ctx context.Context, productLotID uuid.UUID) ([]*entity.BatchTraceability, error) {
	var traces []*entity.BatchTraceability
	err := r.db.WithContext(ctx).
		Where("product_lot_id = ?", productLotID).
		Find(&traces).Error
	return traces, err
}

func (r *traceabilityRepository) GetByWorkOrder(ctx context.Context, woID uuid.UUID) ([]*entity.BatchTraceability, error) {
	var traces []*entity.BatchTraceability
	err := r.db.WithContext(ctx).
		Where("work_order_id = ?", woID).
		Find(&traces).Error
	return traces, err
}

func (r *traceabilityRepository) GetByMaterialLot(ctx context.Context, materialLotID uuid.UUID) ([]*entity.BatchTraceability, error) {
	var traces []*entity.BatchTraceability
	err := r.db.WithContext(ctx).
		Where("material_lot_id = ?", materialLotID).
		Find(&traces).Error
	return traces, err
}

func (r *traceabilityRepository) UpdateProductLot(ctx context.Context, woID uuid.UUID, productLotID uuid.UUID, productLotNumber string) error {
	return r.db.WithContext(ctx).
		Model(&entity.BatchTraceability{}).
		Where("work_order_id = ?", woID).
		Updates(map[string]interface{}{
			"product_lot_id":     productLotID,
			"product_lot_number": productLotNumber,
		}).Error
}
