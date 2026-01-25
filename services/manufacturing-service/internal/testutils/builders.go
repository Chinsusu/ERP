package testutils

import (
	"time"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// BOMBuilder helps creating BOM entities for testing
type BOMBuilder struct {
	bom *entity.BOM
}

func NewBOMBuilder() *BOMBuilder {
	return &BOMBuilder{
		bom: &entity.BOM{
			ID:                   uuid.New(),
			BOMNumber:            "BOM-" + uuid.New().String()[:8],
			ProductID:            uuid.New(),
			Version:              1,
			Name:                 "Test BOM",
			Status:               entity.BOMStatusApproved,
			BatchSize:            100,
			ConfidentialityLevel: entity.ConfidentialityInternal,
		},
	}
}

func (b *BOMBuilder) WithProductID(id uuid.UUID) *BOMBuilder {
	b.bom.ProductID = id
	return b
}

func (b *BOMBuilder) WithStatus(status entity.BOMStatus) *BOMBuilder {
	b.bom.Status = status
	return b
}

func (b *BOMBuilder) WithItems(items []entity.BOMLineItem) *BOMBuilder {
	b.bom.Items = items
	return b
}

func (b *BOMBuilder) Build() *entity.BOM {
	return b.bom
}

// WorkOrderBuilder helps creating WorkOrder entities for testing
type WorkOrderBuilder struct {
	wo *entity.WorkOrder
}

func NewWorkOrderBuilder() *WorkOrderBuilder {
	return &WorkOrderBuilder{
		wo: &entity.WorkOrder{
			ID:              uuid.New(),
			WONumber:        "WO-" + uuid.New().String()[:8],
			WODate:          time.Now(),
			ProductID:       uuid.New(),
			BOMID:           uuid.New(),
			Status:          entity.WOStatusPlanned,
			Priority:        entity.WOPriorityNormal,
			PlannedQuantity: 100,
			UOMID:           uuid.New(),
		},
	}
}

func (b *WorkOrderBuilder) WithBOM(bomID uuid.UUID) *WorkOrderBuilder {
	b.wo.BOMID = bomID
	return b
}

func (b *WorkOrderBuilder) WithStatus(status entity.WOStatus) *WorkOrderBuilder {
	b.wo.Status = status
	return b
}

func (b *WorkOrderBuilder) WithItems(items []entity.WOLineItem) *WorkOrderBuilder {
	b.wo.Items = items
	return b
}

func (b *WorkOrderBuilder) Build() *entity.WorkOrder {
	return b.wo
}

// QCInspectionBuilder helps creating QC entities for testing
type QCInspectionBuilder struct {
	ins *entity.QCInspection
}

func NewQCInspectionBuilder() *QCInspectionBuilder {
	return &QCInspectionBuilder{
		ins: &entity.QCInspection{
			ID:               uuid.New(),
			InspectionNumber: "QC-" + uuid.New().String()[:8],
			Result:           entity.InspectionResultPending,
		},
	}
}

func (b *QCInspectionBuilder) Build() *entity.QCInspection {
	return b.ins
}
