package repository

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// BOMRepository defines BOM repository interface
type BOMRepository interface {
	Create(ctx context.Context, bom *entity.BOM) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.BOM, error)
	GetByNumber(ctx context.Context, bomNumber string) (*entity.BOM, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.BOM, error)
	GetActiveBOMForProduct(ctx context.Context, productID uuid.UUID) (*entity.BOM, error)
	List(ctx context.Context, filter BOMFilter) ([]*entity.BOM, int64, error)
	Update(ctx context.Context, bom *entity.BOM) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Line items
	CreateLineItem(ctx context.Context, item *entity.BOMLineItem) error
	GetLineItems(ctx context.Context, bomID uuid.UUID) ([]*entity.BOMLineItem, error)
	UpdateLineItem(ctx context.Context, item *entity.BOMLineItem) error
	DeleteLineItem(ctx context.Context, id uuid.UUID) error
	
	// Versioning
	CreateVersion(ctx context.Context, version *entity.BOMVersion) error
	GetVersions(ctx context.Context, bomID uuid.UUID) ([]*entity.BOMVersion, error)
}

// BOMFilter for filtering BOMs
type BOMFilter struct {
	ProductID *uuid.UUID
	Status    *entity.BOMStatus
	Search    string
	Page      int
	PageSize  int
}

// WorkOrderRepository defines work order repository interface
type WorkOrderRepository interface {
	Create(ctx context.Context, wo *entity.WorkOrder) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WorkOrder, error)
	GetByNumber(ctx context.Context, woNumber string) (*entity.WorkOrder, error)
	List(ctx context.Context, filter WOFilter) ([]*entity.WorkOrder, int64, error)
	Update(ctx context.Context, wo *entity.WorkOrder) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Line items
	CreateLineItems(ctx context.Context, items []*entity.WOLineItem) error
	GetLineItems(ctx context.Context, woID uuid.UUID) ([]*entity.WOLineItem, error)
	UpdateLineItem(ctx context.Context, item *entity.WOLineItem) error
	
	// Material issues
	CreateMaterialIssue(ctx context.Context, issue *entity.WOMaterialIssue) error
	GetMaterialIssues(ctx context.Context, woID uuid.UUID) ([]*entity.WOMaterialIssue, error)
	
	// Number generation
	GenerateWONumber(ctx context.Context) (string, error)
	GenerateIssueNumber(ctx context.Context) (string, error)
}

// WOFilter for filtering work orders
type WOFilter struct {
	ProductID  *uuid.UUID
	BOMID      *uuid.UUID
	Status     *entity.WOStatus
	Priority   *entity.WOPriority
	DateFrom   *string
	DateTo     *string
	Search     string
	Page       int
	PageSize   int
}

// QCRepository defines QC repository interface
type QCRepository interface {
	// Checkpoints
	GetCheckpoints(ctx context.Context) ([]*entity.QCCheckpoint, error)
	GetCheckpointByID(ctx context.Context, id uuid.UUID) (*entity.QCCheckpoint, error)
	GetCheckpointsByType(ctx context.Context, cpType entity.CheckpointType) ([]*entity.QCCheckpoint, error)
	
	// Inspections
	CreateInspection(ctx context.Context, inspection *entity.QCInspection) error
	GetInspectionByID(ctx context.Context, id uuid.UUID) (*entity.QCInspection, error)
	GetInspectionByNumber(ctx context.Context, number string) (*entity.QCInspection, error)
	ListInspections(ctx context.Context, filter QCFilter) ([]*entity.QCInspection, int64, error)
	UpdateInspection(ctx context.Context, inspection *entity.QCInspection) error
	
	// Inspection items
	CreateInspectionItems(ctx context.Context, items []*entity.QCInspectionItem) error
	GetInspectionItems(ctx context.Context, inspectionID uuid.UUID) ([]*entity.QCInspectionItem, error)
	
	// Number generation
	GenerateInspectionNumber(ctx context.Context) (string, error)
}

// QCFilter for filtering inspections
type QCFilter struct {
	InspectionType *entity.CheckpointType
	Result         *entity.InspectionResult
	ReferenceType  *entity.ReferenceType
	ReferenceID    *uuid.UUID
	DateFrom       *string
	DateTo         *string
	Page           int
	PageSize       int
}

// NCRRepository defines NCR repository interface
type NCRRepository interface {
	Create(ctx context.Context, ncr *entity.NCR) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.NCR, error)
	GetByNumber(ctx context.Context, ncrNumber string) (*entity.NCR, error)
	List(ctx context.Context, filter NCRFilter) ([]*entity.NCR, int64, error)
	Update(ctx context.Context, ncr *entity.NCR) error
	
	// Number generation
	GenerateNCRNumber(ctx context.Context) (string, error)
}

// NCRFilter for filtering NCRs
type NCRFilter struct {
	Status   *entity.NCRStatus
	Severity *entity.NCRSeverity
	NCType   *entity.NCType
	DateFrom *string
	DateTo   *string
	Page     int
	PageSize int
}

// TraceabilityRepository defines traceability repository interface
type TraceabilityRepository interface {
	Create(ctx context.Context, trace *entity.BatchTraceability) error
	CreateBatch(ctx context.Context, traces []*entity.BatchTraceability) error
	
	// Backward trace: product lot → material lots
	GetByProductLot(ctx context.Context, productLotID uuid.UUID) ([]*entity.BatchTraceability, error)
	GetByWorkOrder(ctx context.Context, woID uuid.UUID) ([]*entity.BatchTraceability, error)
	
	// Forward trace: material lot → product lots
	GetByMaterialLot(ctx context.Context, materialLotID uuid.UUID) ([]*entity.BatchTraceability, error)
	
	// Update product lot after WO completion
	UpdateProductLot(ctx context.Context, woID uuid.UUID, productLotID uuid.UUID, productLotNumber string) error
}
