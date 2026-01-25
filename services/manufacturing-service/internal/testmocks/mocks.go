package testmocks

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/infrastructure/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockBOMRepository
type MockBOMRepository struct {
	mock.Mock
}

func (m *MockBOMRepository) Create(ctx context.Context, bom *entity.BOM) error {
	args := m.Called(ctx, bom)
	return args.Error(0)
}

func (m *MockBOMRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.BOM, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BOM), args.Error(1)
}

func (m *MockBOMRepository) GetByNumber(ctx context.Context, num string) (*entity.BOM, error) {
	args := m.Called(ctx, num)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BOM), args.Error(1)
}

func (m *MockBOMRepository) GetActiveBOMForProduct(ctx context.Context, productID uuid.UUID) (*entity.BOM, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BOM), args.Error(1)
}

func (m *MockBOMRepository) Update(ctx context.Context, bom *entity.BOM) error {
	args := m.Called(ctx, bom)
	return args.Error(0)
}

func (m *MockBOMRepository) List(ctx context.Context, filter repository.BOMFilter) ([]*entity.BOM, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.BOM), args.Get(1).(int64), args.Error(2)
}

func (m *MockBOMRepository) GetByProductID(ctx context.Context, id uuid.UUID) ([]*entity.BOM, error) { return nil, nil }
func (m *MockBOMRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockBOMRepository) CreateLineItem(ctx context.Context, item *entity.BOMLineItem) error { return nil }
func (m *MockBOMRepository) GetLineItems(ctx context.Context, bomID uuid.UUID) ([]*entity.BOMLineItem, error) { return nil, nil }
func (m *MockBOMRepository) UpdateLineItem(ctx context.Context, item *entity.BOMLineItem) error { return nil }
func (m *MockBOMRepository) DeleteLineItem(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockBOMRepository) CreateVersion(ctx context.Context, v *entity.BOMVersion) error { return nil }
func (m *MockBOMRepository) GetVersions(ctx context.Context, id uuid.UUID) ([]*entity.BOMVersion, error) { return nil, nil }

// MockWorkOrderRepository
type MockWorkOrderRepository struct {
	mock.Mock
}

func (m *MockWorkOrderRepository) Create(ctx context.Context, wo *entity.WorkOrder) error {
	args := m.Called(ctx, wo)
	return args.Error(0)
}

func (m *MockWorkOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.WorkOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WorkOrder), args.Error(1)
}

func (m *MockWorkOrderRepository) Update(ctx context.Context, wo *entity.WorkOrder) error {
	args := m.Called(ctx, wo)
	return args.Error(0)
}

func (m *MockWorkOrderRepository) GenerateWONumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockWorkOrderRepository) CreateLineItems(ctx context.Context, items []*entity.WOLineItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockWorkOrderRepository) GetByNumber(ctx context.Context, num string) (*entity.WorkOrder, error) { return nil, nil }
func (m *MockWorkOrderRepository) List(ctx context.Context, filter repository.WOFilter) ([]*entity.WorkOrder, int64, error) { return nil, 0, nil }
func (m *MockWorkOrderRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockWorkOrderRepository) GetLineItems(ctx context.Context, woID uuid.UUID) ([]*entity.WOLineItem, error) { return nil, nil }
func (m *MockWorkOrderRepository) UpdateLineItem(ctx context.Context, item *entity.WOLineItem) error { return nil }
func (m *MockWorkOrderRepository) CreateMaterialIssue(ctx context.Context, issue *entity.WOMaterialIssue) error { return nil }
func (m *MockWorkOrderRepository) GetMaterialIssues(ctx context.Context, woID uuid.UUID) ([]*entity.WOMaterialIssue, error) { return nil, nil }
func (m *MockWorkOrderRepository) GenerateIssueNumber(ctx context.Context) (string, error) { return "", nil }

// MockTraceabilityRepository
type MockTraceabilityRepository struct {
	mock.Mock
}

func (m *MockTraceabilityRepository) CreateBatch(ctx context.Context, traces []*entity.BatchTraceability) error {
	args := m.Called(ctx, traces)
	return args.Error(0)
}

func (m *MockTraceabilityRepository) GetByProductLot(ctx context.Context, lotID uuid.UUID) ([]*entity.BatchTraceability, error) {
	args := m.Called(ctx, lotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BatchTraceability), args.Error(1)
}

func (m *MockTraceabilityRepository) GetByMaterialLot(ctx context.Context, lotID uuid.UUID) ([]*entity.BatchTraceability, error) {
	args := m.Called(ctx, lotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BatchTraceability), args.Error(1)
}

func (m *MockTraceabilityRepository) Create(ctx context.Context, trace *entity.BatchTraceability) error { return nil }
func (m *MockTraceabilityRepository) GetByWorkOrder(ctx context.Context, woID uuid.UUID) ([]*entity.BatchTraceability, error) { return nil, nil }
func (m *MockTraceabilityRepository) UpdateProductLot(ctx context.Context, woID uuid.UUID, lotID uuid.UUID, lotNum string) error { return nil }

// MockQCRepository
type MockQCRepository struct {
	mock.Mock
}

func (m *MockQCRepository) CreateInspection(ctx context.Context, ins *entity.QCInspection) error {
	args := m.Called(ctx, ins)
	return args.Error(0)
}

func (m *MockQCRepository) GenerateInspectionNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockQCRepository) GetCheckpointByID(ctx context.Context, id uuid.UUID) (*entity.QCCheckpoint, error) { return nil, nil }
func (m *MockQCRepository) GetCheckpoints(ctx context.Context) ([]*entity.QCCheckpoint, error) { return nil, nil }
func (m *MockQCRepository) GetCheckpointsByType(ctx context.Context, cpType entity.CheckpointType) ([]*entity.QCCheckpoint, error) {
	args := m.Called(ctx, cpType)
	return args.Get(0).([]*entity.QCCheckpoint), args.Error(1)
}
func (m *MockQCRepository) GetInspectionByID(ctx context.Context, id uuid.UUID) (*entity.QCInspection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.QCInspection), args.Error(1)
}
func (m *MockQCRepository) GetInspectionByNumber(ctx context.Context, num string) (*entity.QCInspection, error) { return nil, nil }
func (m *MockQCRepository) ListInspections(ctx context.Context, filter repository.QCFilter) ([]*entity.QCInspection, int64, error) { return nil, 0, nil }
func (m *MockQCRepository) UpdateInspection(ctx context.Context, ins *entity.QCInspection) error {
	args := m.Called(ctx, ins)
	return args.Error(0)
}
func (m *MockQCRepository) CreateInspectionItems(ctx context.Context, items []*entity.QCInspectionItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}
func (m *MockQCRepository) GetInspectionItems(ctx context.Context, id uuid.UUID) ([]*entity.QCInspectionItem, error) { return nil, nil }

// MockNCRRepository
type MockNCRRepository struct {
	mock.Mock
}

func (m *MockNCRRepository) Create(ctx context.Context, ncr *entity.NCR) error {
	args := m.Called(ctx, ncr)
	return args.Error(0)
}

func (m *MockNCRRepository) GenerateNCRNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockNCRRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.NCR, error) { return nil, nil }
func (m *MockNCRRepository) GetByNumber(ctx context.Context, num string) (*entity.NCR, error) { return nil, nil }
func (m *MockNCRRepository) List(ctx context.Context, filter repository.NCRFilter) ([]*entity.NCR, int64, error) { return nil, 0, nil }
func (m *MockNCRRepository) Update(ctx context.Context, ncr *entity.NCR) error { return nil }

// MockEventPublisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) PublishBOMCreated(e event.BOMEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishBOMApproved(e event.BOMEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishWOCreated(e event.WOEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishWOReleased(e event.WOEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishWOStarted(e event.WOEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishWOCompleted(e event.WOCompletedEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishQCPassed(e event.QCEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishQCFailed(e event.QCEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishNCRCreated(e event.NCREvent) error {
	args := m.Called(e)
	return args.Error(0)
}
